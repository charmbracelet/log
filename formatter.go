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

// resolveLogValuer resolves values that implement slog.LogValuer to their
// LogValue, mirroring what the JSON formatter already does. A value passed
// directly satisfies slog.LogValuer, while the slog.Handler path wraps it in a
// slog.Value first; both are handled. The resolved slog.Value is returned so it
// renders the same way a slog.Value passed directly would, leaving values that
// are already resolved untouched. Values that carry no LogValue are returned
// unchanged.
func resolveLogValuer(value any) any {
	switch v := value.(type) {
	case slogLogValuer:
		return v.LogValue().Resolve()
	case slogValue:
		return v.Resolve()
	}
	return value
}

var (
	// TimestampKey is the key for the timestamp.
	TimestampKey = "time"
	// MessageKey is the key for the message.
	MessageKey = "msg"
	// LevelKey is the key for the level.
	LevelKey = "level"
	// CallerKey is the key for the caller.
	CallerKey = "caller"
	// PrefixKey is the key for the prefix.
	PrefixKey = "prefix"
)
