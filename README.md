# Log

<p>
  <img src="https://user-images.githubusercontent.com/25087/219742757-c8afe0d9-608a-4845-a555-ef59c0af9ebc.png#gh-light-mode-only" width="359" />
  <img src="https://user-images.githubusercontent.com/25087/219743408-3d7bef51-1409-40c0-8159-acc6e52f078e.png#gh-dark-mode-only" width="359" />
  <br>
  <a href="https://github.com/charmbracelet/log/releases"><img src="https://img.shields.io/github/release/charmbracelet/log.svg" alt="Latest Release"></a>
  <a href="https://pkg.go.dev/github.com/charmbracelet/log?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="Go Docs"></a>
  <a href="https://github.com/charmbracelet/log/actions"><img src="https://github.com/charmbracelet/log/workflows/build/badge.svg" alt="Build Status"></a>
</p>

A minimal and colorful Go logging library. 🪵

![Demo](./demo.gif)

It provides a leveled structured human readable logger with a small API. Unlike
[standard `log`][stdlog], the Charm logger provides customizable colorful human
readable logging with batteries included.

- Uses [lipgloss][lipgloss] to style and colorize the output.
- Beautiful human readable format.
- Ability to customize the time stamp format.
- Skips caller frames and marks function as helpers.
- Leveled logging with the ability to turn off logging altogether.
- Store and retrieve logger in and from context.
- Standard log Adapter.

## Usage

The Charm logger comes with a global package-wise logger with timestamps turned
on and the logging level set to `info`.

```go
log.Debug("cookie 🍪") // won't print anything
log.Info("Hello World!") // 2023/01/04 10:04:06 INFO Hello World!
```

All logging levels accept optional key/value pairs to be printed along with a
message.

```go
err := fmt.Errorf("too much sugar")
log.Error("failed to bake cookies", "err", err, "butter", "1 cup")
// 2023/01/04 10:04:06 ERROR failed to bake cookies err="too much sugar" butter="1 cup"
```

You can use `log.Print()` to print messages without a level prefix.

```go
log.Print("Baking 101") // 2023/01/04 10:04:06 Baking 101
```

### New loggers

Use `New()` to create new loggers.

```go
logger := log.New()
if butter {
    logger.Warn("chewy!", "butter", true) // WARN chewy! butter=true
}
```

### Options

You can customize the logger with options. Use `log.WithCaller()` to enable
printing source location. `log.WithTimestamp()` prints the timestamp of each
log.

```go
logger := log.New(log.WithTimestamp(), log.WithTimeFormat(time.Kitchen),
    log.WithCaller(), log.WithPrefix("baking 🍪"))
logger.Info("Starting oven!", "degree", 375)
// 10:00AM INFO <cookies/oven.go:56> baking 🍪: Starting oven! degree=375
time.Sleep(10 * time.Minute)
logger.Info("Finished baking")
// 10:10AM INFO <cookies/oven.go:60> baking 🍪: Finished baking
```

Use `log.SetFormatter()` or `log.WithFormatter()` to change the output format.
Available options are:

- `log.TextFormatter` (_default_)
- `log.JSONFormatter`
- `log.LogfmtFormatter`

> **Note** styling only affects the `TextFormatter`.

For a list of available options, refer to [options.go](./options.go).

Set the logger level and options.

```go
logger.SetReportTimestamp(false)
logger.SetReportCaller(false)
logger.SetLevel(log.DebugLevel)
logger.Debug("Preparing batch 2...") // DEBUG baking 🍪: Preparing batch 2...
```

### Sub-logger

Create sub-loggers with their specific fields.

```go
batch2 := logger.With("batch", 2, "chocolateChips", true)
batch2.Debug("Adding chocolate chips")
// DEBUG <cookies/oven.go:68> baking 🍪: Adding chocolate chips batch=2 chocolateChips=true
```

### Format Messages

You can use `fmt.Sprintf()` to format messages.

```go
for item := 1; i <= 100; i++ {
    log.Info(fmt.Sprintf("Baking %d/100...", item))
    // INFO Baking 1/100...
    // INFO Baking 2/100...
    // INFO Baking 3/100...
    // ...
    // INFO Baking 100/100...
}
```

Or arguments:

```go
for temp := 375; temp <= 400; temp++ {
    log.Info("Increasing temperature", "degree", fmt.Sprintf("%d°F", temp))
    // INFO Increasing temperature degree=375°F
    // INFO Increasing temperature degree=376°F
    // INFO Increasing temperature degree=377°F
    // ...
    // INFO Increasing temperature degree=400°F
}
```

### Helper Functions

Skip caller frames in helper functions. Similar to what you can do with
`testing.TB().Helper()`.

```go
function startOven(degree int) {
    log.Helper()
    log.Info("Starting oven", "degree", degree)
}

log.SetReportCaller(true)
startOven(400) // INFO <cookies/oven.go:123> Starting oven degree=400
```

This will use the _caller_ function (`startOven`) line number instead of the
logging function (`log.Info`) to report the source location.

### Standard Log Adapter

Some Go libraries, especially the ones in the standard library, will only accept
the [standard logger][stdlog] interface. For instance, the HTTP Server from
`net/http` will only take a `*log.Logger` for its `ErrorLog` field.

For this, you can use the standard log adapter which simply wraps the logger in
a `*log.Logger` interface.

```go
stdlog := log.New(log.WithPrefix("http")).StandardLog(log.StandardLogOption{
    ForceLevel: log.ErrorLevel,
})
s := &http.Server{
    Addr:     ":8080",
    Handler:  handler,
    ErrorLog: stdlog,
}
stdlog.Printf("Failed to make bake request, %s", fmt.Errorf("temperature is to low"))
// ERROR http: Failed to make bake request, temperature is to low
```

[lipgloss]: https://github.com/charmbracelet/lipgloss
[stdlog]: https://pkg.go.dev/log

## License

[MIT](https://github.com/charmbracelet/log/raw/master/LICENSE)

---

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="the Charm logo" src="https://stuff.charm.sh/charm-badge.jpg" width="400"></a>

Charm热爱开源 • Charm loves open source • نحنُ نحب المصادر المفتوحة
