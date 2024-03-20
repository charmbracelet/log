package log

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func _zeroTime(time.Time) time.Time {
	return time.Time{}.AddDate(1, 0, 0)
}

func TestNilStyles(t *testing.T) {
	st := DefaultStyles()
	l := New(io.Discard)
	l.SetStyles(nil)
	assert.Equal(t, st, l.styles)
}

func TestTextCaller(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf)
	logger.SetReportCaller(true)
	// We calculate the caller offset based on the caller line number.
	_, file, line, _ := runtime.Caller(0)
	cases := []struct {
		name     string
		expected string
		msg      string
		kvs      []interface{}
		f        func(msg interface{}, kvs ...interface{})
	}{
		{
			name:     "simple caller",
			expected: fmt.Sprintf("INFO <log/%s:%d> info\n", filepath.Base(file), line+14),
			msg:      "info",
			kvs:      nil,
			f: func(msg interface{}, kvs ...interface{}) {
				logger.Info(msg, kvs...)
			},
		},
		{
			name:     "helper caller",
			expected: fmt.Sprintf("INFO <log/%s:%d> info\n", filepath.Base(file), line+58),
			msg:      "info",
			kvs:      nil,
			f: func(msg interface{}, kvs ...interface{}) {
				logger.Helper()
				logger.Info(msg, kvs...)
			},
		},
		{
			name:     "nested helper caller",
			expected: fmt.Sprintf("INFO <log/%s:%d> info\n", filepath.Base(file), line+37),
			msg:      "info",
			kvs:      nil,
			f: func(msg interface{}, kvs ...interface{}) {
				fun := func(msg interface{}, kvs ...interface{}) {
					logger.Helper()
					logger.Info(msg, kvs...)
				}
				fun(msg, kvs...)
			},
		},
		{
			name:     "double nested helper caller",
			expected: fmt.Sprintf("INFO <log/%s:%d> info\n", filepath.Base(file), line+58),
			msg:      "info",
			kvs:      nil,
			f: func(msg interface{}, kvs ...interface{}) {
				logger.Helper()
				fun := func(msg interface{}, kvs ...interface{}) {
					logger.Helper()
					logger.Info(msg, kvs...)
				}
				fun(msg, kvs...)
			},
		},
	}
	for _, c := range cases {
		buf.Reset()
		t.Run(c.name, func(t *testing.T) {
			c.f(c.msg, c.kvs...)
			assert.Equal(t, c.expected, buf.String())
		})
	}
}

func TestTextLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf)
	cases := []struct {
		name     string
		expected string
		msg      string
		kvs      []interface{}
		f        func(msg interface{}, kvs ...interface{})
	}{
		{
			name:     "simple message",
			expected: "INFO info\n",
			msg:      "info",
			kvs:      nil,
			f:        logger.Info,
		},
		{
			name:     "ignored message",
			expected: "",
			msg:      "this is a debug message",
			kvs:      nil,
			f:        logger.Debug,
		},
		{
			name:     "message with keyvals",
			expected: "INFO info key1=val1 key2=val2\n",
			msg:      "info",
			kvs:      []interface{}{"key1", "val1", "key2", "val2"},
			f:        logger.Info,
		},
		{
			name:     "error message with keyvals",
			expected: "ERRO info key1=val1 key2=val2\n",
			msg:      "info",
			kvs:      []interface{}{"key1", "val1", "key2", "val2"},
			f:        logger.Error,
		},
		{
			name:     "error message with multiline",
			expected: "ERRO info\n  key1=\n  │ val1\n  │ val2\n",
			msg:      "info",
			kvs:      []interface{}{"key1", "val1\nval2"},
			f:        logger.Error,
		},
		{
			name:     "odd number of keyvals",
			expected: "ERRO info key1=val1 key2=val2 key3=\"missing value\"\n",
			msg:      "info",
			kvs:      []interface{}{"key1", "val1", "key2", "val2", "key3"},
			f:        logger.Error,
		},
		{
			name:     "error field",
			expected: "ERRO info key1=\"error value\"\n",
			msg:      "info",
			kvs:      []interface{}{"key1", errors.New("error value")},
			f:        logger.Error,
		},
		{
			name:     "struct field",
			expected: "ERRO info key1={foo:bar}\n",
			msg:      "info",
			kvs:      []interface{}{"key1", struct{ foo string }{foo: "bar"}},
			f:        logger.Error,
		},
		{
			name:     "struct field quoted",
			expected: "ERRO info key1=\"{foo:bar baz}\"\n",
			msg:      "info",
			kvs:      []interface{}{"key1", struct{ foo string }{foo: "bar baz"}},
			f:        logger.Error,
		},
		{
			name:     "slice of strings",
			expected: "ERRO info key1=\"[foo bar]\"\n",
			msg:      "info",
			kvs:      []interface{}{"key1", []string{"foo", "bar"}},
			f:        logger.Error,
		},
		{
			name:     "slice of structs",
			expected: "ERRO info key1=\"[{foo:bar} {foo:baz}]\"\n",
			msg:      "info",
			kvs:      []interface{}{"key1", []struct{ foo string }{{foo: "bar"}, {foo: "baz"}}},
			f:        logger.Error,
		},
		{
			name:     "slice of errors",
			expected: "ERRO info key1=\"[error value1 error value2]\"\n",
			msg:      "info",
			kvs:      []interface{}{"key1", []error{errors.New("error value1"), errors.New("error value2")}},
			f:        logger.Error,
		},
		{
			name:     "map of strings",
			expected: "ERRO info key1=\"map[baz:qux foo:bar]\"\n",
			msg:      "info",
			kvs:      []interface{}{"key1", map[string]string{"foo": "bar", "baz": "qux"}},
			f:        logger.Error,
		},
	}
	for _, c := range cases {
		buf.Reset()
		t.Run(c.name, func(t *testing.T) {
			c.f(c.msg, c.kvs...)
			assert.Equal(t, c.expected, buf.String())
		})
	}
}

func TestTextHelper(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf)
	logger.SetReportCaller(true)
	helper := func() {
		logger.Helper()
		logger.Info("helper func")
	}

	helper()
	_, file, line, ok := runtime.Caller(0)
	require.True(t, ok)
	assert.Equal(t, fmt.Sprintf("INFO <log/%s:%d> helper func\n", filepath.Base(file), line-1), buf.String())
}

func TestTextFatal(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf)
	logger.SetReportCaller(true)
	if os.Getenv("FATAL") == "1" {
		logger.Fatal("i'm dead")
		logger.Fatalf("bye %s", "bye")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestTextFatal")
	cmd.Env = append(os.Environ(), "FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestTextValueStyles(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf)
	logger.SetColorProfile(termenv.ANSI256)
	lipgloss.SetColorProfile(termenv.ANSI256)
	st := DefaultStyles()
	st.Value = lipgloss.NewStyle().Bold(true)
	st.Values["key3"] = st.Value.Copy().Underline(true)
	logger.SetStyles(st)
	cases := []struct {
		name     string
		expected string
		msg      string
		kvs      []interface{}
		f        func(msg interface{}, kvs ...interface{})
	}{
		{
			name:     "simple message",
			expected: fmt.Sprintf("%s info\n", st.Levels[InfoLevel]),
			msg:      "info",
			kvs:      nil,
			f:        logger.Info,
		},
		{
			name:     "ignored message",
			expected: "",
			msg:      "this is a debug message",
			kvs:      nil,
			f:        logger.Debug,
		},
		{
			name: "message with keyvals",
			expected: fmt.Sprintf(
				"%s info %s%s%s %s%s%s\n",
				st.Levels[InfoLevel],
				st.Key.Render("key1"), st.Separator.Render(separator), st.Value.Render("val1"),
				st.Key.Render("key2"), st.Separator.Render(separator), st.Value.Render("val2"),
			),
			msg: "info",
			kvs: []interface{}{"key1", "val1", "key2", "val2"},
			f:   logger.Info,
		},
		{
			name: "error message with multiline",
			expected: fmt.Sprintf(
				"%s info\n  %s%s\n%s%s\n%s%s\n",
				st.Levels[ErrorLevel],
				st.Key.Render("key1"), st.Separator.Render(separator),
				st.Separator.Render(indentSeparator), st.Value.Render("val1"),
				st.Separator.Render(indentSeparator), st.Value.Render("val2"),
			),
			msg: "info",
			kvs: []interface{}{"key1", "val1\nval2"},
			f:   logger.Error,
		},
		{
			name: "error message with keyvals",
			expected: fmt.Sprintf(
				"%s info %s%s%s %s%s%s\n",
				st.Levels[ErrorLevel],
				st.Key.Render("key1"), st.Separator.Render(separator), st.Value.Render("val1"),
				st.Key.Render("key2"), st.Separator.Render(separator), st.Value.Render("val2"),
			),
			msg: "info",
			kvs: []interface{}{"key1", "val1", "key2", "val2"},
			f:   logger.Error,
		},
		{
			name: "odd number of keyvals",
			expected: fmt.Sprintf(
				"%s info %s%s%s %s%s%s %s%s%s\n",
				st.Levels[ErrorLevel],
				st.Key.Render("key1"), st.Separator.Render(separator), st.Value.Render("val1"),
				st.Key.Render("key2"), st.Separator.Render(separator), st.Value.Render("val2"),
				st.Key.Render("key3"), st.Separator.Render(separator), st.Values["key3"].Render(`"missing value"`),
			),
			msg: "info",
			kvs: []interface{}{"key1", "val1", "key2", "val2", "key3"},
			f:   logger.Error,
		},
		{
			name: "error field",
			expected: fmt.Sprintf(
				"%s info %s%s%s\n",
				st.Levels[ErrorLevel],
				st.Key.Render("key1"), st.Separator.Render(separator), st.Value.Render(`"error value"`),
			),
			msg: "info",
			kvs: []interface{}{"key1", errors.New("error value")},
			f:   logger.Error,
		},
		{
			name: "struct field",
			expected: fmt.Sprintf(
				"%s info %s%s%s\n",
				st.Levels[InfoLevel],
				st.Key.Render("key1"), st.Separator.Render(separator), st.Value.Render("{foo:bar}"),
			),
			msg: "info",
			kvs: []interface{}{"key1", struct{ foo string }{foo: "bar"}},
			f:   logger.Info,
		},
		{
			name: "struct field quoted",
			expected: fmt.Sprintf(
				"%s info %s%s%s\n",
				st.Levels[InfoLevel],
				st.Key.Render("key1"), st.Separator.Render(separator), st.Value.Render(`"{foo:bar baz}"`),
			),
			msg: "info",
			kvs: []interface{}{"key1", struct{ foo string }{foo: "bar baz"}},
			f:   logger.Info,
		},
		{
			name: "slice of strings",
			expected: fmt.Sprintf(
				"%s info %s%s%s\n",
				st.Levels[ErrorLevel],
				st.Key.Render("key1"), st.Separator.Render(separator), st.Value.Render(`"[foo bar]"`),
			),
			msg: "info",
			kvs: []interface{}{"key1", []string{"foo", "bar"}},
			f:   logger.Error,
		},
		{
			name: "slice of structs",
			expected: fmt.Sprintf(
				"%s info %s%s%s\n",
				st.Levels[ErrorLevel],
				st.Key.Render("key1"), st.Separator.Render(separator), st.Value.Render(`"[{foo:bar} {foo:baz}]"`),
			),
			msg: "info",
			kvs: []interface{}{"key1", []struct{ foo string }{{foo: "bar"}, {foo: "baz"}}},
			f:   logger.Error,
		},
		{
			name: "slice of errors",
			expected: fmt.Sprintf(
				"%s info %s%s%s\n",
				st.Levels[ErrorLevel],
				st.Key.Render("key1"), st.Separator.Render(separator), st.Value.Render(`"[error value1 error value2]"`),
			),
			msg: "info",
			kvs: []interface{}{"key1", []error{errors.New("error value1"), errors.New("error value2")}},
			f:   logger.Error,
		},
		{
			name: "map of strings",
			expected: fmt.Sprintf(
				"%s info %s%s%s\n",
				st.Levels[ErrorLevel],
				st.Key.Render("key1"), st.Separator.Render(separator), st.Value.Render(`"map[baz:qux foo:bar]"`),
			),
			msg: "info",
			kvs: []interface{}{"key1", map[string]string{"foo": "bar", "baz": "qux"}},
			f:   logger.Error,
		},
	}
	for _, c := range cases {
		buf.Reset()
		t.Run(c.name, func(t *testing.T) {
			c.f(c.msg, c.kvs...)
			assert.Equal(t, c.expected, buf.String())
		})
	}
}

func TestColorProfile(t *testing.T) {
	cases := []termenv.Profile{
		termenv.Ascii,
		termenv.ANSI,
		termenv.ANSI256,
		termenv.TrueColor,
	}
	l := New(io.Discard)
	for _, p := range cases {
		l.SetColorProfile(p)
		assert.Equal(t, p, l.re.ColorProfile())
	}
}

func TestCustomLevelStyle(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf)
	st := DefaultStyles()
	lvl := Level(1234)
	st.Levels[lvl] = lipgloss.NewStyle().Bold(true).SetString("FUNKY")
	l.SetStyles(st)
	l.SetLevel(lvl)
	l.Log(lvl, "foobar")
	assert.Equal(t, "FUNKY foobar\n", buf.String())
}
