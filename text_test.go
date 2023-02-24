package log

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func _zeroTime() time.Time {
	return time.Time{}
}

func TestTextCaller(t *testing.T) {
	var buf bytes.Buffer
	logger := New(WithOutput(&buf), WithCaller())
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
	logger := New(WithOutput(&buf))
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
			expected: "ERROR info key1=val1 key2=val2\n",
			msg:      "info",
			kvs:      []interface{}{"key1", "val1", "key2", "val2"},
			f:        logger.Error,
		},
		{
			name:     "error message with multiline",
			expected: "ERROR info\n  key1=\n  │ val1\n  │ val2\n",
			msg:      "info",
			kvs:      []interface{}{"key1", "val1\nval2"},
			f:        logger.Error,
		},
		{
			name:     "odd number of keyvals",
			expected: "ERROR info key1=val1 key2=val2 key3=\"missing value\"\n",
			msg:      "info",
			kvs:      []interface{}{"key1", "val1", "key2", "val2", "key3"},
			f:        logger.Error,
		},
		{
			name:     "error field",
			expected: "ERROR info key1=\"error value\"\n",
			msg:      "info",
			kvs:      []interface{}{"key1", errors.New("error value")},
			f:        logger.Error,
		},
		{
			name:     "struct field",
			expected: "ERROR info key1={foo:bar}\n",
			msg:      "info",
			kvs:      []interface{}{"key1", struct{ foo string }{foo: "bar"}},
			f:        logger.Error,
		},
		{
			name:     "struct field quoted",
			expected: "ERROR info key1=\"{foo:bar baz}\"\n",
			msg:      "info",
			kvs:      []interface{}{"key1", struct{ foo string }{foo: "bar baz"}},
			f:        logger.Error,
		},
		{
			name:     "slice of strings",
			expected: "ERROR info key1=\"[foo bar]\"\n",
			msg:      "info",
			kvs:      []interface{}{"key1", []string{"foo", "bar"}},
			f:        logger.Error,
		},
		{
			name:     "slice of structs",
			expected: "ERROR info key1=\"[{foo:bar} {foo:baz}]\"\n",
			msg:      "info",
			kvs:      []interface{}{"key1", []struct{ foo string }{{foo: "bar"}, {foo: "baz"}}},
			f:        logger.Error,
		},
		{
			name:     "slice of errors",
			expected: "ERROR info key1=\"[error value1 error value2]\"\n",
			msg:      "info",
			kvs:      []interface{}{"key1", []error{errors.New("error value1"), errors.New("error value2")}},
			f:        logger.Error,
		},
		{
			name:     "map of strings",
			expected: "ERROR info key1=\"map[baz:qux foo:bar]\"\n",
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
	logger := New(WithOutput(&buf), WithCaller())

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
	logger := New(WithOutput(&buf), WithCaller())
	if os.Getenv("FATAL") == "1" {
		logger.Fatal("i'm dead")
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
	logger := New(WithOutput(&buf)).(*logger)
	logger.noStyles = false
	ValueStyle = lipgloss.NewStyle().Bold(true)
	cases := []struct {
		name     string
		expected string
		msg      string
		kvs      []interface{}
		f        func(msg interface{}, kvs ...interface{})
	}{
		{
			name:     "simple message",
			expected: fmt.Sprintf("%s info\n", InfoLevelStyle.Render("INFO")),
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
				InfoLevelStyle.Render("INFO"),
				KeyStyle.Render("key1"), SeparatorStyle.Render("="), ValueStyle.Render("val1"),
				KeyStyle.Render("key2"), SeparatorStyle.Render("="), ValueStyle.Render("val2"),
			),
			msg: "info",
			kvs: []interface{}{"key1", "val1", "key2", "val2"},
			f:   logger.Info,
		},
		{
			name: "error message with multiline",
			expected: fmt.Sprintf(
				"%s info\n  %s%s\n%s%s\n%s%s\n",
				ErrorLevelStyle.Render("ERRO"),
				KeyStyle.Render("key1"), SeparatorStyle.Render("="),
				SeparatorStyle.Render("  │ "), ValueStyle.Render("val1"),
				SeparatorStyle.Render("  │ "), ValueStyle.Render("val2"),
			),
			msg: "info",
			kvs: []interface{}{"key1", "val1\nval2"},
			f:   logger.Error,
		},
		{
			name: "error message with keyvals",
			expected: fmt.Sprintf(
				"%s info %s%s%s %s%s%s\n",
				ErrorLevelStyle.Render("ERRO"),
				KeyStyle.Render("key1"), SeparatorStyle.Render("="), ValueStyle.Render("val1"),
				KeyStyle.Render("key2"), SeparatorStyle.Render("="), ValueStyle.Render("val2"),
			),
			msg: "info",
			kvs: []interface{}{"key1", "val1", "key2", "val2"},
			f:   logger.Error,
		},
		{
			name: "odd number of keyvals",
			expected: fmt.Sprintf(
				"%s info %s%s%s %s%s%s %s%s\"%s\"\n",
				ErrorLevelStyle.Render("ERRO"),
				KeyStyle.Render("key1"), SeparatorStyle.Render("="), ValueStyle.Render("val1"),
				KeyStyle.Render("key2"), SeparatorStyle.Render("="), ValueStyle.Render("val2"),
				KeyStyle.Render("key3"), SeparatorStyle.Render("="), ValueStyle.Render("missing value"),
			),
			msg: "info",
			kvs: []interface{}{"key1", "val1", "key2", "val2", "key3"},
			f:   logger.Error,
		},
		{
			name: "error field",
			expected: fmt.Sprintf(
				"%s info %s%s\"%s\"\n",
				ErrorLevelStyle.Render("ERRO"),
				KeyStyle.Render("key1"), SeparatorStyle.Render("="), ValueStyle.Render("error value"),
			),
			msg: "info",
			kvs: []interface{}{"key1", errors.New("error value")},
			f:   logger.Error,
		},
		{
			name: "struct field",
			expected: fmt.Sprintf(
				"%s info %s%s%s\n",
				InfoLevelStyle.Render("INFO"),
				KeyStyle.Render("key1"), SeparatorStyle.Render("="), ValueStyle.Render("{foo:bar}"),
			),
			msg: "info",
			kvs: []interface{}{"key1", struct{ foo string }{foo: "bar"}},
			f:   logger.Info,
		},
		{
			name: "struct field quoted",
			expected: fmt.Sprintf(
				"%s info %s%s\"%s\"\n",
				InfoLevelStyle.Render("INFO"),
				KeyStyle.Render("key1"), SeparatorStyle.Render("="), ValueStyle.Render("{foo:bar baz}"),
			),
			msg: "info",
			kvs: []interface{}{"key1", struct{ foo string }{foo: "bar baz"}},
			f:   logger.Info,
		},
		{
			name: "slice of strings",
			expected: fmt.Sprintf(
				"%s info %s%s\"%s\"\n",
				ErrorLevelStyle.Render("ERRO"),
				KeyStyle.Render("key1"), SeparatorStyle.Render("="), ValueStyle.Render("[foo bar]"),
			),
			msg: "info",
			kvs: []interface{}{"key1", []string{"foo", "bar"}},
			f:   logger.Error,
		},
		{
			name: "slice of structs",
			expected: fmt.Sprintf(
				"%s info %s%s\"%s\"\n",
				ErrorLevelStyle.Render("ERRO"),
				KeyStyle.Render("key1"), SeparatorStyle.Render("="), ValueStyle.Render("[{foo:bar} {foo:baz}]"),
			),
			msg: "info",
			kvs: []interface{}{"key1", []struct{ foo string }{{foo: "bar"}, {foo: "baz"}}},
			f:   logger.Error,
		},
		{
			name: "slice of errors",
			expected: fmt.Sprintf(
				"%s info %s%s\"%s\"\n",
				ErrorLevelStyle.Render("ERRO"),
				KeyStyle.Render("key1"), SeparatorStyle.Render("="), ValueStyle.Render("[error value1 error value2]"),
			),
			msg: "info",
			kvs: []interface{}{"key1", []error{errors.New("error value1"), errors.New("error value2")}},
			f:   logger.Error,
		},
		{
			name: "map of strings",
			expected: fmt.Sprintf(
				"%s info %s%s\"%s\"\n",
				ErrorLevelStyle.Render("ERRO"),
				KeyStyle.Render("key1"), SeparatorStyle.Render("="), ValueStyle.Render("map[baz:qux foo:bar]"),
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
