# Adapter Template Architecture

Tài liệu kỹ thuật chi tiết về kiến trúc của package `adapter-template`.

## 📐 Architectural Overview

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Your Application                          │
│                                                              │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐            │
│  │   HTTP     │  │   gRPC     │  │   Kafka    │            │
│  │  Adapter   │  │  Adapter   │  │  Adapter   │  ...       │
│  └──────┬─────┘  └──────┬─────┘  └──────┬─────┘            │
│         │                │                │                  │
│         └────────────────┴────────────────┘                  │
│                          │                                   │
│         ┌────────────────▼────────────────┐                 │
│         │   adapter-template Package      │                 │
│         │                                  │                 │
│         │  ┌──────────────────────────┐   │                 │
│         │  │  BaseAdapter[T]          │   │                 │
│         │  │  AdapterLifecycle        │   │                 │
│         │  │  ICoreController         │   │                 │
│         │  │  RegisterRouter          │   │                 │
│         │  └──────────────────────────┘   │                 │
│         └──────────────┬───────────────────┘                 │
│                        │                                     │
│         ┌──────────────▼───────────────┐                    │
│         │      Uber Fx Framework       │                    │
│         │   (Dependency Injection)     │                    │
│         └──────────────────────────────┘                    │
└─────────────────────────────────────────────────────────────┘
```

---

## 🏛️ Design Patterns

### 1. Template Method Pattern

**Intent**: Define skeleton of algorithm, let subclasses override specific steps.

**Implementation**:

```
┌─────────────────────────────────┐
│     AdapterLifecycle            │
│     (Template Interface)        │
├─────────────────────────────────┤
│ + OnStart(ctx) error    ◄────┐  │
│ + OnStop(ctx) error      ◄──┐│  │
└─────────────────────────────┘││  │
              ▲                 ││  │
              │                 ││  │
    ┌─────────┴──────────┐      ││  │
    │                    │      ││  │
┌───┴────────┐  ┌────────┴───┐  ││  │
│ GinAdapter │  │ GrpcAdapter│  ││  │
├────────────┤  ├────────────┤  ││  │
│ OnStart()  │──┘  OnStart()  │──┘  │
│ OnStop()   │     OnStop()   │     │
└────────────┘  └────────────┘     │
     │                │             │
     └────────┬───────┘             │
              │                     │
        ┌─────▼──────┐              │
        │ BaseTemplate│◄─────────────┘
        │   (Hooks)  │
        └────────────┘
```

**Code Flow**:

```go
// Template definition (interface)
type AdapterLifecycle interface {
    OnStart(ctx context.Context) error  // Step 1
    OnStop(ctx context.Context) error   // Step 2
}

// Template execution (BaseTemplate)
func BaseTemplate(lc fx.Lifecycle, impl AdapterLifecycle) {
    lc.Append(fx.Hook{
        OnStart: impl.OnStart,  // Call step 1
        OnStop:  impl.OnStop,   // Call step 2
    })
}

// Concrete implementation
type GinAdapter struct {
    BaseAdapter[Config]
}

func (g *GinAdapter) OnStart(ctx context.Context) error {
    // Custom implementation of step 1
    RegisterRouters(...)
    StartServer(...)
}

func (g *GinAdapter) OnStop(ctx context.Context) error {
    // Custom implementation of step 2
    StopServer(...)
}
```

**Sequence Diagram**:

```
Application    Fx          BaseTemplate    GinAdapter
    │           │                │              │
    │  Start    │                │              │
    │──────────►│                │              │
    │           │  Append Hook   │              │
    │           │───────────────►│              │
    │           │                │  OnStart()   │
    │           │                │─────────────►│
    │           │                │              │
    │           │                │   Register   │
    │           │                │   Routes     │
    │           │                │◄─────────────│
    │           │                │              │
    │           │                │   Start      │
    │           │                │   Server     │
    │           │                │◄─────────────│
    │  Running  │                │              │
    │◄──────────│                │              │
    │           │                │              │
    │  Stop     │                │              │
    │──────────►│                │              │
    │           │  OnStop Hook   │              │
    │           │───────────────►│              │
    │           │                │  OnStop()    │
    │           │                │─────────────►│
    │           │                │              │
    │           │                │   Shutdown   │
    │           │                │◄─────────────│
    │  Stopped  │                │              │
    │◄──────────│                │              │
```

---

### 2. Factory Pattern

**Intent**: Create objects without specifying exact class.

**Implementation**:

```
┌────────────────────────────────────┐
│         ForRoot Function           │
│         (Factory Method)           │
├────────────────────────────────────┤
│ + ForRoot(config...) fx.Option    │
└────────────┬───────────────────────┘
             │
             │ creates
             ▼
┌────────────────────────────────────┐
│        Fx Module                   │
├────────────────────────────────────┤
│ • Provides Dependencies            │
│ • Registers Lifecycle              │
│ • Injects Controllers              │
└────────────────────────────────────┘
```

**Code Example**:

```go
func ForRoot(port int, controllerGroup string) fx.Option {
    return fx.Module("gin-adapter",
        // Provide dependencies
        fx.Provide(
            func() *gin.Engine { return gin.Default() },
            func() int { return port },
            fx.Annotate(
                NewGinAdapter,
                fx.ParamTags(``, ``, `group:"controllers"`),
            ),
        ),
        // Register lifecycle
        fx.Invoke(func(lc fx.Lifecycle, adapter *GinAdapter) {
            adapter.RegisterLifecycle(lc, adapter)
        }),
    )
}
```

**Benefits**:
- ✅ Encapsulates complex Fx configuration
- ✅ Provides clean API for users
- ✅ Allows different adapter types with same interface

---

### 3. Reflection-based Registry Pattern

**Intent**: Auto-discover and register components at runtime.

**Architecture**:

```
┌───────────────────────────────────────────────────────────┐
│                   RegisterRouter                          │
│                                                            │
│  1. Reflect Controller                                    │
│     ┌──────────────────────────────────────────┐         │
│     │  value := reflect.ValueOf(controller)    │         │
│     │  valueType := value.Type()               │         │
│     └──────────────────────────────────────────┘         │
│                        │                                  │
│  2. Iterate Methods    ▼                                  │
│     ┌──────────────────────────────────────────┐         │
│     │  for i := 0; i < NumMethod(); i++ {     │         │
│     │      method := value.Method(i)           │         │
│     │      if isValid(method) {                │         │
│     │          call(method)                    │         │
│     │      }                                    │         │
│     │  }                                        │         │
│     └──────────────────────────────────────────┘         │
│                        │                                  │
│  3. Validate Signature ▼                                  │
│     ┌──────────────────────────────────────────┐         │
│     │  func(context.Context) ✓                 │         │
│     │  func(ctx, string) ✗                     │         │
│     │  func() ✗                                │         │
│     └──────────────────────────────────────────┘         │
│                        │                                  │
│  4. Call Method        ▼                                  │
│     ┌──────────────────────────────────────────┐         │
│     │  method.Call([]Value{ctx})               │         │
│     │  → route registered!                     │         │
│     └──────────────────────────────────────────┘         │
└───────────────────────────────────────────────────────────┘
```

**Validation Flow**:

```go
func isValidDynamicMethod(methodType reflect.Type) bool {
    // Must have exactly 1 input
    if methodType.NumIn() != 1 {
        return false  // ✗ No params or multiple params
    }
    
    // Must have 0 outputs
    if methodType.NumOut() != 0 {
        return false  // ✗ Has return values
    }
    
    // Input must be context.Context
    ctxType := reflect.TypeOf((*context.Context)(nil)).Elem()
    if methodType.In(0) != ctxType {
        return false  // ✗ Wrong parameter type
    }
    
    return true  // ✓ Valid signature
}
```

**Registration Timeline**:

```
Controller Methods (Declaration Order):
┌─────────────────────────────────────┐
│ func GetUsers(ctx)    ← Declared 1st │
│ func CreateUser(ctx)  ← Declared 2nd │
│ func DeleteUser(ctx)  ← Declared 3rd │
└─────────────────────────────────────┘
                │
                │ Reflection sorts alphabetically
                ▼
┌─────────────────────────────────────┐
│ func CreateUser(ctx)  ← Called 1st  │
│ func DeleteUser(ctx)  ← Called 2nd  │
│ func GetUsers(ctx)    ← Called 3rd  │
└─────────────────────────────────────┘
```

---

## 🔄 Component Interaction

### Component Dependency Graph

```
┌──────────────────────────────────────────────────────────────┐
│                                                               │
│   User Code (main.go, controllers)                           │
│                                                               │
└───────────────┬──────────────────────────────────────────────┘
                │
                │ uses
                ▼
┌──────────────────────────────────────────────────────────────┐
│                                                               │
│   Adapter Implementation (GinAdapter, GrpcAdapter)           │
│   ┌─────────────────────────────────────────────┐            │
│   │  • Embeds BaseAdapter[Config]               │            │
│   │  • Implements AdapterLifecycle              │            │
│   │  • Calls RegisterRouters                    │            │
│   └─────────────────────────────────────────────┘            │
│                                                               │
└───────────────┬──────────────────────────────────────────────┘
                │
                │ depends on
                ▼
┌──────────────────────────────────────────────────────────────┐
│                                                               │
│   adapter-template Package                                   │
│   ┌─────────────────────────────────────────────┐            │
│   │  BaseAdapter[T]      (Generic Template)    │            │
│   │  AdapterLifecycle    (Interface)           │            │
│   │  ICoreController     (Marker Interface)    │            │
│   │  RegisterRouter      (Reflection Engine)   │            │
│   │  AsRoute             (Fx Helper)           │            │
│   └─────────────────────────────────────────────┘            │
│                                                               │
└───────────────┬──────────────────────────────────────────────┘
                │
                │ depends on
                ▼
┌──────────────────────────────────────────────────────────────┐
│                                                               │
│   External Dependencies                                       │
│   • go.uber.org/fx       (Dependency Injection)              │
│   • context              (Context propagation)               │
│   • reflect              (Runtime type inspection)           │
│   • fmt                  (Error formatting)                  │
│                                                               │
└──────────────────────────────────────────────────────────────┘
```

### Data Flow

```
┌─────────────────────────────────────────────────────────────────┐
│  Phase 1: Application Initialization                            │
└─────────────────────────────────────────────────────────────────┘
    │
    │  fx.New(...)
    │  ├─ ForRoot(8080, "controllers")  ← Factory creates module
    │  └─ ControllerModules...          ← Controllers provided
    │
    ▼
┌─────────────────────────────────────────────────────────────────┐
│  Phase 2: Dependency Resolution                                 │
└─────────────────────────────────────────────────────────────────┘
    │
    │  Fx builds dependency graph:
    │  ┌──────────────────────────────────────┐
    │  │  Controllers (group injection)       │
    │  │    ├─ UserController                 │
    │  │    ├─ HealthController               │
    │  │    └─ ProductController              │
    │  └──────────────────────────────────────┘
    │           │
    │           │ injected into
    │           ▼
    │  ┌──────────────────────────────────────┐
    │  │  Adapter (with controllers)          │
    │  │    Config.Controllers = [...]        │
    │  └──────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────────────────────────────────┐
│  Phase 3: Lifecycle Registration                                │
└─────────────────────────────────────────────────────────────────┘
    │
    │  adapter.RegisterLifecycle(lc, adapter)
    │    └─> BaseTemplate(lc, adapter)
    │          └─> lc.Append(fx.Hook{...})
    │
    ▼
┌─────────────────────────────────────────────────────────────────┐
│  Phase 4: Application Start (OnStart)                           │
└─────────────────────────────────────────────────────────────────┘
    │
    │  Fx calls: adapter.OnStart(ctx)
    │    │
    │    ├─> RegisterRouters(controllers, ctx)
    │    │     │
    │    │     ├─> RegisterRouter(UserController, ctx)
    │    │     │     ├─ Reflect methods
    │    │     │     ├─ Call CreateUser(ctx)  → route registered
    │    │     │     ├─ Call GetUsers(ctx)    → route registered
    │    │     │     └─ Call UpdateUser(ctx)  → route registered
    │    │     │
    │    │     ├─> RegisterRouter(HealthController, ctx)
    │    │     │     └─ Call Health(ctx)      → route registered
    │    │     │
    │    │     └─> Return nil or error
    │    │
    │    └─> Start server/service
    │
    ▼
┌─────────────────────────────────────────────────────────────────┐
│  Phase 5: Running                                                │
└─────────────────────────────────────────────────────────────────┘
    │
    │  Application handles requests...
    │
    ▼
┌─────────────────────────────────────────────────────────────────┐
│  Phase 6: Application Stop (OnStop)                             │
└─────────────────────────────────────────────────────────────────┘
    │
    │  Fx calls: adapter.OnStop(ctx)
    │    └─> Graceful shutdown
    │          └─> Cleanup resources
    │
    ▼
   Exit
```

---

## 🧩 Module Structure

### Package Organization

```
adapter-template/
│
├── doc.go                    # Package-level documentation
├── README.md                 # User guide
├── ARCHITECTURE.md          # This file
│
├── Core Components
│   ├── base_adapter.go       # BaseAdapter + AdapterLifecycle
│   ├── dynamic_controller.go # ICoreController interface
│   ├── fx_annotations.go     # AsRoute helper
│   └── router_registry.go    # RegisterRouter logic
│
├── Tests
│   └── dynamic_controller_test.go  # Comprehensive tests
│
└── Documentation
    └── examples/
        ├── README.md
        ├── simple_adapter.go      # Basic example
        └── validation_adapter.go  # Advanced example
```

### API Surface

**Public Exports** (5 functions + 3 types):

```go
// Functions
func BaseTemplate(lc fx.Lifecycle, impl AdapterLifecycle)
func (b *BaseAdapter[T]) RegisterLifecycle(...)
func RegisterRouter(controller ICoreController, ctx context.Context) error
func RegisterRouters(controllers []ICoreController, ctx context.Context) error
func AsRoute(f any, groupTag string, annotation ...fx.Annotation) any

// Types
type AdapterLifecycle interface { ... }
type BaseAdapter[T any] struct { ... }
type ICoreController interface {}
```

**Private Functions** (1 function):

```go
func isValidDynamicMethod(methodType reflect.Type) bool
```

---

## 🔒 Design Decisions

### Decision 1: Empty ICoreController Interface

**Options Considered**:
```go
// Option A: Empty interface (chosen)
type ICoreController interface {}

// Option B: Marker method
type ICoreController interface {
    MarkAsCoreController()
}

// Option C: Registry method
type ICoreController interface {
    RegisterRoutes(ctx context.Context)
}
```

**Chosen**: Option A

**Rationale**:
- ✅ Maximum flexibility
- ✅ No boilerplate for users
- ✅ Reflection does the work
- ⚠️ Trade-off: Runtime validation only

---

### Decision 2: Fail-Fast Error Handling

**Options Considered**:
```go
// Option A: Fail-fast (chosen)
for _, controller := range controllers {
    if err := RegisterRouter(controller, ctx); err != nil {
        return err  // Stop immediately
    }
}

// Option B: Error aggregation
var errs []error
for _, controller := range controllers {
    if err := RegisterRouter(controller, ctx); err != nil {
        errs = append(errs, err)  // Continue
    }
}
return errors.Join(errs...)
```

**Chosen**: Option A

**Rationale**:
- ✅ Production safety - don't start with partial registration
- ✅ Clear failure point
- ✅ Easier debugging
- ⚠️ Trade-off: Less forgiving for non-critical errors

---

### Decision 3: Context Fallback to Background

**Code**:
```go
if ctx == nil {
    ctx = context.Background()
}
```

**Rationale**:
- ✅ Safe default behavior
- ✅ No nil pointer panics
- ✅ Allows nil-safe calling
- ⚠️ May hide bugs where context should be required

---

## 📊 Performance Characteristics

### Time Complexity

| Operation | Complexity | Notes |
|-----------|-----------|-------|
| RegisterRouter | O(N) | N = methods in controller |
| RegisterRouters | O(C × M) | C = controllers, M = avg methods |
| isValidDynamicMethod | O(1) | Type comparison only |
| Reflection overhead | ~5-10μs | Per method call |

### Space Complexity

| Structure | Space | Notes |
|-----------|-------|-------|
| BaseAdapter[T] | O(1) + sizeof(T) | Just config field |
| Method reflection | O(N) | Temporary during registration |
| Error strings | O(E) | E = number of errors |

### Benchmark Estimates

```
Scenario: 100 controllers × 10 methods each = 1000 method calls

Time breakdown:
  Reflection setup:       ~100μs
  Method validation:      ~50μs
  Method calls:           ~500μs
  Total:                  ~650μs

Memory:
  Reflection metadata:    ~10KB (temporary)
  Error strings:          ~1KB (if errors)
```

**Conclusion**: Negligible overhead for startup (one-time cost).

---

## 🔮 Extension Points

### Adding New Adapter Types

```go
// 1. Define config
type MyConfig struct {
    Setting1 string
    Controllers []ICoreController
}

// 2. Create adapter
type MyAdapter struct {
    BaseAdapter[MyConfig]
}

// 3. Implement lifecycle
func (m *MyAdapter) OnStart(ctx context.Context) error {
    RegisterRouters(m.Config.Controllers, ctx)
    // Start your service
}

// 4. Create factory
func ForRoot(...) fx.Option { ... }
```

### Custom Controller Validation

```go
// Wrap RegisterRouter with custom logic
func RegisterValidatedRouter(controller ICoreController, ctx context.Context) error {
    // Pre-validation
    if validator, ok := controller.(interface{ Validate() error }); ok {
        if err := validator.Validate(); err != nil {
            return err
        }
    }
    
    // Standard registration
    return RegisterRouter(controller, ctx)
}
```

---

## 🎓 Learning Resources

- [Template Method Pattern](https://refactoring.guru/design-patterns/template-method)
- [Factory Pattern](https://refactoring.guru/design-patterns/factory-method)
- [Uber Fx Tutorial](https://uber-go.github.io/fx/get-started/)
- [Go Reflection](https://go.dev/blog/laws-of-reflection)

---

**Maintained by**: adapter-template team  
**Last Updated**: 2025-10-31
