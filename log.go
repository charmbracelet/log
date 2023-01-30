package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

// LoggerOption is an option for a logger.
type LoggerOption = func(*logger)

var _ Logger = &logger{}

// logger is a logger that implements Logger.
type logger struct {
	w  io.Writer
	b  bytes.Buffer
	mu *sync.RWMutex

	level        Level
	prefix       string
	timeFunc     TimeFunction
	timeFormat   string
	callerOffset int

	caller    bool
	noStyles  bool
	timestamp bool

	keyvals []interface{}

	helpers sync.Map

	styles Styles
}

// New returns a new logger. It uses os.Stderr as the default output.
func New(opts ...LoggerOption) Logger {
	l := &logger{
		b:      bytes.Buffer{},
		mu:     &sync.RWMutex{},
		level:  InfoLevel,
		styles: DefaultStyles(),
	}

	for _, opt := range opts {
		opt(l)
	}

	if l.w == nil {
		l.w = os.Stderr
	}

	if l.timeFunc == nil {
		l.timeFunc = time.Now
	}

	if l.timeFormat == "" {
		l.timeFormat = DefaultTimeFormat
	}

	if !isTerminal(l.w) {
		l.noStyles = true
	}

	return l
}

func writeIndent(w io.Writer, str string, indent string, newline bool) {
	// kindly borrowed from hclog
	for {
		nl := strings.IndexByte(str, '\n')
		if nl == -1 {
			if str != "" {
				w.Write([]byte(indent))
				writeEscapedForOutput(w, str, false)
				if newline {
					w.Write([]byte{'\n'})
				}
			}
			return
		}

		w.Write([]byte(indent))
		writeEscapedForOutput(w, str[:nl], false)
		w.Write([]byte{'\n'})
		str = str[nl+1:]
	}
}

func needsEscaping(str string) bool {
	for _, b := range str {
		if !unicode.IsPrint(b) || b == '"' {
			return true
		}
	}

	return false
}

const (
	lowerhex = "0123456789abcdef"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func writeEscapedForOutput(w io.Writer, str string, escapeQuotes bool) {
	// kindly borrowed from hclog
	if !needsEscaping(str) {
		w.Write([]byte(str))
		return
	}

	bb := bufPool.Get().(*bytes.Buffer)
	bb.Reset()

	defer bufPool.Put(bb)

	for _, r := range str {
		if escapeQuotes && r == '"' {
			bb.WriteString(`\"`)
		} else if unicode.IsPrint(r) {
			bb.WriteRune(r)
		} else {
			switch r {
			case '\a':
				bb.WriteString(`\a`)
			case '\b':
				bb.WriteString(`\b`)
			case '\f':
				bb.WriteString(`\f`)
			case '\n':
				bb.WriteString(`\n`)
			case '\r':
				bb.WriteString(`\r`)
			case '\t':
				bb.WriteString(`\t`)
			case '\v':
				bb.WriteString(`\v`)
			default:
				switch {
				case r < ' ':
					bb.WriteString(`\x`)
					bb.WriteByte(lowerhex[byte(r)>>4])
					bb.WriteByte(lowerhex[byte(r)&0xF])
				case !utf8.ValidRune(r):
					r = 0xFFFD
					fallthrough
				case r < 0x10000:
					bb.WriteString(`\u`)
					for s := 12; s >= 0; s -= 4 {
						bb.WriteByte(lowerhex[r>>uint(s)&0xF])
					}
				default:
					bb.WriteString(`\U`)
					for s := 28; s >= 0; s -= 4 {
						bb.WriteByte(lowerhex[r>>uint(s)&0xF])
					}
				}
			}
		}
	}

	w.Write(bb.Bytes())
}

// isNormal indicates if the rune is one allowed to exist as an unquoted
// string value. This is a subset of ASCII, `-` through `~`.
func isNormal(r rune) bool {
	return '-' <= r && r <= '~'
}

// needsQuoting returns false if all the runes in string are normal, according
// to isNormal.
func needsQuoting(str string) bool {
	for _, r := range str {
		if !isNormal(r) {
			return true
		}
	}

	return false
}

const (
	separator       = "="
	indentSeparator = "  │ "
)

func (l *logger) log(level Level, skip int, msg interface{}, keyvals ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	defer l.b.Reset()
	t := l.timeFunc()

	// skip logging if writer is discard
	if l.w == io.Discard {
		return
	}
	// check if the level is allowed
	if l.level > level {
		return
	}

	s := l.styles
	if l.timestamp {
		ts := t.Format(l.timeFormat)
		if !l.noStyles {
			ts = s.Timestamp.Render(ts)
		}
		l.b.WriteString(ts)
		l.b.WriteByte(' ')
	}

	if level != noLevel {
		lvl := strings.ToUpper(level.String())
		if !l.noStyles {
			lvl = s.Level(level).String()
		}
		l.b.WriteString(lvl)
		l.b.WriteByte(' ')
	}

	if l.caller {
		// Call stack is log.Error -> log.log (2)
		file, line, _ := l.fillLoc(l.callerOffset + skip + 2)
		caller := fmt.Sprintf("<%s:%d>", trimCallerPath(file), line)
		if !l.noStyles {
			caller = s.Caller.Render(caller)
		}
		l.b.WriteString(caller)
		l.b.WriteByte(' ')
	}

	if l.prefix != "" {
		prefix := l.prefix + ":"
		if !l.noStyles {
			prefix = s.Prefix.Render(prefix)
		}
		l.b.WriteString(prefix)
		l.b.WriteByte(' ')
	}

	if msg != nil {
		m := fmt.Sprint(msg)
		if !l.noStyles {
			m = s.Message.Render(m)
		}
		l.b.WriteString(m)
	}

	keyvals = append(l.keyvals, keyvals...)
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "MISSING_VALUE")
	}

	sep := separator
	indentSep := indentSeparator
	if !l.noStyles {
		sep = s.Separator.Render(sep)
		indentSep = s.Separator.Render(indentSep)
	}
	for i := 0; i < len(keyvals); i += 2 {
		moreKeys := i < len(keyvals)-2
		key := fmt.Sprint(keyvals[i])
		val := fmt.Sprint(keyvals[i+1])
		raw := val == ""
		if raw {
			val = `""`
		}
		if key == "" {
			key = "MISSING_KEY"
		}
		if !l.noStyles {
			key = s.Key.Render(key)
			val = s.Value.Render(val)
		}

		// Values may contain multiple lines, and that format
		// is preserved, with each line prefixed with a "  | "
		// to show it's part of a collection of lines.
		//
		// Values may also need quoting, if not all the runes
		// in the value string are "normal", like if they
		// contain ANSI escape sequences.
		if strings.Contains(val, "\n") {
			l.b.WriteString("\n  ")
			l.b.WriteString(key)
			l.b.WriteString(sep + "\n")
			writeIndent(&l.b, val, indentSep, moreKeys)
			// If there are more keyvals, separate them with a space.
			if moreKeys {
				l.b.WriteByte(' ')
			}
		} else if !raw && needsQuoting(val) {
			l.b.WriteByte(' ')
			l.b.WriteString(key)
			l.b.WriteString(sep)
			l.b.WriteByte('"')
			writeEscapedForOutput(&l.b, val, true)
			l.b.WriteByte('"')
		} else {
			l.b.WriteByte(' ')
			l.b.WriteString(key)
			l.b.WriteString(sep)
			l.b.WriteString(val)
		}
	}

	l.b.WriteByte('\n')

	l.w.Write(l.b.Bytes())
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

// EnableTimestamp enables printing the timestamp.
func (l *logger) EnableTimestamp() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.timestamp = true
}

// DisableTimestamp disables printing the timestamp.
func (l *logger) DisableTimestamp() {
	l.timestamp = false
	l.mu.Lock()
	defer l.mu.Unlock()
}

// EnableCaller enables printing the caller.
func (l *logger) EnableCaller() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.caller = true
}

// DisableCaller disables printing the caller.
func (l *logger) DisableCaller() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.caller = false
}

// EnableStyles enables colored output.
func (l *logger) EnableStyles() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.noStyles = false
}

// DisableStyles disables colored output.
func (l *logger) DisableStyles() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.noStyles = true
}

// GetLevel returns the current level.
func (l *logger) GetLevel() Level {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.level
}

// SetLevel sets the current level.
func (l *logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
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
	l.log(DebugLevel, 0, msg, keyvals...)
}

// Info prints an info message.
func (l *logger) Info(msg interface{}, keyvals ...interface{}) {
	l.log(InfoLevel, 0, msg, keyvals...)
}

// Warn prints a warning message.
func (l *logger) Warn(msg interface{}, keyvals ...interface{}) {
	l.log(WarnLevel, 0, msg, keyvals...)
}

// Error prints an error message.
func (l *logger) Error(msg interface{}, keyvals ...interface{}) {
	l.log(ErrorLevel, 0, msg, keyvals...)
}

// Fatal prints a fatal message and exits.
func (l *logger) Fatal(msg interface{}, keyvals ...interface{}) {
	l.log(FatalLevel, 0, msg, keyvals...)
	os.Exit(1)
}

// Print prints a message with no level.
func (l *logger) Print(msg interface{}, keyvals ...interface{}) {
	l.log(noLevel, 0, msg, keyvals...)
}
