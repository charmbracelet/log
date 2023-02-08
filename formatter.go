package log

// Formatter is a formatter for log messages.
type Formatter uint8

const (
	// TextFormatter is a formatter that formats log messages as text. Suitable for
	// console output and log files.
	TextFormatter Formatter = iota
	// JSONFormatter is a formatter that formats log messages as JSON.
	JSONFormatter
	// LogfmtFormatter is a formatter that formats log messages as logfmt.
	LogfmtFormatter
)

const (
	tsKey     = "ts"
	msgKey    = "msg"
	lvlKey    = "lvl"
	callerKey = "caller"
	prefixKey = "prefix"
)
