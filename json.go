package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (l *Logger) jsonFormatter(keyvals ...interface{}) {
	jw := &jsonWriter{w: &l.b, r: l.re, s: l.styles.Separator}
	jw.start()

	i := 0
	for i < len(keyvals) {
		switch kv := keyvals[i].(type) {
		case slogAttr:
			l.jsonFormatterRoot(jw, kv.Key, kv.Value)
			i++
		default:
			if i+1 < len(keyvals) {
				l.jsonFormatterRoot(jw, keyvals[i], keyvals[i+1])
			}
			i += 2
		}
	}

	jw.end()
	l.b.WriteRune('\n')
}

func (l *Logger) jsonFormatterRoot(jw *jsonWriter, key, value any) {
	switch key {
	case TimestampKey:
		if t, ok := value.(time.Time); ok {
			jw.objectItem(l.styles.Key, TimestampKey, l.styles.Timestamp, t.Format(l.timeFormat))
		}
	case LevelKey:
		if level, ok := value.(Level); ok {
			ls, ok := l.styles.Levels[level]
			if ok {
				jw.objectItem(l.styles.Key, LevelKey, ls, level.String())
			}
		}
	case CallerKey:
		if caller, ok := value.(string); ok {
			jw.objectItem(l.styles.Key, CallerKey, l.styles.Caller, caller)
		}
	case PrefixKey:
		if prefix, ok := value.(string); ok {
			jw.objectItem(l.styles.Key, PrefixKey, l.styles.Prefix, prefix)
		}
	case MessageKey:
		if msg := value; msg != nil {
			jw.objectItem(l.styles.Key, MessageKey, l.styles.Message, fmt.Sprint(msg))
		}
	default:
		l.jsonFormatterItem(jw, 0, l.styles.Key, key, l.styles.Value, value)
	}
}

func (l *Logger) jsonFormatterItem(
	jw *jsonWriter, d int, ks lipgloss.Style, anyKey any, vs lipgloss.Style, value any,
) {
	var key string
	switch k := anyKey.(type) {
	case fmt.Stringer:
		key = k.String()
	case error:
		key = k.Error()
	default:
		key = fmt.Sprint(k)
	}

	// override styles based on root key
	if d == 0 {
		if s, ok := l.styles.Keys[key]; ok {
			ks = s
		}
		if s, ok := l.styles.Values[key]; ok {
			vs = s
		}
	}

	jw.objectKey(ks, key)

	switch v := value.(type) {
	case error:
		jw.objectValue(vs, v.Error())
	case slogLogValuer:
		l.writeSlogValue(jw, d, ks, vs, v.LogValue())
	case slogValue:
		l.writeSlogValue(jw, d, ks, vs, v.Resolve())
	case fmt.Stringer:
		jw.objectValue(vs, v.String())
	default:
		jw.objectValue(vs, v)
	}
}

func (l *Logger) writeSlogValue(jw *jsonWriter, depth int, ks, vs lipgloss.Style, v slogValue) {
	switch v.Kind() {
	case slogKindGroup:
		jw.start()
		for _, attr := range v.Group() {
			l.jsonFormatterItem(jw, depth+1, ks, attr.Key, vs, attr.Value)
		}
		jw.end()
	default:
		jw.objectValue(vs, v.Any())
	}
}

type jsonWriter struct {
	w *bytes.Buffer
	r *lipgloss.Renderer
	s lipgloss.Style
	d int
}

func (w *jsonWriter) start() {
	objectStart := w.s.Renderer(w.r).Render("{")
	w.w.WriteString(objectStart)
	w.d = 0
}

func (w *jsonWriter) end() {
	objectEnd := w.s.Renderer(w.r).Render("}")
	w.w.WriteString(objectEnd)
}

func (w *jsonWriter) objectItem(
	ks lipgloss.Style, key string,
	vs lipgloss.Style, value any,
) {
	w.objectKey(ks, key)
	w.objectValue(vs, value)
}

func (w *jsonWriter) objectKey(s lipgloss.Style, key string) {
	if w.d > 0 {
		itemSep := w.s.Renderer(w.r).Render(",")
		w.w.WriteString(itemSep)
	}
	w.d++

	pos := w.w.Len()
	err := w.writeEncoded(key)
	if err != nil {
		w.w.Truncate(pos)
		w.w.WriteString(`"invalid key"`)
	}

	// re-apply value with style
	w.renderStyle(s, pos)

	valSep := w.s.Renderer(w.r).Render(`:`)
	w.w.WriteString(valSep)
}

func (w *jsonWriter) objectValue(s lipgloss.Style, value any) {
	pos := w.w.Len()
	err := w.writeEncoded(value)
	if err != nil {
		w.w.Truncate(pos)
		w.w.WriteString(`"invalid value"`)
	}

	// re-apply value with style
	w.renderStyle(s, pos)
}

// renderStyle applies the given style to the string at the given position.
func (w *jsonWriter) renderStyle(st lipgloss.Style, pos int) {
	s := w.w.String()[pos:]

	// manually apply quotes
	sep := ""
	if len(s) > 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1] // apply style within quotes
		sep = w.s.Renderer(w.r).Render(`"`)
	} else if st.String() != "" {
		sep = w.s.Renderer(w.r).Render(`"`)
	}

	// render with style
	s = st.Renderer(w.r).Render(s)

	// rewind
	w.w.Truncate(pos)

	// re-apply with colors
	w.w.WriteString(sep)
	w.w.WriteString(s)
	w.w.WriteString(sep)
}

func (w *jsonWriter) writeEncoded(v any) error {
	e := json.NewEncoder(w.w)
	e.SetEscapeHTML(false)
	if err := e.Encode(v); err != nil {
		return err
	}

	// trailing \n added by json.Encode
	b := w.w.Bytes()
	if len(b) > 0 && b[len(b)-1] == '\n' {
		w.w.Truncate(w.w.Len() - 1)
	}
	return nil
}
