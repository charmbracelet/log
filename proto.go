package log

import (
	"fmt"
	"io"

	"github.com/go-kit/log/level"
)

var (
	msgKey = "message"
	tsKey  = "ts"
	errKey = "error"
	lvlKey = fmt.Sprint(level.Key())
)

type Logger interface {
	// SetOutput sets the output destination for the logger.
	SetOutput(w io.Writer)
	// SetLevel sets the minimum level for the logger.
	SetLevel(level Level)
	// SetOptions sets the logger options.
	SetOptions(opts ...Option)
	// SetFields sets the logger fields.
	SetFields(keyvals ...interface{})

	With(keyvals ...interface{}) Logger
	WithError(err error) Logger

	Debug(v ...interface{})
	Print(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})

	Debugln(v ...interface{})
	Println(v ...interface{})
	Infoln(v ...interface{})
	Warnln(v ...interface{})
	Errorln(v ...interface{})
	Fatalln(v ...interface{})

	Debugf(format string, v ...interface{})
	Printf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}
