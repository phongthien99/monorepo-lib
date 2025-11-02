package interceptor

// NextFunc represents a function that processes a request and returns a result.
// It can be either a business handler or the next interceptor in the chain.
type NextFunc[M any] func(ctx *UniversalContext[M]) (any, error)

// InterceptorFunc is a function type that implements the Interceptor interface.
// Allows using plain functions as interceptors.
type InterceptorFunc[M any] func(ctx *UniversalContext[M], next NextFunc[M]) (any, error)

// Intercept implements the Interceptor interface for InterceptorFunc.
func (f InterceptorFunc[M]) Intercept(ctx *UniversalContext[M], next NextFunc[M]) (any, error) {
	return f(ctx, next)
}
