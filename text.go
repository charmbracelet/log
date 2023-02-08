package log

import (
	"fmt"
	"strings"
	"time"
)

const (
	separator       = "="
	indentSeparator = "  │ "
)

func (l *logger) textFormatter(keyvals ...interface{}) {
	s := l.styles
	for i := 0; i < len(keyvals); i += 2 {
		switch keyvals[i].(string) {
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
			if level, ok := keyvals[i+1].(Level); ok && level != noLevel {
				lvl := strings.ToUpper(level.String())
				if !l.noStyles {
					lvl = s.Level(level).String()
				}
				l.b.WriteString(lvl)
				l.b.WriteByte(' ')
			}
		case callerKey:
			if caller, ok := keyvals[i+1].(string); ok {
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
			val := fmt.Sprint(keyvals[i+1])
			raw := val == ""
			if raw {
				val = `""`
			}
			if key == "" {
				key = "MISSING_KEY"
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
