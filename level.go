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
	// OffLevel is the off level.
	OffLevel
)

// String returns the string representation of the level.
func (l Level) String() string {
	return [...]string{"debug", "info", "warn", "error"}[l]
}
