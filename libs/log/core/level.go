package core

// Level represents the log level
type Level int

const (
	// DebugLevel is for debug messages
	DebugLevel Level = iota - 1
	// InfoLevel is for informational messages (default)
	InfoLevel
	// WarnLevel is for warning messages
	WarnLevel
	// ErrorLevel is for error messages
	ErrorLevel
	// DPanicLevel is for messages that should panic in development
	DPanicLevel
	// PanicLevel is for messages that will panic
	PanicLevel
	// FatalLevel is for fatal messages that will exit the program
	FatalLevel
)

// String returns the string representation of the log level
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case DPanicLevel:
		return "DPANIC"
	case PanicLevel:
		return "PANIC"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}
