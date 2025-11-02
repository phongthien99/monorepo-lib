package interceptor

import "context"

// UniversalContext carries request information across all adapters.
// Use context.WithValue for storing additional data.
type UniversalContext[M any] struct {
	context.Context
	Protocol string // "http", "grpc", "kafka", etc.
	Method   string // Route, RPC method, or topic name
	Meta     M      // Adapter-specific metadata
}

// NewUniversalContext creates a new UniversalContext.
func NewUniversalContext[M any](
	ctx context.Context,
	protocol, method string,
	meta M,
) *UniversalContext[M] {
	if ctx == nil {
		ctx = context.Background()
	}

	return &UniversalContext[M]{
		Context:  ctx,
		Protocol: protocol,
		Method:   method,
		Meta:     meta,
	}
}
