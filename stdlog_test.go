package log

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStdLog(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf)
	cases := []struct {
		f        func(l *log.Logger)
		name     string
		expected string
	}{
		{
			name:     "simple",
			expected: "INFO info\n",
			f:        func(l *log.Logger) { l.Print("info") },
		},
		{
			name:     "without level",
			expected: "INFO coffee\n",
			f:        func(l *log.Logger) { l.Print("coffee") },
		},
		{
			name:     "error level",
			expected: "ERRO coffee\n",
			f:        func(l *log.Logger) { l.Print("ERROR coffee") },
		},
	}
	for _, c := range cases {
		buf.Reset()
		t.Run(c.name, func(t *testing.T) {
			l.SetOutput(&buf)
			l.SetTimeFunction(_zeroTime)
			c.f(l.StandardLog())
			assert.Equal(t, c.expected, buf.String())
		})
	}
}

func TestStdLog_forceLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf)
	cases := []struct {
		name     string
		expected string
		level    Level
	}{
		{
			name:     "debug",
			expected: "",
			level:    DebugLevel,
		},
		{
			name:     "info",
			expected: "INFO coffee\n",
			level:    InfoLevel,
		},
		{
			name:     "error",
			expected: "ERRO coffee\n",
			level:    ErrorLevel,
		},
	}
	for _, c := range cases {
		buf.Reset()
		t.Run(c.name, func(t *testing.T) {
			l := logger.StandardLog(StandardLogOptions{ForceLevel: c.level})
			l.Print("coffee")
			assert.Equal(t, c.expected, buf.String())
		})
	}
}

func TestStdLog_writer(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf)
	logger.SetReportCaller(true)
	_, file, line, ok := runtime.Caller(0)
	require.True(t, ok)
	cases := []struct {
		name     string
		expected string
		level    Level
	}{
		{
			name:     "debug",
			expected: "",
			level:    DebugLevel,
		},
		{
			name:     "info",
			expected: fmt.Sprintf("INFO <log/%s:%d> coffee\n", filepath.Base(file), line+27),
			level:    InfoLevel,
		},
		{
			name:     "error",
			expected: fmt.Sprintf("ERRO <log/%s:%d> coffee\n", filepath.Base(file), line+27),
			level:    ErrorLevel,
		},
	}
	for _, c := range cases {
		buf.Reset()
		t.Run(c.name, func(t *testing.T) {
			l := log.New(logger.StandardLog(StandardLogOptions{ForceLevel: c.level}).Writer(), "", 0)
			l.Print("coffee")
			assert.Equal(t, c.expected, buf.String())
		})
	}
}
