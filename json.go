package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

func (l *Logger) jsonFormatter(keyvals ...interface{}) {
	jw := jsonWriter{w: &l.b}
	jw.start()
	for i := 0; i < len(keyvals); i += 2 {
		switch keyvals[i] {
		case TimestampKey:
			if t, ok := keyvals[i+1].(time.Time); ok {
				jw.write(TimestampKey, t.Format(l.timeFormat))
			}
		case LevelKey:
			if level, ok := keyvals[i+1].(Level); ok {
				jw.write(LevelKey, level.String())
			}
		case CallerKey:
			if caller, ok := keyvals[i+1].(string); ok {
				jw.write(CallerKey, caller)
			}
		case PrefixKey:
			if prefix, ok := keyvals[i+1].(string); ok {
				jw.write(PrefixKey, prefix)
			}
		case MessageKey:
			if msg := keyvals[i+1]; msg != nil {
				jw.write(MessageKey, fmt.Sprint(msg))
			}
		default:
			var (
				key string
				val interface{}
			)
			switch k := keyvals[i].(type) {
			case fmt.Stringer:
				key = k.String()
			case error:
				key = k.Error()
			default:
				key = fmt.Sprint(k)
			}
			switch v := keyvals[i+1].(type) {
			case error:
				val = v.Error()
			case fmt.Stringer:
				val = v.String()
			default:
				val = v
			}
			jw.write(key, val)
		}
	}
	jw.end()
}

type jsonWriter struct {
	w *bytes.Buffer
}

func (w *jsonWriter) start() {
	w.w.WriteRune('{')
}

func (w *jsonWriter) end() {
	w.w.WriteRune('}')
	w.w.WriteRune('\n')
}

func (w *jsonWriter) write(key string, value any) {
	// store pos if we need to rewind
	pos := w.w.Len()

	// add separator when buffer is longer than '{'
	if w.w.Len() > 1 {
		w.w.WriteRune(',')
	}

	err := w.writeEncoded(key)
	if err != nil {
		w.w.Truncate(pos)
		return
	}
	w.w.WriteRune(':')

	pos = w.w.Len()
	err = w.writeEncoded(value)
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
