package adaptertemplate

import (
	"context"

	"go.uber.org/fx"
)

// AdapterLifecycle định nghĩa hành vi của module (Template Method interface)
type AdapterLifecycle interface {
	OnStart(ctx context.Context) error
	OnStop(ctx context.Context) error
}

// BaseTemplate đăng ký OnStart/OnStop với Fx lifecycle
//
// Parameters:
//   - lc: Fx lifecycle để đăng ký hooks
//   - impl: Implementation của AdapterLifecycle interface
//
// Panics:
//   - Nếu lc hoặc impl là nil
func BaseTemplate(lc fx.Lifecycle, impl AdapterLifecycle) {
	if lc == nil {
		panic("fx.Lifecycle cannot be nil")
	}
	if impl == nil {
		panic("AdapterLifecycle implementation cannot be nil")
	}

	lc.Append(fx.Hook{
		OnStart: impl.OnStart,
		OnStop:  impl.OnStop,
	})
}

// BaseAdapter generic: gom Config + hỗ trợ lifecycle chung
type BaseAdapter[T any] struct {
	Config T
}

// RegisterLifecycle đăng ký adapter lifecycle với Fx
// Method này add validation layer trên BaseTemplate
//
// Parameters:
//   - lc: Fx lifecycle để đăng ký hooks
//   - impl: Implementation của AdapterLifecycle interface
//
// Panics:
//   - Nếu lc hoặc impl là nil
func (b *BaseAdapter[T]) RegisterLifecycle(lc fx.Lifecycle, impl AdapterLifecycle) {
	BaseTemplate(lc, impl)
}
