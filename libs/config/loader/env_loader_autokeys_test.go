package loader

import (
	"os"
	"testing"
)

func TestEnvLoader_WithAutoKeys(t *testing.T) {
	type AutoKeysConfig struct {
		Server struct {
			Host string `mapstructure:"host"`
			Port int    `mapstructure:"port"`
		} `mapstructure:"server"`
		Database struct {
			Host string `mapstructure:"host"`
			Port int    `mapstructure:"port"`
		} `mapstructure:"database"`
	}

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

	// Use WithAutoKeys instead of manually listing keys
	loader := NewEnvLoader("APP").WithAutoKeys(AutoKeysConfig{})
	cfg := &AutoKeysConfig{}

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify all fields loaded
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

func TestEnvLoader_WithAutoKeys_Pointer(t *testing.T) {
	type PointerConfig struct {
		Name string `mapstructure:"name"`
		Port int    `mapstructure:"port"`
	}

	os.Setenv("APP_NAME", "test-app")
	os.Setenv("APP_PORT", "8080")
	defer func() {
		os.Unsetenv("APP_NAME")
		os.Unsetenv("APP_PORT")
	}()

	// Test with pointer type
	loader := NewEnvLoader("APP").WithAutoKeys(&PointerConfig{})
	cfg := &PointerConfig{}

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Name != "test-app" {
		t.Errorf("Expected name=test-app, got %s", cfg.Name)
	}

	if cfg.Port != 8080 {
		t.Errorf("Expected port=8080, got %d", cfg.Port)
	}
}

func TestEnvLoader_WithAutoKeys_DeepNested(t *testing.T) {
	type DeepNestedConfig struct {
		App struct {
			Server struct {
				HTTP struct {
					Host string `mapstructure:"host"`
					Port int    `mapstructure:"port"`
				} `mapstructure:"http"`
			} `mapstructure:"server"`
		} `mapstructure:"app"`
	}

	os.Setenv("APP_APP_SERVER_HTTP_HOST", "api.example.com")
	os.Setenv("APP_APP_SERVER_HTTP_PORT", "443")
	defer func() {
		os.Unsetenv("APP_APP_SERVER_HTTP_HOST")
		os.Unsetenv("APP_APP_SERVER_HTTP_PORT")
	}()

	loader := NewEnvLoader("APP").WithAutoKeys(DeepNestedConfig{})
	cfg := &DeepNestedConfig{}

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.App.Server.HTTP.Host != "api.example.com" {
		t.Errorf("Expected host=api.example.com, got %s", cfg.App.Server.HTTP.Host)
	}

	if cfg.App.Server.HTTP.Port != 443 {
		t.Errorf("Expected port=443, got %d", cfg.App.Server.HTTP.Port)
	}
}

func TestEnvLoader_WithAutoKeys_Partial(t *testing.T) {
	type PartialConfig struct {
		Server struct {
			Host string `mapstructure:"host"`
			Port int    `mapstructure:"port"`
		} `mapstructure:"server"`
	}

	// Only set one env var
	os.Setenv("APP_SERVER_HOST", "partial-host")
	defer os.Unsetenv("APP_SERVER_HOST")

	loader := NewEnvLoader("APP").WithAutoKeys(PartialConfig{})
	cfg := &PartialConfig{}

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Set env var should be loaded
	if cfg.Server.Host != "partial-host" {
		t.Errorf("Expected server.host=partial-host, got %s", cfg.Server.Host)
	}

	// Unset env var should be zero value
	if cfg.Server.Port != 0 {
		t.Errorf("Expected server.port=0, got %d", cfg.Server.Port)
	}
}
