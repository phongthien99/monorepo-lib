package core

import (
	"fmt"
	"reflect"
)

// DefaultMerge is the default merge strategy using reflection.
// Merges src into dst, only overriding non-zero values.
//
// Rules:
//   - Struct fields: merge recursively, non-zero values override
//   - Slices: override entirely if src slice is not empty
//   - Maps: deep merge keys
//   - Pointers: merge recursively if src is not nil
//   - Primitives: override if src is not zero value
//
// Example:
//
//	dst := &AppConfig{Server: ServerConfig{Host: "localhost", Port: 8080}}
//	src := &AppConfig{Server: ServerConfig{Port: 9090}}
//	DefaultMerge(dst, src)
//	// Result: dst.Server.Host = "localhost", dst.Server.Port = 9090
func DefaultMerge[T any](dst, src *T) error {
	dstVal := reflect.ValueOf(dst).Elem()
	srcVal := reflect.ValueOf(src).Elem()

	return deepMerge(dstVal, srcVal)
}

// deepMerge recursively merges src into dst using reflection.
func deepMerge(dst, src reflect.Value) error {
	if dst.Type() != src.Type() {
		return fmt.Errorf("type mismatch: %v != %v", dst.Type(), src.Type())
	}

	switch src.Kind() {
	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			srcField := src.Field(i)
			dstField := dst.Field(i)

			if !dstField.CanSet() {
				continue
			}

			if !srcField.IsZero() {
				if err := deepMerge(dstField, srcField); err != nil {
					return fmt.Errorf("field %s: %w", src.Type().Field(i).Name, err)
				}
			}
		}

	case reflect.Slice:
		if src.Len() > 0 {
			dst.Set(src)
		}

	case reflect.Map:
		if !src.IsZero() {
			if dst.IsNil() {
				dst.Set(reflect.MakeMap(src.Type()))
			}
			for _, key := range src.MapKeys() {
				srcValue := src.MapIndex(key)
				dstValue := dst.MapIndex(key)

				if dstValue.IsValid() && !dstValue.IsZero() {
					if srcValue.Kind() == reflect.Map || srcValue.Kind() == reflect.Struct {
						merged := reflect.New(srcValue.Type()).Elem()
						merged.Set(dstValue)
						if err := deepMerge(merged, srcValue); err != nil {
							return err
						}
						dst.SetMapIndex(key, merged)
					} else {
						dst.SetMapIndex(key, srcValue)
					}
				} else {
					dst.SetMapIndex(key, srcValue)
				}
			}
		}

	case reflect.Ptr:
		if !src.IsNil() {
			if dst.IsNil() {
				dst.Set(reflect.New(src.Type().Elem()))
			}
			if err := deepMerge(dst.Elem(), src.Elem()); err != nil {
				return err
			}
		}

	default:
		if !src.IsZero() {
			dst.Set(src)
		}
	}

	return nil
}

// ShallowMerge is an alternative merge strategy - overrides entire struct.
// Useful when deep merge is not needed, only full config replacement.
//
// Example:
//
//	cfg := config.New[AppConfig](loaders...).
//	    WithMerge(core.ShallowMerge[AppConfig])
func ShallowMerge[T any](dst, src *T) error {
	srcVal := reflect.ValueOf(src).Elem()
	dstVal := reflect.ValueOf(dst).Elem()

	if !srcVal.IsZero() {
		dstVal.Set(srcVal)
	}
	return nil
}
