package log

import (
	"io"
	"log"
)

var defaultLogger = New(WithTimestamp()).(*logger)

// Default returns the default logger. The default logger comes with timestamp enabled.
func Default() Logger {
	return defaultLogger
}

// EnableTimestamp enables timestamps for the default logger.
func EnableTimestamp() {
	defaultLogger.EnableTimestamp()
}

// DisableTimestamp disables timestamps for the default logger.
func DisableTimestamp() {
	defaultLogger.DisableTimestamp()
}

// EnableCaller enables caller for the default logger.
func EnableCaller() {
	defaultLogger.EnableCaller()
}

// DisableCaller disables caller for the default logger.
func DisableCaller() {
	defaultLogger.DisableCaller()
}

// EnableStyles enables colored output for the default logger.
func EnableStyles() {
	defaultLogger.EnableStyles()
}

// DisableStyles disables colored output for the default logger.
func DisableStyles() {
	defaultLogger.DisableStyles()
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

// SetPrefix sets the prefix for the default logger.
func SetPrefix(prefix string) {
	defaultLogger.SetPrefix(prefix)
}

// GetPrefix returns the prefix for the default logger.
func GetPrefix() string {
	return defaultLogger.GetPrefix()
}

// With returns a new logger with the given keyvals.
func With(keyvals ...interface{}) Logger {
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
	defaultLogger.log(DebugLevel, 0, msg, keyvals...)
}

// Info logs an info message.
func Info(msg interface{}, keyvals ...interface{}) {
	defaultLogger.log(InfoLevel, 0, msg, keyvals...)
}

// Warn logs a warning message.
func Warn(msg interface{}, keyvals ...interface{}) {
	defaultLogger.log(WarnLevel, 0, msg, keyvals...)
}

// Error logs an error message.
func Error(msg interface{}, keyvals ...interface{}) {
	defaultLogger.log(ErrorLevel, 0, msg, keyvals...)
}

// Fatal logs a fatal message and exit.
func Fatal(msg interface{}, keyvals ...interface{}) {
	defaultLogger.log(FatalLevel, 0, msg, keyvals...)
}

// Print logs a message with no level.
func Print(msg interface{}, keyvals ...interface{}) {
	defaultLogger.log(noLevel, 0, msg, keyvals...)
}

// StandardLogger returns a standard logger from the default logger.
func StandardLogger() *log.Logger {
	return defaultLogger.StandardLogger()
}
