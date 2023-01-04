package log

// Level is a logging level.
type Level int32

// Supported log levels. The default level is info.
const (
	// LevelDebug is the debug level.
	LevelDebug Level = iota
	// LevelInfo is the info level.
	LevelInfo
	// LevelWarn is the warn level.
	LevelWarn
	// LevelError is the error level.
	LevelError
	// LevelOff is the off level.
	LevelOff
)

// String returns the string representation of the level.
func (l Level) String() string {
	return [...]string{"debug", "info", "warn", "error"}[l]
}
