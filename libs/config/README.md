# Config Library

A flexible, type-safe configuration management library for Go with support for multiple configuration sources and customizable merge strategies.

## Features

- **Type-Safe**: Generic-based API ensures compile-time type safety
- **Multiple Sources**: Load configuration from files, environment variables, and command-line flags
- **Flexible Merging**: Deep merge by default, with support for custom merge strategies
- **Validation**: Built-in validation support with composable validators
- **Priority System**: Control configuration precedence through loader order
- **Zero Dependencies** (core): Core functionality has minimal dependencies

## Installation

```bash
go get github.com/phongthien99/monorepo-lib/libs/config
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/phongthien99/monorepo-lib/libs/config"
    "github.com/phongthien99/monorepo-lib/libs/config/loader"
)

// Define your configuration structure
type AppConfig struct {
    Server struct {
        Host string `mapstructure:"host"`
        Port int    `mapstructure:"port"`
    } `mapstructure:"server"`
    Database struct {
        URL      string `mapstructure:"url"`
        MaxConns int    `mapstructure:"max_conns"`
    } `mapstructure:"database"`
}

func main() {
    // Create loaders (order matters - last loader has highest priority)
    fileLoader := loader.NewFileLoader("config.yaml", "yaml")
    envLoader := loader.NewEnvLoader("APP").WithAutoKeys(AppConfig{})
    
    // Create config with loaders
    cfg := config.New[AppConfig](
        fileLoader,  // Lowest priority
        envLoader,   // Highest priority - will override file config
    )
    
    // Load configuration
    if err := cfg.Load(); err != nil {
        log.Fatal(err)
    }
    
    // Get configuration
    appConfig := cfg.Get()
    fmt.Printf("Server: %s:%d\n", appConfig.Server.Host, appConfig.Server.Port)
}
```

## Loaders

### File Loader

Load configuration from files. Supports JSON, YAML, TOML, Properties, and HCL formats.

```go
// YAML file
fileLoader := loader.NewFileLoader("config.yaml", "yaml")

// JSON file
jsonLoader := loader.NewFileLoader("config.json", "json")

// TOML file
tomlLoader := loader.NewFileLoader("config.toml", "toml")
```

**Example config.yaml:**
```yaml
server:
  host: localhost
  port: 8080
database:
  url: postgres://localhost/mydb
  max_conns: 10
```

### Environment Variable Loader

Load configuration from environment variables with automatic key mapping.

```go
// With prefix "APP" - reads APP_SERVER_HOST, APP_SERVER_PORT, etc.
envLoader := loader.NewEnvLoader("APP").WithAutoKeys(AppConfig{})

// Without prefix - reads SERVER_HOST, SERVER_PORT, etc.
envLoader := loader.NewEnvLoader("").WithAutoKeys(AppConfig{})

// Manual key specification
envLoader := loader.NewEnvLoader("APP").
    WithKeys("server.host", "server.port", "database.url")
```

**Environment variable mapping:**
- Prefix is automatically uppercased: `"app"` → `"APP_"`
- Dots (`.`) are converted to underscores (`_`): `server.host` → `SERVER_HOST`
- Full example: With prefix `"APP"`, the key `server.host` maps to `APP_SERVER_HOST`

**Example:**
```bash
export APP_SERVER_HOST=0.0.0.0
export APP_SERVER_PORT=9090
export APP_DATABASE_URL=postgres://prod-db/mydb
```

### Command-Line Flag Loader

Load configuration from command-line flags using pflag.

```go
import "github.com/spf13/pflag"

// Define flags
pflag.String("server.host", "localhost", "Server host")
pflag.Int("server.port", 8080, "Server port")
pflag.Parse()

// Create loader
flagLoader := loader.NewFlagLoader(nil) // nil uses global pflag.CommandLine
```

**Usage:**
```bash
./myapp --server.host=0.0.0.0 --server.port=9090
```

## Merge Strategies

### Default Merge (Deep Merge)

By default, the library uses deep merge strategy that intelligently combines configurations:

```go
cfg := config.New[AppConfig](loaders...) // Uses DefaultMerge
```

**Merge rules:**
- **Struct fields**: Merged recursively, non-zero values override
- **Slices**: Completely replaced if source slice is not empty
- **Maps**: Deep merge of keys
- **Pointers**: Merged recursively if source is not nil
- **Primitives**: Overridden if source is not zero value

**Example:**
```go
// File config:
{
    "server": {"host": "localhost", "port": 8080},
    "database": {"url": "postgres://localhost/db"}
}

// Environment: APP_SERVER_PORT=9090
// Result: {"host": "localhost", "port": 9090, "url": "postgres://localhost/db"}
```

### Shallow Merge

Replace entire struct instead of deep merging:

```go
cfg := config.New[AppConfig](loaders...).
    WithMerge(config.ShallowMerge[AppConfig])
```

### Custom Merge Strategy

Define your own merge logic:

```go
func customMerge(dst, src *AppConfig) error {
    // Your custom merge logic
    if src.Server.Port != 0 {
        dst.Server.Port = src.Server.Port
    }
    // ... more custom logic
    return nil
}

cfg := config.New[AppConfig](loaders...).
    WithMerge(customMerge)
```

## Validation

### Basic Validation

```go
type AppConfigValidator struct{}

func (v *AppConfigValidator) Validate(cfg *AppConfig) error {
    if cfg.Server.Port < 1024 {
        return fmt.Errorf("port must be >= 1024")
    }
    if cfg.Database.MaxConns <= 0 {
        return fmt.Errorf("max_conns must be positive")
    }
    return nil
}

cfg := config.New[AppConfig](loaders...).
    WithValidator(&AppConfigValidator{})
```

### Function-Based Validation

```go
validateFunc := func(cfg *AppConfig) error {
    if cfg.Server.Port < 1024 {
        return fmt.Errorf("port must be >= 1024")
    }
    return nil
}

cfg := config.New[AppConfig](loaders...).
    WithValidator(config.ValidatorFunc[AppConfig](validateFunc))
```

### Composite Validation

Combine multiple validators:

```go
portValidator := config.ValidatorFunc[AppConfig](func(cfg *AppConfig) error {
    if cfg.Server.Port < 1024 {
        return fmt.Errorf("invalid port")
    }
    return nil
})

dbValidator := config.ValidatorFunc[AppConfig](func(cfg *AppConfig) error {
    if cfg.Database.MaxConns <= 0 {
        return fmt.Errorf("invalid max_conns")
    }
    return nil
})

validator := config.NewCompositeValidator(portValidator, dbValidator)

cfg := config.New[AppConfig](loaders...).
    WithValidator(validator)
```

## Configuration Priority

Loaders are processed in order, with later loaders having higher priority:

```go
cfg := config.New[AppConfig](
    fileLoader,    // Priority 1 (lowest)
    envLoader,     // Priority 2
    flagLoader,    // Priority 3 (highest)
)
```

In this example:
1. File configuration is loaded first
2. Environment variables override file values
3. Command-line flags override both file and environment values

## Complete Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/phongthien99/monorepo-lib/libs/config"
    "github.com/phongthien99/monorepo-lib/libs/config/loader"
    "github.com/spf13/pflag"
)

type AppConfig struct {
    Server struct {
        Host string `mapstructure:"host"`
        Port int    `mapstructure:"port"`
    } `mapstructure:"server"`
    Database struct {
        URL      string `mapstructure:"url"`
        MaxConns int    `mapstructure:"max_conns"`
    } `mapstructure:"database"`
    LogLevel string `mapstructure:"log_level"`
}

func main() {
    // Define command-line flags
    pflag.String("server.host", "localhost", "Server host")
    pflag.Int("server.port", 8080, "Server port")
    pflag.String("log_level", "info", "Log level")
    pflag.Parse()
    
    // Create loaders
    fileLoader := loader.NewFileLoader("config.yaml", "yaml")
    envLoader := loader.NewEnvLoader("APP").WithAutoKeys(AppConfig{})
    flagLoader := loader.NewFlagLoader(nil)
    
    // Create validator
    validator := config.ValidatorFunc[AppConfig](func(cfg *AppConfig) error {
        if cfg.Server.Port < 1024 || cfg.Server.Port > 65535 {
            return fmt.Errorf("port must be between 1024 and 65535")
        }
        if cfg.Database.MaxConns <= 0 {
            return fmt.Errorf("max_conns must be positive")
        }
        return nil
    })
    
    // Create and load config
    cfg := config.New[AppConfig](
        fileLoader,
        envLoader,
        flagLoader,
    ).WithValidator(validator)
    
    if err := cfg.Load(); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Use configuration
    appConfig := cfg.Get()
    fmt.Printf("Server running on %s:%d\n", 
        appConfig.Server.Host, 
        appConfig.Server.Port)
    fmt.Printf("Database: %s (max %d connections)\n",
        appConfig.Database.URL,
        appConfig.Database.MaxConns)
    fmt.Printf("Log level: %s\n", appConfig.LogLevel)
}
```

**Run with:**
```bash
# Use defaults from config.yaml
./myapp

# Override with environment variables
APP_SERVER_PORT=9090 ./myapp

# Override with flags
./myapp --server.port=9090 --log_level=debug

# Combine all sources
APP_DATABASE_URL=postgres://prod/db ./myapp --server.port=9090
```

## Advanced Usage

### Method Chaining

```go
cfg := config.New[AppConfig](loaders...).
    WithMerge(customMergeFunc).
    WithValidator(validator).
    Load()
```

### Getting Configuration

```go
// Get by value
appConfig := cfg.Get()

// Get by pointer (useful for modifications)
appConfigPtr := cfg.GetPtr()
```

### Custom Struct Tags

The library uses `mapstructure` tags for field mapping:

```go
type Config struct {
    ServerHost string `mapstructure:"server_host"` // Maps to server_host
    ServerPort int    `mapstructure:"port"`        // Maps to port
    Internal   string `mapstructure:"-"`           // Ignored
}
```

## Error Handling

The library provides detailed error messages:

```go
if err := cfg.Load(); err != nil {
    // Errors include context about which loader failed
    log.Printf("Config error: %v", err)
}
```

Common errors:
- `loader[N] failed`: Loader at index N failed to load
- `merge loader[N] failed`: Failed to merge data from loader N
- `config validation failed`: Validation failed after loading

## Best Practices

1. **Order loaders by priority**: Place lowest priority first, highest priority last
2. **Use AutoKeys with environment loader**: Automatically extracts all keys from your struct
3. **Always validate**: Add validation to catch configuration errors early
4. **Use struct tags**: Use `mapstructure` tags for clear field mapping
5. **Handle errors**: Always check errors from `Load()`

## License

[Your License Here]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
