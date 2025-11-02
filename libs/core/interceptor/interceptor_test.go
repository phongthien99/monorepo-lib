package interceptor

import (
	"context"
	"errors"
	"testing"
)

func TestChain_EmptyInterceptors(t *testing.T) {
	handler := func(ctx *UniversalContext[TestMeta]) (any, error) {
		return "result", nil
	}

	pipeline := Chain(handler)

	ctx := NewUniversalContext[TestMeta](nil, "test", "method", TestMeta{})
	result, err := pipeline(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "result" {
		t.Errorf("Expected 'result', got %v", result)
	}
}

func TestChain_SingleInterceptor(t *testing.T) {
	var calls []string

	interceptor1 := InterceptorFunc[TestMeta](func(ctx *UniversalContext[TestMeta], next NextFunc[TestMeta]) (any, error) {
		calls = append(calls, "before-1")
		result, err := next(ctx)
		calls = append(calls, "after-1")
		return result, err
	})

	handler := func(ctx *UniversalContext[TestMeta]) (any, error) {
		calls = append(calls, "handler")
		return "result", nil
	}

	pipeline := Chain(handler, interceptor1)
	ctx := NewUniversalContext[TestMeta](nil, "test", "method", TestMeta{})
	result, err := pipeline(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "result" {
		t.Errorf("Expected 'result', got %v", result)
	}

	expectedCalls := []string{"before-1", "handler", "after-1"}
	if !equalSlices(calls, expectedCalls) {
		t.Errorf("Expected calls %v, got %v", expectedCalls, calls)
	}
}

func TestChain_MultipleInterceptors(t *testing.T) {
	var calls []string

	interceptor1 := InterceptorFunc[TestMeta](func(ctx *UniversalContext[TestMeta], next NextFunc[TestMeta]) (any, error) {
		calls = append(calls, "before-1")
		result, err := next(ctx)
		calls = append(calls, "after-1")
		return result, err
	})

	interceptor2 := InterceptorFunc[TestMeta](func(ctx *UniversalContext[TestMeta], next NextFunc[TestMeta]) (any, error) {
		calls = append(calls, "before-2")
		result, err := next(ctx)
		calls = append(calls, "after-2")
		return result, err
	})

	interceptor3 := InterceptorFunc[TestMeta](func(ctx *UniversalContext[TestMeta], next NextFunc[TestMeta]) (any, error) {
		calls = append(calls, "before-3")
		result, err := next(ctx)
		calls = append(calls, "after-3")
		return result, err
	})

	handler := func(ctx *UniversalContext[TestMeta]) (any, error) {
		calls = append(calls, "handler")
		return "result", nil
	}

	pipeline := Chain(handler, interceptor1, interceptor2, interceptor3)
	ctx := NewUniversalContext[TestMeta](nil, "test", "method", TestMeta{})
	result, err := pipeline(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "result" {
		t.Errorf("Expected 'result', got %v", result)
	}

	expectedCalls := []string{
		"before-1", "before-2", "before-3",
		"handler",
		"after-3", "after-2", "after-1",
	}
	if !equalSlices(calls, expectedCalls) {
		t.Errorf("Expected calls %v, got %v", expectedCalls, calls)
	}
}

func TestChain_FailFast(t *testing.T) {
	var calls []string
	expectedErr := errors.New("interceptor error")

	interceptor1 := InterceptorFunc[TestMeta](func(ctx *UniversalContext[TestMeta], next NextFunc[TestMeta]) (any, error) {
		calls = append(calls, "before-1")
		result, err := next(ctx)
		calls = append(calls, "after-1")
		return result, err
	})

	interceptor2 := InterceptorFunc[TestMeta](func(ctx *UniversalContext[TestMeta], next NextFunc[TestMeta]) (any, error) {
		calls = append(calls, "before-2")
		return nil, expectedErr
	})

	interceptor3 := InterceptorFunc[TestMeta](func(ctx *UniversalContext[TestMeta], next NextFunc[TestMeta]) (any, error) {
		calls = append(calls, "before-3")
		result, err := next(ctx)
		calls = append(calls, "after-3")
		return result, err
	})

	handler := func(ctx *UniversalContext[TestMeta]) (any, error) {
		calls = append(calls, "handler")
		return "result", nil
	}

	pipeline := Chain(handler, interceptor1, interceptor2, interceptor3)
	ctx := NewUniversalContext[TestMeta](nil, "test", "method", TestMeta{})
	result, err := pipeline(ctx)

	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	expectedCalls := []string{"before-1", "before-2", "after-1"}
	if !equalSlices(calls, expectedCalls) {
		t.Errorf("Expected calls %v (fail fast), got %v", expectedCalls, calls)
	}
}

func TestChain_ShortCircuit(t *testing.T) {
	var calls []string

	interceptor1 := InterceptorFunc[TestMeta](func(ctx *UniversalContext[TestMeta], next NextFunc[TestMeta]) (any, error) {
		calls = append(calls, "interceptor-1")
		return next(ctx)
	})

	interceptor2 := InterceptorFunc[TestMeta](func(ctx *UniversalContext[TestMeta], next NextFunc[TestMeta]) (any, error) {
		calls = append(calls, "interceptor-2-short-circuit")
		return "short-circuit-result", nil
	})

	handler := func(ctx *UniversalContext[TestMeta]) (any, error) {
		calls = append(calls, "handler")
		return "handler-result", nil
	}

	pipeline := Chain(handler, interceptor1, interceptor2)
	ctx := NewUniversalContext[TestMeta](nil, "test", "method", TestMeta{})
	result, err := pipeline(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "short-circuit-result" {
		t.Errorf("Expected 'short-circuit-result', got %v", result)
	}

	expectedCalls := []string{"interceptor-1", "interceptor-2-short-circuit"}
	if !equalSlices(calls, expectedCalls) {
		t.Errorf("Expected calls %v, got %v", expectedCalls, calls)
	}
}

func TestChain_ModifyResult(t *testing.T) {
	interceptor1 := InterceptorFunc[TestMeta](func(ctx *UniversalContext[TestMeta], next NextFunc[TestMeta]) (any, error) {
		result, err := next(ctx)
		if err != nil {
			return nil, err
		}
		return result.(string) + "-modified-1", nil
	})

	interceptor2 := InterceptorFunc[TestMeta](func(ctx *UniversalContext[TestMeta], next NextFunc[TestMeta]) (any, error) {
		result, err := next(ctx)
		if err != nil {
			return nil, err
		}
		return result.(string) + "-modified-2", nil
	})

	handler := func(ctx *UniversalContext[TestMeta]) (any, error) {
		return "original", nil
	}

	pipeline := Chain(handler, interceptor1, interceptor2)
	ctx := NewUniversalContext[TestMeta](nil, "test", "method", TestMeta{})
	result, err := pipeline(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "original-modified-2-modified-1" {
		t.Errorf("Expected 'original-modified-2-modified-1', got %v", result)
	}
}

func TestChain_AccessMeta(t *testing.T) {
	interceptor := InterceptorFunc[TestMeta](func(ctx *UniversalContext[TestMeta], next NextFunc[TestMeta]) (any, error) {
		if ctx.Meta.UserID == "" {
			return nil, errors.New("unauthorized")
		}
		if ctx.Meta.Role != "admin" {
			return nil, errors.New("forbidden")
		}
		return next(ctx)
	})

	handler := func(ctx *UniversalContext[TestMeta]) (any, error) {
		return "success", nil
	}

	pipeline := Chain(handler, interceptor)

	// Test with valid meta
	ctx := NewUniversalContext(nil, "http", "GET /admin", TestMeta{UserID: "user123", Role: "admin"})
	result, err := pipeline(ctx)
	if err != nil {
		t.Errorf("Expected no error for valid meta, got %v", err)
	}
	if result != "success" {
		t.Errorf("Expected 'success', got %v", result)
	}

	// Test with invalid role
	ctx = NewUniversalContext(nil, "http", "GET /admin", TestMeta{UserID: "user123", Role: "user"})
	result, err = pipeline(ctx)
	if err == nil || err.Error() != "forbidden" {
		t.Errorf("Expected 'forbidden' error, got %v", err)
	}

	// Test with empty userID
	ctx = NewUniversalContext(nil, "http", "GET /admin", TestMeta{UserID: "", Role: "admin"})
	result, err = pipeline(ctx)
	if err == nil || err.Error() != "unauthorized" {
		t.Errorf("Expected 'unauthorized' error, got %v", err)
	}
}

func TestChain_ContextCancellation(t *testing.T) {
	interceptor := InterceptorFunc[TestMeta](func(ctx *UniversalContext[TestMeta], next NextFunc[TestMeta]) (any, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return next(ctx)
		}
	})

	handler := func(ctx *UniversalContext[TestMeta]) (any, error) {
		return "result", nil
	}

	pipeline := Chain(handler, interceptor)

	// Test with cancelled context
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()

	ctx := NewUniversalContext(cancelCtx, "http", "GET /", TestMeta{})
	result, err := pipeline(ctx)

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

func TestChain_PassContextValues(t *testing.T) {
	interceptor1 := InterceptorFunc[TestMeta](func(ctx *UniversalContext[TestMeta], next NextFunc[TestMeta]) (any, error) {
		ctx.Context = context.WithValue(ctx.Context, "key1", "value1")
		return next(ctx)
	})

	interceptor2 := InterceptorFunc[TestMeta](func(ctx *UniversalContext[TestMeta], next NextFunc[TestMeta]) (any, error) {
		ctx.Context = context.WithValue(ctx.Context, "key2", "value2")
		return next(ctx)
	})

	handler := func(ctx *UniversalContext[TestMeta]) (any, error) {
		val1 := ctx.Value("key1")
		val2 := ctx.Value("key2")
		return map[string]string{"key1": val1.(string), "key2": val2.(string)}, nil
	}

	pipeline := Chain(handler, interceptor1, interceptor2)
	ctx := NewUniversalContext(context.Background(), "http", "GET /", TestMeta{})
	result, err := pipeline(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	resultMap := result.(map[string]string)
	if resultMap["key1"] != "value1" {
		t.Errorf("Expected key1='value1', got %v", resultMap["key1"])
	}
	if resultMap["key2"] != "value2" {
		t.Errorf("Expected key2='value2', got %v", resultMap["key2"])
	}
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
