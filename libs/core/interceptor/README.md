# Interceptor

A minimal, type-safe AOP (Aspect-Oriented Programming) interceptor framework for Go using generics.

## Features

- **Minimal Core**: Only 30 lines for chain logic
- **Type-Safe**: Generic metadata `[M any]` with compile-time checking
- **Framework Agnostic**: Bridge pattern for Gin, Echo, gRPC, Kafka, etc.
- **Zero Overhead**: No reflection, no allocations in hot path
- **FailFast**: Simple error handling - pipeline stops on first error
- **Go Idiomatic**: Embedded context, error wrapping, functional composition

## Installation

```bash
go get github.com/phongthien99/monorepo-lib/libs/core
```

Then import in your code:

```go
import "github.com/phongthien99/monorepo-lib/libs/core/interceptor"
```

## Quick Start

### 1. Define Your Metadata

```go
type GinMeta struct {
    RawQuery    string
    ContentType string
    UserAgent   string
}
```

### 2. Create Interceptors

```go
// Logging interceptor
var loggingInterceptor = interceptor.InterceptorFunc[GinMeta](
    func(ctx *interceptor.UniversalContext[GinMeta], next interceptor.NextFunc[GinMeta]) (any, error) {
        start := time.Now()
        result, err := next(ctx)
        log.Printf("[%s] %s - %v", ctx.Protocol, ctx.Method, time.Since(start))
        return result, err
    },
)

// Auth interceptor
var authInterceptor = interceptor.InterceptorFunc[GinMeta](
    func(ctx *interceptor.UniversalContext[GinMeta], next interceptor.NextFunc[GinMeta]) (any, error) {
        token := ctx.Meta.UserAgent
        if token == "" {
            return nil, errors.New("unauthorized")
        }
        return next(ctx)
    },
)
```

### 3. Build Chain

```go
handler := func(ctx *interceptor.UniversalContext[GinMeta]) (any, error) {
    return "Hello World", nil
}

pipeline := interceptor.Chain(handler, loggingInterceptor, authInterceptor)

// Execute
ctx := interceptor.NewUniversalContext(context.Background(), "http", "/api/users", GinMeta{})
result, err := pipeline(ctx)
```

## Bridge Pattern for Frameworks

### Gin Example

```go
type GinBridge struct {
    *interceptor.BaseBridge[GinMeta, *gin.Context]
}

func NewGinBridge() *GinBridge {
    return &GinBridge{
        BaseBridge: &interceptor.BaseBridge[GinMeta, *gin.Context]{
            Protocol: "http",
            ExtractMetaFn: func(c *gin.Context) GinMeta {
                return GinMeta{
                    RawQuery:    c.Request.URL.RawQuery,
                    ContentType: c.ContentType(),
                    UserAgent:   c.GetHeader("User-Agent"),
                }
            },
            GetMethodFn: func(c *gin.Context) string {
                return c.Request.Method + " " + c.FullPath()
            },
            OnErrorFn: func(c *gin.Context, err error) {
                c.JSON(500, gin.H{"error": err.Error()})
            },
        },
    }
}
```

### Global Middleware

```go
func (b *GinBridge) GlobalMiddleware(interceptors ...interceptor.Interceptor[GinMeta]) gin.HandlerFunc {
    return func(c *gin.Context) {
        handler := func(ctx *interceptor.UniversalContext[GinMeta]) (any, error) {
            c.Next()
            return nil, nil
        }

        uCtx := b.CreateUniversalContext(c)
        uCtx.Context = c.Request.Context()

        pipeline := interceptor.Chain(handler, interceptors...)
        _, err := pipeline(uCtx)

        if err != nil {
            b.OnError(c, err)
            c.Abort()
        }
    }
}
```

## Core API

### Types

```go
// Core interface - implement this for custom interceptors
type Interceptor[M any] interface {
    Intercept(ctx *UniversalContext[M], next NextFunc[M]) (any, error)
}

// Function type for simple interceptors
type InterceptorFunc[M any] func(ctx *UniversalContext[M], next NextFunc[M]) (any, error)

// Handler/next function
type NextFunc[M any] func(ctx *UniversalContext[M]) (any, error)

// Context carrying request info
type UniversalContext[M any] struct {
    context.Context
    Protocol string
    Method   string
    Meta     M
}
```

### Functions

```go
// Chain composes interceptors into pipeline
func Chain[M any](handler NextFunc[M], interceptors ...Interceptor[M]) NextFunc[M]

// Create new context
func NewUniversalContext[M any](ctx context.Context, protocol, method string, meta M) *UniversalContext[M]
```

## Advanced Usage

### Context Values

Use standard Go context for storing data:

```go
func authInterceptor(ctx *UniversalContext[GinMeta], next NextFunc[GinMeta]) (any, error) {
    userID := getUserID(ctx.Meta.UserAgent)
    
    // Store in context
    ctx.Context = context.WithValue(ctx.Context, "userID", userID)
    
    return next(ctx)
}

func businessHandler(ctx *UniversalContext[GinMeta]) (any, error) {
    // Retrieve from context
    userID := ctx.Value("userID").(string)
    return fmt.Sprintf("User: %s", userID), nil
}
```

### Custom Interceptor Type

```go
type TimingInterceptor struct {
    threshold time.Duration
}

func (t *TimingInterceptor) Intercept(ctx *UniversalContext[GinMeta], next NextFunc[GinMeta]) (any, error) {
    start := time.Now()
    result, err := next(ctx)
    
    duration := time.Since(start)
    if duration > t.threshold {
        log.Printf("SLOW: %s took %v", ctx.Method, duration)
    }
    
    return result, err
}
```

### Error Handling

```go
import "errors"

func validationInterceptor(ctx *UniversalContext[GinMeta], next NextFunc[GinMeta]) (any, error) {
    if ctx.Meta.ContentType != "application/json" {
        return nil, interceptor.NewInterceptorError("validation", errors.New("invalid content type"))
    }
    return next(ctx)
}

// Check error type
result, err := pipeline(ctx)
if err != nil {
    var interceptorErr *interceptor.InterceptorError
    if errors.As(err, &interceptorErr) {
        log.Printf("Interceptor '%s' failed: %v", interceptorErr.InterceptorName, interceptorErr.Err)
    }
}
```

## Integration with Registry

For dynamic interceptor selection based on rules, use the optional registry module:

```go
import "github.com/phongthien99/monorepo-lib/libs/interceptor-registry"

reg := registry.NewRegistry[GinMeta]()

reg.Global(loggingInterceptor)
reg.ForProtocol("http", authInterceptor)
reg.ForMethod("/api/*", rateLimitInterceptor)

// Resolve interceptors for specific context
interceptors := reg.Resolve(ctx, "/api/users")
pipeline := interceptor.Chain(handler, interceptors...)
```

## Design Patterns

- **Chain of Responsibility**: Sequential interceptor execution
- **Bridge Pattern**: Framework-agnostic integration
- **Template Method**: BaseBridge with customizable hooks
- **Strategy Pattern**: InterceptorResolver for selection logic

## Performance

- **Zero allocations** in Chain loop (closure reuse)
- **No reflection** (pure generics)
- **No locks** (immutable chain)
- **~5-10ns** per interceptor (function call overhead only)

## Best Practices

1. **Keep interceptors focused** - Each interceptor should do one thing
2. **Always call next()** - Unless intentionally short-circuiting
3. **Use context.WithValue()** - For passing data between interceptors
4. **Fail fast** - Return errors immediately, don't swallow them
5. **Order matters** - First in slice = first to execute

## Examples

See [cmd/gin-example](../../../cmd/gin-example) for complete working examples with:
- Multiple controllers
- Global and route-specific interceptors
- Error handling
- Metrics collection

## License

MIT
