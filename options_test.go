package log

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptions(t *testing.T) {
	opts := Options{
		Level:        ErrorLevel,
		ReportCaller: true,
		Fields:       []interface{}{"foo", "bar"},
	}
	logger := NewWithOptions(io.Discard, opts)
	require.Equal(t, ErrorLevel, logger.GetLevel())
	require.True(t, logger.reportCaller)
	require.False(t, logger.reportTimestamp)
	require.Equal(t, []interface{}{"foo", "bar"}, logger.fields)
	require.Equal(t, TextFormatter, logger.formatter)
	require.Equal(t, DefaultTimeFormat, logger.timeFormat)
	require.NotNil(t, logger.timeFunc)
}

func TestCallerFormatter(t *testing.T) {
	var buf bytes.Buffer
	l := NewWithOptions(&buf, Options{ReportCaller: true})
	frames := l.frames(0)
	frame, _ := frames.Next()
	file, line, fn := frame.File, frame.Line, frame.Function
	hi := func() { l.Info("hi") }
	cases := []struct {
		name     string
		expected string
		format   CallerFormatter
	}{
		{
			name:     "short caller formatter",
			expected: fmt.Sprintf("INFO <log/options_test.go:%d> hi\n", line+3),
			format:   ShortCallerFormatter,
		},
		{
			name:     "long caller formatter",
			expected: fmt.Sprintf("INFO <%s:%d> hi\n", file, line+3),
			format:   LongCallerFormatter,
		},
		{
			name:     "foo caller formatter",
			expected: "INFO <foo> hi\n",
			format: func(file string, line int, fn string) string {
				return "foo"
			},
		},
		{
			name:     "custom caller formatter",
			expected: fmt.Sprintf("INFO <%s:%d:%s.func1> hi\n", file, line+3, fn),
			format: func(file string, line int, fn string) string {
				return fmt.Sprintf("%s:%d:%s", file, line, fn)
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf.Reset()
			l.callerFormatter = c.format
			hi()
			require.Equal(t, c.expected, buf.String())
		})
	}
}
