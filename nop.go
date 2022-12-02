package log

import "github.com/go-kit/log"

// NewNopLogger returns a new logger that discards all log events.
func NewNopLogger() Logger {
	logger := log.NewNopLogger()
	output := &Output{
		logger: logger,
	}
	return output
}
