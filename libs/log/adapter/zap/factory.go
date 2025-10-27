package zap

import (
	"fmt"

	"github.com/phongthien99/monorepo-lib/libs/log/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewDevelopment creates a development logger (human-friendly, colorful)
func NewDevelopment() (core.ISugaredLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return NewZapAdapterFromLogger(logger, core.DebugLevel), nil
}

// NewProduction creates a production logger (JSON, optimized for performance)
func NewProduction() (core.ISugaredLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return NewZapAdapterFromLogger(logger, core.InfoLevel), nil
}

// NewExample creates an example logger (for testing/examples)
func NewExample() core.ISugaredLogger {
	logger := zap.NewExample()
	return NewZapAdapterFromLogger(logger, core.DebugLevel)
}

// NewNop creates a no-op logger (for testing)
func NewNop() core.ISugaredLogger {
	return NewZapAdapterFromLogger(zap.NewNop(), core.InfoLevel)
}

// Config holds the configuration for creating a custom logger
type Config struct {
	Level            core.Level
	Development      bool
	Encoding         string // "json" or "console"
	OutputPaths      []string
	ErrorOutputPaths []string
}

// NewWithConfig creates a logger with custom configuration
func NewWithConfig(cfg Config) (core.ISugaredLogger, error) {
	// Validate and set defaults
	if cfg.Encoding == "" {
		cfg.Encoding = "json"
	}
	if cfg.Encoding != "json" && cfg.Encoding != "console" {
		return nil, fmt.Errorf("invalid encoding: %s (must be 'json' or 'console')", cfg.Encoding)
	}
	if len(cfg.OutputPaths) == 0 {
		cfg.OutputPaths = []string{"stdout"}
	}
	if len(cfg.ErrorOutputPaths) == 0 {
		cfg.ErrorOutputPaths = []string{"stderr"}
	}

	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(coreToZapLevel(cfg.Level)),
		Development:      cfg.Development,
		Encoding:         cfg.Encoding,
		OutputPaths:      cfg.OutputPaths,
		ErrorOutputPaths: cfg.ErrorOutputPaths,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return NewZapAdapterFromLogger(logger, cfg.Level), nil
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		Level:            core.InfoLevel,
		Development:      false,
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// DevelopmentConfig returns a development-friendly configuration
func DevelopmentConfig() Config {
	return Config{
		Level:            core.DebugLevel,
		Development:      true,
		Encoding:         "console",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}
