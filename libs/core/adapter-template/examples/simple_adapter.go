package examples

import (
	"context"
	"fmt"
	"log"

	adaptertemplate "github.com/phongthien99/monorepo-lib/libs/core/adapter-template"
	"go.uber.org/fx"
)

// Example 1: Simple Adapter
// This example shows the minimal code needed to create an adapter

// SimpleConfig holds configuration for the simple adapter
type SimpleConfig struct {
	Name        string
	Controllers []adaptertemplate.ICoreController
}

// SimpleAdapter demonstrates a minimal adapter implementation
type SimpleAdapter struct {
	adaptertemplate.BaseAdapter[SimpleConfig]
}

// NewSimpleAdapter creates a new simple adapter instance
func NewSimpleAdapter(name string, controllers []adaptertemplate.ICoreController) *SimpleAdapter {
	return &SimpleAdapter{
		BaseAdapter: adaptertemplate.BaseAdapter[SimpleConfig]{
			Config: SimpleConfig{
				Name:        name,
				Controllers: controllers,
			},
		},
	}
}

// OnStart implements AdapterLifecycle.OnStart
func (s *SimpleAdapter) OnStart(ctx context.Context) error {
	log.Printf("🚀 Starting %s adapter", s.Config.Name)

	// Register all dynamic controllers
	if err := adaptertemplate.RegisterRouters(s.Config.Controllers, ctx); err != nil {
		return fmt.Errorf("failed to register controllers: %w", err)
	}

	log.Printf("✅ %s adapter started successfully", s.Config.Name)
	return nil
}

// OnStop implements AdapterLifecycle.OnStop
func (s *SimpleAdapter) OnStop(ctx context.Context) error {
	log.Printf("🧹 Stopping %s adapter", s.Config.Name)
	// Cleanup logic here
	log.Printf("✅ %s adapter stopped successfully", s.Config.Name)
	return nil
}

// ForRoot creates an Fx module for the simple adapter
func ForRoot(name string, controllerGroup string) fx.Option {
	if controllerGroup == "" {
		controllerGroup = "simpleControllers"
	}

	return fx.Module("simple-adapter",
		fx.Provide(
			func() string { return name },
			fx.Annotate(
				NewSimpleAdapter,
				fx.ParamTags(``, fmt.Sprintf(`group:"%s"`, controllerGroup)),
			),
		),
		fx.Invoke(func(lc fx.Lifecycle, adapter *SimpleAdapter) {
			adapter.RegisterLifecycle(lc, adapter)
		}),
	)
}

// Example Controller

// PrintController demonstrates a simple controller
type PrintController struct {
	prefix string
}

var _ adaptertemplate.ICoreController = (*PrintController)(nil)

// NewPrintController creates a new print controller
func NewPrintController(prefix string) adaptertemplate.ICoreController {
	return &PrintController{prefix: prefix}
}

// HelloWorld will be auto-called by RegisterRouter
func (p *PrintController) HelloWorld(ctx context.Context) {
	log.Printf("%s: Hello World registered!", p.prefix)
}

// GoodbyeWorld will be auto-called by RegisterRouter
func (p *PrintController) GoodbyeWorld(ctx context.Context) {
	log.Printf("%s: Goodbye World registered!", p.prefix)
}

// PrintControllerModule is the Fx module for PrintController
var PrintControllerModule = fx.Module("print-controller",
	fx.Provide(
		fx.Annotate(
			func() adaptertemplate.ICoreController {
				return NewPrintController("PRINT")
			},
			fx.ResultTags(`group:"simpleControllers"`),
		),
	),
)

// Usage example:
//
//	func main() {
//	    app := fx.New(
//	        ForRoot("MySimpleAdapter", "simpleControllers"),
//	        PrintControllerModule,
//	    )
//	    app.Run()
//	}
