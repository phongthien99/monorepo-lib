package core

import (
	"fmt"
	"testing"
)

type AppConfig struct {
	Server struct {
		Host string
		Port int
	}
	Database struct {
		Host string
		Port int
	}
}

// MockLoader for testing
type MockLoader struct {
	data AppConfig
	err  error
}

func (m *MockLoader) Load(dst *AppConfig) error {
	if m.err != nil {
		return m.err
	}
	*dst = m.data
	return nil
}

func TestConfig_New(t *testing.T) {
	cfg := New[AppConfig]()

	if cfg == nil {
		t.Fatal("Expected non-nil config")
	}

	if cfg.mergeFunc == nil {
		t.Error("Expected default merge function to be set")
	}
}

func TestConfig_WithMerge(t *testing.T) {
	customMergeCalled := false
	customMerge := func(dst, src *AppConfig) error {
		customMergeCalled = true
		return nil
	}

	loader := &MockLoader{
		data: AppConfig{},
	}

	cfg := New[AppConfig](loader).WithMerge(customMerge)
	cfg.Load()

	if !customMergeCalled {
		t.Error("Expected custom merge to be called")
	}
}

func TestConfig_LoadSingleLoader(t *testing.T) {
	loader := &MockLoader{
		data: AppConfig{},
	}
	loader.data.Server.Host = "localhost"
	loader.data.Server.Port = 8080

	cfg := New[AppConfig](loader)
	if err := cfg.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	result := cfg.Get()

	if result.Server.Host != "localhost" {
		t.Errorf("Expected host=localhost, got %s", result.Server.Host)
	}

	if result.Server.Port != 8080 {
		t.Errorf("Expected port=8080, got %d", result.Server.Port)
	}
}

func TestConfig_LoadMultipleLoaders(t *testing.T) {
	// Loader 1: base config
	loader1 := &MockLoader{
		data: AppConfig{},
	}
	loader1.data.Server.Host = "localhost"
	loader1.data.Server.Port = 8080
	loader1.data.Database.Host = "dbhost"
	loader1.data.Database.Port = 5432

	// Loader 2: override server.port
	loader2 := &MockLoader{
		data: AppConfig{},
	}
	loader2.data.Server.Port = 9090

	// Loader 3: override database.host
	loader3 := &MockLoader{
		data: AppConfig{},
	}
	loader3.data.Database.Host = "prod-db"

	cfg := New[AppConfig](loader1, loader2, loader3)
	if err := cfg.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	result := cfg.Get()

	// server.host from loader1 (not overridden)
	if result.Server.Host != "localhost" {
		t.Errorf("Expected server.host=localhost, got %s", result.Server.Host)
	}

	// server.port from loader2 (overridden)
	if result.Server.Port != 9090 {
		t.Errorf("Expected server.port=9090, got %d", result.Server.Port)
	}

	// database.host from loader3 (overridden)
	if result.Database.Host != "prod-db" {
		t.Errorf("Expected database.host=prod-db, got %s", result.Database.Host)
	}

	// database.port from loader1 (not overridden)
	if result.Database.Port != 5432 {
		t.Errorf("Expected database.port=5432, got %d", result.Database.Port)
	}
}

func TestConfig_LoadError(t *testing.T) {
	loader := &MockLoader{
		err: fmt.Errorf("load error"),
	}

	cfg := New[AppConfig](loader)
	if err := cfg.Load(); err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestConfig_MergeError(t *testing.T) {
	loader := &MockLoader{
		data: AppConfig{},
	}

	customMerge := func(dst, src *AppConfig) error {
		return fmt.Errorf("merge error")
	}

	cfg := New[AppConfig](loader).WithMerge(customMerge)
	if err := cfg.Load(); err == nil {
		t.Error("Expected merge error, got nil")
	}
}

func TestConfig_GetBeforeLoad(t *testing.T) {
	cfg := New[AppConfig]()

	// Get before Load should return zero value
	result := cfg.Get()

	if result.Server.Host != "" || result.Server.Port != 0 {
		t.Error("Expected zero value before Load")
	}
}

func TestConfig_GetPtr(t *testing.T) {
	loader := &MockLoader{
		data: AppConfig{},
	}
	loader.data.Server.Host = "localhost"

	cfg := New[AppConfig](loader)
	cfg.Load()

	ptr := cfg.GetPtr()

	if ptr == nil {
		t.Fatal("Expected non-nil pointer")
	}

	if ptr.Server.Host != "localhost" {
		t.Errorf("Expected host=localhost, got %s", ptr.Server.Host)
	}

	// Modify through pointer
	ptr.Server.Port = 9999

	// Should affect internal data
	if cfg.Get().Server.Port != 9999 {
		t.Error("Expected pointer modification to affect internal data")
	}
}

func TestConfig_EmptyLoaders(t *testing.T) {
	cfg := New[AppConfig]()
	if err := cfg.Load(); err != nil {
		t.Fatalf("Load with empty loaders should not fail: %v", err)
	}

	result := cfg.Get()

	// Should be zero value
	if result.Server.Host != "" || result.Server.Port != 0 {
		t.Error("Expected zero value for empty loaders")
	}
}

func TestConfig_MethodChaining(t *testing.T) {
	loader := &MockLoader{
		data: AppConfig{},
	}

	customMerge := func(dst, src *AppConfig) error {
		return DefaultMerge(dst, src)
	}

	// Test method chaining
	cfg := New[AppConfig](loader).
		WithMerge(customMerge)

	if err := cfg.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should work
	_ = cfg.Get()
}

func TestConfig_LoadOrder(t *testing.T) {
	loadOrder := []int{}

	// Custom merge to track load order
	customMerge := func(dst, src *AppConfig) error {
		if src.Server.Port != 0 {
			loadOrder = append(loadOrder, src.Server.Port)
		}
		return DefaultMerge(dst, src)
	}

	loader1 := &MockLoader{data: AppConfig{}}
	loader1.data.Server.Port = 1

	loader2 := &MockLoader{data: AppConfig{}}
	loader2.data.Server.Port = 2

	loader3 := &MockLoader{data: AppConfig{}}
	loader3.data.Server.Port = 3

	cfg := New[AppConfig](loader1, loader2, loader3).WithMerge(customMerge)
	cfg.Load()

	// Should be loaded in order
	if len(loadOrder) != 3 {
		t.Fatalf("Expected 3 loads, got %d", len(loadOrder))
	}

	if loadOrder[0] != 1 || loadOrder[1] != 2 || loadOrder[2] != 3 {
		t.Errorf("Expected load order [1,2,3], got %v", loadOrder)
	}
}
