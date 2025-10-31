# Adapter Template Examples

This directory contains practical examples demonstrating different use cases of the `adapter-template` package.

## Examples Overview

### 1. Simple Adapter (`simple_adapter.go`)

**Purpose**: Minimal example showing basic adapter creation

**Features**:
- Simple configuration
- Basic lifecycle implementation
- Dynamic controller registration
- Minimal boilerplate

**Use When**:
- Getting started with adapter-template
- Creating proof-of-concept adapters
- Learning the basic patterns

**Code Highlights**:
```go
type SimpleAdapter struct {
    adaptertemplate.BaseAdapter[SimpleConfig]
}

func (s *SimpleAdapter) OnStart(ctx context.Context) error {
    return adaptertemplate.RegisterRouters(s.Config.Controllers, ctx)
}
```

---

### 2. Validated Adapter (`validation_adapter.go`)

**Purpose**: Production-ready adapter with comprehensive validation

**Features**:
- Config validation with `Validate()` method
- Error handling at multiple levels
- Graceful shutdown with context checking
- Defense-in-depth validation strategy

**Use When**:
- Building production adapters
- Need robust error handling
- Configuration validation is critical
- Graceful shutdown is required

**Code Highlights**:
```go
type ValidatedConfig struct {
    Port int
    ServiceName string
}

func (c *ValidatedConfig) Validate() error {
    if c.Port <= 0 || c.Port > 65535 {
        return fmt.Errorf("invalid port: %d", c.Port)
    }
    return nil
}
```

---

## Running Examples

These are code examples meant to be copied and adapted for your use case. They are not standalone runnable programs.

To use an example:

1. **Copy the relevant example code** to your project
2. **Modify** the adapter to fit your needs (HTTP, gRPC, Kafka, etc.)
3. **Add your business logic** in OnStart/OnStop
4. **Create controllers** for your use case
5. **Wire with Fx** in your main.go

## Example Structure

Each example follows this pattern:

```
Example File
├── Config Struct           - Configuration definition
├── Adapter Struct         - Embeds BaseAdapter[Config]
├── OnStart/OnStop         - Lifecycle implementation
├── ForRoot Function       - Fx module factory
└── Example Controller     - Sample controller usage
```

## Integration with Real Adapters

### HTTP Example

See `libs/http/ginfx/` for a complete HTTP adapter implementation:

```go
// Based on simple_adapter.go pattern
type GinAdapter struct {
    adaptertemplate.BaseAdapter[Config]
    server *http.Server  // Adapter-specific state
}

func (g *GinAdapter) OnStart(ctx context.Context) error {
    // Register routes
    adaptertemplate.RegisterRouters(g.Config.Controllers, ctx)
    
    // Start HTTP server
    go g.server.ListenAndServe()
    return nil
}
```

### gRPC Example (Conceptual)

```go
// Based on validation_adapter.go pattern
type GrpcConfig struct {
    Port        int
    Controllers []adaptertemplate.ICoreController
}

func (c *GrpcConfig) Validate() error {
    if c.Port <= 0 {
        return fmt.Errorf("invalid port")
    }
    return nil
}

type GrpcAdapter struct {
    adaptertemplate.BaseAdapter[GrpcConfig]
    server *grpc.Server
}

func (g *GrpcAdapter) OnStart(ctx context.Context) error {
    if err := g.Config.Validate(); err != nil {
        return err
    }
    
    // Register gRPC services via controllers
    adaptertemplate.RegisterRouters(g.Config.Controllers, ctx)
    
    // Start gRPC server
    go g.server.Serve(listener)
    return nil
}
```

## Best Practices from Examples

### From Simple Adapter

✅ **Keep it minimal**: Don't add complexity until needed

```go
// Good: Simple and clear
func (s *SimpleAdapter) OnStart(ctx context.Context) error {
    return adaptertemplate.RegisterRouters(s.Config.Controllers, ctx)
}

// Avoid: Premature optimization
func (s *SimpleAdapter) OnStart(ctx context.Context) error {
    // Validate
    // Cache
    // Monitor
    // ... (unnecessary for simple use case)
}
```

### From Validated Adapter

✅ **Validate early**: Fail fast during construction

```go
func NewValidatedAdapter(config Config) (*Adapter, error) {
    if err := config.Validate(); err != nil {
        return nil, err  // Fail at creation time
    }
    return &Adapter{...}
}
```

✅ **Defense in depth**: Validate at multiple layers

```go
func (v *ValidatedAdapter) OnStart(ctx context.Context) error {
    // Re-validate even if validated in constructor
    if err := v.Config.Validate(); err != nil {
        return err
    }
    // ...
}
```

## Common Patterns

### Pattern 1: Stateful Adapter

```go
type StatefulAdapter struct {
    adaptertemplate.BaseAdapter[Config]
    server  *http.Server  // Managed state
    db      *sql.DB       // External resources
}

func (s *StatefulAdapter) OnStop(ctx context.Context) error {
    // Cleanup managed resources
    if s.server != nil {
        s.server.Shutdown(ctx)
    }
    if s.db != nil {
        s.db.Close()
    }
    return nil
}
```

### Pattern 2: Multi-Protocol Adapter

```go
type MultiAdapter struct {
    adaptertemplate.BaseAdapter[Config]
    httpServer *http.Server
    grpcServer *grpc.Server
}

func (m *MultiAdapter) OnStart(ctx context.Context) error {
    // Start multiple protocols
    go m.httpServer.ListenAndServe()
    go m.grpcServer.Serve(listener)
    return nil
}
```

### Pattern 3: Controller Groups

```go
// Different controller groups for different purposes
fx.New(
    MyAdapter.ForRoot(8080, "publicControllers"),
    MyAdapter.ForRoot(9090, "adminControllers"),
    
    // Public controllers
    PublicHealthModule,
    PublicAPIModule,
    
    // Admin controllers
    AdminDashboardModule,
    AdminMetricsModule,
)
```

## Troubleshooting

### Issue: Controllers not registered

**Symptom**: Methods not called during OnStart

**Solutions**:
1. Check method signature is exactly `func(context.Context)`
2. Verify controller implements `ICoreController`
3. Ensure controller added to correct Fx group
4. Check for panics in controller methods

### Issue: Config validation fails

**Symptom**: Adapter fails to start with validation error

**Solutions**:
1. Implement `Validate()` method
2. Call `Validate()` in constructor
3. Check all required fields have values
4. Add detailed error messages

### Issue: Graceful shutdown timeout

**Symptom**: OnStop takes too long

**Solutions**:
1. Check context deadline in OnStop
2. Add timeout to resource cleanup
3. Use `select` to respect context cancellation

```go
func (a *Adapter) OnStop(ctx context.Context) error {
    done := make(chan struct{})
    go func() {
        // Cleanup logic
        close(done)
    }()
    
    select {
    case <-done:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

## Further Reading

- [Main README](../README.md) - Full package documentation
- [libs/http/ginfx](../../../http/ginfx) - Production HTTP adapter
- [Fx Documentation](https://uber-go.github.io/fx/) - Fx framework guide

## Contributing Examples

Want to add an example? Please ensure:

1. ✅ Code compiles and follows Go conventions
2. ✅ Includes comprehensive comments
3. ✅ Demonstrates a specific use case
4. ✅ Follows existing example structure
5. ✅ Updates this README with description

---

**Happy Coding!** 🚀
