package log

import (
	"fmt"
	"io"
	stdlog "log"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// Logger is the main type in the log package.
type Output struct {
	logger log.Logger
	w      io.Writer
}

var _ log.Logger = (*Output)(nil)
var _ Logger = (*Output)(nil)

// Option is a functional option type for configuring a Logger.
type Option func(*Output)

// New returns a new Logger.
func New(w io.Writer, opts ...Option) *Output {
	var logger log.Logger
	if w == nil {
		w = os.Stderr
	}
	w = log.NewSyncWriter(w)
	if IsTerminal(w) {
		logger = newColorLogger(w, log.NewLogfmtLogger, DefaultStyles())
	} else {
		logger = log.NewLogfmtLogger(w)
	}
	output := &Output{
		logger: logger,
	}
	output.SetLevel(DebugLevel)
	output.SetOptions(opts...)
	return output
}

// SetOutput implements the Logger interface.
func (l *Output) SetOutput(w io.Writer) {
	l.logger = New(w)
}

// SetLevel implements the Logger interface.
func (l *Output) SetLevel(lvl Level) {
	stdlog.Printf("SetLevel %s", lvl)
	l.logger = level.NewFilter(l.logger, levelOption(lvl))
}

// SetOptions implements Logger.
func (l *Output) SetOptions(opts ...Option) {
	for _, opt := range opts {
		opt(l)
	}
}

// SetFields implements Logger.
func (l *Output) SetFields(keyvals ...interface{}) {
	l.logger = log.With(l.logger, keyvals...)
}

// Log implements the log.Logger interface.
func (l *Output) Log(keyvals ...interface{}) error {
	return l.logger.Log(keyvals...)
}

// WithFields implements Logger.
func WithFields(keyvals ...interface{}) Option {
	return func(l *Output) {
		l.SetFields(keyvals...)
	}
}

// WithTimestamp returns a logger option that adds a timestamp to each log event.
func WithTimestamp() Option {
	return WithTimestampFormat(DefaultTimestampFormat)
}

// WithTimestampUTC returns a logger option that adds a UTC timestamp to each log event.
func WithTimestampUTC() Option {
	return WithTimestampFormat(DefaultTimestampFormat)
}

// WithTimestampUTCFormat returns a logger option that adds a UTC timestamp to
// each log event.
func WithTimestampUTCFormat(format string) Option {
	return func(l *Output) {
		l.logger = log.With(l.logger, tsKey, log.TimestampFormat(
			func() time.Time { return time.Now().UTC() },
			format,
		))
	}
}

// WithTimestampFormat returns a logger option that adds a timestamp to each log event.
func WithTimestampFormat(format string) Option {
	return func(l *Output) {
		l.logger = log.With(l.logger, tsKey,
			log.TimestampFormat(time.Now, format))
	}
}

// WithLevel returns a logger option that sets the log level.
func WithLevel(lvl Level) Option {
	return func(l *Output) {
		l.SetLevel(lvl)
	}
}

// With returns a new Logger with the given keyvals set.
func (l *Output) With(keyvals ...interface{}) Logger {
	return &Output{
		logger: log.With(l.logger, keyvals...),
	}
}

// WithError returns a new Logger with the given error.
func (l *Output) WithError(err error) Logger {
	return &Output{
		logger: log.With(l.logger, errKey, err),
	}
}

// Debug implements Logger
func (l *Output) Debug(v ...interface{}) {
	level.Debug(l.logger).Log(msgKey, fmt.Sprint(v...))
}

// Debugf implements Logger
func (l *Output) Debugf(format string, v ...interface{}) {
	l.Debug(fmt.Sprintf(format, v...))
}

// Debugln implements Logger
func (l *Output) Debugln(v ...interface{}) {
	l.Debug(v...)
}

// Error implements Logger
func (l *Output) Error(v ...interface{}) {
	level.Debug(l.logger).Log(msgKey, fmt.Sprint(v...))
}

// Errorf implements Logger
func (l *Output) Errorf(format string, v ...interface{}) {
	l.Error(fmt.Sprintf(format, v...))
}

// Errorln implements Logger
func (l *Output) Errorln(v ...interface{}) {
	l.Error(v...)
}

// Fatal implements Logger
func (l *Output) Fatal(v ...interface{}) {
	l.Error(v...)
	os.Exit(1)
}

// Fatalf implements Logger
func (l *Output) Fatalf(format string, v ...interface{}) {
	l.Fatal(fmt.Sprintf(format, v...))

}

// Fatalln implements Logger
func (l *Output) Fatalln(v ...interface{}) {
	l.Fatal(v...)
}

// Info implements Logger
func (l *Output) Info(v ...interface{}) {
	level.Info(l.logger).Log(msgKey, fmt.Sprint(v...))
}

// Infof implements Logger
func (l *Output) Infof(format string, v ...interface{}) {
	l.Info(fmt.Sprintf(format, v...))
}

// Infoln implements Logger
func (l *Output) Infoln(v ...interface{}) {
	l.Info(v...)
}

// Print implements Logger
func (l *Output) Print(v ...interface{}) {
	l.Info(v...)
}

// Printf implements Logger
func (l *Output) Printf(format string, v ...interface{}) {
	l.Info(fmt.Sprintf(format, v...))
}

// Println implements Logger
func (l *Output) Println(v ...interface{}) {
	l.Info(v...)
}

// Warn implements Logger
func (l *Output) Warn(v ...interface{}) {
	level.Warn(l.logger).Log(msgKey, fmt.Sprint(v...))
}

// Warnf implements Logger
func (l *Output) Warnf(format string, v ...interface{}) {
	l.Warn(fmt.Sprintf(format, v...))
}

// Warnln implements Logger
func (l *Output) Warnln(v ...interface{}) {
	l.Warn(v...)
}
