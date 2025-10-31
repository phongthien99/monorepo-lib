# Adapter Template

A powerful Go framework for building pluggable adapters with dynamic controller registration and Fx lifecycle management.

[![Go Version](https://img.shields.io/badge/go-1.23+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## 🎯 Overview

`adapter-template` provides a clean, reusable pattern for creating adapters (HTTP, gRPC, Kafka, Cron, etc.) with:

- ✨ **Template Method Pattern** - Standardized lifecycle management
- 🏭 **Factory Pattern** - Easy Fx module creation
- 🔮 **Reflection-based Registration** - Auto-discover controller methods
- 🔄 **Fail-Fast Error Handling** - Production-ready resilience
- 📦 **Zero External Dependencies** - Only stdlib + Fx

## 📦 Installation

```bash
go get github.com/phongthien99/monorepo-lib/libs/core/adapter-template
```

## 🚀 Quick Start

### 1. Create Your Adapter

```go
package myadapter

import (
    "context"
    adaptertemplate "github.com/phongthien99/monorepo-lib/libs/core/adapter-template"
)

// Config for your adapter
type Config struct {
    Port        int
    Controllers []adaptertemplate.ICoreController
}

// MyAdapter embeds BaseAdapter for lifecycle management
type MyAdapter struct {
    adaptertemplate.BaseAdapter[Config]
}

// OnStart implements AdapterLifecycle
func (m *MyAdapter) OnStart(ctx context.Context) error {
    // Register all dynamic controllers
    if err := adaptertemplate.RegisterRouters(m.Config.Controllers, ctx); err != nil {
        return err
    }
    
    // Start your service
    log.Printf("Starting service on port %d", m.Config.Port)
    return nil
}

// OnStop implements AdapterLifecycle
func (m *MyAdapter) OnStop(ctx context.Context) error {
    log.Println("Stopping service...")
    return nil
}
```

### 2. Create a Factory Function

```go
func ForRoot(port int, controllerGroup string) fx.Option {
    return fx.Module("myadapter",
        fx.Provide(
            func() int { return port },
            fx.Annotate(
                NewMyAdapter,
                fx.ParamTags(``, fmt.Sprintf(`group:"%s"`, controllerGroup)),
            ),
        ),
        fx.Invoke(func(lc fx.Lifecycle, adapter *MyAdapter) {
            adapter.RegisterLifecycle(lc, adapter)
        }),
    )
}
```

### 3. Create Controllers

```go
type UserController struct {
    router *gin.Engine
}

var _ adaptertemplate.ICoreController = (*UserController)(nil)

func NewUserController(router *gin.Engine) adaptertemplate.ICoreController {
    return &UserController{router: router}
}

// Auto-registered! Methods with signature func(context.Context) are called automatically
func (u *UserController) GetUsers(ctx context.Context) {
    u.router.GET("/users", func(c *gin.Context) {
        c.JSON(200, gin.H{"users": []string{"Alice", "Bob"}})
    })
}

func (u *UserController) CreateUser(ctx context.Context) {
    u.router.POST("/users", func(c *gin.Context) {
        c.JSON(201, gin.H{"status": "created"})
    })
}

// Fx module
var UserModule = fx.Module("user-controller",
    fx.Provide(
        adaptertemplate.AsRoute(NewUserController, "controllers"),
    ),
)
```

### 4. Wire Everything Together

```go
func main() {
    app := fx.New(
        myadapter.ForRoot(8080, "controllers"),
        UserModule,
    )
    app.Run()
}
```

## 🏗️ Architecture

### Core Components

```
┌─────────────────────────────────────────┐
│         AdapterLifecycle                │
│  ┌────────────────────────────────┐    │
│  │ OnStart(ctx) error             │    │
│  │ OnStop(ctx) error              │    │
│  └────────────────────────────────┘    │
└─────────────────────────────────────────┘
                   ▲
                   │ implements
                   │
┌─────────────────────────────────────────┐
│      BaseAdapter[Config]                │
│  ┌────────────────────────────────┐    │
│  │ Config: T                      │    │
│  │ RegisterLifecycle(...)         │    │
│  └────────────────────────────────┘    │
└─────────────────────────────────────────┘
                   ▲
                   │ embeds
                   │
┌─────────────────────────────────────────┐
│         Your Adapter                    │
│  (GinAdapter, GrpcAdapter, etc.)        │
└─────────────────────────────────────────┘
```

### Dynamic Controller Registration Flow

```
Application Start
       │
       ▼
   OnStart(ctx)
       │
       ▼
RegisterRouters(controllers, ctx)
       │
       ├─► RegisterRouter(controller1)
       │   ├─► Reflect methods
       │   ├─► Validate signature: func(context.Context)
       │   ├─► Call method1(ctx)
       │   ├─► Call method2(ctx)  ← Routes registered here
       │   └─► ...
       │
       ├─► RegisterRouter(controller2)
       │   └─► ...
       │
       └─► Return error if any panic
```

## 📚 API Reference

### Types

#### AdapterLifecycle

```go
type AdapterLifecycle interface {
    OnStart(ctx context.Context) error
    OnStop(ctx context.Context) error
}
```

Defines lifecycle hooks for adapters.

#### BaseAdapter[T any]

```go
type BaseAdapter[T any] struct {
    Config T
}

func (b *BaseAdapter[T]) RegisterLifecycle(lc fx.Lifecycle, impl AdapterLifecycle)
```

Generic template for adapter configuration and lifecycle registration.

#### ICoreController

```go
type ICoreController interface {}
```

Marker interface for controllers. Controllers implementing this interface can be registered dynamically.

### Functions

#### BaseTemplate

```go
func BaseTemplate(lc fx.Lifecycle, impl AdapterLifecycle)
```

Registers lifecycle hooks with Fx. Panics if `lc` or `impl` is nil.

#### RegisterRouter

```go
func RegisterRouter(controller ICoreController, ctx context.Context) error
```

Registers a single controller by calling all methods matching signature `func(context.Context)`.

**Behavior:**
- Methods are called in **alphabetical order** (reflection behavior)
- Panics are recovered and returned as errors
- **Fail-fast**: stops at first error

#### RegisterRouters

```go
func RegisterRouters(controllers []ICoreController, ctx context.Context) error
```

Registers multiple controllers. Stops at first error (fail-fast).

#### AsRoute

```go
func AsRoute(f any, groupTag string, annotation ...fx.Annotation) any
```

Helper to annotate controller constructors for Fx group injection.

**Example:**
```go
fx.Provide(
    AsRoute(NewUserController, "controllers"),
)
```

## 🎨 Design Patterns

### Template Method Pattern

The `AdapterLifecycle` interface defines the template:

```go
type AdapterLifecycle interface {
    OnStart(ctx context.Context) error   // Template step 1
    OnStop(ctx context.Context) error    // Template step 2
}
```

Each adapter implements these steps differently:

```go
// HTTP Adapter
func (h *HttpAdapter) OnStart(ctx context.Context) error {
    registerRoutes()
    startServer()
}

// gRPC Adapter
func (g *GrpcAdapter) OnStart(ctx context.Context) error {
    registerServices()
    startGrpcServer()
}

// Kafka Adapter
func (k *KafkaAdapter) OnStart(ctx context.Context) error {
    subscribeTopics()
    startConsumer()
}
```

### Factory Pattern

`ForRoot` functions create pre-configured Fx modules:

```go
func ForRoot(port int, controllerGroup string) fx.Option {
    return fx.Module("adapter",
        fx.Provide(...),  // Dependencies
        fx.Invoke(...),   // Lifecycle registration
    )
}

// Usage:
fx.New(
    HttpAdapter.ForRoot(8080, "httpControllers"),
    GrpcAdapter.ForRoot(9090, "grpcControllers"),
).Run()
```

### Reflection-based Auto-Discovery

Controllers don't need explicit route registration:

```go
// ❌ Manual registration (tedious):
func (u *UserController) Register(router *gin.Engine) {
    router.GET("/users", u.GetUsers)
    router.POST("/users", u.CreateUser)
    router.PUT("/users/:id", u.UpdateUser)
    router.DELETE("/users/:id", u.DeleteUser)
}

// ✅ Auto-discovery (clean):
func (u *UserController) GetUsers(ctx context.Context) {
    u.router.GET("/users", handler)
}
// Just implement the method - it's auto-registered!
```

## ⚡ Performance

### Reflection Overhead

Reflection is used **only during application startup** (OnStart phase):

| Controllers | Methods | Registration Time |
|-------------|---------|-------------------|
| 10          | 100     | ~50-100μs         |
| 100         | 1000    | ~500μs-1ms        |

**No runtime performance impact** - routes are registered once during startup.

### Complexity

- **Time**: O(N × M) where N = controllers, M = methods per controller
- **Space**: O(1) - no caching, minimal allocations

## 🔐 Best Practices

### 1. Immutable Config

```go
// ✅ Good: Create config once
config := Config{Port: 8080}
adapter := NewAdapter(config)

// ❌ Bad: Mutate after creation (race condition risk)
adapter.Config.Port = 9090  // Dangerous if accessed concurrently
```

### 2. Fail-Fast Validation

```go
func (m *MyAdapter) OnStart(ctx context.Context) error {
    // Validate config first
    if m.Config.Port == 0 {
        return fmt.Errorf("port cannot be zero")
    }
    
    // Then register controllers
    return adaptertemplate.RegisterRouters(m.Config.Controllers, ctx)
}
```

### 3. Controller Separation of Concerns

```go
// ✅ Good: Controller only registers routes
func (u *UserController) GetUsers(ctx context.Context) {
    u.router.GET("/users", u.handleGetUsers)
}

func (u *UserController) handleGetUsers(c *gin.Context) {
    // Business logic here
}

// ❌ Bad: Mixing registration and business logic
func (u *UserController) GetUsers(ctx context.Context) {
    users := u.db.FindAll()  // Business logic in registration method
    u.router.GET("/users", ...)
}
```

### 4. Error Handling

```go
func (m *MyAdapter) OnStart(ctx context.Context) error {
    // Don't panic - return errors
    if err := adaptertemplate.RegisterRouters(m.Config.Controllers, ctx); err != nil {
        return fmt.Errorf("failed to register routes: %w", err)
    }
    return nil
}
```

## ⚠️ Important Notes

### Method Iteration Order

Methods are called in **alphabetical order**, not declaration order:

```go
type Controller struct{}

func (c *Controller) Zebra(ctx context.Context)    {} // Called 3rd
func (c *Controller) Alpha(ctx context.Context)    {} // Called 1st
func (c *Controller) MyMethod(ctx context.Context) {} // Called 2nd
```

This is a Go reflection behavior. If order matters, use explicit registration instead of dynamic controllers.

### Context Propagation

The context passed to `RegisterRouters` is propagated to all controller methods:

```go
func (m *MyAdapter) OnStart(ctx context.Context) error {
    // ctx from Fx lifecycle (may have timeout, values)
    return adaptertemplate.RegisterRouters(m.Config.Controllers, ctx)
}

func (u *UserController) GetUsers(ctx context.Context) {
    // Same ctx as above
    deadline, ok := ctx.Deadline()
}
```

### Thread Safety

`BaseAdapter.Config` is **not protected by mutex**. Treat it as immutable after construction.

## 📖 Examples

### Complete Example: HTTP Adapter

See [`libs/http/ginfx`](../../http/ginfx) for a full implementation:

- GinAdapter implementing AdapterLifecycle
- ForRoot factory with configurable controller groups
- Health and User controllers with dynamic registration
- Graceful shutdown with context timeout

### Running the Example

```bash
cd cmd/gin-example
go run main.go
```

Then visit:
- http://localhost:8080/health
- http://localhost:8080/users

## 🧪 Testing

Run tests:

```bash
cd libs/core/adapter-template
go test -v
```

Test coverage:

```bash
go test -cover
# coverage: 79.1% of statements
```

## 🤝 Contributing

Contributions welcome! Please:

1. Add tests for new features
2. Update documentation
3. Follow existing code style
4. Ensure `go vet` and `gofmt` pass

## 📄 License

MIT License - see LICENSE file

## 🙏 Acknowledgments

- Inspired by NestJS module system
- Built on [Uber Fx](https://github.com/uber-go/fx)
- Uses Go 1.23+ generics

## 📚 Further Reading

- [Fx Documentation](https://uber-go.github.io/fx/)
- [Template Method Pattern](https://refactoring.guru/design-patterns/template-method)
- [Reflection in Go](https://go.dev/blog/laws-of-reflection)

---

**Made with ❤️ for the Go community**
