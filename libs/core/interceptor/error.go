package interceptor

import "fmt"

// InterceptorError wraps an error from interceptor execution.
type InterceptorError struct {
	InterceptorName string
	Err             error
}

// Error implements the error interface.
func (e *InterceptorError) Error() string {
	return fmt.Sprintf("interceptor[%s]: %v", e.InterceptorName, e.Err)
}

// Unwrap returns the underlying error for errors.Is() and errors.As().
func (e *InterceptorError) Unwrap() error {
	return e.Err
}

// NewInterceptorError creates a new InterceptorError.
func NewInterceptorError(name string, err error) error {
	if err == nil {
		return nil
	}

	return &InterceptorError{
		InterceptorName: name,
		Err:             err,
	}
}
