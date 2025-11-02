package interceptor

import (
	"errors"
	"testing"
)

func TestInterceptorError_Error(t *testing.T) {
	originalErr := errors.New("original error")
	interceptorErr := &InterceptorError{
		InterceptorName: "AuthInterceptor",
		Err:             originalErr,
	}

	expected := "interceptor[AuthInterceptor]: original error"
	if interceptorErr.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, interceptorErr.Error())
	}
}

func TestInterceptorError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	interceptorErr := &InterceptorError{
		InterceptorName: "LoggingInterceptor",
		Err:             originalErr,
	}

	unwrapped := interceptorErr.Unwrap()
	if unwrapped != originalErr {
		t.Errorf("Expected unwrapped error to be original error")
	}
}

func TestInterceptorError_ErrorsIs(t *testing.T) {
	originalErr := errors.New("database connection failed")
	interceptorErr := &InterceptorError{
		InterceptorName: "DBInterceptor",
		Err:             originalErr,
	}

	if !errors.Is(interceptorErr, originalErr) {
		t.Error("errors.Is should match original error")
	}
}

func TestInterceptorError_ErrorsAs(t *testing.T) {
	originalErr := &InterceptorError{
		InterceptorName: "ValidationInterceptor",
		Err:             errors.New("validation failed"),
	}

	wrappedErr := &InterceptorError{
		InterceptorName: "OuterInterceptor",
		Err:             originalErr,
	}

	var target *InterceptorError
	if !errors.As(wrappedErr, &target) {
		t.Error("errors.As should match InterceptorError type")
	}

	if target.InterceptorName != "OuterInterceptor" {
		t.Errorf("Expected 'OuterInterceptor', got '%s'", target.InterceptorName)
	}
}

func TestNewInterceptorError_WithError(t *testing.T) {
	originalErr := errors.New("test error")
	err := NewInterceptorError("TestInterceptor", originalErr)

	if err == nil {
		t.Error("Expected non-nil error")
	}

	var interceptorErr *InterceptorError
	if !errors.As(err, &interceptorErr) {
		t.Error("Expected error to be of type InterceptorError")
	}

	if interceptorErr.InterceptorName != "TestInterceptor" {
		t.Errorf("Expected name 'TestInterceptor', got '%s'", interceptorErr.InterceptorName)
	}

	if !errors.Is(err, originalErr) {
		t.Error("Expected error to wrap original error")
	}
}

func TestNewInterceptorError_WithNilError(t *testing.T) {
	err := NewInterceptorError("TestInterceptor", nil)

	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}

func TestInterceptorError_ChainedErrors(t *testing.T) {
	baseErr := errors.New("base error")
	err1 := NewInterceptorError("Interceptor1", baseErr)
	err2 := NewInterceptorError("Interceptor2", err1)
	err3 := NewInterceptorError("Interceptor3", err2)

	// Test error message contains outermost interceptor
	if !errors.Is(err3, baseErr) {
		t.Error("errors.Is should match base error through chain")
	}

	// Test unwrapping chain
	var interceptorErr *InterceptorError
	if !errors.As(err3, &interceptorErr) {
		t.Error("errors.As should match InterceptorError")
	}
	if interceptorErr.InterceptorName != "Interceptor3" {
		t.Errorf("Expected outermost interceptor 'Interceptor3', got '%s'", interceptorErr.InterceptorName)
	}
}

type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

func TestInterceptorError_WithCustomError(t *testing.T) {
	customErr := &CustomError{Code: 404, Message: "not found"}
	interceptorErr := &InterceptorError{
		InterceptorName: "RouterInterceptor",
		Err:             customErr,
	}

	unwrapped := interceptorErr.Unwrap()
	if customErrUnwrapped, ok := unwrapped.(*CustomError); !ok {
		t.Error("Expected unwrapped error to be CustomError")
	} else {
		if customErrUnwrapped.Code != 404 {
			t.Errorf("Expected code 404, got %d", customErrUnwrapped.Code)
		}
		if customErrUnwrapped.Message != "not found" {
			t.Errorf("Expected message 'not found', got '%s'", customErrUnwrapped.Message)
		}
	}
}

func TestInterceptorError_InChainExecution(t *testing.T) {
	expectedErr := errors.New("auth failed")

	authInterceptor := InterceptorFunc[TestMeta](func(ctx *UniversalContext[TestMeta], next NextFunc[TestMeta]) (any, error) {
		return nil, NewInterceptorError("AuthInterceptor", expectedErr)
	})

	handler := func(ctx *UniversalContext[TestMeta]) (any, error) {
		return "success", nil
	}

	pipeline := Chain(handler, authInterceptor)
	ctx := NewUniversalContext[TestMeta](nil, "http", "GET /", TestMeta{})
	result, err := pipeline(ctx)

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	var interceptorErr *InterceptorError
	if !errors.As(err, &interceptorErr) {
		t.Fatalf("Expected InterceptorError, got %T", err)
	}

	if interceptorErr.InterceptorName != "AuthInterceptor" {
		t.Errorf("Expected name 'AuthInterceptor', got '%s'", interceptorErr.InterceptorName)
	}

	if !errors.Is(err, expectedErr) {
		t.Error("Expected to match original error")
	}
}
