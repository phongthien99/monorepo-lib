package adaptertemplate

import (
	"context"
	"fmt"
	"reflect"
)

// RegisterRouter sử dụng reflection để tự động gọi tất cả methods
// của controller có signature: func(context.Context)
//
// Parameters:
//   - controller: Instance của controller (phải implement ICoreController interface)
//   - ctx: Context được truyền vào mỗi method call. Nếu nil, sẽ dùng context.Background()
//
// Returns:
//   - error: Error ngay khi có method bị panic hoặc lỗi (fail-fast)
//
// Behavior:
//   - Quét tất cả exported methods của controller
//   - Chỉ gọi methods có đúng signature: func(context.Context)
//   - Skip methods không đúng signature
//   - Truyền context vào mỗi method call
//   - Recover từ panic và return error ngay lập tức
//   - DỪNG NGAY khi 1 method fail (fail-fast pattern)
//
// Example:
//
//	type UserController struct {
//	    router *gin.Engine
//	}
//
//	func (u *UserController) GetUsers(ctx context.Context) {
//	    u.router.GET("/users", handler)
//	}
//
//	func (u *UserController) CreateUser(ctx context.Context) {
//	    u.router.POST("/users", handler)
//	}
//
//	// RegisterRouter sẽ tự động gọi cả GetUsers và CreateUser
//	// Nếu 1 method fail, sẽ dừng ngay và return error
//	ctx := context.WithTimeout(context.Background(), 5*time.Second)
//	if err := RegisterRouter(controller, ctx); err != nil {
//	    log.Fatalf("Failed to register routes: %v", err)
//	}
func RegisterRouter(controller ICoreController, ctx context.Context) error {
	if controller == nil {
		return nil
	}

	// Sử dụng context.Background() nếu ctx nil
	if ctx == nil {
		ctx = context.Background()
	}

	value := reflect.ValueOf(controller)
	valueType := value.Type()

	// Iterate qua tất cả methods của controller
	for i := 0; i < value.NumMethod(); i++ {
		method := value.Method(i)
		methodType := method.Type()
		methodName := valueType.Method(i).Name

		// Validate method signature: func(context.Context)
		if !isValidDynamicMethod(methodType) {
			// Skip methods không đúng signature
			continue
		}

		// Recover từ panic và return error ngay lập tức
		var panicErr error
		func() {
			defer func() {
				if r := recover(); r != nil {
					panicErr = fmt.Errorf("method %s.%s panicked: %v",
						valueType.String(), methodName, r)
				}
			}()

			// Gọi method với context được truyền vào
			method.Call([]reflect.Value{reflect.ValueOf(ctx)})
		}()

		// Fail-fast: dừng ngay khi có panic
		if panicErr != nil {
			return panicErr
		}
	}

	return nil
}

// RegisterRouters là helper để register nhiều controllers cùng lúc
// Useful cho việc register batch controllers từ Fx group
//
// Parameters:
//   - controllers: Danh sách controllers cần register (phải implement ICoreController interface)
//   - ctx: Context được truyền vào mỗi controller. Nếu nil, sẽ dùng context.Background()
//
// Returns:
//   - error: Error ngay khi có controller bị lỗi (fail-fast)
//
// Behavior:
//   - Register controllers theo thứ tự
//   - DỪNG NGAY khi có 1 controller fail (fail-fast pattern)
//   - Return error chứa thông tin controller index bị lỗi
//
// Example:
//
//	func (h *HttpAdapter) OnStart(ctx context.Context) error {
//	    // Register tất cả dynamic controllers với context từ Fx lifecycle
//	    if err := RegisterRouters(h.Config.DynamicControllers, ctx); err != nil {
//	        log.Fatalf("Failed to register controllers: %v", err)
//	        return err
//	    }
//	    return nil
//	}
func RegisterRouters(controllers []ICoreController, ctx context.Context) error {
	// Sử dụng context.Background() nếu ctx nil
	if ctx == nil {
		ctx = context.Background()
	}

	for i, controller := range controllers {
		if err := RegisterRouter(controller, ctx); err != nil {
			// Fail-fast: dừng ngay và return error với controller index
			return fmt.Errorf("controller[%d]: %w", i, err)
		}
	}

	return nil
}

// isValidDynamicMethod kiểm tra method có đúng signature: func(context.Context) không
func isValidDynamicMethod(methodType reflect.Type) bool {
	// Method phải có đúng 1 input parameter
	if methodType.NumIn() != 1 {
		return false
	}

	// Method không có return value
	if methodType.NumOut() != 0 {
		return false
	}

	// Input parameter phải là context.Context
	ctxType := reflect.TypeOf((*context.Context)(nil)).Elem()
	if methodType.In(0) != ctxType {
		return false
	}

	return true
}
