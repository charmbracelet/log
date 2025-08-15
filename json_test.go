package log

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/require"
)

func TestJsonCustomLevelWithStyle(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf)
	styles := DefaultStyles()
	Levels[int(Critical)] = Critical
	styles.Levels[int(Critical)] = lipgloss.NewStyle().
		SetString(strings.ToUpper(Critical.String())).
		Bold(true).
		MaxWidth(4).
		Foreground(lipgloss.Color("134"))
	l.SetStyles(styles)
	l.SetLevel(InfoLevel)
	l.SetFormatter(JSONFormatter)
	l.Logf(Critical, "foo")
	require.Equal(t, "{\"level\":\"crit\",\"msg\":\"foo\"}\n", buf.String())
}

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
			expected: "{\"level\":\"info\",\"msg\":\"info\"}\n",
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
			expected: "{\"level\":\"error\",\"msg\":\"info\"}\n",
			msg:      "info",
			kvs:      nil,
			f:        l.Error,
		},
		{
			name:     "multiline message",
			expected: "{\"level\":\"error\",\"msg\":\"info\\ninfo\"}\n",
			msg:      "info\ninfo",
			kvs:      nil,
			f:        l.Error,
		},
		{
			name:     "multiline kvs",
			expected: "{\"level\":\"error\",\"msg\":\"info\",\"multiline\":\"info\\ninfo\"}\n",
			msg:      "info",
			kvs:      []interface{}{"multiline", "info\ninfo"},
			f:        l.Error,
		},
		{
			name:     "odd number of kvs",
			expected: "{\"level\":\"error\",\"msg\":\"info\",\"foo\":\"bar\",\"baz\":\"missing value\"}\n",
			msg:      "info",
			kvs:      []interface{}{"foo", "bar", "baz"},
			f:        l.Error,
		},
		{
			name:     "error field",
			expected: "{\"level\":\"error\",\"msg\":\"info\",\"error\":\"error message\"}\n",
			msg:      "info",
			kvs:      []interface{}{"error", errors.New("error message")},
			f:        l.Error,
		},
		{
			name:     "struct field",
			expected: "{\"level\":\"info\",\"msg\":\"info\",\"struct\":{}}\n",
			msg:      "info",
			kvs:      []interface{}{"struct", struct{ foo string }{foo: "bar"}},
			f:        l.Info,
		},
		{
			name:     "slice field",
			expected: "{\"level\":\"info\",\"msg\":\"info\",\"slice\":[1,2,3]}\n",
			msg:      "info",
			kvs:      []interface{}{"slice", []int{1, 2, 3}},
			f:        l.Info,
		},
		{
			name:     "slice of structs",
			expected: "{\"level\":\"info\",\"msg\":\"info\",\"slice\":[{},{}]}\n",
			msg:      "info",
			kvs:      []interface{}{"slice", []struct{ foo string }{{foo: "bar"}, {foo: "baz"}}},
			f:        l.Info,
		},
		{
			name:     "slice of strings",
			expected: "{\"level\":\"info\",\"msg\":\"info\",\"slice\":[\"foo\",\"bar\"]}\n",
			msg:      "info",
			kvs:      []interface{}{"slice", []string{"foo", "bar"}},
			f:        l.Info,
		},
		{
			name:     "slice of errors",
			expected: "{\"level\":\"info\",\"msg\":\"info\",\"slice\":[{},{}]}\n",
			msg:      "info",
			kvs:      []interface{}{"slice", []error{errors.New("error message1"), errors.New("error message2")}},
			f:        l.Info,
		},
		{
			name:     "map of strings",
			expected: "{\"level\":\"info\",\"msg\":\"info\",\"map\":{\"a\":\"b\",\"foo\":\"bar\"}}\n",
			msg:      "info",
			kvs:      []interface{}{"map", map[string]string{"a": "b", "foo": "bar"}},
			f:        l.Info,
		},
		{
			name:     "slog any value error type",
			expected: "{\"level\":\"info\",\"msg\":\"info\",\"error\":\"error message\"}\n",
			msg:      "info",
			kvs:      []interface{}{"error", slogAnyValue(fmt.Errorf("error message"))},
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
			expected: fmt.Sprintf("{\"level\":\"info\",\"caller\":\"log/%s:%d\",\"msg\":\"info\"}\n", filepath.Base(file), line+30),
			msg:      "info",
			kvs:      nil,
			f:        l.Info,
		},
		{
			name:     "nested caller",
			expected: fmt.Sprintf("{\"level\":\"info\",\"caller\":\"log/%s:%d\",\"msg\":\"info\"}\n", filepath.Base(file), line+30),
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

func TestJsonTime(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf)
	logger.SetTimeFunction(_zeroTime)
	logger.SetFormatter(JSONFormatter)
	logger.SetReportTimestamp(true)
	logger.Info("info")
	require.Equal(t, "{\"time\":\"0002/01/01 00:00:00\",\"level\":\"info\",\"msg\":\"info\"}\n", buf.String())
}

func TestJsonPrefix(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf)
	logger.SetFormatter(JSONFormatter)
	logger.SetPrefix("my-prefix")
	logger.Info("info")
	require.Equal(t, "{\"level\":\"info\",\"prefix\":\"my-prefix\",\"msg\":\"info\"}\n", buf.String())
}

func TestJsonCustomKey(t *testing.T) {
	var buf bytes.Buffer
	oldTsKey := TimestampKey
	defer func() {
		TimestampKey = oldTsKey
	}()
	TimestampKey = "other-time"
	logger := New(&buf)
	logger.SetTimeFunction(_zeroTime)
	logger.SetFormatter(JSONFormatter)
	logger.SetReportTimestamp(true)
	logger.Info("info")
	require.Equal(t, "{\"other-time\":\"0002/01/01 00:00:00\",\"level\":\"info\",\"msg\":\"info\"}\n", buf.String())
}

func TestJsonWriter(t *testing.T) {
	testCases := []struct {
		name     string
		fn       func(w *jsonWriter)
		expected string
	}{
		{
			"string",
			func(w *jsonWriter) {
				w.start()
				w.objectItem("a", "value")
				w.end()
			},
			`{"a":"value"}`,
		},
		{
			"int",
			func(w *jsonWriter) {
				w.start()
				w.objectItem("a", 123)
				w.end()
			},
			`{"a":123}`,
		},
		{
			"bytes",
			func(w *jsonWriter) {
				w.start()
				w.objectItem("b", []byte{0x0, 0x1})
				w.end()
			},
			`{"b":"AAE="}`,
		},
		{
			"no fields",
			func(w *jsonWriter) {
				w.start()
				w.end()
			},
			`{}`,
		},
		{
			"multiple in asc order",
			func(w *jsonWriter) {
				w.start()
				w.objectItem("a", "value")
				w.objectItem("b", "some-other")
				w.end()
			},
			`{"a":"value","b":"some-other"}`,
		},
		{
			"multiple in desc order",
			func(w *jsonWriter) {
				w.start()
				w.objectItem("b", "some-other")
				w.objectItem("a", "value")
				w.end()
			},
			`{"b":"some-other","a":"value"}`,
		},
		{
			"depth",
			func(w *jsonWriter) {
				w.start()
				w.objectItem("a", map[string]int{"b": 123})
				w.end()
			},
			`{"a":{"b":123}}`,
		},
		{
			"key contains reserved",
			func(w *jsonWriter) {
				w.start()
				w.objectItem("a:\"b", "value")
				w.end()
			},
			`{"a:\"b":"value"}`,
		},
		{
			"pointer",
			func(w *jsonWriter) {
				w.start()
				w.objectItem("a", ptr("pointer"))
				w.end()
			},
			`{"a":"pointer"}`,
		},
		{
			"double-pointer",
			func(w *jsonWriter) {
				w.start()
				w.objectItem("a", ptr(ptr("pointer")))
				w.end()
			},
			`{"a":"pointer"}`,
		},
		{
			"invalid",
			func(w *jsonWriter) {
				w.start()
				w.objectItem("a", invalidJSON{})
				w.end()
			},
			`{"a":"invalid value"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			tc.fn(&jsonWriter{w: &buf})
			require.Equal(t, tc.expected, buf.String())
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}

type invalidJSON struct{}

func (invalidJSON) MarshalJSON() ([]byte, error) {
	return nil, errors.New("invalid json error")
}
