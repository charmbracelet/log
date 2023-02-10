package log

import (
	"io"
	"log"
	"time"
)

// DefaultTimeFormat is the default time format.
const DefaultTimeFormat = "2006/01/02 15:04:05"

// TimeFunction is a function that returns a time.Time.
type TimeFunction = func() time.Time

// NowUTC is a convenient function that returns the
// current time in UTC timezone.
//
// This is to be used as a time function.
// For example:
//
//	log.SetTimeFunction(log.NowUTC)
func NowUTC() time.Time {
	return time.Now().UTC()
}

// Logger is an interface for logging.
type Logger interface {
	// SetLevel sets the allowed level.
	SetLevel(level Level)
	// GetLevel returns the allowed level.
	GetLevel() Level

	// SetPrefix sets the logger prefix. The default is no prefix.
	SetPrefix(prefix string)
	// GetPrefix returns the logger prefix.
	GetPrefix() string

	// SetReportTimestamp sets whether the logger should report the timestamp.
	SetReportTimestamp(bool)
	// SetReportCaller sets whether the logger should report the caller location.
	SetReportCaller(bool)
	// SetTimeFunction sets the time function used to get the time.
	// The default is time.Now.
	//
	// To use UTC time instead of local time set the time
	// function to `NowUTC`.
	SetTimeFunction(f TimeFunction)
	// SetTimeFormat sets the time format. The default is "2006/01/02 15:04:05".
	SetTimeFormat(format string)
	// SetOutput sets the output destination. The default is os.Stderr.
	SetOutput(w io.Writer)
	// SetFormatter sets the formatter. The default is TextFormatter.
	SetFormatter(f Formatter)

	// Helper marks the calling function as a helper
	// and skips it for source location information.
	// It's the equivalent of testing.TB.Helper().
	Helper()

	// With returns a new sub logger with the given key value pairs.
	With(keyval ...interface{}) Logger

	// Debug logs a debug message.
	Debug(msg interface{}, keyval ...interface{})
	// Info logs an info message.
	Info(msg interface{}, keyval ...interface{})
	// Warn logs a warning message.
	Warn(msg interface{}, keyval ...interface{})
	// Error logs an error message.
	Error(msg interface{}, keyval ...interface{})
	// Fatal logs a fatal message.
	Fatal(msg interface{}, keyval ...interface{})
	// Print logs a message with no level.
	Print(msg interface{}, keyval ...interface{})

	// StandardLog returns a standard logger from this logger.
	StandardLog(...StandardLogOption) *log.Logger
}
