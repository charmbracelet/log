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
	"time"
)

var (
	// ErrMissingValue is returned when a key is missing a value.
	ErrMissingValue = fmt.Errorf("missing value")
)

// LoggerOption is an option for a logger.
type LoggerOption = func(*logger)

var _ Logger = &logger{}

// logger is a logger that implements Logger.
type logger struct {
	w  io.Writer
	b  bytes.Buffer
	mu *sync.RWMutex

	isDiscard uint32

	level        int32
	prefix       string
	timeFunc     TimeFunction
	timeFormat   string
	callerOffset int
	formatter    Formatter

	caller    bool
	noStyles  bool
	timestamp bool

	keyvals []interface{}

	helpers sync.Map
}

// New returns a new logger. It uses os.Stderr as the default output.
func New(opts ...LoggerOption) Logger {
	l := &logger{
		b:     bytes.Buffer{},
		mu:    &sync.RWMutex{},
		level: int32(InfoLevel),
	}

	for _, opt := range opts {
		opt(l)
	}

	if l.w == nil {
		l.w = os.Stderr
	}

	l.SetOutput(l.w)
	l.SetLevel(Level(l.level))

	if l.timeFunc == nil {
		l.timeFunc = time.Now
	}

	if l.timeFormat == "" {
		l.timeFormat = DefaultTimeFormat
	}

	return l
}

func (l *logger) log(level Level, msg interface{}, keyvals ...interface{}) {
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
	if l.timestamp {
		kvs = append(kvs, tsKey, l.timeFunc())
	}

	if level != noLevel {
		kvs = append(kvs, lvlKey, level)
	}

	if l.caller {
		// Call stack is log.Error -> log.log (2)
		file, line, _ := l.fillLoc(l.callerOffset + 2)
		caller := fmt.Sprintf("%s:%d", trimCallerPath(file), line)
		kvs = append(kvs, callerKey, caller)
	}

	if l.prefix != "" {
		kvs = append(kvs, prefixKey, l.prefix+":")
	}

	if msg != nil {
		m := fmt.Sprint(msg)
		kvs = append(kvs, msgKey, m)
	}

	// append logger fields
	kvs = append(kvs, l.keyvals...)
	if len(l.keyvals)%2 != 0 {
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
func (l *logger) Helper() {
	l.helper(1)
}

func (l *logger) helper(skip int) {
	_, _, fn := location(skip + 1)
	l.helpers.LoadOrStore(fn, struct{}{})
}

func (l *logger) fillLoc(skip int) (file string, line int, fn string) {
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

// Cleanup a path by returning the last 2 segments of the path only.
func trimCallerPath(path string) string {
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

	// Find the last separator.
	idx := strings.LastIndexByte(path, '/')
	if idx == -1 {
		return path
	}

	// Find the penultimate separator.
	idx = strings.LastIndexByte(path[:idx], '/')
	if idx == -1 {
		return path
	}

	return path[idx+1:]
}

// SetReportTimestamp sets whether the timestamp should be reported.
func (l *logger) SetReportTimestamp(report bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.timestamp = report
}

// SetReportCaller sets whether the caller location should be reported.
func (l *logger) SetReportCaller(report bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.caller = report
}

// GetLevel returns the current level.
func (l *logger) GetLevel() Level {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return Level(l.level)
}

// SetLevel sets the current level.
func (l *logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	atomic.StoreInt32(&l.level, int32(level))
}

// GetPrefix returns the current prefix.
func (l *logger) GetPrefix() string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.prefix
}

// SetPrefix sets the current prefix.
func (l *logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

// SetTimeFormat sets the time format.
func (l *logger) SetTimeFormat(format string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.timeFormat = format
}

// SetTimeFunction sets the time function.
func (l *logger) SetTimeFunction(f TimeFunction) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.timeFunc = f
}

// SetOutput sets the output destination.
func (l *logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.w = w
	var isDiscard uint32 = 0
	if w == ioutil.Discard {
		isDiscard = 1
	}
	atomic.StoreUint32(&l.isDiscard, isDiscard)
	if !isTerminal(w) {
		// This only affects the TextFormatter
		l.noStyles = true
	}
}

// SetFormatter sets the formatter.
func (l *logger) SetFormatter(f Formatter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.formatter = f
}

// With returns a new logger with the given keyvals added.
func (l *logger) With(keyvals ...interface{}) Logger {
	sl := *l
	sl.b = bytes.Buffer{}
	sl.mu = &sync.RWMutex{}
	sl.keyvals = append(l.keyvals, keyvals...)
	return &sl
}

// Debug prints a debug message.
func (l *logger) Debug(msg interface{}, keyvals ...interface{}) {
	l.log(DebugLevel, msg, keyvals...)
}

// Info prints an info message.
func (l *logger) Info(msg interface{}, keyvals ...interface{}) {
	l.log(InfoLevel, msg, keyvals...)
}

// Warn prints a warning message.
func (l *logger) Warn(msg interface{}, keyvals ...interface{}) {
	l.log(WarnLevel, msg, keyvals...)
}

// Error prints an error message.
func (l *logger) Error(msg interface{}, keyvals ...interface{}) {
	l.log(ErrorLevel, msg, keyvals...)
}

// Fatal prints a fatal message and exits.
func (l *logger) Fatal(msg interface{}, keyvals ...interface{}) {
	l.log(FatalLevel, msg, keyvals...)
	os.Exit(1)
}

// Print prints a message with no level.
func (l *logger) Print(msg interface{}, keyvals ...interface{}) {
	l.log(noLevel, msg, keyvals...)
}
