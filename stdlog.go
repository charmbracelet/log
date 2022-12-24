package log

import (
	"log"
	"strings"
)

type stdLogger struct {
	l *logger
}

func (l *stdLogger) Write(p []byte) (n int, err error) {
	str := strings.TrimSpace(string(p))

	switch {
	case strings.HasPrefix(str, "DEBUG"):
		l.l.Debug(strings.TrimSpace(str[5:]))
	case strings.HasPrefix(str, "INFO"):
		l.l.Info(strings.TrimSpace(str[4:]))
	case strings.HasPrefix(str, "WARN"):
		l.l.Warn(strings.TrimSpace(str[4:]))
	case strings.HasPrefix(str, "ERROR"):
		l.l.Error(strings.TrimSpace(str[5:]))
	case strings.HasPrefix(str, "ERR"):
		l.l.Error(strings.TrimSpace(str[3:]))
	default:
		l.l.Info(str)
	}

	return len(p), nil
}

// StandardLogger returns a standard logger from Logger. The returned logger
// can infer log levels from message prefix. Expected prefixes are DEBUG, INFO,
// WARN, ERROR, and ERR.
func (l *logger) StandardLogger() *log.Logger {
	nl := *l
	// The caller stack is
	// log.Printf() -> l.Output() -> l.out.Write(stdLogger.Write)
	nl.callerOffset += 3
	sl := &stdLogger{
		l: &nl,
	}
	return log.New(sl, "", 0)
}
