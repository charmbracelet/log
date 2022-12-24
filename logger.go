package log

import (
	"io"
	"log"
	"time"
)

// DefaultTimeFormat is the default time format.
const DefaultTimeFormat = "2006/01/02 15:04:05"

// TimeFunction is a function that returns a time.Time.
type TimeFunction = func() time.Time

// Logger is an interface for logging.
type Logger interface {
	EnableTimestamp()
	DisableTimestamp()

	EnableCaller()
	DisableCaller()

	SetLevel(level Level)
	GetLevel() Level

	SetPrefix(prefix string)
	GetPrefix() string

	SetTimeFunction(f TimeFunction)
	SetTimeFormat(format string)
	SetOutput(w io.Writer)

	With(keyval ...interface{}) Logger

	Debug(msg interface{}, keyval ...interface{})
	Info(msg interface{}, keyval ...interface{})
	Warn(msg interface{}, keyval ...interface{})
	Error(msg interface{}, keyval ...interface{})

	StandardLogger() *log.Logger
}
