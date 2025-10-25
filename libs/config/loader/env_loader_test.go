package loader

import (
	"os"
	"testing"
)

func TestEnvLoader_Load(t *testing.T) {
	// Set environment variables
	os.Setenv("APP_SERVER_HOST", "localhost")
	os.Setenv("APP_SERVER_PORT", "9090")
	os.Setenv("APP_DATABASE_HOST", "dbhost")
	os.Setenv("APP_DATABASE_PORT", "5432")
	defer func() {
		os.Unsetenv("APP_SERVER_HOST")
		os.Unsetenv("APP_SERVER_PORT")
		os.Unsetenv("APP_DATABASE_HOST")
		os.Unsetenv("APP_DATABASE_PORT")
	}()

	// Need to specify keys for viper to bind
	loader := NewEnvLoader("APP").WithKeys(
		"server.host", "server.port",
		"database.host", "database.port",
	)
	cfg := &TestConfig{}

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Server.Host != "localhost" {
		t.Errorf("Expected server.host=localhost, got %s", cfg.Server.Host)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("Expected server.port=9090, got %d", cfg.Server.Port)
	}

	if cfg.Database.Host != "dbhost" {
		t.Errorf("Expected database.host=dbhost, got %s", cfg.Database.Host)
	}

	if cfg.Database.Port != 5432 {
		t.Errorf("Expected database.port=5432, got %d", cfg.Database.Port)
	}
}

func TestEnvLoader_WithoutPrefix(t *testing.T) {
	// Set environment variables without prefix
	os.Setenv("SERVER_HOST", "localhost")
	os.Setenv("SERVER_PORT", "8080")
	defer func() {
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")
	}()

	loader := NewEnvLoader("").WithKeys("server.host", "server.port")
	cfg := &TestConfig{}

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Server.Host != "localhost" {
		t.Errorf("Expected server.host=localhost, got %s", cfg.Server.Host)
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("Expected server.port=8080, got %d", cfg.Server.Port)
	}
}

func TestEnvLoader_WithKeys(t *testing.T) {
	os.Setenv("APP_SERVER_HOST", "localhost")
	os.Setenv("APP_SERVER_PORT", "9090")
	os.Setenv("APP_DATABASE_HOST", "dbhost")
	defer func() {
		os.Unsetenv("APP_SERVER_HOST")
		os.Unsetenv("APP_SERVER_PORT")
		os.Unsetenv("APP_DATABASE_HOST")
	}()

	// Only bind specific keys
	loader := NewEnvLoader("APP").WithKeys("server.host", "server.port")
	cfg := &TestConfig{}

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Specified keys should be loaded
	if cfg.Server.Host != "localhost" {
		t.Errorf("Expected server.host=localhost, got %s", cfg.Server.Host)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("Expected server.port=9090, got %d", cfg.Server.Port)
	}

	// Note: With AutomaticEnv, all env vars are still accessible
	// WithKeys is more for explicit binding
}

func TestEnvLoader_PartialEnv(t *testing.T) {
	// Only set some env vars
	os.Setenv("APP_SERVER_HOST", "localhost")
	defer os.Unsetenv("APP_SERVER_HOST")

	loader := NewEnvLoader("APP").WithKeys("server.host", "server.port")
	cfg := &TestConfig{}

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Set env var should be loaded
	if cfg.Server.Host != "localhost" {
		t.Errorf("Expected server.host=localhost, got %s", cfg.Server.Host)
	}

	// Unset env vars should be zero value
	if cfg.Server.Port != 0 {
		t.Errorf("Expected server.port=0, got %d", cfg.Server.Port)
	}

	if cfg.Database.Host != "" {
		t.Errorf("Expected database.host to be empty, got %s", cfg.Database.Host)
	}
}

func TestEnvLoader_NoEnvVars(t *testing.T) {
	loader := NewEnvLoader("NONEXISTENT")
	cfg := &TestConfig{}

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// All fields should be zero values
	if cfg.Server.Host != "" || cfg.Server.Port != 0 {
		t.Error("Expected zero values when no env vars set")
	}
}

func TestEnvLoader_TypeConversion(t *testing.T) {
	type TypeTestConfig struct {
		StringVal string `mapstructure:"string_val"`
		IntVal    int    `mapstructure:"int_val"`
		BoolVal   bool   `mapstructure:"bool_val"`
	}

	os.Setenv("APP_STRING_VAL", "hello")
	os.Setenv("APP_INT_VAL", "42")
	os.Setenv("APP_BOOL_VAL", "true")
	defer func() {
		os.Unsetenv("APP_STRING_VAL")
		os.Unsetenv("APP_INT_VAL")
		os.Unsetenv("APP_BOOL_VAL")
	}()

	loader := NewEnvLoader("APP").WithKeys("string_val", "int_val", "bool_val")
	cfg := &TypeTestConfig{}

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.StringVal != "hello" {
		t.Errorf("Expected string_val=hello, got %s", cfg.StringVal)
	}

	if cfg.IntVal != 42 {
		t.Errorf("Expected int_val=42, got %d", cfg.IntVal)
	}

	if cfg.BoolVal != true {
		t.Errorf("Expected bool_val=true, got %v", cfg.BoolVal)
	}
}

func TestEnvLoader_UnderscoreConversion(t *testing.T) {
	// Test that underscore is converted to dot
	os.Setenv("APP_SERVER_HOST", "localhost")
	defer os.Unsetenv("APP_SERVER_HOST")

	loader := NewEnvLoader("APP").WithKeys("server.host")
	cfg := &TestConfig{}

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// APP_SERVER_HOST should map to server.host
	if cfg.Server.Host != "localhost" {
		t.Errorf("Expected server.host=localhost, got %s", cfg.Server.Host)
	}
}
