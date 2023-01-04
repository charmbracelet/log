package log

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdLog(t *testing.T) {
	var buf bytes.Buffer
	cases := []struct {
		name     string
		expected string
		logger   Logger
		f        func(l *log.Logger)
	}{
		{
			name:     "simple",
			expected: "INFO info\n",
			logger:   New(),
			f:        func(l *log.Logger) { l.Print("info") },
		},
		{
			name:     "without level",
			expected: "INFO coffee\n",
			logger:   New(),
			f:        func(l *log.Logger) { l.Print("coffee") },
		},
		{
			name:     "error level",
			expected: "ERROR coffee\n",
			logger:   New(),
			f:        func(l *log.Logger) { l.Print("ERROR coffee") },
		},
	}
	for _, c := range cases {
		buf.Reset()
		t.Run(c.name, func(t *testing.T) {
			c.logger.SetOutput(&buf)
			c.logger.DisableStyles()
			c.logger.SetTimeFunction(_zeroTime)
			c.f(c.logger.StandardLogger())
			assert.Equal(t, c.expected, buf.String())
		})
	}
}
