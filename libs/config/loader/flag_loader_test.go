package loader

import (
	"testing"

	"github.com/spf13/pflag"
)

func TestFlagLoader_Load(t *testing.T) {
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("server.host", "localhost", "Server host")
	flags.Int("server.port", 8080, "Server port")
	flags.String("database.host", "dbhost", "Database host")
	flags.Int("database.port", 5432, "Database port")

	// Parse flags with values
	flags.Parse([]string{
		"--server.host=api.example.com",
		"--server.port=9090",
	})

	loader := NewFlagLoader(flags)
	cfg := &TestConfig{}

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Parsed flags should be loaded
	if cfg.Server.Host != "api.example.com" {
		t.Errorf("Expected server.host=api.example.com, got %s", cfg.Server.Host)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("Expected server.port=9090, got %d", cfg.Server.Port)
	}

	// Unparsed flags should have default values
	if cfg.Database.Host != "dbhost" {
		t.Errorf("Expected database.host=dbhost (default), got %s", cfg.Database.Host)
	}

	if cfg.Database.Port != 5432 {
		t.Errorf("Expected database.port=5432 (default), got %d", cfg.Database.Port)
	}
}

func TestFlagLoader_DefaultValues(t *testing.T) {
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("server.host", "default-host", "Server host")
	flags.Int("server.port", 3000, "Server port")

	// Parse without arguments - should use defaults
	flags.Parse([]string{})

	loader := NewFlagLoader(flags)
	cfg := &TestConfig{}

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Server.Host != "default-host" {
		t.Errorf("Expected server.host=default-host, got %s", cfg.Server.Host)
	}

	if cfg.Server.Port != 3000 {
		t.Errorf("Expected server.port=3000, got %d", cfg.Server.Port)
	}
}

func TestFlagLoader_NilFlagSet(t *testing.T) {
	// Should use pflag.CommandLine by default
	loader := NewFlagLoader(nil)

	if loader.flagSet == nil {
		t.Error("Expected flagSet to be set to CommandLine")
	}

	if loader.flagSet != pflag.CommandLine {
		t.Error("Expected flagSet to be pflag.CommandLine")
	}
}

func TestFlagLoader_NestedStructure(t *testing.T) {
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)

	// Using dot notation for nested structure
	flags.String("server.host", "", "Server host")
	flags.Int("server.port", 0, "Server port")

	flags.Parse([]string{
		"--server.host=nested-host",
		"--server.port=7777",
	})

	loader := NewFlagLoader(flags)
	cfg := &TestConfig{}

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Server.Host != "nested-host" {
		t.Errorf("Expected server.host=nested-host, got %s", cfg.Server.Host)
	}

	if cfg.Server.Port != 7777 {
		t.Errorf("Expected server.port=7777, got %d", cfg.Server.Port)
	}
}

func TestFlagLoader_EmptyFlags(t *testing.T) {
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	// No flags defined
	flags.Parse([]string{})

	loader := NewFlagLoader(flags)
	cfg := &TestConfig{}

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should be zero values
	if cfg.Server.Host != "" || cfg.Server.Port != 0 {
		t.Error("Expected zero values for empty flags")
	}
}

func TestFlagLoader_PartialFlags(t *testing.T) {
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("server.host", "", "Server host")
	// Only define server.host, not server.port

	flags.Parse([]string{"--server.host=partial-host"})

	loader := NewFlagLoader(flags)
	cfg := &TestConfig{}

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Defined flag should be loaded
	if cfg.Server.Host != "partial-host" {
		t.Errorf("Expected server.host=partial-host, got %s", cfg.Server.Host)
	}

	// Undefined flag should be zero value
	if cfg.Server.Port != 0 {
		t.Errorf("Expected server.port=0, got %d", cfg.Server.Port)
	}
}

func TestFlagLoader_TypeConversion(t *testing.T) {
	type TypeTestConfig struct {
		StringVal string `mapstructure:"string_val"`
		IntVal    int    `mapstructure:"int_val"`
		BoolVal   bool   `mapstructure:"bool_val"`
	}

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("string_val", "", "String value")
	flags.Int("int_val", 0, "Int value")
	flags.Bool("bool_val", false, "Bool value")

	flags.Parse([]string{
		"--string_val=hello",
		"--int_val=99",
		"--bool_val=true",
	})

	loader := NewFlagLoader(flags)
	cfg := &TypeTestConfig{}

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.StringVal != "hello" {
		t.Errorf("Expected string_val=hello, got %s", cfg.StringVal)
	}

	if cfg.IntVal != 99 {
		t.Errorf("Expected int_val=99, got %d", cfg.IntVal)
	}

	if cfg.BoolVal != true {
		t.Errorf("Expected bool_val=true, got %v", cfg.BoolVal)
	}
}

func TestFlagLoader_Shorthand(t *testing.T) {
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.StringP("server.host", "h", "localhost", "Server host")
	flags.IntP("server.port", "p", 8080, "Server port")

	// Use shorthand flags
	flags.Parse([]string{"-h", "short-host", "-p", "4000"})

	loader := NewFlagLoader(flags)
	cfg := &TestConfig{}

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Server.Host != "short-host" {
		t.Errorf("Expected server.host=short-host, got %s", cfg.Server.Host)
	}

	if cfg.Server.Port != 4000 {
		t.Errorf("Expected server.port=4000, got %d", cfg.Server.Port)
	}
}
