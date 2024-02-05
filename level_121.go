//go:build go1.21
// +build go1.21

package log

import (
	"log/slog"
)

// fromSlogLevel converts slog.Level to log.Level.
var fromSlogLevel = map[slog.Level]Level{
	slog.LevelDebug: DebugLevel,
	slog.LevelInfo:  InfoLevel,
	slog.LevelWarn:  WarnLevel,
	slog.LevelError: ErrorLevel,
	slog.Level(12):  FatalLevel,
}

var _ slog.Leveler = Level(0)

// Leveler is a dynamic logging leveler.
type Leveler = slog.Leveler

// Level implements slog.Leveler.
func (l Level) Level() slog.Level {
	return slog.Level(l)
}
