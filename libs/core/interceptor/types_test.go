package interceptor

import (
	"errors"
	"testing"
)

type TestMeta struct {
	UserID string
	Role   string
}

func TestInterceptorFunc_Intercept(t *testing.T) {
	called := false
	interceptorFunc := InterceptorFunc[TestMeta](func(ctx *UniversalContext[TestMeta], next NextFunc[TestMeta]) (any, error) {
		called = true
		return next(ctx)
	})

	handler := func(ctx *UniversalContext[TestMeta]) (any, error) {
		return "result", nil
	}

	ctx := NewUniversalContext[TestMeta](nil, "test", "method", TestMeta{})
	result, err := interceptorFunc.Intercept(ctx, handler)

	if !called {
		t.Error("InterceptorFunc was not called")
	}
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "result" {
		t.Errorf("Expected 'result', got %v", result)
	}
}

func TestInterceptorFunc_ReturnsError(t *testing.T) {
	expectedErr := errors.New("test error")
	interceptorFunc := InterceptorFunc[TestMeta](func(ctx *UniversalContext[TestMeta], next NextFunc[TestMeta]) (any, error) {
		return nil, expectedErr
	})

	handler := func(ctx *UniversalContext[TestMeta]) (any, error) {
		return "should not reach", nil
	}

	ctx := NewUniversalContext[TestMeta](nil, "test", "method", TestMeta{})
	result, err := interceptorFunc.Intercept(ctx, handler)

	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

func TestInterceptorFunc_ModifiesResult(t *testing.T) {
	interceptorFunc := InterceptorFunc[TestMeta](func(ctx *UniversalContext[TestMeta], next NextFunc[TestMeta]) (any, error) {
		result, err := next(ctx)
		if err != nil {
			return nil, err
		}
		return result.(string) + "-modified", nil
	})

	handler := func(ctx *UniversalContext[TestMeta]) (any, error) {
		return "original", nil
	}

	ctx := NewUniversalContext[TestMeta](nil, "test", "method", TestMeta{})
	result, err := interceptorFunc.Intercept(ctx, handler)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "original-modified" {
		t.Errorf("Expected 'original-modified', got %v", result)
	}
}

func TestNextFunc_DirectCall(t *testing.T) {
	handler := NextFunc[TestMeta](func(ctx *UniversalContext[TestMeta]) (any, error) {
		return ctx.Meta.UserID, nil
	})

	ctx := NewUniversalContext[TestMeta](nil, "test", "method", TestMeta{UserID: "user123"})
	result, err := handler(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "user123" {
		t.Errorf("Expected 'user123', got %v", result)
	}
}
