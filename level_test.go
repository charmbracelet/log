package log

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLevel(t *testing.T) {
	var level Level
	assert.Equal(t, InfoLevel, level)
}
func TestParseLevel(t *testing.T) {
	testCases := []struct {
		name        string
		level       string
		envVarLevel string
		expLevel    Level
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
		{
			name:        "From os.EnvVar",
			envVarLevel: "DEBUG",
			expLevel:    DebugLevel,
		},
		{
			name:     "Wrong level, set INFO",
			level:    "WRONG_LEVEL",
			expLevel: InfoLevel,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var lvl Level
			if tc.envVarLevel != "" {
				t.Setenv("LOG_LEVEL", tc.envVarLevel)
				lvl = ParseLevel(os.Getenv("LOG_LEVEL"))
			} else {
				lvl = ParseLevel(tc.level)
			}
			assert.Equal(t, tc.expLevel, lvl)
		})
	}
}
