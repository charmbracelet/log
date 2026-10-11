package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/charmbracelet/colorprofile"
	"github.com/stretchr/testify/assert"
)

func TestSubLogger(t *testing.T) {
	var buf bytes.Buffer
	l := New(colorprofile.NewWriter(&buf, os.Environ()))
	cases := []struct {
		name     string
		expected string
		msg      string
		fields   []any
		kvs      []any
	}{
		{
			name:     "sub logger nil fields",
			expected: "INFO  info\n",
			msg:      "info",
			fields:   nil,
			kvs:      nil,
		},
		{
			name:     "sub logger info",
			expected: "INFO  info foo=bar\n",
			msg:      "info",
			fields:   []any{"foo", "bar"},
			kvs:      nil,
		},
		{
			name:     "sub logger info with kvs",
			expected: "INFO  info foo=bar foobar=baz\n",
			msg:      "info",
			fields:   []any{"foo", "bar"},
			kvs:      []any{"foobar", "baz"},
		},
		{
			name:     "emoji",
			expected: "INFO  👍 🐱\n",
			msg:      "👍 🐱",
			fields:   nil,
			kvs:      nil,
		},
	}
	for _, c := range cases {
		buf.Reset()
		t.Run(c.name, func(t *testing.T) {
			l.With(c.fields...).Info(c.msg, c.kvs...)
			assert.Equal(t, c.expected, buf.String())
		})
	}
}

func TestWrongLevel(t *testing.T) {
	var buf bytes.Buffer
	cases := []struct {
		name     string
		expected string
		level    Level
	}{
		{
			name:     "wrong level",
			expected: "",
			level:    Level(999),
		},
		{
			name:     "wrong level negative",
			expected: "INFO  info\n",
			level:    Level(-999),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf.Reset()
			l := New(colorprofile.NewWriter(&buf, os.Environ()))
			l.SetLevel(c.level)
			l.Info("info")
			assert.Equal(t, c.expected, buf.String())
		})
	}
}

func TestLogFormatter(t *testing.T) {
	var buf bytes.Buffer
	l := New(colorprofile.NewWriter(&buf, os.Environ()))
	l.SetLevel(DebugLevel)
	cases := []struct {
		name     string
		format   string
		args     []any
		fun      func(string, ...any)
		expected string
	}{
		{
			name:     "info format",
			format:   "%s %s",
			args:     []any{"foo", "bar"},
			fun:      l.Infof,
			expected: "INFO  foo bar\n",
		},
		{
			name:     "debug format",
			format:   "%s %s",
			args:     []any{"foo", "bar"},
			fun:      l.Debugf,
			expected: "DEBUG foo bar\n",
		},
		{
			name:     "warn format",
			format:   "%s %s",
			args:     []any{"foo", "bar"},
			fun:      l.Warnf,
			expected: "WARN  foo bar\n",
		},
		{
			name:     "error format",
			format:   "%s %s",
			args:     []any{"foo", "bar"},
			fun:      l.Errorf,
			expected: "ERROR foo bar\n",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf.Reset()
			c.fun(c.format, "foo", "bar")
			assert.Equal(t, c.expected, buf.String())
		})
	}
}

func TestEmptyMessage(t *testing.T) {
	var buf bytes.Buffer
	l := New(colorprofile.NewWriter(&buf, os.Environ()))
	cases := []struct {
		name     string
		expected string
		msg      string
		fields   []any
		kvs      []any
	}{
		{
			name:     "empty message nil fields",
			expected: "INFO \n",
			msg:      "",
			fields:   nil,
			kvs:      nil,
		},
		{
			name:     "empty message with fields",
			expected: "INFO  foo=bar\n",
			msg:      "",
			fields:   []any{"foo", "bar"},
			kvs:      nil,
		},
		{
			name:     "empty message with fields & kvs",
			expected: "INFO  foo=bar foobar=baz\n",
			msg:      "",
			fields:   []any{"foo", "bar"},
			kvs:      []any{"foobar", "baz"},
		},
	}
	for _, c := range cases {
		buf.Reset()
		t.Run(c.name, func(t *testing.T) {
			l.With(c.fields...).Info(c.msg, c.kvs...)
			assert.Equal(t, c.expected, buf.String())
		})
	}
}

func TestLogWithPrefix(t *testing.T) {
	var buf bytes.Buffer
	cases := []struct {
		name     string
		expected string
		prefix   string
		msg      string
	}{
		{
			name:     "with prefix",
			expected: "INFO  prefix: info\n",
			prefix:   "prefix",
			msg:      "info",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf.Reset()
			l := New(colorprofile.NewWriter(&buf, os.Environ()))
			l.SetPrefix(c.prefix)
			l.Info(c.msg)
			assert.Equal(t, c.expected, buf.String())
		})
	}
}

func TestLogWithRaceCondition(t *testing.T) {
	w := io.Discard
	cases := []struct {
		name string
	}{
		{
			name: "must be run with -race",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := New(w)

			var done sync.WaitGroup

			longArgs := make([]any, 0, 1000)
			for i := 0; i < 1000; i++ {
				longArgs = append(longArgs, fmt.Sprintf("arg%d", i), fmt.Sprintf("val%d", i))
			}
			l = l.With(longArgs...)

			for i := 0; i < 100; i++ {
				done.Add(1)
				go func() {
					ll := l.With("arg1", "val1", "arg2", "val2")
					ll.Info("kinda long long log message")
					done.Done()
				}()

				done.Add(1)
				go func() {
					l.Info("kinda long log message")
					done.Done()
				}()
			}
			done.Wait()
		})
	}
}

func TestRace(t *testing.T) {
	t.Parallel()

	w := io.Discard
	l := New(w)
	for i := 0; i < 100; i++ {
		t.Run("race", func(t *testing.T) {
			t.Parallel()
			s := l.StandardLog()
			l.Info("foo")
			l.GetLevel()
			l.Print("foo")

			s.Print("foo")
			s.Writer().Write([]byte("bar"))
			s.Output(1, "baz")

			l.SetOutput(w)
			l.Debug("foo")
			l.SetLevel(InfoLevel)
			l.GetPrefix()

			o := l.With("foo", "bar")
			o.Printf("foo %s", "bar")
			o.SetTimeFormat(time.Kitchen)
			o.Warn("foo")
			o.SetOutput(w)
			o.Error("foo")
			o.SetFormatter(JSONFormatter)
		})
	}
}

func TestCustomLevel(t *testing.T) {
	var buf bytes.Buffer
	level500 := Level(500)
	l := New(&buf)
	l.SetLevel(level500)
	l.Logf(level500, "foo")
	assert.Equal(t, "foo\n", buf.String())
}

func TestWrappedLoggerInheritsLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithOptions(&buf, Options{
		Level:           InfoLevel,
		ReportTimestamp: false,
		ReportCaller:    false,
	})

	wrappedLogger := logger.With("component", "auth")

	// Initially InfoLevel: debug message must be suppressed on child.
	wrappedLogger.Debug("auth debug suppressed")
	assert.Empty(t, buf.String())

	// Changing parent level to DebugLevel must propagate to child.
	logger.SetLevel(DebugLevel)
	assert.Equal(t, DebugLevel, wrappedLogger.GetLevel())

	wrappedLogger.Debug("auth debug allowed")
	assert.Equal(t, "DEBUG auth debug allowed component=auth\n", buf.String())

	// Changing parent level to WarnLevel must suppress info messages on child.
	buf.Reset()
	logger.SetLevel(WarnLevel)
	assert.Equal(t, WarnLevel, wrappedLogger.GetLevel())

	wrappedLogger.Info("auth info suppressed")
	assert.Empty(t, buf.String())

	wrappedLogger.Warn("auth warning allowed")
	assert.Equal(t, "WARN  auth warning allowed component=auth\n", buf.String())
}

func TestWrappedLoggerExplicitChildLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithOptions(&buf, Options{
		Level:           InfoLevel,
		ReportTimestamp: false,
		ReportCaller:    false,
	})

	child := logger.With("component", "worker")
	child.SetLevel(WarnLevel)

	// Child override must not mutate parent level.
	assert.Equal(t, InfoLevel, logger.GetLevel())
	assert.Equal(t, WarnLevel, child.GetLevel())

	// Subsequent parent level change must not override explicit child setting.
	logger.SetLevel(DebugLevel)
	assert.Equal(t, DebugLevel, logger.GetLevel())
	assert.Equal(t, WarnLevel, child.GetLevel())

	// Parent permits debug messages while child suppresses them.
	logger.Debug("parent debug allowed")
	assert.Equal(t, "DEBUG parent debug allowed\n", buf.String())

	buf.Reset()
	child.Debug("child debug suppressed")
	assert.Empty(t, buf.String())

	child.Info("child info suppressed")
	assert.Empty(t, buf.String())

	child.Warn("child warn allowed")
	assert.Equal(t, "WARN  child warn allowed component=worker\n", buf.String())
}

func TestWrappedLoggerHierarchy(t *testing.T) {
	var buf bytes.Buffer
	root := NewWithOptions(&buf, Options{
		Level:           InfoLevel,
		ReportTimestamp: false,
		ReportCaller:    false,
	})

	child := root.With("tier", 1)
	grandchild := child.With("tier", 2)

	assert.Equal(t, InfoLevel, root.GetLevel())
	assert.Equal(t, InfoLevel, child.GetLevel())
	assert.Equal(t, InfoLevel, grandchild.GetLevel())

	// Root update propagates down entire hierarchy.
	root.SetLevel(DebugLevel)
	assert.Equal(t, DebugLevel, root.GetLevel())
	assert.Equal(t, DebugLevel, child.GetLevel())
	assert.Equal(t, DebugLevel, grandchild.GetLevel())

	grandchild.Debug("grandchild debug")
	assert.Equal(t, "DEBUG grandchild debug tier=1 tier=2\n", buf.String())

	// Overriding child propagates to grandchild, while root is unaffected.
	buf.Reset()
	child.SetLevel(ErrorLevel)
	assert.Equal(t, DebugLevel, root.GetLevel())
	assert.Equal(t, ErrorLevel, child.GetLevel())
	assert.Equal(t, ErrorLevel, grandchild.GetLevel())

	grandchild.Warn("grandchild warn suppressed")
	assert.Empty(t, buf.String())

	grandchild.Error("grandchild error allowed")
	assert.Equal(t, "ERROR grandchild error allowed tier=1 tier=2\n", buf.String())

	buf.Reset()
	root.Debug("root debug allowed")
	assert.Equal(t, "DEBUG root debug allowed\n", buf.String())
}

func TestWrappedLoggerRace(t *testing.T) {
	w := io.Discard
	logger := New(w)
	child := logger.With("component", "worker")
	grandchild := child.With("task", 42)

	var wg sync.WaitGroup
	levels := []Level{DebugLevel, InfoLevel, WarnLevel, ErrorLevel}

	for i := 0; i < 20; i++ {
		wg.Add(4)
		lvl := levels[i%len(levels)]

		go func() {
			defer wg.Done()
			logger.SetLevel(lvl)
			logger.Debug("logger debug")
			logger.Info("logger info")
		}()

		go func() {
			defer wg.Done()
			child.GetLevel()
			child.Info("child info")
			child.Warn("child warn")
		}()

		go func() {
			defer wg.Done()
			grandchild.GetLevel()
			grandchild.Debug("grandchild debug")
			grandchild.Error("grandchild error")
		}()

		go func(idx int) {
			defer wg.Done()
			if idx%2 == 0 {
				child.SetLevel(lvl)
			}
		}(i)
	}

	wg.Wait()
}

