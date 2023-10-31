package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLevel(t *testing.T) {
	var level Level
	assert.Equal(t, InfoLevel, level)
}

type parseLevelResult struct {
	level Level
	err   error
}

func TestParseLevel(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		result parseLevelResult
	}{
		{
			name:   "Parse debug",
			input:  "debug",
			result: parseLevelResult{DebugLevel, nil},
		},
		{
			name:   "Parse info",
			input:  "Info",
			result: parseLevelResult{InfoLevel, nil},
		},
		{
			name:   "Parse warn",
			input:  "WARN",
			result: parseLevelResult{WarnLevel, nil},
		},
		{
			name:   "Parse error",
			input:  "error",
			result: parseLevelResult{ErrorLevel, nil},
		},
		{
			name:   "Parse fatal",
			input:  "FATAL",
			result: parseLevelResult{FatalLevel, nil},
		},
		{
			name:   "Default",
			input:  "",
			result: parseLevelResult{InfoLevel, parseLevelError("")},
		},
		{
			name:   "Wrong level, set INFO",
			input:  "WRONG_LEVEL",
			result: parseLevelResult{InfoLevel, parseLevelError("WRONG_LEVEL")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lvl, err := ParseLevel(tc.input)
			assert.Equal(t, tc.result, parseLevelResult{lvl, err})
		})
	}
}
