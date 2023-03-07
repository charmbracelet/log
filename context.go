package log

import "context"

// WithContext wraps the given logger in context.
func WithContext(ctx context.Context, logger Logger, keyvals ...interface{}) context.Context {
	if len(keyvals) > 0 {
		logger = logger.With(keyvals...)
	}
	return context.WithValue(ctx, loggerContextKey, &logger)
}

// FromContext returns the logger from the given context.
// This will return the default package logger if no logger
// found in context.
func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerContextKey).(*Logger); ok {
		return *logger
	}
	return defaultLogger
}

// UpdateContext updates the logger in the given context. Returns a boolean if the logger was
// successfully updated. If there's no logger in the context, this will return false.
func UpdateContext(ctx context.Context, fn func(Logger) Logger) bool {
	loggerPtr := ctx.Value(loggerContextKey)
	if loggerPtr == nil {
		return false
	}

	logger, ok := loggerPtr.(*Logger)
	if !ok {
		return false
	}

	*logger = fn(*logger)

	return true
}

type contextKey struct{}

var loggerContextKey = contextKey{}
