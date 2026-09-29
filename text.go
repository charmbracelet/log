package log

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"charm.land/lipgloss/v2"
)

const (
	separator       = "="
	indentSeparator = "  │ "
)

func (l *Logger) writeIndent(w io.Writer, str string, indent string, newline bool, key string) {
	st := l.styles
	str = normalizeLogNewlines(str)

	// kindly borrowed from hclog
	for {
		nl := strings.IndexByte(str, '\n')
		if nl == -1 { //nolint:nestif
			if str != "" {
				_, _ = w.Write([]byte(indent))
				val := escapeStringForOutput(str, false)
				if valueStyle, ok := st.Values[key]; ok {
					val = renderKeepTabs(valueStyle, val)
				} else {
					val = renderKeepTabs(st.Value, val)
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
		val = renderKeepTabs(st.Value, val)
		_, _ = w.Write([]byte(val))
		_, _ = w.Write([]byte{'\n'})
		str = str[nl+1:]
	}
}

// normalizeLogNewlines turns CRLF and leftover CR into LF so Windows line
// endings don't show up as a literal \r in the text formatter.
func normalizeLogNewlines(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.ReplaceAll(s, "\r", "\n")
}

// renderKeepTabs styles s without converting tabs to spaces. Lip Gloss
// otherwise replaces tabs with 4 spaces (or removes them when TabWidth is 0),
// which breaks pretty-printed errors that indent with tabs.
func renderKeepTabs(style lipgloss.Style, s string) string {
	return style.TabWidth(lipgloss.NoTabConversion).Render(s)
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
	New: func() any {
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
				// Keep tabs as tabs so multiline values (pretty-printed
				// errors, source fragments) retain their indentation.
				bb.WriteRune(r)
			case '\v':
				bb.WriteString(`\v`)
			default:
				switch {
				case r < ' ':
					bb.WriteString(`\x`)
					bb.WriteByte(lowerhex[r>>4&0xF])
					bb.WriteByte(lowerhex[r&0xF])
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

func needsQuoting(s string) bool {
	for i := 0; i < len(s); {
		b := s[i]
		if b < utf8.RuneSelf {
			if needsQuotingSet[b] {
				return true
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError || unicode.IsSpace(r) || !unicode.IsPrint(r) {
			return true
		}
		i += size
	}
	return false
}

var needsQuotingSet = [utf8.RuneSelf]bool{
	'"': true,
	'=': true,
}

func init() {
	for i := 0; i < utf8.RuneSelf; i++ {
		r := rune(i)
		if unicode.IsSpace(r) || !unicode.IsPrint(r) {
			needsQuotingSet[i] = true
		}
	}
}

func writeSpace(w io.Writer, first bool) {
	if !first {
		_, _ = w.Write([]byte{' '})
	}
}

func (l *Logger) textFormatter(keyvals ...any) {
	st := l.styles
	lenKeyvals := len(keyvals)

	for i := 0; i < lenKeyvals; i += 2 {
		firstKey := i == 0
		moreKeys := i < lenKeyvals-2

		switch keyvals[i] {
		case TimestampKey:
			if t, ok := keyvals[i+1].(time.Time); ok {
				ts := t.Format(l.timeFormat)
				ts = st.Timestamp.Render(ts)
				writeSpace(&l.b, firstKey)
				l.b.WriteString(ts)
			}
		case LevelKey:
			if level, ok := keyvals[i+1].(Level); ok {
				var lvl string
				lvlStyle, ok := st.Levels[level]
				if !ok {
					continue
				}

				lvl = lvlStyle.String()
				if lvl != "" {
					writeSpace(&l.b, firstKey)
					l.b.WriteString(lvl)
				}
			}
		case CallerKey:
			if caller, ok := keyvals[i+1].(string); ok {
				caller = fmt.Sprintf("<%s>", caller)
				caller = st.Caller.Render(caller)
				writeSpace(&l.b, firstKey)
				l.b.WriteString(caller)
			}
		case PrefixKey:
			if prefix, ok := keyvals[i+1].(string); ok {
				prefix = st.Prefix.Render(prefix + ":")
				writeSpace(&l.b, firstKey)
				l.b.WriteString(prefix)
			}
		case MessageKey:
			if msg := keyvals[i+1]; msg != nil {
				m := normalizeLogNewlines(fmt.Sprint(msg))
				m = renderKeepTabs(st.Message, m)
				writeSpace(&l.b, firstKey)
				l.b.WriteString(m)
			}
		default:
			sep := separator
			indentSep := indentSeparator
			sep = st.Separator.Render(sep)
			indentSep = st.Separator.Render(indentSep)
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
			valueStyle := st.Value
			if vs, ok := st.Values[actualKey]; ok {
				valueStyle = vs
			}
			if keyStyle, ok := st.Keys[key]; ok {
				key = keyStyle.Render(key)
			} else {
				key = st.Key.Render(key)
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
			} else if !raw && needsQuoting(val) {
				writeSpace(&l.b, firstKey)
				l.b.WriteString(key)
				l.b.WriteString(sep)
				l.b.WriteString(renderKeepTabs(valueStyle, fmt.Sprintf(`"%s"`,
					escapeStringForOutput(val, true))))
			} else {
				val = renderKeepTabs(valueStyle, val)
				writeSpace(&l.b, firstKey)
				l.b.WriteString(key)
				l.b.WriteString(sep)
				l.b.WriteString(val)
			}
		}
	}

	// Add a newline to the end of the log message.
	l.b.WriteByte('\n')
}
