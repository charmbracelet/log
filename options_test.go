package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithCallerFormat(t *testing.T) {
	l := New(WithCallerFormat(CallerLong)).(*logger)
	assert.Equal(t, CallerLong, l.callerFormat)
}
