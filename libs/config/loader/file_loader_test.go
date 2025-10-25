package loader

import (
	"os"
	"path/filepath"
	"testing"
)

type TestConfig struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`
	Database struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"database"`
}

func TestFileLoader_LoadJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.json")

	jsonContent := `{
		"server": {
			"host": "localhost",
			"port": 8080
		},
		"database": {
			"host": "dbhost",
			"port": 5432
		}
	}`

	if err := os.WriteFile(configPath, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	loader := NewFileLoader(configPath, "json")
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

	if cfg.Database.Host != "dbhost" {
		t.Errorf("Expected database.host=dbhost, got %s", cfg.Database.Host)
	}

	if cfg.Database.Port != 5432 {
		t.Errorf("Expected database.port=5432, got %d", cfg.Database.Port)
	}
}

func TestFileLoader_LoadYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.yaml")

	yamlContent := `
server:
  host: localhost
  port: 8080
database:
  host: dbhost
  port: 5432
`

	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	loader := NewFileLoader(configPath, "yaml")
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

func TestFileLoader_FileNotFound(t *testing.T) {
	loader := NewFileLoader("/nonexistent/config.json", "json")
	cfg := &TestConfig{}

	if err := loader.Load(cfg); err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestFileLoader_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.json")

	invalidJSON := `{invalid json}`

	if err := os.WriteFile(configPath, []byte(invalidJSON), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	loader := NewFileLoader(configPath, "json")
	cfg := &TestConfig{}

	if err := loader.Load(cfg); err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestFileLoader_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.yaml")

	invalidYAML := `
invalid yaml
  - this is: bad
    - format
`

	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	loader := NewFileLoader(configPath, "yaml")
	cfg := &TestConfig{}

	if err := loader.Load(cfg); err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestFileLoader_PartialConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.json")

	// Only server config, no database
	jsonContent := `{
		"server": {
			"host": "localhost",
			"port": 8080
		}
	}`

	if err := os.WriteFile(configPath, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	loader := NewFileLoader(configPath, "json")
	cfg := &TestConfig{}

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Server should be loaded
	if cfg.Server.Host != "localhost" {
		t.Errorf("Expected server.host=localhost, got %s", cfg.Server.Host)
	}

	// Database should be zero value
	if cfg.Database.Host != "" {
		t.Errorf("Expected database.host to be empty, got %s", cfg.Database.Host)
	}

	if cfg.Database.Port != 0 {
		t.Errorf("Expected database.port to be 0, got %d", cfg.Database.Port)
	}
}

func TestFileLoader_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.json")

	emptyJSON := `{}`

	if err := os.WriteFile(configPath, []byte(emptyJSON), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	loader := NewFileLoader(configPath, "json")
	cfg := &TestConfig{}

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// All fields should be zero values
	if cfg.Server.Host != "" || cfg.Server.Port != 0 {
		t.Error("Expected zero values for empty config file")
	}
}

func TestFileLoader_TOML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.toml")

	tomlContent := `
[server]
host = "localhost"
port = 8080

[database]
host = "dbhost"
port = 5432
`

	if err := os.WriteFile(configPath, []byte(tomlContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	loader := NewFileLoader(configPath, "toml")
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
