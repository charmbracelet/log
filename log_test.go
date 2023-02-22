package log

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubLogger(t *testing.T) {
	cases := []struct {
		name     string
		expected string
		msg      string
		fields   []interface{}
		kvs      []interface{}
		level    string
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
			level:    "WARN",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.level != "" {
				err := os.Setenv("LOG_LEVEL", c.level)
				assert.Nil(t, err)
			}
			var buf bytes.Buffer
			l := New(WithOutput(&buf))
			if c.level != "" {
				l.With(c.fields...).Warn(c.msg, c.kvs...)
				assert.Equal(t, c.expected, buf.String())
				level := l.GetLevel()
				assert.Equal(t, strings.ToLower(c.level), level.String())
			}
			l.With(c.fields...).Info(c.msg, c.kvs...)
			assert.Equal(t, c.expected, buf.String())
		})
		os.Unsetenv("LOG_LEVEL")
	}
}
