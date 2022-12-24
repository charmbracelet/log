package log

import "io"

// WithOutput returns a LoggerOption that sets the output for the logger.
func WithOutput(w io.Writer) LoggerOption {
	return func(l *logger) {
		l.w = w
	}
}

// WithTimeFunction returns a LoggerOption that sets the time function for the logger.
func WithTimeFunction(f TimeFunction) LoggerOption {
	return func(l *logger) {
		l.timeFunc = f
	}
}

// WithTimeFormat returns a LoggerOption that sets the time format for the logger.
func WithTimeFormat(format string) LoggerOption {
	return func(l *logger) {
		l.timeFormat = format
	}
}

// WithLevel returns a LoggerOption that sets the level for the logger.
func WithLevel(level Level) LoggerOption {
	return func(l *logger) {
		l.level = level
	}
}

// WithPrefix returns a LoggerOption that sets the prefix for the logger.
func WithPrefix(prefix string) LoggerOption {
	return func(l *logger) {
		l.prefix = prefix
	}
}

// WithNoColor returns a LoggerOption that disables colors for the logger.
func WithNoColor() LoggerOption {
	return func(l *logger) {
		l.noColor = true
	}
}

// WithTimestamp returns a LoggerOption that enables timestamps for the logger.
func WithTimestamp() LoggerOption {
	return func(l *logger) {
		l.timestamp = true
	}
}

// WithCaller returns a LoggerOption that enables caller for the logger.
func WithCaller() LoggerOption {
	return func(l *logger) {
		l.caller = true
	}
}

// WithFields returns a LoggerOption that sets the fields for the logger.
func WithFields(keyvals ...interface{}) LoggerOption {
	return func(l *logger) {
		l.keyvals = keyvals
	}
}
