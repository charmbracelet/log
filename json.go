package log

import (
	"encoding/json"
	"fmt"
	"time"
)

func (l *logger) jsonFormatter(keyvals ...interface{}) {
	m := make(map[string]interface{}, len(keyvals)/2)
	for i := 0; i < len(keyvals); i += 2 {
		switch keyvals[i].(string) {
		case tsKey:
			if t, ok := keyvals[i+1].(time.Time); ok {
				m[tsKey] = t.Format(l.timeFormat)
			}
		case lvlKey:
			if level, ok := keyvals[i+1].(Level); ok && level != noLevel {
				m[lvlKey] = level.String()
			}
		case callerKey:
			if caller, ok := keyvals[i+1].(string); ok {
				m[callerKey] = caller
			}
		case prefixKey:
			if prefix, ok := keyvals[i+1].(string); ok {
				m[prefixKey] = prefix
			}
		case msgKey:
			if msg := keyvals[i+1]; msg != nil {
				m[msgKey] = fmt.Sprint(msg)
			}
		default:
			var k string
			if key, ok := keyvals[i].(string); ok {
				k = key
			} else {
				k = fmt.Sprint(keyvals[i])
			}
			m[k] = keyvals[i+1]
		}
	}

	json.NewEncoder(l.w).SetEscapeHTML(false)
	_ = json.NewEncoder(l.w).Encode(m)
}
