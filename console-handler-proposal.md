# Proposal: Enable Timestamps by Default in log.New()

## Problem

`log.New()` omits timestamps while `log.Default()` includes them. This inconsistency surprises users.

```go
// log.Default() includes timestamps
log.Info("Hello world!")
// Output: 2023/01/23 14:23:45 INFO Hello world!

// log.New() does not
logger := log.New(os.Stderr)
logger.Info("Hello world!")
// Output: INFO Hello world!
```

## Solution

Enable timestamps by default in `log.New()`:

```go
func New(w io.Writer) *Logger {
    return NewWithOptions(w, Options{
        ReportTimestamp: true,
    })
}
```

This uses `DefaultTimeFormat` (`"2006/01/02 15:04:05"`), matching `log.Default()` behavior.

Users can disable via `NewWithOptions()`:

```go
logger := log.NewWithOptions(os.Stderr, log.Options{
    ReportTimestamp: false,
})
```

## Impact

**Breaking change**: Existing code using `log.New()` will now include timestamps.

**Benefits**:
- Consistent with `log.Default()`
- Matches user expectations
- Aligns with standard library behavior

## API Compatibility with log/slog

The `*Logger` type already implements `slog.Handler`:

```go
handler := log.NewWithOptions(os.Stderr, log.Options{
    ReportTimestamp: true,
    Level:           log.DebugLevel,
})
logger := slog.New(handler)
```

**For a familiar slog-style API, add:**

```go
// HandlerOptions provides slog-style configuration.
type HandlerOptions struct {
    AddSource bool
    Level     Level
}

// NewTextHandler returns a Handler with text formatting.
// Timestamps are enabled by default.
func NewTextHandler(w io.Writer, opts *HandlerOptions) *Logger {
    if opts == nil {
        opts = &HandlerOptions{}
    }
    return NewWithOptions(w, Options{
        ReportTimestamp: true,
        Formatter:       TextFormatter,
        Level:           opts.Level,
        ReportCaller:    opts.AddSource,
    })
}

// NewJSONHandler returns a Handler with JSON formatting.
// Timestamps are enabled by default.
func NewJSONHandler(w io.Writer, opts *HandlerOptions) *Logger {
    if opts == nil {
        opts = &HandlerOptions{}
    }
    return NewWithOptions(w, Options{
        ReportTimestamp: true,
        Formatter:       JSONFormatter,
        Level:           opts.Level,
        ReportCaller:    opts.AddSource,
    })
}
```

**Usage:**

```go
handler := log.NewTextHandler(os.Stderr, &log.HandlerOptions{
    Level: log.DebugLevel,
})
logger := slog.New(handler)
```

## Implementation

- Update `New()` to set `ReportTimestamp: true`
- Add `HandlerOptions`, `NewTextHandler()`, `NewJSONHandler()`
- Update tests
- Document breaking change in release notes
