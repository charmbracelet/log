package log

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var (
	// ErrMissingValue is returned when a key is missing a value.
	ErrMissingValue = fmt.Errorf("missing value")
)

// LoggerOption is an option for a logger.
type LoggerOption = func(*Logger)

// Logger is a Logger that implements Logger.
type Logger struct {
	w  io.Writer
	b  bytes.Buffer
	mu *sync.RWMutex
	re *lipgloss.Renderer

	isDiscard uint32

	level           int32
	prefix          string
	timeFunc        TimeFunction
	timeFormat      string
	callerOffset    int
	callerFormatter CallerFormatter
	formatter       Formatter

	reportCaller    bool
	reportTimestamp bool

	fields []interface{}

	helpers *sync.Map
}

func (l *Logger) log(level Level, msg interface{}, keyvals ...interface{}) {
	if atomic.LoadUint32(&l.isDiscard) != 0 {
		return
	}

	// check if the level is allowed
	if atomic.LoadInt32(&l.level) > int32(level) {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	defer l.b.Reset()

	var kvs []interface{}
	if l.reportTimestamp {
		kvs = append(kvs, TimestampKey, l.timeFunc())
	}

	if level != noLevel {
		kvs = append(kvs, LevelKey, level)
	}

	if l.reportCaller {
		// Call stack is log.Error -> log.log (2)
		file, line, fn := l.fillLoc(l.callerOffset + 2)
		caller := l.callerFormatter(file, line, fn)
		kvs = append(kvs, CallerKey, caller)
	}

	if l.prefix != "" {
		kvs = append(kvs, PrefixKey, l.prefix+":")
	}

	if msg != nil {
		m := fmt.Sprint(msg)
		kvs = append(kvs, MessageKey, m)
	}

	// append logger fields
	kvs = append(kvs, l.fields...)
	if len(l.fields)%2 != 0 {
		kvs = append(kvs, ErrMissingValue)
	}
	// append the rest
	kvs = append(kvs, keyvals...)
	if len(keyvals)%2 != 0 {
		kvs = append(kvs, ErrMissingValue)
	}

	switch l.formatter {
	case LogfmtFormatter:
		l.logfmtFormatter(kvs...)
	case JSONFormatter:
		l.jsonFormatter(kvs...)
	default:
		l.textFormatter(kvs...)
	}

	_, _ = l.w.Write(l.b.Bytes())
}

// Helper marks the calling function as a helper
// and skips it for source location information.
// It's the equivalent of testing.TB.Helper().
func (l *Logger) Helper() {
	l.helper(1)
}

func (l *Logger) helper(skip int) {
	_, _, fn := location(skip + 1)
	l.helpers.LoadOrStore(fn, struct{}{})
}

func (l *Logger) fillLoc(skip int) (file string, line int, fn string) {
	// Copied from testing.T
	const maxStackLen = 50
	var pc [maxStackLen]uintptr

	// Skip two extra frames to account for this function
	// and runtime.Callers itself.
	n := runtime.Callers(skip+2, pc[:])
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		_, helper := l.helpers.Load(frame.Function)
		if !helper || !more {
			// Found a frame that wasn't a helper function.
			// Or we ran out of frames to check.
			return frame.File, frame.Line, frame.Function
		}
	}
}

func location(skip int) (file string, line int, fn string) {
	pc, file, line, _ := runtime.Caller(skip + 1)
	f := runtime.FuncForPC(pc)
	return file, line, f.Name()
}

// Cleanup a path by returning the last n segments of the path only.
func trimCallerPath(path string, n int) string {
	// lovely borrowed from zap
	// nb. To make sure we trim the path correctly on Windows too, we
	// counter-intuitively need to use '/' and *not* os.PathSeparator here,
	// because the path given originates from Go stdlib, specifically
	// runtime.Caller() which (as of Mar/17) returns forward slashes even on
	// Windows.
	//
	// See https://github.com/golang/go/issues/3335
	// and https://github.com/golang/go/issues/18151
	//
	// for discussion on the issue on Go side.

	// Return the full path if n is 0.
	if n <= 0 {
		return path
	}

	// Find the last separator.
	idx := strings.LastIndexByte(path, '/')
	if idx == -1 {
		return path
	}

	for i := 0; i < n-1; i++ {
		// Find the penultimate separator.
		idx = strings.LastIndexByte(path[:idx], '/')
		if idx == -1 {
			return path
		}
	}

	return path[idx+1:]
}

// SetReportTimestamp sets whether the timestamp should be reported.
func (l *Logger) SetReportTimestamp(report bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.reportTimestamp = report
}

// SetReportCaller sets whether the caller location should be reported.
func (l *Logger) SetReportCaller(report bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.reportCaller = report
}

// GetLevel returns the current level.
func (l *Logger) GetLevel() Level {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return Level(l.level)
}

// SetLevel sets the current level.
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	atomic.StoreInt32(&l.level, int32(level))
}

// GetPrefix returns the current prefix.
func (l *Logger) GetPrefix() string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.prefix
}

// SetPrefix sets the current prefix.
func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

// SetTimeFormat sets the time format.
func (l *Logger) SetTimeFormat(format string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.timeFormat = format
}

// SetTimeFunction sets the time function.
func (l *Logger) SetTimeFunction(f TimeFunction) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.timeFunc = f
}

// SetOutput sets the output destination.
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if w == nil {
		w = os.Stderr
	}
	l.w = w
	var isDiscard uint32
	if w == ioutil.Discard {
		isDiscard = 1
	}
	atomic.StoreUint32(&l.isDiscard, isDiscard)
	// Reuse cached renderers
	if v, ok := registry.Load(w); ok {
		l.re = v.(*lipgloss.Renderer)
	} else {
		l.re = lipgloss.NewRenderer(w, termenv.WithColorCache(true))
		registry.Store(w, l.re)
	}
}

// SetFormatter sets the formatter.
func (l *Logger) SetFormatter(f Formatter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.formatter = f
}

// SetCallerFormatter sets the caller formatter.
func (l *Logger) SetCallerFormatter(f CallerFormatter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.callerFormatter = f
}

// With returns a new logger with the given keyvals added.
func (l *Logger) With(keyvals ...interface{}) *Logger {
	sl := *l
	sl.b = bytes.Buffer{}
	sl.mu = &sync.RWMutex{}
	sl.helpers = &sync.Map{}
	sl.fields = append(l.fields, keyvals...)
	return &sl
}

// WithPrefix returns a new logger with the given prefix.
func (l *Logger) WithPrefix(prefix string) *Logger {
	sl := l.With()
	sl.SetPrefix(prefix)
	return sl
}

// Debug prints a debug message.
func (l *Logger) Debug(msg interface{}, keyvals ...interface{}) {
	l.log(DebugLevel, msg, keyvals...)
}

// Info prints an info message.
func (l *Logger) Info(msg interface{}, keyvals ...interface{}) {
	l.log(InfoLevel, msg, keyvals...)
}

// Warn prints a warning message.
func (l *Logger) Warn(msg interface{}, keyvals ...interface{}) {
	l.log(WarnLevel, msg, keyvals...)
}

// Error prints an error message.
func (l *Logger) Error(msg interface{}, keyvals ...interface{}) {
	l.log(ErrorLevel, msg, keyvals...)
}

// Fatal prints a fatal message and exits.
func (l *Logger) Fatal(msg interface{}, keyvals ...interface{}) {
	l.log(FatalLevel, msg, keyvals...)
	os.Exit(1)
}

// Print prints a message with no level.
func (l *Logger) Print(msg interface{}, keyvals ...interface{}) {
	l.log(noLevel, msg, keyvals...)
}

// Debugf prints a debug message with formatting.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(DebugLevel, fmt.Sprintf(format, args...))
}

// Infof prints an info message with formatting.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(InfoLevel, fmt.Sprintf(format, args...))
}

// Warnf prints a warning message with formatting.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(WarnLevel, fmt.Sprintf(format, args...))
}

// Errorf prints an error message with formatting.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(ErrorLevel, fmt.Sprintf(format, args...))
}

// Fatalf prints a fatal message with formatting and exits.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log(FatalLevel, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Printf prints a message with no level and formatting.
func (l *Logger) Printf(format string, args ...interface{}) {
	l.log(noLevel, fmt.Sprintf(format, args...))
}
