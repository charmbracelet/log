package log

import "github.com/go-kit/log/level"

type Level = level.Value

var (
	// DebugLevel is the lowest level of logging.
	DebugLevel Level = level.DebugValue()
	// InfoLevel is the default level of logging.
	InfoLevel Level = level.InfoValue()
	// WarnLevel is the level of logging for warnings.
	WarnLevel Level = level.WarnValue()
	// ErrorLevel is the level of logging for errors.
	ErrorLevel Level = level.ErrorValue()
)

// option returns a level.Option that sets the level to l.
func levelOption(l Level) level.Option {
	return level.Allow(l)
}
