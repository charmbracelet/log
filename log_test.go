package log

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func _zeroTime() time.Time {
	return time.Time{}
}

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := New()
	logger.SetOutput(&buf)
	logger.SetTimeFunction(_zeroTime)
	logger.DisableColors()
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
	}
	for _, c := range cases {
		buf.Reset()
		t.Run(c.name, func(t *testing.T) {
			c.f(c.msg, c.kvs...)
			assert.Equal(t, c.expected, buf.String())
		})
	}
}
