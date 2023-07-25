package log

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	separator       = "="
	indentSeparator = "  │ "
)

func (l *Logger) writeIndent(w io.Writer, str string, indent string, newline bool, key string) {
	// kindly borrowed from hclog
	for {
		nl := strings.IndexByte(str, '\n')
		if nl == -1 {
			if str != "" {
				_, _ = w.Write([]byte(indent))
				val := escapeStringForOutput(str, false)
				if valueStyle, ok := ValueStyles[key]; ok {
					val = valueStyle.Renderer(l.re).Render(val)
				} else {
					val = ValueStyle.Renderer(l.re).Render(val)
				}
				_, _ = w.Write([]byte(val))
				if newline {
					_, _ = w.Write([]byte{'\n'})
				}
			}
			return
		}

		_, _ = w.Write([]byte(indent))
		val := escapeStringForOutput(str[:nl], false)
		val = ValueStyle.Renderer(l.re).Render(val)
		_, _ = w.Write([]byte(val))
		_, _ = w.Write([]byte{'\n'})
		str = str[nl+1:]
	}
}

func needsEscaping(str string) bool {
	for _, b := range str {
		if !unicode.IsPrint(b) || b == '"' {
			return true
		}
	}

	return false
}

const (
	lowerhex = "0123456789abcdef"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(strings.Builder)
	},
}

func escapeStringForOutput(str string, escapeQuotes bool) string {
	// kindly borrowed from hclog
	if !needsEscaping(str) {
		return str
	}

	bb := bufPool.Get().(*strings.Builder)
	bb.Reset()

	defer bufPool.Put(bb)

	for _, r := range str {
		if escapeQuotes && r == '"' {
			bb.WriteString(`\"`)
		} else if unicode.IsPrint(r) {
			bb.WriteRune(r)
		} else {
			switch r {
			case '\a':
				bb.WriteString(`\a`)
			case '\b':
				bb.WriteString(`\b`)
			case '\f':
				bb.WriteString(`\f`)
			case '\n':
				bb.WriteString(`\n`)
			case '\r':
				bb.WriteString(`\r`)
			case '\t':
				bb.WriteString(`\t`)
			case '\v':
				bb.WriteString(`\v`)
			default:
				switch {
				case r < ' ':
					bb.WriteString(`\x`)
					bb.WriteByte(lowerhex[byte(r)>>4])
					bb.WriteByte(lowerhex[byte(r)&0xF])
				case !utf8.ValidRune(r):
					r = 0xFFFD
					fallthrough
				case r < 0x10000:
					bb.WriteString(`\u`)
					for s := 12; s >= 0; s -= 4 {
						bb.WriteByte(lowerhex[r>>uint(s)&0xF])
					}
				default:
					bb.WriteString(`\U`)
					for s := 28; s >= 0; s -= 4 {
						bb.WriteByte(lowerhex[r>>uint(s)&0xF])
					}
				}
			}
		}
	}

	return bb.String()
}

// isNormal indicates if the rune is one allowed to exist as an unquoted
// string value. This is a subset of ASCII, `-` through `~`.
func isNormal(r rune) bool {
	return '-' <= r && r <= '~'
}

// needsQuoting returns false if all the runes in string are normal, according
// to isNormal.
func needsQuoting(str string) bool {
	for _, r := range str {
		if !isNormal(r) {
			return true
		}
	}

	return false
}

func (l *Logger) textFormatter(keyvals ...interface{}) {
	for i := 0; i < len(keyvals); i += 2 {
		switch keyvals[i] {
		case TimestampKey:
			if t, ok := keyvals[i+1].(time.Time); ok {
				ts := t.Format(l.timeFormat)
				ts = TimestampStyle.Renderer(l.re).Render(ts)
				l.b.WriteString(ts)
				l.b.WriteByte(' ')
			}
		case LevelKey:
			if level, ok := keyvals[i+1].(Level); ok {
				lvl := levelStyle(level).Renderer(l.re).String()
				l.b.WriteString(lvl)
				l.b.WriteByte(' ')
			}
		case CallerKey:
			if caller, ok := keyvals[i+1].(string); ok {
				caller = fmt.Sprintf("<%s>", caller)
				caller = CallerStyle.Renderer(l.re).Render(caller)
				l.b.WriteString(caller)
				l.b.WriteByte(' ')
			}
		case PrefixKey:
			if prefix, ok := keyvals[i+1].(string); ok {
				prefix = PrefixStyle.Renderer(l.re).Render(prefix + ":")
				l.b.WriteString(prefix)
				l.b.WriteByte(' ')
			}
		case MessageKey:
			if msg := keyvals[i+1]; msg != nil {
				m := fmt.Sprint(msg)
				m = MessageStyle.Renderer(l.re).Render(m)
				l.b.WriteString(m)
			}
		default:
			sep := separator
			indentSep := indentSeparator
			sep = SeparatorStyle.Renderer(l.re).Render(sep)
			indentSep = SeparatorStyle.Renderer(l.re).Render(indentSep)
			moreKeys := i < len(keyvals)-2
			key := fmt.Sprint(keyvals[i])
			val := fmt.Sprintf("%+v", keyvals[i+1])
			raw := val == ""
			if raw {
				val = `""`
			}
			if key == "" {
				continue
			}
			actualKey := key
			valueStyle := ValueStyle
			if vs, ok := ValueStyles[actualKey]; ok {
				valueStyle = vs
			}
			if keyStyle, ok := KeyStyles[key]; ok {
				key = keyStyle.Renderer(l.re).Render(key)
			} else {
				key = KeyStyle.Renderer(l.re).Render(key)
			}

			// Values may contain multiple lines, and that format
			// is preserved, with each line prefixed with a "  | "
			// to show it's part of a collection of lines.
			//
			// Values may also need quoting, if not all the runes
			// in the value string are "normal", like if they
			// contain ANSI escape sequences.
			if strings.Contains(val, "\n") {
				l.b.WriteString("\n  ")
				l.b.WriteString(key)
				l.b.WriteString(sep + "\n")
				l.writeIndent(&l.b, val, indentSep, moreKeys, actualKey)
				// If there are more keyvals, separate them with a space.
				if moreKeys {
					l.b.WriteByte(' ')
				}
			} else if !raw && needsQuoting(val) {
				l.b.WriteByte(' ')
				l.b.WriteString(key)
				l.b.WriteString(sep)
				l.b.WriteString(valueStyle.Renderer(l.re).Render(fmt.Sprintf(`"%s"`,
					escapeStringForOutput(val, true))))
			} else {
				val = valueStyle.Renderer(l.re).Render(val)
				l.b.WriteByte(' ')
				l.b.WriteString(key)
				l.b.WriteString(sep)
				l.b.WriteString(val)
			}
		}
	}

	// Add a newline to the end of the log message.
	l.b.WriteByte('\n')
}
