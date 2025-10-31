package examples

import (
	"context"
	"fmt"
	"log"

	adaptertemplate "github.com/phongthien99/monorepo-lib/libs/core/adapter-template"
	"go.uber.org/fx"
)

// Example 2: Adapter with Validation
// This example shows how to add validation to your adapter

// ValidatedConfig demonstrates config validation
type ValidatedConfig struct {
	Port        int
	ServiceName string
	MaxRetries  int
	Controllers []adaptertemplate.ICoreController
}

// Validate checks if the config is valid
func (c *ValidatedConfig) Validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid port: %d (must be 1-65535)", c.Port)
	}
	if c.ServiceName == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if c.MaxRetries < 0 {
		return fmt.Errorf("max retries cannot be negative: %d", c.MaxRetries)
	}
	return nil
}

// ValidatedAdapter demonstrates config validation
type ValidatedAdapter struct {
	adaptertemplate.BaseAdapter[ValidatedConfig]
}

// NewValidatedAdapter creates a validated adapter
func NewValidatedAdapter(
	port int,
	serviceName string,
	maxRetries int,
	controllers []adaptertemplate.ICoreController,
) (*ValidatedAdapter, error) {
	config := ValidatedConfig{
		Port:        port,
		ServiceName: serviceName,
		MaxRetries:  maxRetries,
		Controllers: controllers,
	}

	// Validate config before creating adapter
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &ValidatedAdapter{
		BaseAdapter: adaptertemplate.BaseAdapter[ValidatedConfig]{
			Config: config,
		},
	}, nil
}

// OnStart implements AdapterLifecycle with validation
func (v *ValidatedAdapter) OnStart(ctx context.Context) error {
	log.Printf("🚀 Starting %s on port %d (max retries: %d)",
		v.Config.ServiceName,
		v.Config.Port,
		v.Config.MaxRetries,
	)

	// Double-check validation (defense in depth)
	if err := v.Config.Validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	// Register controllers with fail-fast
	if err := adaptertemplate.RegisterRouters(v.Config.Controllers, ctx); err != nil {
		return fmt.Errorf("controller registration failed: %w", err)
	}

	log.Printf("✅ %s started successfully", v.Config.ServiceName)
	return nil
}

// OnStop implements graceful shutdown
func (v *ValidatedAdapter) OnStop(ctx context.Context) error {
	log.Printf("🧹 Stopping %s gracefully...", v.Config.ServiceName)

	// Check if context already cancelled
	select {
	case <-ctx.Done():
		log.Printf("⚠️  Shutdown timeout exceeded, forcing stop")
		return ctx.Err()
	default:
	}

	log.Printf("✅ %s stopped successfully", v.Config.ServiceName)
	return nil
}

// ForRootValidated creates an Fx module with validation
func ForRootValidated(port int, serviceName string, maxRetries int, controllerGroup string) fx.Option {
	if controllerGroup == "" {
		controllerGroup = "validatedControllers"
	}

	return fx.Module("validated-adapter",
		fx.Provide(
			func() int { return port },
			func() string { return serviceName },
			func() int { return maxRetries },
			fx.Annotate(
				NewValidatedAdapter,
				fx.ParamTags(``, ``, ``, fmt.Sprintf(`group:"%s"`, controllerGroup)),
			),
		),
		fx.Invoke(func(lc fx.Lifecycle, adapter *ValidatedAdapter, err error) error {
			if err != nil {
				return fmt.Errorf("adapter creation failed: %w", err)
			}
			adapter.RegisterLifecycle(lc, adapter)
			return nil
		}),
	)
}

// Example: Controller with validation

// ValidatedController demonstrates controller-level validation
type ValidatedController struct {
	minVersion string
}

var _ adaptertemplate.ICoreController = (*ValidatedController)(nil)

// NewValidatedController creates a controller with validation
func NewValidatedController(minVersion string) (adaptertemplate.ICoreController, error) {
	if minVersion == "" {
		return nil, fmt.Errorf("minVersion cannot be empty")
	}
	return &ValidatedController{minVersion: minVersion}, nil
}

// RegisterAPIv1 registers v1 endpoints
func (v *ValidatedController) RegisterAPIv1(ctx context.Context) {
	log.Printf("✅ API v1 registered (min version: %s)", v.minVersion)
}

// RegisterAPIv2 registers v2 endpoints
func (v *ValidatedController) RegisterAPIv2(ctx context.Context) {
	log.Printf("✅ API v2 registered (min version: %s)", v.minVersion)
}

// ValidatedControllerModule is the Fx module with error handling
var ValidatedControllerModule = fx.Module("validated-controller",
	fx.Provide(
		fx.Annotate(
			func() (adaptertemplate.ICoreController, error) {
				return NewValidatedController("1.0.0")
			},
			fx.ResultTags(`group:"validatedControllers"`),
		),
	),
)

// Usage example:
//
//	func main() {
//	    app := fx.New(
//	        ForRootValidated(8080, "MyService", 3, "validatedControllers"),
//	        ValidatedControllerModule,
//	        fx.WithLogger(func() fxevent.Logger {
//	            return &fxevent.ConsoleLogger{W: os.Stdout}
//	        }),
//	    )
//
//	    // Will fail if validation errors occur
//	    if err := app.Err(); err != nil {
//	        log.Fatalf("Application failed to start: %v", err)
//	    }
//
//	    app.Run()
//	}
