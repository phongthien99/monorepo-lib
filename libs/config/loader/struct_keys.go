package loader

import (
	"reflect"
	"strings"
)

// extractStructKeys recursively extracts all keys from a struct using mapstructure tags.
// Returns a flat list of keys in dot notation.
//
// Example:
//
//	type Config struct {
//	    Server struct {
//	        Host string `mapstructure:"host"`
//	        Port int    `mapstructure:"port"`
//	    } `mapstructure:"server"`
//	}
//
//	keys := extractStructKeys(reflect.TypeOf(Config{}), "")
//	// Returns: ["server.host", "server.port"]
func extractStructKeys(t reflect.Type, prefix string) []string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil
	}

	var keys []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if !field.IsExported() {
			continue
		}

		tag := field.Tag.Get("mapstructure")
		if tag == "" {
			tag = strings.ToLower(field.Name)
		}

		if tag == "-" {
			continue
		}

		var fullKey string
		if prefix == "" {
			fullKey = tag
		} else {
			fullKey = prefix + "." + tag
		}

		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		if fieldType.Kind() == reflect.Struct {
			nestedKeys := extractStructKeys(fieldType, fullKey)
			keys = append(keys, nestedKeys...)
		} else {
			keys = append(keys, fullKey)
		}
	}

	return keys
}

// ExtractKeysFromType extracts all config keys from a struct type.
// Accepts any type (value or pointer) and returns all keys in dot notation.
func ExtractKeysFromType(example interface{}) []string {
	t := reflect.TypeOf(example)
	return extractStructKeys(t, "")
}
