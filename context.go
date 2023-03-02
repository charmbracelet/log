package log

import "context"

// WithContext wraps the given logger in context.
func WithContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

// FromContext returns the logger from the given context.
// This will return the default package logger if no logger
// found in context.
func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(loggerContextKey).(*Logger); ok {
		return logger
	}
	return defaultLogger
}

type contextKey struct{}

var loggerContextKey = contextKey{}
