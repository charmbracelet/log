package log

import (
	"bytes"
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

func writeIndent(w io.Writer, str string, indent string, newline bool) {
	// kindly borrowed from hclog
	for {
		nl := strings.IndexByte(str, '\n')
		if nl == -1 {
			if str != "" {
				_, _ = w.Write([]byte(indent))
				writeEscapedForOutput(w, str, false)
				if newline {
					_, _ = w.Write([]byte{'\n'})
				}
			}
			return
		}

		_, _ = w.Write([]byte(indent))
		writeEscapedForOutput(w, str[:nl], false)
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
		return new(bytes.Buffer)
	},
}

func writeEscapedForOutput(w io.Writer, str string, escapeQuotes bool) {
	// kindly borrowed from hclog
	if !needsEscaping(str) {
		_, _ = w.Write([]byte(str))
		return
	}

	bb := bufPool.Get().(*bytes.Buffer)
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

	_, _ = w.Write(bb.Bytes())
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

func (l *logger) textFormatter(keyvals ...interface{}) {
	s := l.styles
	for i := 0; i < len(keyvals); i += 2 {
		switch keyvals[i] {
		case tsKey:
			if t, ok := keyvals[i+1].(time.Time); ok {
				ts := t.Format(l.timeFormat)
				if !l.noStyles {
					ts = s.Timestamp.Render(ts)
				}
				l.b.WriteString(ts)
				l.b.WriteByte(' ')
			}
		case lvlKey:
			if level, ok := keyvals[i+1].(Level); ok {
				lvl := strings.ToUpper(level.String())
				if !l.noStyles {
					lvl = s.Level(level).String()
				}
				l.b.WriteString(lvl)
				l.b.WriteByte(' ')
			}
		case callerKey:
			if caller, ok := keyvals[i+1].(string); ok {
				caller = fmt.Sprintf("<%s>", caller)
				if !l.noStyles {
					caller = s.Caller.Render(caller)
				}
				l.b.WriteString(caller)
				l.b.WriteByte(' ')
			}
		case prefixKey:
			if prefix, ok := keyvals[i+1].(string); ok {
				if !l.noStyles {
					prefix = s.Prefix.Render(prefix)
				}
				l.b.WriteString(prefix)
				l.b.WriteByte(' ')
			}
		case msgKey:
			if msg := keyvals[i+1]; msg != nil {
				m := fmt.Sprint(msg)
				if !l.noStyles {
					m = s.Message.Render(m)
				}
				l.b.WriteString(m)
			}
		default:
			sep := separator
			indentSep := indentSeparator
			if !l.noStyles {
				sep = s.Separator.Render(sep)
				indentSep = s.Separator.Render(indentSep)
			}
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
			if !l.noStyles {
				key = s.Key.Render(key)
				val = s.Value.Render(val)
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
				writeIndent(&l.b, val, indentSep, moreKeys)
				// If there are more keyvals, separate them with a space.
				if moreKeys {
					l.b.WriteByte(' ')
				}
			} else if !raw && needsQuoting(val) {
				l.b.WriteByte(' ')
				l.b.WriteString(key)
				l.b.WriteString(sep)
				l.b.WriteByte('"')
				writeEscapedForOutput(&l.b, val, true)
				l.b.WriteByte('"')
			} else {
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
