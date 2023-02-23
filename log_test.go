package log

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubLogger(t *testing.T) {
	var buf bytes.Buffer
	l := New(WithOutput(&buf))
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
