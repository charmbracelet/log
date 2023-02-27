package log

import "io"

// WithOutput returns a LoggerOption that sets the output for the logger. The
// default is os.Stderr.
func WithOutput(w io.Writer) LoggerOption {
	return func(l *logger) {
		l.w = w
	}
}

// WithTimeFunction returns a LoggerOption that sets the time function for the
// logger. The default is time.Now.
func WithTimeFunction(f TimeFunction) LoggerOption {
	return func(l *logger) {
		l.timeFunc = f
	}
}

// WithTimeFormat returns a LoggerOption that sets the time format for the
// logger. The default is "2006/01/02 15:04:05".
func WithTimeFormat(format string) LoggerOption {
	return func(l *logger) {
		l.timeFormat = format
	}
}

// WithLevel returns a LoggerOption that sets the level for the logger. The
// default is InfoLevel.
func WithLevel(level Level) LoggerOption {
	return func(l *logger) {
		l.level = int32(level)
	}
}

// WithPrefix returns a LoggerOption that sets the prefix for the logger.
func WithPrefix(prefix string) LoggerOption {
	return func(l *logger) {
		l.prefix = prefix
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

// WithFormatter returns a LoggerOption that sets the formatter for the logger.
func WithFormatter(f Formatter) LoggerOption {
	return func(l *logger) {
		l.formatter = f
	}
}
