package log

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
)

func TestGlobal(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetTimeFunction(_zeroTime)
	cases := []struct {
		name     string
		expected string
		msg      string
		kvs      []interface{}
		f        func(msg interface{}, kvs ...interface{})
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
