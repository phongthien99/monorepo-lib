package interceptor

// Interceptor is the core interface for implementing cross-cutting concerns.
// Uses Chain of Responsibility pattern.
// All errors fail fast - the pipeline stops immediately on any error.
type Interceptor[M any] interface {
	// Intercept executes the interceptor's logic.
	// Must call next(ctx) to continue the chain (unless short-circuiting).
	Intercept(ctx *UniversalContext[M], next NextFunc[M]) (any, error)
}

// Chain composes multiple interceptors into a single execution pipeline.
// Execution order: interceptors[0] → interceptors[1] → ... → handler
func Chain[M any](handler NextFunc[M], interceptors ...Interceptor[M]) NextFunc[M] {
	if len(interceptors) == 0 {
		return handler
	}

	// Build pipeline from right to left (last to first)
	for i := len(interceptors) - 1; i >= 0; i-- {
		currentInterceptor := interceptors[i]
		nextFunc := handler

		handler = func(ctx *UniversalContext[M]) (any, error) {
			return currentInterceptor.Intercept(ctx, nextFunc)
		}
	}

	return handler
}
