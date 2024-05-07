//go:build go1.21
// +build go1.21

package log

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSlogSimple(t *testing.T) {
	var buf bytes.Buffer
	h := New(&buf)
	h.SetLevel(DebugLevel)
	l := slog.New(h)
	cases := []struct {
		name     string
		expected string
		msg      string
		attrs    []any
		print    func(string, ...any)
	}{
		{
			name:     "slog debug",
			expected: "DEBU slog debug\n",
			msg:      "slog debug",
			print:    l.Debug,
			attrs:    nil,
		},
		{
			name:     "slog info",
			expected: "INFO slog info\n",
			msg:      "slog info",
			print:    l.Info,
			attrs:    nil,
		},
		{
			name:     "slog warn",
			expected: "WARN slog warn\n",
			msg:      "slog warn",
			print:    l.Warn,
			attrs:    nil,
		},
		{
			name:     "slog error",
			expected: "ERRO slog error\n",
			msg:      "slog error",
			print:    l.Error,
			attrs:    nil,
		},
		{
			name:     "slog error attrs",
			expected: "ERRO slog error foo=bar\n",
			msg:      "slog error",
			print:    l.Error,
			attrs:    []any{"foo", "bar"},
		},
	}

	for _, c := range cases {
		buf.Reset()
		t.Run(c.name, func(t *testing.T) {
			c.print(c.msg, c.attrs...)
			assert.Equal(t, c.expected, buf.String())
		})
	}
}

func TestSlogWith(t *testing.T) {
	var buf bytes.Buffer
	h := New(&buf)
	h.SetLevel(DebugLevel)
	l := slog.New(h).With("a", "b")
	cases := []struct {
		name     string
		expected string
		msg      string
		attrs    []any
		print    func(string, ...any)
	}{
		{
			name:     "slog debug",
			expected: "DEBU slog debug a=b foo=bar\n",
			msg:      "slog debug",
			print:    l.Debug,
			attrs:    []any{"foo", "bar"},
		},
		{
			name:     "slog info",
			expected: "INFO slog info a=b foo=bar\n",
			msg:      "slog info",
			print:    l.Info,
			attrs:    []any{"foo", "bar"},
		},
		{
			name:     "slog warn",
			expected: "WARN slog warn a=b foo=bar\n",
			msg:      "slog warn",
			print:    l.Warn,
			attrs:    []any{"foo", "bar"},
		},
		{
			name:     "slog error",
			expected: "ERRO slog error a=b foo=bar\n",
			msg:      "slog error",
			print:    l.Error,
			attrs:    []any{"foo", "bar"},
		},
	}

	for _, c := range cases {
		buf.Reset()
		t.Run(c.name, func(t *testing.T) {
			c.print(c.msg, c.attrs...)
			assert.Equal(t, c.expected, buf.String())
		})
	}
}

func TestSlogWithGroup(t *testing.T) {
	var buf bytes.Buffer
	h := New(&buf)
	l := slog.New(h).WithGroup("charm").WithGroup("bracelet")
	cases := []struct {
		name     string
		expected string
		msg      string
	}{
		{
			name:     "simple",
			msg:      "message",
			expected: "INFO charm.bracelet: message\n",
		},
		{
			name:     "empty",
			msg:      "",
			expected: "INFO charm.bracelet:\n",
		},
	}
	for _, c := range cases {
		buf.Reset()
		t.Run(c.name, func(t *testing.T) {
			l.Info(c.msg)
			assert.Equal(t, c.expected, buf.String())
		})
	}
}

func TestSlogCustomLevel(t *testing.T) {
	var buf bytes.Buffer
	cases := []struct {
		name     string
		expected string
		level    slog.Level
		minLevel Level
	}{
		{
			name:     "custom level not enabled",
			expected: "",
			level:    slog.Level(500),
			minLevel: Level(600),
		},
		{
			name:     "custom level enabled",
			expected: "foo\n",
			level:    slog.Level(500),
			minLevel: Level(100),
		},
	}
	for _, c := range cases {
		buf.Reset()
		t.Run(c.name, func(t *testing.T) {
			l := New(&buf)
			l.SetLevel(c.minLevel)
			l.Handle(context.Background(), slog.NewRecord(time.Now(), c.level, "foo", 0))
			assert.Equal(t, c.expected, buf.String())
		})
	}
}

type testLogValue struct {
	v slog.Value
}

func (v testLogValue) LogValue() slog.Value {
	return v.v
}

func TestSlogAttr(t *testing.T) {
	cases := []struct {
		name     string
		expected string
		kvs      []interface{}
	}{
		{
			name:     "any",
			expected: `{"level":"info","msg":"message","any":42}` + "\n",
			kvs:      []any{"any", slog.AnyValue(42)},
		},
		{
			name:     "bool",
			expected: `{"level":"info","msg":"message","bool":false}` + "\n",
			kvs:      []any{"bool", slog.BoolValue(false)},
		},
		{
			name:     "duration",
			expected: `{"level":"info","msg":"message","duration":10800000000000}` + "\n",
			kvs:      []any{"duration", slog.DurationValue(3 * time.Hour)},
		},
		{
			name:     "float64",
			expected: `{"level":"info","msg":"message","float64":123}` + "\n",
			kvs:      []any{"float64", slog.Float64Value(123)},
		},
		{
			name:     "string",
			expected: `{"level":"info","msg":"message","string":"hello"}` + "\n",
			kvs:      []any{"string", slog.StringValue("hello")},
		},
		{
			name:     "time",
			expected: `{"level":"info","msg":"message","_time":"1970-01-01T00:00:00Z"}` + "\n",
			kvs:      []any{"_time", slog.TimeValue(time.Unix(0, 0).UTC())},
		},
		{
			name:     "uint64",
			expected: `{"level":"info","msg":"message","uint64":42}` + "\n",
			kvs:      []any{"uint64", slog.Uint64Value(42)},
		},
		{
			name:     "group",
			expected: `{"level":"info","msg":"message","g":{"b":true}}` + "\n",
			kvs:      []any{slog.Group("g", slog.Bool("b", true))},
		},
		{
			name:     "log valuer",
			expected: `{"level":"info","msg":"message","lv":42}` + "\n",
			kvs: []any{
				"lv", testLogValue{slog.AnyValue(42)},
			},
		},
		{
			name:     "log valuer",
			expected: `{"level":"info","msg":"message","lv":{"first":"hello","last":"world"}}` + "\n",
			kvs: []any{
				"lv", testLogValue{slog.GroupValue(
					slog.String("first", "hello"),
					slog.String("last", "world"),
				)},
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			// expect same output from slog and log
			var buf bytes.Buffer
			l := NewWithOptions(&buf, Options{Formatter: JSONFormatter})
			l.Info("message", c.kvs...)
			assert.Equal(t, c.expected, buf.String())

			buf.Truncate(0)
			sl := slog.New(l)
			sl.Info("message", c.kvs...)
			assert.Equal(t, c.expected, buf.String())
		})
	}
}
