package zap

import (
	"github.com/phongthien99/monorepo-lib/libs/log/core"
	"go.uber.org/zap/zapcore"
)

// coreToZapLevel converts our core.Level to zapcore.Level
func coreToZapLevel(level core.Level) zapcore.Level {
	switch level {
	case core.DebugLevel:
		return zapcore.DebugLevel
	case core.InfoLevel:
		return zapcore.InfoLevel
	case core.WarnLevel:
		return zapcore.WarnLevel
	case core.ErrorLevel:
		return zapcore.ErrorLevel
	case core.DPanicLevel:
		return zapcore.DPanicLevel
	case core.PanicLevel:
		return zapcore.PanicLevel
	case core.FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// zapToCoreLevel converts zapcore.Level to our core.Level
func zapToCoreLevel(level zapcore.Level) core.Level {
	switch level {
	case zapcore.DebugLevel:
		return core.DebugLevel
	case zapcore.InfoLevel:
		return core.InfoLevel
	case zapcore.WarnLevel:
		return core.WarnLevel
	case zapcore.ErrorLevel:
		return core.ErrorLevel
	case zapcore.DPanicLevel:
		return core.DPanicLevel
	case zapcore.PanicLevel:
		return core.PanicLevel
	case zapcore.FatalLevel:
		return core.FatalLevel
	default:
		return core.InfoLevel
	}
}
