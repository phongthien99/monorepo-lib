package zap

import (
	"github.com/phongthien99/monorepo-lib/libs/log/core"
	"go.uber.org/zap"
)

// zapAdapter wraps zap.SugaredLogger to implement our ISugaredLogger interface
type zapAdapter struct {
	logger *zap.SugaredLogger
	level  core.Level
}

// NewZapAdapter creates a new adapter that wraps zap.SugaredLogger
func NewZapAdapter(zapLogger *zap.SugaredLogger, level core.Level) core.ISugaredLogger {
	return &zapAdapter{
		logger: zapLogger,
		level:  level,
	}
}

// NewZapAdapterFromLogger creates a new adapter from zap.Logger
func NewZapAdapterFromLogger(zapLogger *zap.Logger, level core.Level) core.ISugaredLogger {
	return &zapAdapter{
		logger: zapLogger.Sugar(),
		level:  level,
	}
}

// IBasicLogger implementation
func (z *zapAdapter) Debug(args ...any) {
	z.logger.Debug(args...)
}

func (z *zapAdapter) Info(args ...any) {
	z.logger.Info(args...)
}

func (z *zapAdapter) Warn(args ...any) {
	z.logger.Warn(args...)
}

func (z *zapAdapter) Error(args ...any) {
	z.logger.Error(args...)
}

func (z *zapAdapter) DPanic(args ...any) {
	z.logger.DPanic(args...)
}

func (z *zapAdapter) Panic(args ...any) {
	z.logger.Panic(args...)
}

func (z *zapAdapter) Fatal(args ...any) {
	z.logger.Fatal(args...)
}

// IFormattedLogger implementation
func (z *zapAdapter) Debugf(template string, args ...any) {
	z.logger.Debugf(template, args...)
}

func (z *zapAdapter) Infof(template string, args ...any) {
	z.logger.Infof(template, args...)
}

func (z *zapAdapter) Warnf(template string, args ...any) {
	z.logger.Warnf(template, args...)
}

func (z *zapAdapter) Errorf(template string, args ...any) {
	z.logger.Errorf(template, args...)
}

func (z *zapAdapter) DPanicf(template string, args ...any) {
	z.logger.DPanicf(template, args...)
}

func (z *zapAdapter) Panicf(template string, args ...any) {
	z.logger.Panicf(template, args...)
}

func (z *zapAdapter) Fatalf(template string, args ...any) {
	z.logger.Fatalf(template, args...)
}

func (z *zapAdapter) Logf(level core.Level, template string, args ...any) {
	z.logger.Logf(coreToZapLevel(level), template, args...)
}

// IStructuredLogger implementation
func (z *zapAdapter) Debugw(msg string, keysAndValues ...any) {
	z.logger.Debugw(msg, keysAndValues...)
}

func (z *zapAdapter) Infow(msg string, keysAndValues ...any) {
	z.logger.Infow(msg, keysAndValues...)
}

func (z *zapAdapter) Warnw(msg string, keysAndValues ...any) {
	z.logger.Warnw(msg, keysAndValues...)
}

func (z *zapAdapter) Errorw(msg string, keysAndValues ...any) {
	z.logger.Errorw(msg, keysAndValues...)
}

func (z *zapAdapter) DPanicw(msg string, keysAndValues ...any) {
	z.logger.DPanicw(msg, keysAndValues...)
}

func (z *zapAdapter) Panicw(msg string, keysAndValues ...any) {
	z.logger.Panicw(msg, keysAndValues...)
}

func (z *zapAdapter) Fatalw(msg string, keysAndValues ...any) {
	z.logger.Fatalw(msg, keysAndValues...)
}

func (z *zapAdapter) Logw(level core.Level, msg string, keysAndValues ...any) {

	z.logger.Logw(coreToZapLevel(level), msg, keysAndValues...)
}

// ILineLogger implementation
func (z *zapAdapter) Debugln(args ...any) {
	z.logger.Debugln(args...)
}

func (z *zapAdapter) Infoln(args ...any) {
	z.logger.Infoln(args...)
}

func (z *zapAdapter) Warnln(args ...any) {
	z.logger.Warnln(args...)
}

func (z *zapAdapter) Errorln(args ...any) {
	z.logger.Errorln(args...)
}

func (z *zapAdapter) DPanicln(args ...any) {
	z.logger.DPanicln(args...)
}

func (z *zapAdapter) Panicln(args ...any) {
	z.logger.Panicln(args...)
}

func (z *zapAdapter) Fatalln(args ...any) {
	z.logger.Fatalln(args...)
}

func (z *zapAdapter) Logln(level core.Level, args ...any) {
	z.logger.Logln(coreToZapLevel(level), args...)
}

// IContextualLogger implementation
func (z *zapAdapter) With(args ...any) core.ISugaredLogger {
	return &zapAdapter{
		logger: z.logger.With(args...),
		level:  z.level,
	}
}

func (z *zapAdapter) WithLazy(args ...any) core.ISugaredLogger {
	return &zapAdapter{
		logger: z.logger.WithLazy(args...),
		level:  z.level,
	}
}

func (z *zapAdapter) Named(name string) core.ISugaredLogger {
	return &zapAdapter{
		logger: z.logger.Named(name),
		level:  z.level,
	}
}

// IContextLogger implementation
func (z *zapAdapter) WithContext(ctx any) core.ISugaredLogger {
	// Extract trace info from context if available
	// This is a simple implementation, can be extended with OpenTelemetry
	return &zapAdapter{
		logger: z.logger,
		level:  z.level,
	}
}

// ILoggerControl implementation
func (z *zapAdapter) Desugar() any {
	return z.logger.Desugar()
}

func (z *zapAdapter) Level() core.Level {
	return z.level
}

func (z *zapAdapter) Sync() error {
	return z.logger.Sync()
}
