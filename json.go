package log

import (
	"encoding/json"
	"fmt"
	"time"
)

func (l *Logger) jsonFormatter(keyvals ...interface{}) {
	m := make(map[string]interface{}, len(keyvals)/2)
	for i := 0; i < len(keyvals); i += 2 {
		switch keyvals[i] {
		case TimestampKey:
			if t, ok := keyvals[i+1].(time.Time); ok {
				m[TimestampKey] = t.Format(l.timeFormat)
			}
		case LevelKey:
			if level, ok := keyvals[i+1].(Level); ok {
				m[LevelKey] = level.String()
			}
		case CallerKey:
			if caller, ok := keyvals[i+1].(string); ok {
				m[CallerKey] = caller
			}
		case PrefixKey:
			if prefix, ok := keyvals[i+1].(string); ok {
				m[PrefixKey] = prefix
			}
		case MessageKey:
			if msg := keyvals[i+1]; msg != nil {
				m[MessageKey] = fmt.Sprint(msg)
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
			m[key] = val
		}
	}

	e := json.NewEncoder(&l.b)
	e.SetEscapeHTML(false)
	_ = e.Encode(m)
}
