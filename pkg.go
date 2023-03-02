package log

import (
	"bytes"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

var defaultLogger = NewWithOptions(os.Stderr, Options{ReportTimestamp: true})

// Default returns the default logger. The default logger comes with timestamp enabled.
func Default() *Logger {
	return defaultLogger
}

// SetDefault sets the default global logger.
func SetDefault(logger *Logger) {
	defaultLogger = logger
}

// New returns a new logger with the default options.
func New(w io.Writer) *Logger {
	return NewWithOptions(w, Options{})
}

// NewWithOptions returns a new logger using the provided options.
func NewWithOptions(w io.Writer, o Options) *Logger {
	l := &Logger{
		b:               bytes.Buffer{},
		mu:              &sync.RWMutex{},
		helpers:         &sync.Map{},
		level:           int32(o.Level),
		reportTimestamp: o.ReportTimestamp,
		reportCaller:    o.ReportCaller,
		prefix:          o.Prefix,
		timeFunc:        o.TimeFunction,
		timeFormat:      o.TimeFormat,
		formatter:       o.Formatter,
		fields:          o.Fields,
	}

	l.SetOutput(w)
	l.SetLevel(Level(l.level))

	if l.timeFunc == nil {
		l.timeFunc = time.Now
	}

	if l.timeFormat == "" {
		l.timeFormat = DefaultTimeFormat
	}

	return l
}

// SetReportTimestamp sets whether to report timestamp for the default logger.
func SetReportTimestamp(report bool) {
	defaultLogger.SetReportTimestamp(report)
}

// SetReportCaller sets whether to report caller location for the default logger.
func SetReportCaller(report bool) {
	defaultLogger.SetReportCaller(report)
}

// SetLevel sets the level for the default logger.
func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}

// GetLevel returns the level for the default logger.
func GetLevel() Level {
	return defaultLogger.GetLevel()
}

// SetTimeFormat sets the time format for the default logger.
func SetTimeFormat(format string) {
	defaultLogger.SetTimeFormat(format)
}

// SetTimeFunction sets the time function for the default logger.
func SetTimeFunction(f TimeFunction) {
	defaultLogger.SetTimeFunction(f)
}

// SetOutput sets the output for the default logger.
func SetOutput(w io.Writer) {
	defaultLogger.SetOutput(w)
}

// SetFormatter sets the formatter for the default logger.
func SetFormatter(f Formatter) {
	defaultLogger.SetFormatter(f)
}

// SetPrefix sets the prefix for the default logger.
func SetPrefix(prefix string) {
	defaultLogger.SetPrefix(prefix)
}

// GetPrefix returns the prefix for the default logger.
func GetPrefix() string {
	return defaultLogger.GetPrefix()
}

// With returns a new logger with the given keyvals.
func With(keyvals ...interface{}) *Logger {
	return defaultLogger.With(keyvals...)
}

// Helper marks the calling function as a helper
// and skips it for source location information.
// It's the equivalent of testing.TB.Helper().
func Helper() {
	// skip this function frame
	defaultLogger.helper(1)
}

// Debug logs a debug message.
func Debug(msg interface{}, keyvals ...interface{}) {
	defaultLogger.log(DebugLevel, msg, keyvals...)
}

// Info logs an info message.
func Info(msg interface{}, keyvals ...interface{}) {
	defaultLogger.log(InfoLevel, msg, keyvals...)
}

// Warn logs a warning message.
func Warn(msg interface{}, keyvals ...interface{}) {
	defaultLogger.log(WarnLevel, msg, keyvals...)
}

// Error logs an error message.
func Error(msg interface{}, keyvals ...interface{}) {
	defaultLogger.log(ErrorLevel, msg, keyvals...)
}

// Fatal logs a fatal message and exit.
func Fatal(msg interface{}, keyvals ...interface{}) {
	defaultLogger.log(FatalLevel, msg, keyvals...)
	os.Exit(1)
}

// Print logs a message with no level.
func Print(msg interface{}, keyvals ...interface{}) {
	defaultLogger.log(noLevel, msg, keyvals...)
}

// StandardLog returns a standard logger from the default logger.
func StandardLog(opts ...StandardLogOptions) *log.Logger {
	return defaultLogger.StandardLog(opts...)
}
