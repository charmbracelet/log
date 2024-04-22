//go:build go1.21
// +build go1.21

package log

import (
	"bytes"
	"context"
	"testing"
	"time"

	"log/slog"

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
