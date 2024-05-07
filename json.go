package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

func (l *Logger) jsonFormatter(keyvals ...interface{}) {
	jw := &jsonWriter{w: &l.b}
	jw.start()

	for i := 0; i < len(keyvals); i += 2 {
		l.jsonFormatterKeyVal(jw, keyvals[i], keyvals[i+1])
	}

	jw.end()
	l.b.WriteRune('\n')
}

func (l *Logger) jsonFormatterKeyVal(jw *jsonWriter, anyKey, value any) {
	switch anyKey {
	case TimestampKey:
		if t, ok := value.(time.Time); ok {
			jw.objectItem(TimestampKey, t.Format(l.timeFormat))
		}
	case LevelKey:
		if level, ok := value.(Level); ok {
			jw.objectItem(LevelKey, level.String())
		}
	case CallerKey:
		if caller, ok := value.(string); ok {
			jw.objectItem(CallerKey, caller)
		}
	case PrefixKey:
		if prefix, ok := value.(string); ok {
			jw.objectItem(PrefixKey, prefix)
		}
	case MessageKey:
		if msg := value; msg != nil {
			jw.objectItem(MessageKey, fmt.Sprint(msg))
		}
	default:
		switch k := anyKey.(type) {
		case fmt.Stringer:
			jw.objectKey(k.String())
		case error:
			jw.objectKey(k.Error())
		default:
			jw.objectKey(fmt.Sprint(k))
		}
		switch v := value.(type) {
		case error:
			jw.objectValue(v.Error())
		case fmt.Stringer:
			jw.objectValue(v.String())
		default:
			jw.objectValue(v)
		}
	}
}

type jsonWriter struct {
	w *bytes.Buffer
	d int
}

func (w *jsonWriter) start() {
	w.w.WriteRune('{')
	w.d = 0
}

func (w *jsonWriter) end() {
	w.w.WriteRune('}')
}

func (w *jsonWriter) objectItem(key string, value any) {
	w.objectKey(key)
	w.objectValue(value)
}

func (w *jsonWriter) objectKey(key string) {
	if w.d > 0 {
		w.w.WriteRune(',')
	}
	w.d++

	pos := w.w.Len()
	err := w.writeEncoded(key)
	if err != nil {
		w.w.Truncate(pos)
		w.w.WriteString(`"invalid key"`)
	}
	w.w.WriteRune(':')
}

func (w *jsonWriter) objectValue(value any) {
	pos := w.w.Len()
	err := w.writeEncoded(value)
	if err != nil {
		w.w.Truncate(pos)
		w.w.WriteString(`"invalid value"`)
	}
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
