package log

import (
	"fmt"

	"github.com/go-kit/log/level"
)

var (
	tsKey  = "timestamp"
	msgKey = "message"
	errKey = "error"
	lvlKey = fmt.Sprint(level.Key())
)

type Logger interface {
	SetLevel(lvl Level)
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
