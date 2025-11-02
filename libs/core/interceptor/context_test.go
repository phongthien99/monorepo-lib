package interceptor

import (
	"context"
	"testing"
	"time"
)

func TestNewUniversalContext(t *testing.T) {
	meta := TestMeta{UserID: "user123", Role: "admin"}
	ctx := NewUniversalContext(nil, "http", "GET /api/users", meta)

	if ctx.Protocol != "http" {
		t.Errorf("Expected protocol 'http', got '%s'", ctx.Protocol)
	}
	if ctx.Method != "GET /api/users" {
		t.Errorf("Expected method 'GET /api/users', got '%s'", ctx.Method)
	}
	if ctx.Meta.UserID != "user123" {
		t.Errorf("Expected UserID 'user123', got '%s'", ctx.Meta.UserID)
	}
	if ctx.Meta.Role != "admin" {
		t.Errorf("Expected Role 'admin', got '%s'", ctx.Meta.Role)
	}
	if ctx.Context == nil {
		t.Error("Expected Context to be set to background context")
	}
}

func TestNewUniversalContext_WithParentContext(t *testing.T) {
	parentCtx := context.WithValue(context.Background(), "key", "value")
	meta := TestMeta{UserID: "user123"}

	ctx := NewUniversalContext(parentCtx, "grpc", "GetUser", meta)

	if ctx.Context != parentCtx {
		t.Error("Expected parent context to be preserved")
	}
	if ctx.Value("key") != "value" {
		t.Error("Expected to inherit parent context values")
	}
}

func TestNewUniversalContext_NilContext(t *testing.T) {
	ctx := NewUniversalContext[TestMeta](nil, "http", "GET /", TestMeta{})

	if ctx.Context == nil {
		t.Error("Expected context to default to background")
	}
}

func TestUniversalContext_ContextMethods(t *testing.T) {
	deadline := time.Now().Add(5 * time.Second)
	parentCtx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	ctx := NewUniversalContext[TestMeta](parentCtx, "http", "GET /", TestMeta{})

	// Test Deadline()
	ctxDeadline, ok := ctx.Deadline()
	if !ok {
		t.Error("Expected deadline to be set")
	}
	if !ctxDeadline.Equal(deadline) {
		t.Errorf("Expected deadline %v, got %v", deadline, ctxDeadline)
	}

	// Test Done()
	if ctx.Done() == nil {
		t.Error("Expected Done() channel to be available")
	}

	// Test Err() before cancel
	if ctx.Err() != nil {
		t.Error("Expected no error before cancellation")
	}

	// Cancel and test
	cancel()
	<-ctx.Done()

	if ctx.Err() == nil {
		t.Error("Expected error after cancellation")
	}
}

func TestUniversalContext_WithValue(t *testing.T) {
	ctx := NewUniversalContext[TestMeta](context.Background(), "http", "GET /", TestMeta{})

	// Store value using context.WithValue
	ctx.Context = context.WithValue(ctx.Context, "userID", "user123")
	ctx.Context = context.WithValue(ctx.Context, "requestID", "req456")

	// Retrieve values
	if userID := ctx.Value("userID"); userID != "user123" {
		t.Errorf("Expected 'user123', got %v", userID)
	}
	if requestID := ctx.Value("requestID"); requestID != "req456" {
		t.Errorf("Expected 'req456', got %v", requestID)
	}
	if unknown := ctx.Value("unknown"); unknown != nil {
		t.Errorf("Expected nil for unknown key, got %v", unknown)
	}
}

func TestUniversalContext_GenericMeta(t *testing.T) {
	type CustomMeta struct {
		TraceID string
		SpanID  string
	}

	meta := CustomMeta{TraceID: "trace123", SpanID: "span456"}
	ctx := NewUniversalContext(context.Background(), "kafka", "user.created", meta)

	if ctx.Meta.TraceID != "trace123" {
		t.Errorf("Expected TraceID 'trace123', got '%s'", ctx.Meta.TraceID)
	}
	if ctx.Meta.SpanID != "span456" {
		t.Errorf("Expected SpanID 'span456', got '%s'", ctx.Meta.SpanID)
	}
}

func TestUniversalContext_EmptyMeta(t *testing.T) {
	ctx := NewUniversalContext(context.Background(), "http", "GET /", struct{}{})

	if ctx.Protocol != "http" {
		t.Error("Context should work with empty struct meta")
	}
}
