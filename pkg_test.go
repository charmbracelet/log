package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
)

func TestDefaultRace(t *testing.T) {
	l := Default()
	t.Cleanup(func() {
		SetDefault(l)
	})

	for i := 0; i < 2; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			SetDefault(New(io.Discard))
			Default().Info("foo")
		})
	}
}

func TestGlobal(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetTimeFunction(_zeroTime)
	cases := []struct {
		name     string
		expected string
		msg      string
		kvs      []any
		f        func(msg any, kvs ...any)
	}{
		{
			name:     "default logger info with timestamp",
			expected: "0002/01/01 00:00:00 INFO info\n",
			msg:      "info",
			kvs:      nil,
			f:        Info,
		},
		{
			name:     "default logger debug with timestamp",
			expected: "",
			msg:      "info",
			kvs:      nil,
			f:        Debug,
		},
		{
			name:     "default logger error with timestamp",
			expected: "0002/01/01 00:00:00 ERRO info\n",
			msg:      "info",
			kvs:      nil,
			f:        Error,
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

func TestPrint(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(FatalLevel)
	SetTimeFunction(_zeroTime)
	SetReportTimestamp(true)
	SetReportCaller(false)
	SetTimeFormat(DefaultTimeFormat)
	SetColorProfile(termenv.ANSI)
	Error("error")
	Print("print")
	assert.Equal(t, "0002/01/01 00:00:00 print\n", buf.String())
}

func TestPrintf(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(FatalLevel)
	SetTimeFunction(_zeroTime)
	SetReportTimestamp(true)
	SetReportCaller(false)
	SetTimeFormat(DefaultTimeFormat)
	Errorf("error")
	Printf("print")
	assert.Equal(t, "0002/01/01 00:00:00 print\n", buf.String())
}

func TestFatal(t *testing.T) {
	SetReportTimestamp(true)
	SetReportCaller(false)
	SetTimeFormat(DefaultTimeFormat)
	if os.Getenv("FATAL") == "1" {
		Fatal("i'm dead")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestFatal")
	cmd.Env = append(os.Environ(), "FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestFatalf(t *testing.T) {
	SetReportTimestamp(true)
	SetReportCaller(false)
	SetTimeFormat(DefaultTimeFormat)
	if os.Getenv("FATAL") == "1" {
		Fatalf("i'm %s", "dead")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalf")
	cmd.Env = append(os.Environ(), "FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestDebugf(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(DebugLevel)
	SetTimeFunction(_zeroTime)
	SetReportTimestamp(true)
	SetReportCaller(true)
	SetTimeFormat(DefaultTimeFormat)
	_, file, line, _ := runtime.Caller(0)
	Debugf("debug %s", "foo")
	assert.Equal(t, fmt.Sprintf("0002/01/01 00:00:00 DEBU <log/%s:%d> debug foo\n", filepath.Base(file), line+1), buf.String())
}

func TestInfof(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(InfoLevel)
	SetReportTimestamp(false)
	SetReportCaller(false)
	SetTimeFormat(DefaultTimeFormat)
	Infof("info %s", "foo")
	assert.Equal(t, "INFO info foo\n", buf.String())
}

func TestWarnf(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(WarnLevel)
	SetReportCaller(false)
	SetReportTimestamp(true)
	SetTimeFunction(_zeroTime)
	SetTimeFormat(DefaultTimeFormat)
	Warnf("warn %s", "foo")
	assert.Equal(t, "0002/01/01 00:00:00 WARN warn foo\n", buf.String())
}

func TestErrorf(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(ErrorLevel)
	SetReportCaller(false)
	SetReportTimestamp(true)
	SetTimeFunction(_zeroTime)
	SetTimeFormat(time.Kitchen)
	Errorf("error %s", "foo")
	assert.Equal(t, "12:00AM ERRO error foo\n", buf.String())
}

func TestWith(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(InfoLevel)
	SetReportCaller(false)
	SetReportTimestamp(true)
	SetTimeFunction(_zeroTime)
	SetTimeFormat(DefaultTimeFormat)
	With("foo", "bar").Info("info")
	assert.Equal(t, "0002/01/01 00:00:00 INFO info foo=bar\n", buf.String())
}

func TestGetLevel(t *testing.T) {
	SetLevel(InfoLevel)
	assert.Equal(t, InfoLevel, GetLevel())
}

func TestPrefix(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(WarnLevel)
	SetReportCaller(false)
	SetReportTimestamp(false)
	SetPrefix("prefix")
	Warn("info")
	assert.Equal(t, "WARN prefix: info\n", buf.String())
	assert.Equal(t, "prefix", GetPrefix())
	SetPrefix("")
}

func TestFormatter(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(InfoLevel)
	SetReportCaller(false)
	SetReportTimestamp(false)
	SetFormatter(JSONFormatter)
	Info("info")
	assert.Equal(t, "{\"level\":\"info\",\"msg\":\"info\"}\n", buf.String())
}

func TestWithPrefix(t *testing.T) {
	l := WithPrefix("test")
	assert.Equal(t, "test", l.prefix)
}

func TestGlobalCustomLevel(t *testing.T) {
	var buf bytes.Buffer
	lvl := Level(-1)
	SetOutput(&buf)
	SetLevel(lvl)
	SetReportCaller(false)
	SetReportTimestamp(false)
	SetFormatter(JSONFormatter)
	Log(lvl, "info")
	Logf(lvl, "hey %s", "you")
	assert.Equal(t, "{\"msg\":\"info\"}\n{\"msg\":\"hey you\"}\n", buf.String())
}

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf)
	l.SetTimeFunction(_zeroTime)
	l.Info("hello")
	assert.Equal(t, "0002/01/01 00:00:00 INFO hello\n", buf.String())
	assert.True(t, l.reportTimestamp)
}

func TestNewTextHandler(t *testing.T) {
	var buf bytes.Buffer
	l := NewTextHandler(&buf, &HandlerOptions{Level: DebugLevel})
	l.SetTimeFunction(_zeroTime)

	l.Debug("debug msg")
	assert.Equal(t, "0002/01/01 00:00:00 DEBU debug msg\n", buf.String())
	assert.True(t, l.reportTimestamp)
	assert.Equal(t, TextFormatter, l.formatter)
	assert.Equal(t, DebugLevel, l.GetLevel())
}

func TestNewTextHandlerWithSource(t *testing.T) {
	var buf bytes.Buffer
	l := NewTextHandler(&buf, &HandlerOptions{AddSource: true})
	l.SetTimeFunction(_zeroTime)
	_, file, line, _ := runtime.Caller(0)
	l.Info("hello")
	assert.Equal(t, fmt.Sprintf("0002/01/01 00:00:00 INFO <log/%s:%d> hello\n", filepath.Base(file), line+1), buf.String())
	assert.True(t, l.reportCaller)
}

func TestNewJSONHandler(t *testing.T) {
	var buf bytes.Buffer
	l := NewJSONHandler(&buf, &HandlerOptions{Level: WarnLevel})
	l.SetTimeFunction(_zeroTime)

	l.Warn("warn msg")
	assert.Equal(t, "{\"time\":\"0002/01/01 00:00:00\",\"level\":\"warn\",\"msg\":\"warn msg\"}\n", buf.String())
	assert.True(t, l.reportTimestamp)
	assert.Equal(t, JSONFormatter, l.formatter)
	assert.Equal(t, WarnLevel, l.GetLevel())
}

func TestNewJSONHandlerWithSource(t *testing.T) {
	var buf bytes.Buffer
	l := NewJSONHandler(&buf, &HandlerOptions{AddSource: true})
	l.SetTimeFunction(_zeroTime)
	_, file, line, _ := runtime.Caller(0)
	l.Info("hello")
	expected := fmt.Sprintf("{\"time\":\"0002/01/01 00:00:00\",\"level\":\"info\",\"caller\":\"%s:%d\",\"msg\":\"hello\"}\n", trimCallerPath(file, 2), line+1)
	assert.Equal(t, expected, buf.String())
	assert.True(t, l.reportCaller)
}

func TestHandlerOptionsNil(t *testing.T) {
	var buf bytes.Buffer

	// NewTextHandler with nil opts
	tl := NewTextHandler(&buf, nil)
	tl.SetTimeFunction(_zeroTime)
	tl.Info("text")
	assert.Equal(t, "0002/01/01 00:00:00 INFO text\n", buf.String())
	assert.Equal(t, InfoLevel, tl.GetLevel())
	assert.False(t, tl.reportCaller)

	buf.Reset()

	// NewJSONHandler with nil opts
	jl := NewJSONHandler(&buf, nil)
	jl.SetTimeFunction(_zeroTime)
	jl.Info("json")
	assert.Equal(t, "{\"time\":\"0002/01/01 00:00:00\",\"level\":\"info\",\"msg\":\"json\"}\n", buf.String())
	assert.Equal(t, InfoLevel, jl.GetLevel())
	assert.False(t, jl.reportCaller)
}
