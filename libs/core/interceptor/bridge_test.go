package interceptor

import (
	"context"
	"errors"
	"testing"
)

type MockNativeContext struct {
	Path   string
	Method string
	UserID string
}

type MockMeta struct {
	RequestPath string
	UserID      string
}

func TestBaseBridge_ExtractMeta(t *testing.T) {
	bridge := &BaseBridge[MockMeta, *MockNativeContext]{
		ExtractMetaFn: func(nc *MockNativeContext) MockMeta {
			return MockMeta{
				RequestPath: nc.Path,
				UserID:      nc.UserID,
			}
		},
	}

	nativeCtx := &MockNativeContext{Path: "/api/users", UserID: "user123"}
	meta := bridge.ExtractMeta(nativeCtx)

	if meta.RequestPath != "/api/users" {
		t.Errorf("Expected RequestPath '/api/users', got '%s'", meta.RequestPath)
	}
	if meta.UserID != "user123" {
		t.Errorf("Expected UserID 'user123', got '%s'", meta.UserID)
	}
}

func TestBaseBridge_ExtractMeta_NilFunction(t *testing.T) {
	bridge := &BaseBridge[MockMeta, *MockNativeContext]{}

	nativeCtx := &MockNativeContext{Path: "/api/users"}
	meta := bridge.ExtractMeta(nativeCtx)

	var zeroMeta MockMeta
	if meta != zeroMeta {
		t.Errorf("Expected zero value meta, got %v", meta)
	}
}

func TestBaseBridge_CreateUniversalContext(t *testing.T) {
	bridge := &BaseBridge[MockMeta, *MockNativeContext]{
		Protocol: "http",
		ExtractMetaFn: func(nc *MockNativeContext) MockMeta {
			return MockMeta{RequestPath: nc.Path, UserID: nc.UserID}
		},
		GetMethodFn: func(nc *MockNativeContext) string {
			return nc.Method + " " + nc.Path
		},
	}

	nativeCtx := &MockNativeContext{Path: "/api/users", Method: "GET", UserID: "user123"}
	uCtx := bridge.CreateUniversalContext(nativeCtx)

	if uCtx.Protocol != "http" {
		t.Errorf("Expected protocol 'http', got '%s'", uCtx.Protocol)
	}
	if uCtx.Method != "GET /api/users" {
		t.Errorf("Expected method 'GET /api/users', got '%s'", uCtx.Method)
	}
	if uCtx.Meta.RequestPath != "/api/users" {
		t.Errorf("Expected RequestPath '/api/users', got '%s'", uCtx.Meta.RequestPath)
	}
	if uCtx.Meta.UserID != "user123" {
		t.Errorf("Expected UserID 'user123', got '%s'", uCtx.Meta.UserID)
	}
}

func TestBaseBridge_CreateUniversalContext_NoMethodFn(t *testing.T) {
	bridge := &BaseBridge[MockMeta, *MockNativeContext]{
		Protocol: "grpc",
		ExtractMetaFn: func(nc *MockNativeContext) MockMeta {
			return MockMeta{RequestPath: nc.Path}
		},
	}

	nativeCtx := &MockNativeContext{Path: "/GetUser"}
	uCtx := bridge.CreateUniversalContext(nativeCtx)

	if uCtx.Method != "" {
		t.Errorf("Expected empty method, got '%s'", uCtx.Method)
	}
}

func TestBaseBridge_OnSuccess(t *testing.T) {
	var capturedResult any
	var capturedContext *MockNativeContext

	bridge := &BaseBridge[MockMeta, *MockNativeContext]{
		OnSuccessFn: func(nc *MockNativeContext, result any) {
			capturedContext = nc
			capturedResult = result
		},
	}

	nativeCtx := &MockNativeContext{Path: "/api/users"}
	bridge.OnSuccess(nativeCtx, "success result")

	if capturedContext != nativeCtx {
		t.Error("Expected OnSuccess to receive native context")
	}
	if capturedResult != "success result" {
		t.Errorf("Expected result 'success result', got %v", capturedResult)
	}
}

func TestBaseBridge_OnSuccess_NilFunction(t *testing.T) {
	bridge := &BaseBridge[MockMeta, *MockNativeContext]{}

	nativeCtx := &MockNativeContext{}
	bridge.OnSuccess(nativeCtx, "result")
}

func TestBaseBridge_OnError(t *testing.T) {
	var capturedError error
	var capturedContext *MockNativeContext

	bridge := &BaseBridge[MockMeta, *MockNativeContext]{
		OnErrorFn: func(nc *MockNativeContext, err error) {
			capturedContext = nc
			capturedError = err
		},
	}

	nativeCtx := &MockNativeContext{Path: "/api/users"}
	testErr := errors.New("test error")
	bridge.OnError(nativeCtx, testErr)

	if capturedContext != nativeCtx {
		t.Error("Expected OnError to receive native context")
	}
	if capturedError != testErr {
		t.Errorf("Expected error %v, got %v", testErr, capturedError)
	}
}

func TestBaseBridge_OnError_NilFunction(t *testing.T) {
	bridge := &BaseBridge[MockMeta, *MockNativeContext]{}

	nativeCtx := &MockNativeContext{}
	bridge.OnError(nativeCtx, errors.New("error"))
}

func TestSimpleResolver_Resolve(t *testing.T) {
	interceptor1 := InterceptorFunc[MockMeta](func(ctx *UniversalContext[MockMeta], next NextFunc[MockMeta]) (any, error) {
		return next(ctx)
	})
	interceptor2 := InterceptorFunc[MockMeta](func(ctx *UniversalContext[MockMeta], next NextFunc[MockMeta]) (any, error) {
		return next(ctx)
	})

	resolver := &SimpleResolver[MockMeta]{
		Interceptors: []Interceptor[MockMeta]{interceptor1, interceptor2},
	}

	ctx := NewUniversalContext[MockMeta](nil, "http", "GET /", MockMeta{})
	interceptors := resolver.Resolve(ctx, "/api/users")

	if len(interceptors) != 2 {
		t.Errorf("Expected 2 interceptors, got %d", len(interceptors))
	}
}

func TestSimpleResolver_ResolveEmpty(t *testing.T) {
	resolver := &SimpleResolver[MockMeta]{
		Interceptors: []Interceptor[MockMeta]{},
	}

	ctx := NewUniversalContext[MockMeta](nil, "http", "GET /", MockMeta{})
	interceptors := resolver.Resolve(ctx, "/api/users")

	if len(interceptors) != 0 {
		t.Errorf("Expected 0 interceptors, got %d", len(interceptors))
	}
}

func TestExecutePipeline_Success(t *testing.T) {
	var calls []string
	var onSuccessCalled bool
	var successResult any

	bridge := &BaseBridge[MockMeta, *MockNativeContext]{
		Protocol: "http",
		ExtractMetaFn: func(nc *MockNativeContext) MockMeta {
			return MockMeta{RequestPath: nc.Path}
		},
		GetMethodFn: func(nc *MockNativeContext) string {
			return nc.Method + " " + nc.Path
		},
		OnSuccessFn: func(nc *MockNativeContext, result any) {
			onSuccessCalled = true
			successResult = result
		},
	}

	interceptor1 := InterceptorFunc[MockMeta](func(ctx *UniversalContext[MockMeta], next NextFunc[MockMeta]) (any, error) {
		calls = append(calls, "interceptor1")
		return next(ctx)
	})

	resolver := &SimpleResolver[MockMeta]{
		Interceptors: []Interceptor[MockMeta]{interceptor1},
	}

	handler := func(ctx *UniversalContext[MockMeta]) (any, error) {
		calls = append(calls, "handler")
		return "success", nil
	}

	nativeCtx := &MockNativeContext{Path: "/api/users", Method: "GET"}
	result, err := ExecutePipeline(bridge, resolver, nativeCtx, "/api/users", handler)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "success" {
		t.Errorf("Expected 'success', got %v", result)
	}
	if !onSuccessCalled {
		t.Error("Expected OnSuccess to be called")
	}
	if successResult != "success" {
		t.Errorf("Expected OnSuccess to receive 'success', got %v", successResult)
	}

	expectedCalls := []string{"interceptor1", "handler"}
	if !equalSlices(calls, expectedCalls) {
		t.Errorf("Expected calls %v, got %v", expectedCalls, calls)
	}
}

func TestExecutePipeline_Error(t *testing.T) {
	var onErrorCalled bool
	var capturedError error
	expectedErr := errors.New("pipeline error")

	bridge := &BaseBridge[MockMeta, *MockNativeContext]{
		Protocol: "http",
		ExtractMetaFn: func(nc *MockNativeContext) MockMeta {
			return MockMeta{}
		},
		OnErrorFn: func(nc *MockNativeContext, err error) {
			onErrorCalled = true
			capturedError = err
		},
	}

	interceptor1 := InterceptorFunc[MockMeta](func(ctx *UniversalContext[MockMeta], next NextFunc[MockMeta]) (any, error) {
		return nil, expectedErr
	})

	resolver := &SimpleResolver[MockMeta]{
		Interceptors: []Interceptor[MockMeta]{interceptor1},
	}

	handler := func(ctx *UniversalContext[MockMeta]) (any, error) {
		return "success", nil
	}

	nativeCtx := &MockNativeContext{}
	result, err := ExecutePipeline(bridge, resolver, nativeCtx, "/api/users", handler)

	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
	if !onErrorCalled {
		t.Error("Expected OnError to be called")
	}
	if capturedError != expectedErr {
		t.Errorf("Expected OnError to receive error %v, got %v", expectedErr, capturedError)
	}
}

func TestExecutePipeline_NoInterceptors(t *testing.T) {
	bridge := &BaseBridge[MockMeta, *MockNativeContext]{
		Protocol:      "http",
		ExtractMetaFn: func(nc *MockNativeContext) MockMeta { return MockMeta{} },
	}

	resolver := &SimpleResolver[MockMeta]{
		Interceptors: []Interceptor[MockMeta]{},
	}

	handler := func(ctx *UniversalContext[MockMeta]) (any, error) {
		return "direct result", nil
	}

	nativeCtx := &MockNativeContext{}
	result, err := ExecutePipeline(bridge, resolver, nativeCtx, "/api/users", handler)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "direct result" {
		t.Errorf("Expected 'direct result', got %v", result)
	}
}

func TestExecutePipeline_ContextPropagation(t *testing.T) {
	bridge := &BaseBridge[MockMeta, *MockNativeContext]{
		Protocol: "http",
		ExtractMetaFn: func(nc *MockNativeContext) MockMeta {
			return MockMeta{UserID: nc.UserID}
		},
	}

	interceptor1 := InterceptorFunc[MockMeta](func(ctx *UniversalContext[MockMeta], next NextFunc[MockMeta]) (any, error) {
		ctx.Context = context.WithValue(ctx.Context, "interceptor-key", "interceptor-value")
		return next(ctx)
	})

	resolver := &SimpleResolver[MockMeta]{
		Interceptors: []Interceptor[MockMeta]{interceptor1},
	}

	handler := func(ctx *UniversalContext[MockMeta]) (any, error) {
		value := ctx.Value("interceptor-key")
		if value == nil {
			return nil, errors.New("context value not propagated")
		}
		return value, nil
	}

	nativeCtx := &MockNativeContext{UserID: "user123"}
	result, err := ExecutePipeline(bridge, resolver, nativeCtx, "/api/users", handler)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "interceptor-value" {
		t.Errorf("Expected 'interceptor-value', got %v", result)
	}
}
