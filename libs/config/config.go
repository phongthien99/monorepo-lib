package config

import (
	"github.com/phongthien99/monorepo-lib/libs/config/core"
)

// Config re-exports core.Config so users can use config.Config[T]
type Config[T any] = core.Config[T]

// Loader re-exports core.Loader so users can use config.Loader[T]
type Loader[T any] = core.Loader[T]

// MergeFunc re-exports core.MergeFunc so users can define custom merge functions
type MergeFunc[T any] = core.MergeFunc[T]

// Validator re-exports core.Validator so users can define custom validators
type Validator[T any] = core.Validator[T]

// ValidatorFunc re-exports core.ValidatorFunc - function adapter for Validator
type ValidatorFunc[T any] = core.ValidatorFunc[T]

// New re-exports core.New to create a new Config with default merge strategy
func New[T any](loaders ...Loader[*T]) *Config[T] {
	return core.New[T](loaders...)
}

// NewCompositeValidator re-exports core.NewCompositeValidator
func NewCompositeValidator[T any](validators ...Validator[T]) *core.CompositeValidator[T] {
	return core.NewCompositeValidator[T](validators...)
}

// DefaultMerge re-exports core.DefaultMerge - deep merge strategy
func DefaultMerge[T any](dst, src *T) error {
	return core.DefaultMerge(dst, src)
}

// ShallowMerge re-exports core.ShallowMerge - shallow merge strategy
func ShallowMerge[T any](dst, src *T) error {
	return core.ShallowMerge(dst, src)
}
