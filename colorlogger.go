package log

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/go-kit/log"
)

// LevelStyle defines the colors for each level.
type LevelStyle struct {
	Name      string
	Level     lipgloss.Style
	Message   lipgloss.Style
	Keys      lipgloss.Style
	Values    lipgloss.Style
	Timestamp lipgloss.Style
}

// Styles is the default styles map.
type Styles struct {
	Debug LevelStyle
	Info  LevelStyle
	Warn  LevelStyle
	Error LevelStyle
}

// DefaultStyles returns the default styles.
func DefaultStyles() *Styles {
	return &Styles{
		Debug: LevelStyle{
			Name:      "DBG",
			Level:     lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true),
			Message:   lipgloss.NewStyle(),
			Keys:      lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
			Values:    lipgloss.NewStyle(),
			Timestamp: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		},
		Info: LevelStyle{
			Name:      "INF",
			Level:     lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true),
			Message:   lipgloss.NewStyle(),
			Keys:      lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
			Values:    lipgloss.NewStyle(),
			Timestamp: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		},
		Warn: LevelStyle{
			Name:      "WRN",
			Level:     lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true),
			Message:   lipgloss.NewStyle(),
			Keys:      lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
			Values:    lipgloss.NewStyle(),
			Timestamp: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		},
		Error: LevelStyle{
			Name:      "ERR",
			Level:     lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true),
			Message:   lipgloss.NewStyle(),
			Keys:      lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
			Values:    lipgloss.NewStyle(),
			Timestamp: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		},
	}
}

// newColorLogger returns a Logger which writes colored logs to w. ANSI color
// codes for the colors returned by color are added to the formatted output
// from the Logger returned by newLogger and the combined result written to w.
func newColorLogger(w io.Writer, newLogger func(io.Writer) log.Logger, styles *Styles) log.Logger {
	if styles == nil {
		styles = DefaultStyles()
	}
	return &colorLogger{
		w:             w,
		newLogger:     newLogger,
		styles:        styles,
		bufPool:       sync.Pool{New: func() interface{} { return &loggerBuf{} }},
		noColorLogger: newLogger(w),
	}
}

type colorLogger struct {
	w             io.Writer
	newLogger     func(io.Writer) log.Logger
	styles        *Styles
	bufPool       sync.Pool
	noColorLogger log.Logger
}

func (l *colorLogger) Log(keyvals ...interface{}) error {
	lb := l.getLoggerBuf()
	defer l.putLoggerBuf(lb)
	var ts string
	var msg string
	var err string
	var lvl string
	keys := make([]interface{}, 0, len(keyvals)/2)
	values := make([]interface{}, 0, len(keyvals)/2)

	for i := 0; i < len(keyvals); i += 2 {
		key := keyvals[i]
		switch key {
		case tsKey:
			ts = fmt.Sprint(keyvals[i+1])
		case msgKey:
			msg = fmt.Sprint(keyvals[i+1])
		case errKey:
			err = fmt.Sprint(keyvals[i+1])
		case lvlKey:
			lvl = fmt.Sprint(keyvals[i+1])
		default:
			keys = append(keys, key)
			values = append(values, keyvals[i+1])
		}
	}

	var styles *LevelStyle
	var name string
	switch lvl {
	case "debug":
		styles = &l.styles.Debug
		name = styles.Name
	case "info":
		styles = &l.styles.Info
		name = styles.Name
	case "warn":
		styles = &l.styles.Warn
		name = styles.Name
	case "error":
		styles = &l.styles.Error
		name = styles.Name
	}

	fmt.Fprintf(lb.buf, "%s %s %s ",
		styles.Timestamp.Render(fmt.Sprint(ts)),
		styles.Level.Render(name),
		styles.Message.Render(msg),
	)

	if err != "" {
		fmt.Fprintf(lb.buf, "%s", styles.Keys.Render(err))
	}

	for i := 0; i < len(keys); i++ {
		switch keys[i] {
		case msgKey, tsKey, errKey, lvlKey:
			continue
		}
		fmt.Fprintf(lb.buf, " %s=%s",
			styles.Keys.Render(fmt.Sprint(keys[i])),
			styles.Values.Render(fmt.Sprint(values[i])),
		)
	}
	fmt.Fprintln(lb.buf)
	if _, err := io.Copy(l.w, lb.buf); err != nil {
		return err
	}
	return nil
}

type loggerBuf struct {
	buf    *bytes.Buffer
	logger log.Logger
}

func (l *colorLogger) getLoggerBuf() *loggerBuf {
	lb := l.bufPool.Get().(*loggerBuf)
	if lb.buf == nil {
		lb.buf = &bytes.Buffer{}
		lb.logger = l.newLogger(lb.buf)
	} else {
		lb.buf.Reset()
	}
	return lb
}

func (l *colorLogger) putLoggerBuf(cb *loggerBuf) {
	l.bufPool.Put(cb)
}
