package core

// Validator defines an interface for validating config after loading.
type Validator[T any] interface {
	// Validate checks if the config is valid.
	// Returns nil if valid, error otherwise.
	Validate(*T) error
}

// ValidatorFunc is a function adapter for the Validator interface.
// Allows using a function as a Validator.
//
// Example:
//
//	validateFunc := func(cfg *AppConfig) error {
//	    if cfg.Server.Port < 1024 {
//	        return fmt.Errorf("port must be >= 1024")
//	    }
//	    return nil
//	}
//
//	cfg := config.New[AppConfig](loaders...).
//	    WithValidator(core.ValidatorFunc[AppConfig](validateFunc))
type ValidatorFunc[T any] func(*T) error

// Validate implements the Validator interface.
func (f ValidatorFunc[T]) Validate(cfg *T) error {
	return f(cfg)
}

// CompositeValidator combines multiple validators.
// All validators must pass for validation to succeed.
//
// Example:
//
//	validator := core.NewCompositeValidator(
//	    portValidator,
//	    hostValidator,
//	    databaseValidator,
//	)
type CompositeValidator[T any] struct {
	validators []Validator[T]
}

// NewCompositeValidator creates a new CompositeValidator.
func NewCompositeValidator[T any](validators ...Validator[T]) *CompositeValidator[T] {
	return &CompositeValidator[T]{
		validators: validators,
	}
}

// Validate runs all validators in order.
// Returns the first error encountered, or nil if all pass.
func (c *CompositeValidator[T]) Validate(cfg *T) error {
	for i, validator := range c.validators {
		if err := validator.Validate(cfg); err != nil {
			if len(c.validators) > 1 {
				return &ValidationError{
					ValidatorIndex: i,
					Cause:          err,
				}
			}
			return err
		}
	}
	return nil
}

// ValidationError wraps validation errors with context.
type ValidationError struct {
	ValidatorIndex int
	Cause          error
}

func (e *ValidationError) Error() string {
	return e.Cause.Error()
}

func (e *ValidationError) Unwrap() error {
	return e.Cause
}
