package core

import (
	"fmt"
	"testing"
)

type ValidatedConfig struct {
	Server struct {
		Host string
		Port int
	}
	Database struct {
		Host string
		Port int
	}
}

// MockLoader for ValidatedConfig
type ValidatedMockLoader struct {
	data ValidatedConfig
	err  error
}

func (m *ValidatedMockLoader) Load(dst *ValidatedConfig) error {
	if m.err != nil {
		return m.err
	}
	*dst = m.data
	return nil
}

// Sample validator
type ServerValidator struct{}

func (v *ServerValidator) Validate(cfg *ValidatedConfig) error {
	if cfg.Server.Port < 1024 || cfg.Server.Port > 65535 {
		return fmt.Errorf("server port must be between 1024 and 65535")
	}
	if cfg.Server.Host == "" {
		return fmt.Errorf("server host cannot be empty")
	}
	return nil
}

func TestConfig_WithValidator_Success(t *testing.T) {
	// Mock loader
	loader := &ValidatedMockLoader{
		data: ValidatedConfig{},
	}
	loader.data.Server.Host = "localhost"
	loader.data.Server.Port = 8080

	validator := &ServerValidator{}

	cfg := New[ValidatedConfig](loader).WithValidator(validator)
	if err := cfg.Load(); err != nil {
		t.Fatalf("Load should succeed with valid config: %v", err)
	}

	result := cfg.Get()
	if result.Server.Host != "localhost" {
		t.Errorf("Expected host=localhost, got %s", result.Server.Host)
	}
}

func TestConfig_WithValidator_Fail(t *testing.T) {
	// Mock loader with invalid config
	loader := &ValidatedMockLoader{
		data: ValidatedConfig{},
	}
	loader.data.Server.Host = "localhost"
	loader.data.Server.Port = 80 // Invalid: < 1024

	validator := &ServerValidator{}

	cfg := New[ValidatedConfig](loader).WithValidator(validator)
	if err := cfg.Load(); err == nil {
		t.Error("Load should fail with invalid config")
	}
}

func TestConfig_WithValidator_EmptyHost(t *testing.T) {
	loader := &ValidatedMockLoader{
		data: ValidatedConfig{},
	}
	loader.data.Server.Port = 8080
	// Host is empty

	validator := &ServerValidator{}

	cfg := New[ValidatedConfig](loader).WithValidator(validator)
	if err := cfg.Load(); err == nil {
		t.Error("Load should fail with empty host")
	} else if err.Error() != "config validation failed: server host cannot be empty" {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestValidatorFunc(t *testing.T) {
	loader := &ValidatedMockLoader{
		data: ValidatedConfig{},
	}
	loader.data.Server.Host = "localhost"
	loader.data.Server.Port = 8080

	// Use ValidatorFunc
	validateFunc := ValidatorFunc[ValidatedConfig](func(cfg *ValidatedConfig) error {
		if cfg.Server.Port != 8080 {
			return fmt.Errorf("expected port 8080")
		}
		return nil
	})

	cfg := New[ValidatedConfig](loader).WithValidator(validateFunc)
	if err := cfg.Load(); err != nil {
		t.Fatalf("Load should succeed: %v", err)
	}
}

func TestValidatorFunc_Fail(t *testing.T) {
	loader := &ValidatedMockLoader{
		data: ValidatedConfig{},
	}
	loader.data.Server.Port = 9090

	validateFunc := ValidatorFunc[ValidatedConfig](func(cfg *ValidatedConfig) error {
		if cfg.Server.Port != 8080 {
			return fmt.Errorf("expected port 8080, got %d", cfg.Server.Port)
		}
		return nil
	})

	cfg := New[ValidatedConfig](loader).WithValidator(validateFunc)
	if err := cfg.Load(); err == nil {
		t.Error("Load should fail with wrong port")
	}
}

func TestCompositeValidator_AllPass(t *testing.T) {
	loader := &ValidatedMockLoader{
		data: ValidatedConfig{},
	}
	loader.data.Server.Host = "localhost"
	loader.data.Server.Port = 8080
	loader.data.Database.Host = "dbhost"
	loader.data.Database.Port = 5432

	// Validator 1: check server
	validator1 := ValidatorFunc[ValidatedConfig](func(cfg *ValidatedConfig) error {
		if cfg.Server.Port < 1024 {
			return fmt.Errorf("server port too low")
		}
		return nil
	})

	// Validator 2: check database
	validator2 := ValidatorFunc[ValidatedConfig](func(cfg *ValidatedConfig) error {
		if cfg.Database.Host == "" {
			return fmt.Errorf("database host empty")
		}
		return nil
	})

	composite := NewCompositeValidator(validator1, validator2)

	cfg := New[ValidatedConfig](loader).WithValidator(composite)
	if err := cfg.Load(); err != nil {
		t.Fatalf("Load should succeed when all validators pass: %v", err)
	}
}

func TestCompositeValidator_FirstFail(t *testing.T) {
	loader := &ValidatedMockLoader{
		data: ValidatedConfig{},
	}
	loader.data.Server.Port = 80 // Will fail first validator

	validator1 := ValidatorFunc[ValidatedConfig](func(cfg *ValidatedConfig) error {
		if cfg.Server.Port < 1024 {
			return fmt.Errorf("server port too low")
		}
		return nil
	})

	validator2 := ValidatorFunc[ValidatedConfig](func(cfg *ValidatedConfig) error {
		return nil // Would pass
	})

	composite := NewCompositeValidator(validator1, validator2)

	cfg := New[ValidatedConfig](loader).WithValidator(composite)
	if err := cfg.Load(); err == nil {
		t.Error("Load should fail when first validator fails")
	}
}

func TestCompositeValidator_SecondFail(t *testing.T) {
	loader := &ValidatedMockLoader{
		data: ValidatedConfig{},
	}
	loader.data.Server.Port = 8080 // Will pass first validator
	// Database.Host is empty - will fail second validator

	validator1 := ValidatorFunc[ValidatedConfig](func(cfg *ValidatedConfig) error {
		if cfg.Server.Port < 1024 {
			return fmt.Errorf("server port too low")
		}
		return nil
	})

	validator2 := ValidatorFunc[ValidatedConfig](func(cfg *ValidatedConfig) error {
		if cfg.Database.Host == "" {
			return fmt.Errorf("database host empty")
		}
		return nil
	})

	composite := NewCompositeValidator(validator1, validator2)

	cfg := New[ValidatedConfig](loader).WithValidator(composite)
	if err := cfg.Load(); err == nil {
		t.Error("Load should fail when second validator fails")
	}
}

func TestConfig_NoValidator(t *testing.T) {
	loader := &ValidatedMockLoader{
		data: ValidatedConfig{},
	}
	loader.data.Server.Port = 80 // Would be invalid if validator was set

	cfg := New[ValidatedConfig](loader)
	// No validator set
	if err := cfg.Load(); err != nil {
		t.Fatalf("Load should succeed without validator: %v", err)
	}
}

func TestConfig_WithValidator_MethodChaining(t *testing.T) {
	loader := &ValidatedMockLoader{
		data: ValidatedConfig{},
	}
	loader.data.Server.Host = "localhost"
	loader.data.Server.Port = 8080

	// Test method chaining with both WithMerge and WithValidator
	cfg := New[ValidatedConfig](loader).
		WithMerge(DefaultMerge[ValidatedConfig]).
		WithValidator(&ServerValidator{})

	if err := cfg.Load(); err != nil {
		t.Fatalf("Load should succeed: %v", err)
	}
}
