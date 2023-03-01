package log

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"

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
			expected: "0001/01/01 00:00:00 INFO info\n",
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
			expected: "0001/01/01 00:00:00 ERROR info\n",
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

func TestFullCallerSource(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetCallerFormat(CallerLong)
	SetReportCaller(true)
	SetReportTimestamp(false)
	_, file, line, _ := runtime.Caller(0)
	Error("error")
	assert.Equal(t, fmt.Sprintf("ERROR <%s:%d> error\n", file, line+1), buf.String())
	SetReportTimestamp(true)
	SetReportCaller(false)
	SetCallerFormat(CallerShort)
}

func TestPrint(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(FatalLevel)
	SetTimeFunction(_zeroTime)
	Error("error")
	Print("print")
	assert.Equal(t, "0001/01/01 00:00:00 print\n", buf.String())
}

func TestFatal(t *testing.T) {
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
