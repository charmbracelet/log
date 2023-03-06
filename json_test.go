package log

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJson(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf)
	l.SetFormatter(JSONFormatter)
	cases := []struct {
		name     string
		expected string
		msg      string
		kvs      []interface{}
		f        func(msg interface{}, kvs ...interface{})
	}{
		{
			name:     "default logger info with timestamp",
			expected: "{\"lvl\":\"info\",\"msg\":\"info\"}\n",
			msg:      "info",
			kvs:      nil,
			f:        l.Info,
		},
		{
			name:     "default logger debug with timestamp",
			expected: "",
			msg:      "info",
			kvs:      nil,
			f:        l.Debug,
		},
		{
			name:     "default logger error with timestamp",
			expected: "{\"lvl\":\"error\",\"msg\":\"info\"}\n",
			msg:      "info",
			kvs:      nil,
			f:        l.Error,
		},
		{
			name:     "multiline message",
			expected: "{\"lvl\":\"error\",\"msg\":\"info\\ninfo\"}\n",
			msg:      "info\ninfo",
			kvs:      nil,
			f:        l.Error,
		},
		{
			name:     "multiline kvs",
			expected: "{\"lvl\":\"error\",\"msg\":\"info\",\"multiline\":\"info\\ninfo\"}\n",
			msg:      "info",
			kvs:      []interface{}{"multiline", "info\ninfo"},
			f:        l.Error,
		},
		{
			name:     "odd number of kvs",
			expected: "{\"baz\":\"missing value\",\"foo\":\"bar\",\"lvl\":\"error\",\"msg\":\"info\"}\n",
			msg:      "info",
			kvs:      []interface{}{"foo", "bar", "baz"},
			f:        l.Error,
		},
		{
			name:     "error field",
			expected: "{\"error\":\"error message\",\"lvl\":\"error\",\"msg\":\"info\"}\n",
			msg:      "info",
			kvs:      []interface{}{"error", errors.New("error message")},
			f:        l.Error,
		},
		{
			name:     "struct field",
			expected: "{\"lvl\":\"info\",\"msg\":\"info\",\"struct\":{}}\n",
			msg:      "info",
			kvs:      []interface{}{"struct", struct{ foo string }{foo: "bar"}},
			f:        l.Info,
		},
		{
			name:     "slice field",
			expected: "{\"lvl\":\"info\",\"msg\":\"info\",\"slice\":[1,2,3]}\n",
			msg:      "info",
			kvs:      []interface{}{"slice", []int{1, 2, 3}},
			f:        l.Info,
		},
		{
			name:     "slice of structs",
			expected: "{\"lvl\":\"info\",\"msg\":\"info\",\"slice\":[{},{}]}\n",
			msg:      "info",
			kvs:      []interface{}{"slice", []struct{ foo string }{{foo: "bar"}, {foo: "baz"}}},
			f:        l.Info,
		},
		{
			name:     "slice of strings",
			expected: "{\"lvl\":\"info\",\"msg\":\"info\",\"slice\":[\"foo\",\"bar\"]}\n",
			msg:      "info",
			kvs:      []interface{}{"slice", []string{"foo", "bar"}},
			f:        l.Info,
		},
		{
			name:     "slice of errors",
			expected: "{\"lvl\":\"info\",\"msg\":\"info\",\"slice\":[{},{}]}\n",
			msg:      "info",
			kvs:      []interface{}{"slice", []error{errors.New("error message1"), errors.New("error message2")}},
			f:        l.Info,
		},
		{
			name:     "map of strings",
			expected: "{\"lvl\":\"info\",\"map\":{\"a\":\"b\",\"foo\":\"bar\"},\"msg\":\"info\"}\n",
			msg:      "info",
			kvs:      []interface{}{"map", map[string]string{"a": "b", "foo": "bar"}},
			f:        l.Info,
		},
	}
	for _, c := range cases {
		buf.Reset()
		t.Run(c.name, func(t *testing.T) {
			c.f(c.msg, c.kvs...)
			require.Equal(t, c.expected, buf.String())
		})
	}
}

func TestJsonCaller(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf)
	l.SetFormatter(JSONFormatter)
	l.SetReportCaller(true)
	l.SetLevel(DebugLevel)
	_, file, line, _ := runtime.Caller(0)
	cases := []struct {
		name     string
		expected string
		msg      string
		kvs      []interface{}

		f func(msg interface{}, kvs ...interface{})
	}{
		{
			name:     "simple caller",
			expected: fmt.Sprintf("{\"caller\":\"log/%s:%d\",\"lvl\":\"info\",\"msg\":\"info\"}\n", filepath.Base(file), line+30),
			msg:      "info",
			kvs:      nil,
			f:        l.Info,
		},
		{
			name:     "nested caller",
			expected: fmt.Sprintf("{\"caller\":\"log/%s:%d\",\"lvl\":\"info\",\"msg\":\"info\"}\n", filepath.Base(file), line+30),
			msg:      "info",
			kvs:      nil,
			f: func(msg interface{}, kvs ...interface{}) {
				l.Helper()
				l.Info(msg, kvs...)
			},
		},
	}
	for _, c := range cases {
		buf.Reset()
		t.Run(c.name, func(t *testing.T) {
			c.f(c.msg, c.kvs...)
			require.Equal(t, c.expected, buf.String())
		})
	}
}

func TestJsonCustomKey(t *testing.T) {
	var buf bytes.Buffer
	oldTsKey := TimestampKey
	defer func() {
		TimestampKey = oldTsKey
	}()
	TimestampKey = "time"
	logger := New(&buf)
	logger.SetTimeFunction(_zeroTime)
	logger.SetFormatter(JSONFormatter)
	logger.SetReportTimestamp(true)
	logger.Info("info")
	require.Equal(t, "{\"lvl\":\"info\",\"msg\":\"info\",\"time\":\"0001/01/01 00:00:00\"}\n", buf.String())
}
