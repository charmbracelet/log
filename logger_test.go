package log

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubLogger(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf)
	cases := []struct {
		name     string
		expected string
		msg      string
		fields   []interface{}
		kvs      []interface{}
	}{
		{
			name:     "sub logger nil fields",
			expected: "INFO info\n",
			msg:      "info",
			fields:   nil,
			kvs:      nil,
		},
		{
			name:     "sub logger info",
			expected: "INFO info foo=bar\n",
			msg:      "info",
			fields:   []interface{}{"foo", "bar"},
			kvs:      nil,
		},
		{
			name:     "sub logger info with kvs",
			expected: "INFO info foo=bar foobar=baz\n",
			msg:      "info",
			fields:   []interface{}{"foo", "bar"},
			kvs:      []interface{}{"foobar", "baz"},
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
			expected: "INFO info\n",
			level:    Level(-999),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf.Reset()
			l := New(&buf)
			l.SetLevel(c.level)
			l.Info("info")
			assert.Equal(t, c.expected, buf.String())
		})
	}
}

func TestLogFormatter(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf)
	l.SetLevel(DebugLevel)
	cases := []struct {
		name     string
		format   string
		args     []interface{}
		fun      func(string, ...interface{})
		expected string
	}{
		{
			name:     "info format",
			format:   "%s %s",
			args:     []interface{}{"foo", "bar"},
			fun:      l.Infof,
			expected: "INFO foo bar\n",
		},
		{
			name:     "debug format",
			format:   "%s %s",
			args:     []interface{}{"foo", "bar"},
			fun:      l.Debugf,
			expected: "DEBU foo bar\n",
		},
		{
			name:     "warn format",
			format:   "%s %s",
			args:     []interface{}{"foo", "bar"},
			fun:      l.Warnf,
			expected: "WARN foo bar\n",
		},
		{
			name:     "error format",
			format:   "%s %s",
			args:     []interface{}{"foo", "bar"},
			fun:      l.Errorf,
			expected: "ERRO foo bar\n",
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
			expected: "INFO prefix: info\n",
			prefix:   "prefix",
			msg:      "info",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf.Reset()
			l := New(&buf)
			l.SetPrefix(c.prefix)
			l.Info(c.msg)
			assert.Equal(t, c.expected, buf.String())
		})
	}
}
