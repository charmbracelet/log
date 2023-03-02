package log

import (
	"time"
)

// DefaultTimeFormat is the default time format.
const DefaultTimeFormat = "2006/01/02 15:04:05"

// TimeFunction is a function that returns a time.Time.
type TimeFunction = func() time.Time

// NowUTC is a convenient function that returns the
// current time in UTC timezone.
//
// This is to be used as a time function.
// For example:
//
//	log.SetTimeFunction(log.NowUTC)
func NowUTC() time.Time {
	return time.Now().UTC()
}

// Options is the options for the logger.
type Options struct {
	// TimeFunction is the time function for the logger. The default is time.Now.
	TimeFunction TimeFunction
	// TimeFormat is the time format for the logger. The default is "2006/01/02 15:04:05".
	TimeFormat string
	// Level is the level for the logger. The default is InfoLevel.
	Level Level
	// Prefix is the prefix for the logger. The default is no prefix.
	Prefix string
	// ReportTimestamp is whether the logger should report the timestamp. The default is false.
	ReportTimestamp bool
	// ReportCaller is whether the logger should report the caller location. The default is false.
	ReportCaller bool
	// Fields is the fields for the logger. The default is no fields.
	Fields []interface{}
	// Formatter is the formatter for the logger. The default is TextFormatter.
	Formatter Formatter
}
