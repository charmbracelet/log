package log

import (
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
		name     string
		input    string
		expected Level
		error    error
	}{
		{
			name:     "Parse debug",
			input:    "debug",
			expected: DebugLevel,
			error:    nil,
		},
		{
			name:     "Parse info",
			input:    "Info",
			expected: InfoLevel,
			error:    nil,
		},
		{
			name:     "Parse warn",
			input:    "WARN",
			expected: WarnLevel,
			error:    nil,
		},
		{
			name:     "Parse error",
			input:    "error",
			expected: ErrorLevel,
			error:    nil,
		},
		{
			name:     "Parse fatal",
			input:    "FATAL",
			expected: FatalLevel,
			error:    nil,
		},
		{
			name:     "Default",
			input:    "",
			expected: InfoLevel,
			error:    fmt.Errorf("%w: %q", ErrInvalidLevel, ""),
		},
		{
			name:     "Wrong level, set INFO",
			input:    "WRONG_LEVEL",
			expected: InfoLevel,
			error:    fmt.Errorf("%w: %q", ErrInvalidLevel, "WRONG_LEVEL"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lvl, err := ParseLevel(tc.input)
			assert.Equal(t, tc.expected, lvl)
			assert.Equal(t, tc.error, err)
		})
	}
}
