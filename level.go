package log

// Level is a logging level.
type Level int32

const (
	// DebugLevel is the debug level.
	DebugLevel Level = iota
	// InfoLevel is the info level.
	InfoLevel
	// WarnLevel is the warn level.
	WarnLevel
	// ErrorLevel is the error level.
	ErrorLevel
	// FatalLevel is the fatal level.
	FatalLevel
	// NoLevel is the no level.
	NoLevel
	// OffLevel is the off level.
	OffLevel
)

// String returns the string representation of the level.
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case OffLevel:
		return "off"
	default:
		return ""
	}
}
