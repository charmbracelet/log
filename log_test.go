package log

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubLogger(t *testing.T) {
	t.Setenv("LOG_LEVEL", "WARN")
	cases := []struct {
		name     string
		expected string
		msg      string
		fields   []interface{}
		kvs      []interface{}
		level    string
		isEnvVar bool
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
		{
			name:     "Log level from env variable",
			expected: "WARN log level from env\n",
			msg:      "log level from env",
			level:    os.Getenv("LOG_LEVEL"),
			isEnvVar: true,
		},
		{
			name:     "Log level from string",
			expected: "ERROR log level from string\n",
			msg:      "log level from string",
			level:    "ERROR",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			l := New(WithOutput(&buf))
			if c.level != "" {
				l = New(WithOutput(&buf), WithLevelFromString(c.level))
				if c.isEnvVar {
					l.With(c.fields...).Warn(c.msg, c.kvs...)
				} else {
					l.With(c.fields...).Error(c.msg, c.kvs...)
				}
				assert.Equal(t, c.expected, buf.String())
				level := l.GetLevel()
				assert.Equal(t, strings.ToLower(c.level), level.String())
			} else {
				l.With(c.fields...).Info(c.msg, c.kvs...)
				assert.Equal(t, c.expected, buf.String())
			}
		})
	}
}
