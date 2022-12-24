package log

import "io"

var defaultLogger = New(WithTimestamp())

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

// Debug logs a debug message.
func Debug(msg interface{}, keyvals ...interface{}) {
	defaultLogger.Debug(msg, keyvals...)
}

// Info logs an info message.
func Info(msg interface{}, keyvals ...interface{}) {
	defaultLogger.Info(msg, keyvals...)
}

// Warn logs a warning message.
func Warn(msg interface{}, keyvals ...interface{}) {
	defaultLogger.Warn(msg, keyvals...)
}

// Error logs an error message.
func Error(msg interface{}, keyvals ...interface{}) {
	defaultLogger.Error(msg, keyvals...)
}
