package log

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogfmt(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf)
	l.SetFormatter(LogfmtFormatter)
	cases := []struct {
		name     string
		expected string
		msg      string
		kvs      []interface{}
		f        func(msg interface{}, kvs ...interface{})
	}{
		{
			name:     "simple",
			expected: "lvl=info msg=info\n",
			msg:      "info",
			kvs:      nil,
			f:        l.Info,
		},
		{
			name:     "ignored message",
			expected: "",
			msg:      "info",
			kvs:      nil,
			f:        l.Debug,
		},
		{
			name:     "message with keyvals",
			expected: "lvl=info msg=info foo=bar\n",
			msg:      "info",
			kvs:      []interface{}{"foo", "bar"},
			f:        l.Info,
		},
		{
			name:     "message with multiline keyvals",
			expected: "lvl=info msg=info foo=\"bar\\nbaz\"\n",
			msg:      "info",
			kvs:      []interface{}{"foo", "bar\nbaz"},
			f:        l.Info,
		},
		{
			name:     "multiline message",
			expected: "lvl=info msg=\"info\\nfoo\"\n",
			msg:      "info\nfoo",
			kvs:      nil,
			f:        l.Info,
		},
		{
			name:     "message with error",
			expected: "lvl=info msg=info err=\"foo: bar\"\n",
			msg:      "info",
			kvs:      []interface{}{"err", errors.New("foo: bar")},
			f:        l.Info,
		},
		{
			name:     "odd number of keyvals",
			expected: "lvl=info msg=info foo=bar baz=\"missing value\"\n",
			msg:      "info",
			kvs:      []interface{}{"foo", "bar", "baz"},
			f:        l.Info,
		},
		{
			name:     "struct field",
			expected: "lvl=info msg=info foo=\"{bar:foo bar}\"\n",
			msg:      "info",
			kvs:      []interface{}{"foo", struct{ bar string }{"foo bar"}},
			f:        l.Info,
		},
		{
			name:     "multiple struct fields",
			expected: "lvl=info msg=info foo={bar:baz} foo2={bar:baz}\n",
			msg:      "info",
			kvs:      []interface{}{"foo", struct{ bar string }{"baz"}, "foo2", struct{ bar string }{"baz"}},
			f:        l.Info,
		},
		{
			name:     "slice of structs",
			expected: "lvl=info msg=info foo=\"[{bar:baz} {bar:baz}]\"\n",
			msg:      "info",
			kvs:      []interface{}{"foo", []struct{ bar string }{{"baz"}, {"baz"}}},
			f:        l.Info,
		},
		{
			name:     "slice of strings",
			expected: "lvl=info msg=info foo=\"[bar baz]\"\n",
			msg:      "info",
			kvs:      []interface{}{"foo", []string{"bar", "baz"}},
			f:        l.Info,
		},
		{
			name:     "slice of errors",
			expected: "lvl=info msg=info foo=\"[error1 error2]\"\n",
			msg:      "info",
			kvs:      []interface{}{"foo", []error{errors.New("error1"), errors.New("error2")}},
			f:        l.Info,
		},
		{
			name:     "map of strings",
			expected: "lvl=info msg=info foo=map[bar:baz]\n",
			msg:      "info",
			kvs:      []interface{}{"foo", map[string]string{"bar": "baz"}},
			f:        l.Info,
		},
		{
			name:     "map with interface",
			expected: "lvl=info msg=info foo=map[bar:baz]\n",
			msg:      "info",
			kvs:      []interface{}{"foo", map[string]interface{}{"bar": "baz"}},
			f:        l.Info,
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
