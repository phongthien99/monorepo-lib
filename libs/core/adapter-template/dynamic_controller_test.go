package adaptertemplate

import (
	"context"
	"reflect"
	"testing"
)

// Mock controller cho testing - implements ICoreController interface
type testController struct {
	getMethodCalled     bool
	postMethodCalled    bool
	invalidMethodCalled bool
	shouldPanic         bool
}

// Ensure testController implements ICoreController interface
var _ ICoreController = (*testController)(nil)

// Valid method: func(context.Context)
func (t *testController) GetUsers(ctx context.Context) {
	t.getMethodCalled = true
}

// Valid method: func(context.Context)
func (t *testController) CreateUser(ctx context.Context) {
	t.postMethodCalled = true
	if t.shouldPanic {
		panic("intentional panic for testing")
	}
}

// Invalid method: không có parameter context.Context
func (t *testController) InvalidMethod() {
	t.invalidMethodCalled = true
}

// Invalid method: có return value
func (t *testController) InvalidMethodWithReturn(ctx context.Context) string {
	return "should not be called"
}

// Invalid method: có multiple parameters
func (t *testController) InvalidMethodMultiParams(ctx context.Context, s string) {
}

// Private method: không được export
func (t *testController) privateMethod(ctx context.Context) {
}

func TestRegisterRouter(t *testing.T) {
	controller := &testController{}

	// Execute
	err := RegisterRouter(controller, nil)

	// Verify: No error
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify: Chỉ valid methods được gọi
	if !controller.getMethodCalled {
		t.Error("GetUsers should be called")
	}

	if !controller.postMethodCalled {
		t.Error("CreateUser should be called")
	}

	if controller.invalidMethodCalled {
		t.Error("InvalidMethod should NOT be called")
	}
}

func TestRegisterRouter_WithContext(t *testing.T) {
	controller := &testController{}
	ctx := context.Background()

	// Execute với context
	err := RegisterRouter(controller, ctx)

	// Verify: No error
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify
	if !controller.getMethodCalled || !controller.postMethodCalled {
		t.Error("Methods should be called with provided context")
	}
}

func TestRegisterRouter_Panic(t *testing.T) {
	controller := &testController{shouldPanic: true}

	// Execute - should recover from panic and return error immediately (fail-fast)
	err := RegisterRouter(controller, nil)

	// Verify: Should have error
	if err == nil {
		t.Error("Expected error from panic, got nil")
	}

	// Note: Reflection methods iterate in ALPHABETICAL order
	// CreateUser (C) comes before GetUsers (G)
	// So CreateUser panics first, and GetUsers is never called (fail-fast)

	// Verify: CreateUser was attempted (panic occurred)
	if !controller.postMethodCalled {
		t.Error("CreateUser should be attempted (and panic)")
	}

	// Verify: GetUsers should NOT be called because CreateUser panicked first (fail-fast)
	if controller.getMethodCalled {
		t.Error("GetUsers should NOT be called due to fail-fast after CreateUser panic")
	}
}

func TestRegisterRouter_NilController(t *testing.T) {
	// Should not panic and return nil error
	err := RegisterRouter(nil, nil)
	if err != nil {
		t.Errorf("Expected nil error for nil controller, got: %v", err)
	}

	err = RegisterRouter(nil, context.Background())
	if err != nil {
		t.Errorf("Expected nil error for nil controller, got: %v", err)
	}
}

func TestRegisterRouters(t *testing.T) {
	controller1 := &testController{}
	controller2 := &testController{}

	controllers := []ICoreController{controller1, controller2}

	// Execute
	err := RegisterRouters(controllers, nil)

	// Verify: No error
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify: Cả 2 controllers đều được register
	if !controller1.getMethodCalled || !controller1.postMethodCalled {
		t.Error("Controller1 methods should be called")
	}

	if !controller2.getMethodCalled || !controller2.postMethodCalled {
		t.Error("Controller2 methods should be called")
	}
}

func TestRegisterRouters_WithContext(t *testing.T) {
	controller1 := &testController{}
	controller2 := &testController{}

	controllers := []ICoreController{controller1, controller2}
	ctx := context.Background()

	// Execute với context
	err := RegisterRouters(controllers, ctx)

	// Verify: No error
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify: Cả 2 controllers đều được register
	if !controller1.getMethodCalled || !controller1.postMethodCalled {
		t.Error("Controller1 methods should be called")
	}

	if !controller2.getMethodCalled || !controller2.postMethodCalled {
		t.Error("Controller2 methods should be called")
	}
}

func TestRegisterRouters_WithErrors(t *testing.T) {
	controller1 := &testController{}
	controller2 := &testController{shouldPanic: true} // This will cause error
	controller3 := &testController{}

	controllers := []ICoreController{controller1, controller2, controller3}

	// Execute - fail-fast behavior
	err := RegisterRouters(controllers, nil)

	// Verify: Should have error from controller2
	if err == nil {
		t.Error("Expected error from controller2 panic, got nil")
	}

	// Verify: Controller1 should be registered successfully (processed before controller2)
	if !controller1.getMethodCalled || !controller1.postMethodCalled {
		t.Error("Controller1 should be registered successfully")
	}

	// Verify: Controller2's CreateUser was attempted (and panicked)
	// Note: Methods iterate alphabetically, so CreateUser (C) is called before GetUsers (G)
	if !controller2.postMethodCalled {
		t.Error("Controller2.CreateUser should be attempted (and panic)")
	}

	// Verify: Controller2's GetUsers should NOT be called (fail-fast after CreateUser panic)
	if controller2.getMethodCalled {
		t.Error("Controller2.GetUsers should NOT be called due to fail-fast")
	}

	// Verify: Controller3 should NOT be registered (fail-fast stops at controller2 error)
	if controller3.getMethodCalled || controller3.postMethodCalled {
		t.Error("Controller3 should NOT be registered due to fail-fast behavior")
	}
}

func TestIsValidDynamicMethod(t *testing.T) {
	// Create test controller
	controller := &testController{}
	value := reflect.ValueOf(controller)

	tests := []struct {
		methodName string
		want       bool
	}{
		{"GetUsers", true},                  // Valid: func(context.Context)
		{"CreateUser", true},                // Valid: func(context.Context)
		{"InvalidMethod", false},            // Invalid: no parameters
		{"InvalidMethodWithReturn", false},  // Invalid: has return value
		{"InvalidMethodMultiParams", false}, // Invalid: multiple params
	}

	for _, tt := range tests {
		t.Run(tt.methodName, func(t *testing.T) {
			// Find method by name using reflect.Value.MethodByName
			// This gives us the method WITHOUT receiver in the type signature
			method := value.MethodByName(tt.methodName)
			if !method.IsValid() {
				t.Fatalf("Method %s not found", tt.methodName)
			}

			// Get method type (without receiver)
			methodType := method.Type()

			// Validate
			got := isValidDynamicMethod(methodType)
			if got != tt.want {
				t.Errorf("isValidDynamicMethod(%s) = %v, want %v", tt.methodName, got, tt.want)
			}
		})
	}
}
