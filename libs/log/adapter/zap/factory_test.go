package zap

import (
	"testing"

	"github.com/phongthien99/monorepo-lib/libs/log/core"
)

func TestNewDevelopment(t *testing.T) {
	logger, err := NewDevelopment()
	if err != nil {
		t.Fatalf("NewDevelopment() error = %v", err)
	}
	if logger == nil {
		t.Error("NewDevelopment() returned nil logger")
	}
	if logger.Level() != core.DebugLevel {
		t.Errorf("NewDevelopment() level = %v, want %v", logger.Level(), core.DebugLevel)
	}

	// Test logging
	logger.Info("test development logger")
}

func TestNewProduction(t *testing.T) {
	logger, err := NewProduction()
	if err != nil {
		t.Fatalf("NewProduction() error = %v", err)
	}
	if logger == nil {
		t.Error("NewProduction() returned nil logger")
	}
	if logger.Level() != core.InfoLevel {
		t.Errorf("NewProduction() level = %v, want %v", logger.Level(), core.InfoLevel)
	}

	// Test logging
	logger.Info("test production logger")
}

func TestNewExample(t *testing.T) {
	logger := NewExample()
	if logger == nil {
		t.Error("NewExample() returned nil logger")
	}
	if logger.Level() != core.DebugLevel {
		t.Errorf("NewExample() level = %v, want %v", logger.Level(), core.DebugLevel)
	}

	// Test logging
	logger.Info("test example logger")
}

func TestNewNop(t *testing.T) {
	logger := NewNop()
	if logger == nil {
		t.Error("NewNop() returned nil logger")
	}

	// Nop logger should not output anything
	logger.Info("this should not be output")
}

func TestNewWithConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name:    "Default config",
			config:  DefaultConfig(),
			wantErr: false,
		},
		{
			name:    "Development config",
			config:  DevelopmentConfig(),
			wantErr: false,
		},
		{
			name: "Custom config with JSON encoding",
			config: Config{
				Level:            core.WarnLevel,
				Development:      false,
				Encoding:         "json",
				OutputPaths:      []string{"stdout"},
				ErrorOutputPaths: []string{"stderr"},
			},
			wantErr: false,
		},
		{
			name: "Custom config with console encoding",
			config: Config{
				Level:            core.ErrorLevel,
				Development:      true,
				Encoding:         "console",
				OutputPaths:      []string{"stdout"},
				ErrorOutputPaths: []string{"stderr"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewWithConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWithConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && logger == nil {
				t.Error("NewWithConfig() returned nil logger")
			}
			if !tt.wantErr && logger.Level() != tt.config.Level {
				t.Errorf("NewWithConfig() level = %v, want %v", logger.Level(), tt.config.Level)
			}

			// Test logging
			if logger != nil {
				logger.Info("test custom config logger")
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Level != core.InfoLevel {
		t.Errorf("DefaultConfig().Level = %v, want %v", config.Level, core.InfoLevel)
	}
	if config.Development != false {
		t.Error("DefaultConfig().Development should be false")
	}
	if config.Encoding != "json" {
		t.Errorf("DefaultConfig().Encoding = %v, want json", config.Encoding)
	}
	if len(config.OutputPaths) == 0 {
		t.Error("DefaultConfig().OutputPaths should not be empty")
	}
	if len(config.ErrorOutputPaths) == 0 {
		t.Error("DefaultConfig().ErrorOutputPaths should not be empty")
	}
}

func TestDevelopmentConfig(t *testing.T) {
	config := DevelopmentConfig()

	if config.Level != core.DebugLevel {
		t.Errorf("DevelopmentConfig().Level = %v, want %v", config.Level, core.DebugLevel)
	}
	if config.Development != true {
		t.Error("DevelopmentConfig().Development should be true")
	}
	if config.Encoding != "console" {
		t.Errorf("DevelopmentConfig().Encoding = %v, want console", config.Encoding)
	}
	if len(config.OutputPaths) == 0 {
		t.Error("DevelopmentConfig().OutputPaths should not be empty")
	}
	if len(config.ErrorOutputPaths) == 0 {
		t.Error("DevelopmentConfig().ErrorOutputPaths should not be empty")
	}
}

func TestNewWithConfig_AllLevels(t *testing.T) {
	levels := []core.Level{
		core.DebugLevel,
		core.InfoLevel,
		core.WarnLevel,
		core.ErrorLevel,
		core.DPanicLevel,
		core.PanicLevel,
		core.FatalLevel,
	}

	for _, level := range levels {
		t.Run(level.String(), func(t *testing.T) {
			config := Config{
				Level:            level,
				Development:      false,
				Encoding:         "json",
				OutputPaths:      []string{"stdout"},
				ErrorOutputPaths: []string{"stderr"},
			}

			logger, err := NewWithConfig(config)
			if err != nil {
				t.Fatalf("NewWithConfig() error = %v", err)
			}
			if logger == nil {
				t.Error("NewWithConfig() returned nil logger")
			}
			if logger.Level() != level {
				t.Errorf("Logger level = %v, want %v", logger.Level(), level)
			}
		})
	}
}

func BenchmarkNewDevelopment(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger, err := NewDevelopment()
		if err != nil {
			b.Fatal(err)
		}
		_ = logger.Sync()
	}
}

func BenchmarkNewProduction(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger, err := NewProduction()
		if err != nil {
			b.Fatal(err)
		}
		_ = logger.Sync()
	}
}

func BenchmarkNewExample(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger := NewExample()
		_ = logger.Sync()
	}
}
