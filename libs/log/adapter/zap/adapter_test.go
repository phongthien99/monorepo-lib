package zap

import (
	"testing"

	"github.com/phongthien99/monorepo-lib/libs/log/core"
)

func TestZapAdapter_BasicLogger(t *testing.T) {
	logger := NewExample()

	// Test basic logging methods (just make sure they don't panic)
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")
}

func TestZapAdapter_FormattedLogger(t *testing.T) {
	logger := NewExample()

	// Test formatted logging
	logger.Debugf("debug: %s", "test")
	logger.Infof("info: %d", 123)
	logger.Warnf("warn: %v", true)
	logger.Errorf("error: %s", "test error")
}

func TestZapAdapter_StructuredLogger(t *testing.T) {
	logger := NewExample()

	// Test structured logging
	logger.Debugw("debug with fields", "key1", "value1", "key2", 123)
	logger.Infow("info with fields", "user", "john", "age", 30)
	logger.Warnw("warn with fields", "status", "warning")
	logger.Errorw("error with fields", "error", "something went wrong")
}

func TestZapAdapter_LineLogger(t *testing.T) {
	logger := NewExample()

	// Test ln-style logging
	logger.Debugln("debug", "line")
	logger.Infoln("info", "line")
	logger.Warnln("warn", "line")
	logger.Errorln("error", "line")
}

func TestZapAdapter_ContextualLogger(t *testing.T) {
	logger := NewExample()

	// Test contextual logging
	childLogger := logger.With("component", "test")
	childLogger.(core.ISugaredLogger).Info("message with context")

	namedLogger := logger.Named("myapp")
	namedLogger.(core.ISugaredLogger).Info("message with name")

	lazyLogger := logger.WithLazy("lazy", "value")
	lazyLogger.(core.ISugaredLogger).Info("message with lazy field")
}

func TestZapAdapter_Logf(t *testing.T) {
	logger := NewExample()

	// Test Logf with different levels
	logger.Logf(core.DebugLevel, "debug: %s", "test")
	logger.Logf(core.InfoLevel, "info: %d", 123)
	logger.Logf(core.WarnLevel, "warn: %v", true)
	logger.Logf(core.ErrorLevel, "error: %s", "test")
}

func TestZapAdapter_Logw(t *testing.T) {
	logger := NewExample()

	// Test Logw with different levels
	logger.Logw(core.DebugLevel, "debug message", "key", "value")
	logger.Logw(core.InfoLevel, "info message", "count", 42)
	logger.Logw(core.WarnLevel, "warn message", "status", "warning")
	logger.Logw(core.ErrorLevel, "error message", "error", "test error")
}

func TestZapAdapter_Logln(t *testing.T) {
	logger := NewExample()

	// Test Logln with different levels
	logger.Logln(core.DebugLevel, "debug", "line")
	logger.Logln(core.InfoLevel, "info", "line")
	logger.Logln(core.WarnLevel, "warn", "line")
	logger.Logln(core.ErrorLevel, "error", "line")
}

func TestZapAdapter_Level(t *testing.T) {
	tests := []struct {
		name          string
		createLogger  func() core.ISugaredLogger
		expectedLevel core.Level
	}{
		{
			name: "Development logger has DebugLevel",
			createLogger: func() core.ISugaredLogger {
				logger, _ := NewDevelopment()
				return logger
			},
			expectedLevel: core.DebugLevel,
		},
		{
			name: "Production logger has InfoLevel",
			createLogger: func() core.ISugaredLogger {
				logger, _ := NewProduction()
				return logger
			},
			expectedLevel: core.InfoLevel,
		},
		{
			name: "Example logger has DebugLevel",
			createLogger: func() core.ISugaredLogger {
				return NewExample()
			},
			expectedLevel: core.DebugLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := tt.createLogger()
			if logger.Level() != tt.expectedLevel {
				t.Errorf("Level() = %v, want %v", logger.Level(), tt.expectedLevel)
			}
		})
	}
}

func TestZapAdapter_Sync(t *testing.T) {
	logger := NewExample()

	// Test sync (may return error for stdout, which is acceptable)
	_ = logger.Sync()
}

func TestZapAdapter_Desugar(t *testing.T) {
	logger := NewExample()

	// Test desugar returns non-nil
	underlying := logger.Desugar()
	if underlying == nil {
		t.Error("Desugar() should not return nil")
	}
}

func TestCoreToZapLevel(t *testing.T) {
	tests := []struct {
		name      string
		coreLevel core.Level
	}{
		{"DebugLevel", core.DebugLevel},
		{"InfoLevel", core.InfoLevel},
		{"WarnLevel", core.WarnLevel},
		{"ErrorLevel", core.ErrorLevel},
		{"DPanicLevel", core.DPanicLevel},
		{"PanicLevel", core.PanicLevel},
		{"FatalLevel", core.FatalLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zapLevel := coreToZapLevel(tt.coreLevel)
			converted := zapToCoreLevel(zapLevel)
			if converted != tt.coreLevel {
				t.Errorf("Level conversion round-trip failed: got %v, want %v", converted, tt.coreLevel)
			}
		})
	}
}

func TestZapAdapter_WithPreservesLevel(t *testing.T) {
	logger := NewExample()
	originalLevel := logger.Level()

	// Test that With preserves level
	childLogger := logger.With("key", "value")
	if childLogger.(core.ISugaredLogger).Level() != originalLevel {
		t.Error("With() should preserve the logger level")
	}

	// Test that Named preserves level
	namedLogger := logger.Named("test")
	if namedLogger.(core.ISugaredLogger).Level() != originalLevel {
		t.Error("Named() should preserve the logger level")
	}

	// Test that WithLazy preserves level
	lazyLogger := logger.WithLazy("lazy", "value")
	if lazyLogger.(core.ISugaredLogger).Level() != originalLevel {
		t.Error("WithLazy() should preserve the logger level")
	}
}

func BenchmarkZapAdapter_Info(b *testing.B) {
	logger := NewExample()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message")
	}
}

func BenchmarkZapAdapter_Infof(b *testing.B) {
	logger := NewExample()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infof("benchmark message: %d", i)
	}
}

func BenchmarkZapAdapter_Infow(b *testing.B) {
	logger := NewExample()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infow("benchmark message", "iteration", i)
	}
}
