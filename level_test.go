package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLevel(t *testing.T) {
	var level Level
	assert.Equal(t, InfoLevel, level)
}
func TestParse(t *testing.T) {
	testCases := []struct {
		name     string
		level    string
		expLevel Level
	}{
		{
			name:     "Parse debug",
			level:    "debug",
			expLevel: DebugLevel,
		},
		{
			name:     "Parse info",
			level:    "Info",
			expLevel: InfoLevel,
		},
		{
			name:     "Parse warn",
			level:    "WARN",
			expLevel: WarnLevel,
		},
		{
			name:     "Parse error",
			level:    "error",
			expLevel: ErrorLevel,
		},
		{
			name:     "Parse fatal",
			level:    "FATAL",
			expLevel: FatalLevel,
		},
		{
			name:     "Default",
			level:    "",
			expLevel: InfoLevel,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lvl := ParseLevel(tc.level)
			assert.Equal(t, tc.expLevel, lvl)
		})
	}
}
