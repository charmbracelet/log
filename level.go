package log

import (
	"fmt"
	"strings"
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

// ErrInvalidLevel is an error returned when parsing an invalid level string.
var ErrInvalidLevel = errors.New("invalid level")
// ParseLevel converts level in string to Level type. Default level is InfoLevel.
func ParseLevel(level string) (Level, error) {
	switch strings.ToLower(level) {
	case DebugLevel.String():
		return DebugLevel, nil
	case InfoLevel.String():
		return InfoLevel, nil
	case WarnLevel.String():
		return WarnLevel, nil
	case ErrorLevel.String():
		return ErrorLevel, nil
	case FatalLevel.String():
		return FatalLevel, nil
	default:
		return InfoLevel, fmt.Errorf("%w: %q", ErrInvalidLevel, level)
	}
}
