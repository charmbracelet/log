package log

import (
	"io"
	"os"
	"time"
)

var DefaultTimestampFormat = time.Kitchen

var defaultLogger = NewDefault(os.Stderr)

// NewDefault returns a new default logger.
func NewDefault(w io.Writer, opts ...Option) Logger {
	return New(w, append(opts, WithTimestamp(), WithLevel(InfoLevel))...)
}

// Default returns the default logger.
func Default() Logger {
	return defaultLogger
}

// SetOutput sets the default logger's output.
func SetOutput(w io.Writer) {
	defaultLogger.SetOutput(w)
}

// SetLevel sets the default logger's level.
func SetLevel(lvl Level) {
	defaultLogger.SetLevel(lvl)
}

// SetOptions sets the default logger's options.
func SetOptions(opts ...Option) {
	defaultLogger.SetOptions(opts...)
}

// SetFields sets the default logger's fields.
func SetFields(keyvals ...interface{}) {
	defaultLogger.SetFields(keyvals...)
}

// With returns a new logger with the given keyvals.
func With(keyvals ...interface{}) Logger {
	return defaultLogger.With(keyvals...)
}

// WithError returns a new logger with the given error.
func WithError(err error) Logger {
	return defaultLogger.WithError(err)
}

// Debug logs a debug message.
func Debug(v ...interface{}) {
	defaultLogger.Debug(v...)
}

// Print logs an into message.
func Print(v ...interface{}) {
	defaultLogger.Print(v...)
}

// Info logs an info message.
func Info(v ...interface{}) {
	defaultLogger.Info(v...)
}

// Warn logs a warning message.
func Warn(v ...interface{}) {
	defaultLogger.Warn(v...)
}

// Error logs an error message.
func Error(v ...interface{}) {
	defaultLogger.Error(v...)
}

// Fatal logs a fatal message.
func Fatal(v ...interface{}) {
	defaultLogger.Fatal(v...)
}

// Debugln logs a debug message.
func Debugln(v ...interface{}) {
	defaultLogger.Debugln(v...)
}

// Println logs an info message.
func Println(v ...interface{}) {
	defaultLogger.Println(v...)
}

// Infoln logs an info message.
func Infoln(v ...interface{}) {
	defaultLogger.Infoln(v...)
}

// Warnln logs a warning message.
func Warnln(v ...interface{}) {
	defaultLogger.Warnln(v...)
}

// Errorln logs an error message.
func Errorln(v ...interface{}) {
	defaultLogger.Errorln(v...)
}

// Fatalln logs a fatal message.
func Fatalln(v ...interface{}) {
	defaultLogger.Fatalln(v...)
}

// Debugf logs a debug message.
func Debugf(format string, v ...interface{}) {
	defaultLogger.Debugf(format, v...)
}

// Printf logs an info message.
func Printf(format string, v ...interface{}) {
	defaultLogger.Printf(format, v...)
}

// Infof logs an info message.
func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

// Warnf logs a warning message.
func Warnf(format string, v ...interface{}) {
	defaultLogger.Warnf(format, v...)
}

// Errorf logs an error message.
func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

// Fatalf logs a fatal message.
func Fatalf(format string, v ...interface{}) {
	defaultLogger.Fatalf(format, v...)
}
