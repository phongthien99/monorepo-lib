# Log Library

A flexible, interface-based logging library for Go that follows SOLID principles. This library provides a unified logging interface with multiple adapter implementations.

## Features

- ✅ **Interface Segregation Principle (ISP)** - Small, focused interfaces
- ✅ **Multiple logging styles** - Basic, formatted, structured, and line-based logging
- ✅ **Context support** - For distributed tracing
- ✅ **Multiple adapters** - Currently supports Zap, easy to add more
- ✅ **Type-safe** - Compile-time type checking
- ✅ **Flexible configuration** - Config struct or functional options

## Installation

```bash
go get github.com/phongthien99/monorepo-lib/libs/log
```

## Quick Start

### Basic Usage

```go
package main

import (
    "github.com/phongthien99/monorepo-lib/libs/log/adapter/zap"
)

func main() {
    // Create a development logger
    logger, _ := zap.NewDevelopment()
    defer logger.Sync()

    // Basic logging
    logger.Info("Hello, World!")
    logger.Debug("Debug message")
    logger.Warn("Warning message")
    logger.Error("Error message")
}
```

### Formatted Logging

```go
logger.Infof("User %s logged in at %d", username, timestamp)
logger.Errorf("Failed to connect: %v", err)
```

### Structured Logging (Recommended)

```go
logger.Infow("User login",
    "username", "john",
    "ip", "192.168.1.1",
    "timestamp", time.Now(),
)

logger.Errorw("Database error",
    "error", err,
    "query", sql,
    "duration_ms", elapsed,
)
```

### Contextual Logging

```go
// Add fields to all subsequent logs
userLogger := logger.With("user_id", 123, "session", "abc")
userLogger.Info("User action") // Includes user_id and session

// Named loggers
apiLogger := logger.Named("api")
apiLogger.Info("API request") // Includes logger name
```

### Context Support

```go
// For distributed tracing
ctx := context.Background()
ctxLogger := logger.WithContext(ctx)
ctxLogger.Info("Request processed")
```

## Factory Functions

### Pre-configured Loggers

```go
// Development: console output, debug level, human-readable
devLogger, _ := zap.NewDevelopment()

// Production: JSON output, info level, optimized
prodLogger, _ := zap.NewProduction()

// Example: for testing
exampleLogger := zap.NewExample()

// No-op: discards all logs
nopLogger := zap.NewNop()
```

### Custom Configuration

```go
import "github.com/phongthien99/monorepo-lib/libs/log/core"

cfg := zap.Config{
    Level:            core.WarnLevel,
    Development:      false,
    Encoding:         "json",
    OutputPaths:      []string{"stdout", "/var/log/app.log"},
    ErrorOutputPaths: []string{"stderr"},
}

logger, _ := zap.NewWithConfig(cfg)
```

### Functional Options (Recommended)

```go
logger, _ := zap.NewWithOptions(
    zap.WithLevel(core.DebugLevel),
    zap.WithConsoleEncoding(),
    zap.WithOutputPaths("stdout", "/tmp/app.log"),
)

// Or start with development config
logger, _ := zap.NewDevelopmentWithOptions(
    zap.WithLevel(core.InfoLevel),
)

// Or production config
logger, _ := zap.NewProductionWithOptions(
    zap.WithConsoleEncoding(),
)
```

## Log Levels

```go
core.DebugLevel   // Detailed information for debugging
core.InfoLevel    // General informational messages
core.WarnLevel    // Warning messages
core.ErrorLevel   // Error messages
core.DPanicLevel  // Panic in development, error in production
core.PanicLevel   // Panic
core.FatalLevel   // Fatal, then os.Exit(1)
```

## Interfaces

The library uses Interface Segregation Principle (ISP) to provide focused interfaces:

### Basic Interfaces

- `IBasicLogger` - Simple logging (Debug, Info, Warn, Error, etc.)
- `IFormattedLogger` - Printf-style logging (Debugf, Infof, etc.)
- `IStructuredLogger` - Structured logging with key-value pairs (Debugw, Infow, etc.)
- `ILineLogger` - Line-style logging (Debugln, Infoln, etc.)
- `IContextualLogger` - Contextual logging (With, WithLazy, Named)
- `IContextLogger` - Context support (WithContext)
- `ILoggerControl` - Logger control (Level, Sync, Desugar)

### Composite Interfaces

- `ILogger` - Basic + Formatted + Control (minimal logger)
- `IFullLogger` - Basic + Formatted + Structured + Control
- `ISugaredLogger` - All interfaces (complete logger)

## Dependency Injection

Use specific interfaces based on your needs:

```go
type UserService struct {
    logger core.ILogger // Only needs basic + formatted logging
}

type MetricsService struct {
    logger core.IStructuredLogger // Only needs structured logging
}

type APIHandler struct {
    logger core.ISugaredLogger // Needs all features
}
```

## Best Practices

### 1. Use Structured Logging in Production

```go
// ❌ Bad
logger.Infof("User %s performed action %s", user, action)

// ✅ Good
logger.Infow("User action",
    "user", user,
    "action", action,
)
```

### 2. Add Context to Loggers

```go
// Create a child logger with context
requestLogger := logger.With(
    "request_id", requestID,
    "method", r.Method,
    "path", r.URL.Path,
)

// Use throughout request handling
requestLogger.Info("Processing request")
requestLogger.Errorw("Failed to process", "error", err)
```

### 3. Use Named Loggers for Components

```go
apiLogger := logger.Named("api")
dbLogger := logger.Named("database")
cacheLogger := logger.Named("cache")
```

### 4. Always Sync Before Exit

```go
func main() {
    logger, _ := zap.NewProduction()
    defer logger.Sync() // Flush buffered logs

    // Your code
}
```

## Adding New Adapters

To add support for other logging libraries (logrus, zerolog, etc.):

1. Create `adapter/<library>/adapter.go`
2. Implement `core.ISugaredLogger` interface
3. Create factory functions in `adapter/<library>/factory.go`

Example structure:
```
adapter/
├── zap/
│   ├── adapter.go
│   ├── factory.go
│   ├── options.go
│   └── converter.go
└── logrus/     # Your new adapter
    ├── adapter.go
    └── factory.go
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific adapter tests
go test ./adapter/zap -v
```

## License

MIT License

## Contributing

Contributions are welcome! Please ensure:
- All tests pass
- Code follows Go conventions
- New features include tests
- SOLID principles are maintained
