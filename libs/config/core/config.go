package core

import "fmt"

// MergeFunc defines the function signature for merge strategies.
// dst: destination (current merge result)
// src: source (new data from loader)
type MergeFunc[T any] func(dst, src *T) error

// Config manages configuration with type-safe generics and configurable merge strategy.
type Config[T any] struct {
	loaders   []Loader[*T]
	mergeFunc MergeFunc[T]
	validator Validator[T]
	data      T
}

// New creates a new Config with default merge strategy.
//
// Loaders are merged in array order:
//   - First loader has lowest priority
//   - Last loader has highest priority (will override previous loaders)
//
// Default merge strategy: DefaultMerge (deep merge with reflection)
//
// Example:
//
//	cfg := config.New[AppConfig](
//	    fileLoader,   // Lowest priority
//	    envLoader,    // Medium priority
//	    flagLoader,   // Highest priority
//	)
func New[T any](loaders ...Loader[*T]) *Config[T] {
	return &Config[T]{
		loaders:   loaders,
		mergeFunc: DefaultMerge[T],
	}
}

// WithMerge sets a custom merge function.
// Returns *Config[T] to support method chaining.
//
// Example:
//
//	cfg := config.New[AppConfig](loaders...).
//	    WithMerge(customMergeFunc)
//
//	// Or use ShallowMerge
//	cfg := config.New[AppConfig](loaders...).
//	    WithMerge(core.ShallowMerge[AppConfig])
func (c *Config[T]) WithMerge(mergeFn MergeFunc[T]) *Config[T] {
	c.mergeFunc = mergeFn
	return c
}

// WithValidator sets a validator function.
// Validator will be called after loading and merging config.
// Returns *Config[T] to support method chaining.
//
// Example:
//
//	type AppConfigValidator struct{}
//
//	func (v *AppConfigValidator) Validate(cfg *AppConfig) error {
//	    if cfg.Server.Port < 1024 {
//	        return fmt.Errorf("port must be >= 1024")
//	    }
//	    return nil
//	}
//
//	cfg := config.New[AppConfig](loaders...).
//	    WithValidator(&AppConfigValidator{})
func (c *Config[T]) WithValidator(validator Validator[T]) *Config[T] {
	c.validator = validator
	return c
}

// Load executes loading and merging of all config sources.
//
// Process:
//  1. Initialize accumulated result (zero value)
//  2. Loop through all loaders in order
//  3. Each loader fills data into temp struct
//  4. Merge temp into accumulated using merge strategy
//  5. Validate config if validator is set
//  6. Store accumulated result
//
// Returns error if:
//   - Any loader fails during Load()
//   - Merge function fails
//   - Validation fails
func (c *Config[T]) Load() error {
	accumulated := new(T)

	for i, loader := range c.loaders {
		temp := new(T)

		if err := loader.Load(temp); err != nil {
			return fmt.Errorf("loader[%d] failed: %w", i, err)
		}

		if err := c.mergeFunc(accumulated, temp); err != nil {
			return fmt.Errorf("merge loader[%d] failed: %w", i, err)
		}
	}

	if c.validator != nil {
		if err := c.validator.Validate(accumulated); err != nil {
			return fmt.Errorf("config validation failed: %w", err)
		}
	}

	c.data = *accumulated
	return nil
}

// Get returns the typed config data.
// Must call Load() before Get(), otherwise returns zero value of T.
func (c *Config[T]) Get() T {
	return c.data
}

// GetPtr returns a pointer to config data.
// Useful when you need to modify config or pass by reference.
func (c *Config[T]) GetPtr() *T {
	return &c.data
}
