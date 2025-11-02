package interceptor

// Bridge connects the interceptor system with a specific framework or adapter.
// Each framework (Gin, Echo, gRPC, etc.) should implement this interface.
type Bridge[M any, NativeCtx any] interface {
	// ExtractMeta extracts adapter-specific metadata from native context.
	// Example for Gin: Extract *gin.Context into GinMeta
	ExtractMeta(nativeCtx NativeCtx) M

	// CreateUniversalContext creates UniversalContext from native context.
	// Typically calls NewUniversalContext with extracted metadata.
	CreateUniversalContext(nativeCtx NativeCtx) *UniversalContext[M]

	// OnSuccess handles successful execution (optional hook).
	// Called after pipeline executes without error.
	OnSuccess(nativeCtx NativeCtx, result any)

	// OnError handles error from pipeline (optional hook).
	// Called when pipeline returns an error.
	OnError(nativeCtx NativeCtx, err error)
}

// BaseBridge provides default implementation for Bridge interface.
// Frameworks can embed this and override specific methods.
type BaseBridge[M any, NativeCtx any] struct {
	Protocol      string
	ExtractMetaFn func(NativeCtx) M
	GetMethodFn   func(NativeCtx) string
	OnSuccessFn   func(NativeCtx, any)
	OnErrorFn     func(NativeCtx, error)
}

// ExtractMeta implements Bridge interface.
func (b *BaseBridge[M, NativeCtx]) ExtractMeta(nativeCtx NativeCtx) M {
	if b.ExtractMetaFn != nil {
		return b.ExtractMetaFn(nativeCtx)
	}
	var zero M
	return zero
}

// CreateUniversalContext implements Bridge interface.
func (b *BaseBridge[M, NativeCtx]) CreateUniversalContext(nativeCtx NativeCtx) *UniversalContext[M] {
	meta := b.ExtractMeta(nativeCtx)
	method := ""
	if b.GetMethodFn != nil {
		method = b.GetMethodFn(nativeCtx)
	}

	return NewUniversalContext(
		nil, // Context will be set by framework
		b.Protocol,
		method,
		meta,
	)
}

// OnSuccess implements Bridge interface (default: no-op).
func (b *BaseBridge[M, NativeCtx]) OnSuccess(nativeCtx NativeCtx, result any) {
	if b.OnSuccessFn != nil {
		b.OnSuccessFn(nativeCtx, result)
	}
}

// OnError implements Bridge interface (default: no-op).
func (b *BaseBridge[M, NativeCtx]) OnError(nativeCtx NativeCtx, err error) {
	if b.OnErrorFn != nil {
		b.OnErrorFn(nativeCtx, err)
	}
}

// InterceptorResolver resolves which interceptors to apply.
// This can be a simple slice, or a Registry implementation.
type InterceptorResolver[M any] interface {
	// Resolve returns interceptors to apply for this context.
	Resolve(ctx *UniversalContext[M], handlerKey string) []Interceptor[M]
}

// SimpleResolver is a basic resolver using a static list.
type SimpleResolver[M any] struct {
	Interceptors []Interceptor[M]
}

// Resolve implements InterceptorResolver.
func (s *SimpleResolver[M]) Resolve(ctx *UniversalContext[M], handlerKey string) []Interceptor[M] {
	return s.Interceptors
}

// ExecutePipeline is a helper to execute interceptor pipeline with a bridge.
// This provides the standard flow: Extract → Resolve → Chain → Execute
func ExecutePipeline[M any, NativeCtx any](
	bridge Bridge[M, NativeCtx],
	resolver InterceptorResolver[M],
	nativeCtx NativeCtx,
	handlerKey string,
	businessHandler NextFunc[M],
) (any, error) {
	// 1. Create UniversalContext from native context
	uCtx := bridge.CreateUniversalContext(nativeCtx)

	// 2. Resolve interceptors
	interceptors := resolver.Resolve(uCtx, handlerKey)

	// 3. Build and execute pipeline
	pipeline := Chain(businessHandler, interceptors...)
	result, err := pipeline(uCtx)

	// 4. Call hooks
	if err != nil {
		bridge.OnError(nativeCtx, err)
	} else {
		bridge.OnSuccess(nativeCtx, result)
	}

	return result, err
}
