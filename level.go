package log

import (
	"strings"

	"golang.org/x/exp/slog"
)

// Level is a logging level.
type Level int32

const (
	// DebugLevel is the debug level.
	DebugLevel Level = iota - 1
	// InfoLevel is the info level.
	InfoLevel
	// WarnLevel is the warn level.
	WarnLevel
	// ErrorLevel is the error level.
	ErrorLevel
	// FatalLevel is the fatal level.
	FatalLevel
	// noLevel is used with log.Print.
	noLevel
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
	default:
		return ""
	}
}

// ParseLevel converts level in string to Level type. Default level is InfoLevel.
func ParseLevel(level string) Level {
	switch strings.ToLower(level) {
	case DebugLevel.String():
		return DebugLevel
	case InfoLevel.String():
		return InfoLevel
	case WarnLevel.String():
		return WarnLevel
	case ErrorLevel.String():
		return ErrorLevel
	case FatalLevel.String():
		return FatalLevel
	default:
		return InfoLevel
	}
}

// fromSlogLevel converts slog.Level to log.Level.
var fromSlogLevel = map[slog.Level]Level{
	slog.LevelDebug: DebugLevel,
	slog.LevelInfo:  InfoLevel,
	slog.LevelWarn:  WarnLevel,
	slog.LevelError: ErrorLevel,
	slog.Level(12):  FatalLevel,
}

// toSlogLevel converts log.Level to slog.Level.
var toSlogLevel = map[Level]slog.Level{
	DebugLevel: slog.LevelDebug,
	InfoLevel:  slog.LevelInfo,
	WarnLevel:  slog.LevelWarn,
	ErrorLevel: slog.LevelError,
	FatalLevel: slog.Level(12),
}
