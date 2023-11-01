package log

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLevel(t *testing.T) {
	var level Level
	assert.Equal(t, InfoLevel, level)
}

func TestParseLevel(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		result Level
		error  error
	}{
		{
			name:   "Parse debug",
			input:  "debug",
			result: DebugLevel,
			error:  nil,
		},
		{
			name:   "Parse info",
			input:  "Info",
			result: InfoLevel,
			error:  nil,
		},
		{
			name:   "Parse warn",
			input:  "WARN",
			result: WarnLevel,
			error:  nil,
		},
		{
			name:   "Parse error",
			input:  "error",
			result: ErrorLevel,
			error:  nil,
		},
		{
			name:   "Parse fatal",
			input:  "FATAL",
			result: FatalLevel,
			error:  nil,
		},
		{
			name:   "Default",
			input:  "",
			result: InfoLevel,
			error:  fmt.Errorf("%w: %q", errors.New("invalid level"), ""),
		},
		{
			name:   "Wrong level, set INFO",
			input:  "WRONG_LEVEL",
			result: InfoLevel,
			error:  fmt.Errorf("%w: %q", errors.New("invalid level"), "WRONG_LEVEL"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lvl, err := ParseLevel(tc.input)
			assert.Equal(t, tc.result, lvl)
			assert.Equal(t, tc.error, err)
		})
	}
}
