package core

// IBasicLogger provides basic logging methods with simple arguments
type IBasicLogger interface {
	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	DPanic(args ...any)
	Panic(args ...any)
	Fatal(args ...any)
}

// IFormattedLogger provides Printf-style formatted logging
type IFormattedLogger interface {
	Debugf(template string, args ...any)
	Infof(template string, args ...any)
	Warnf(template string, args ...any)
	Errorf(template string, args ...any)
	DPanicf(template string, args ...any)
	Panicf(template string, args ...any)
	Fatalf(template string, args ...any)
	Logf(level Level, template string, args ...any)
}

// IStructuredLogger provides structured logging with key-value pairs
type IStructuredLogger interface {
	Debugw(msg string, keysAndValues ...any)
	Infow(msg string, keysAndValues ...any)
	Warnw(msg string, keysAndValues ...any)
	Errorw(msg string, keysAndValues ...any)
	DPanicw(msg string, keysAndValues ...any)
	Panicw(msg string, keysAndValues ...any)
	Fatalw(msg string, keysAndValues ...any)
	Logw(level Level, msg string, keysAndValues ...any)
}

// ILineLogger provides ln-style logging (adds newline)
type ILineLogger interface {
	Debugln(args ...any)
	Infoln(args ...any)
	Warnln(args ...any)
	Errorln(args ...any)
	DPanicln(args ...any)
	Panicln(args ...any)
	Fatalln(args ...any)
	Logln(level Level, args ...any)
}

// IContextualLogger provides contextual logging capabilities
type IContextualLogger interface {
	With(args ...any) ISugaredLogger
	WithLazy(args ...any) ISugaredLogger
	Named(name string) ISugaredLogger
}

// IContextLogger provides context-aware logging
type IContextLogger interface {
	// WithContext adds context to the logger for distributed tracing
	WithContext(ctx any) ISugaredLogger
}

// ILoggerControl provides logger control and configuration
type ILoggerControl interface {
	Desugar() any // Returns the underlying logger
	Level() Level
	Sync() error
}

// ILogger is a minimal interface that combines basic and formatted logging
type ILogger interface {
	IBasicLogger
	IFormattedLogger
	ILoggerControl
}

// IFullLogger combines basic, formatted, and structured logging
type IFullLogger interface {
	IBasicLogger
	IFormattedLogger
	IStructuredLogger
	ILoggerControl
}

// ISugaredLogger is the complete interface combining all logging styles
type ISugaredLogger interface {
	IBasicLogger
	IFormattedLogger
	IStructuredLogger
	ILineLogger
	IContextualLogger
	IContextLogger
	ILoggerControl
}
