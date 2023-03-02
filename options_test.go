package log

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptions(t *testing.T) {
	opts := Options{
		Level:        ErrorLevel,
		ReportCaller: true,
		Fields:       []interface{}{"foo", "bar"},
	}
	logger := NewWithOptions(ioutil.Discard, opts)
	require.Equal(t, ErrorLevel, logger.GetLevel())
	require.True(t, logger.reportCaller)
	require.False(t, logger.reportTimestamp)
	require.Equal(t, []interface{}{"foo", "bar"}, logger.fields)
	require.Equal(t, TextFormatter, logger.formatter)
	require.Equal(t, DefaultTimeFormat, logger.timeFormat)
	require.NotNil(t, logger.timeFunc)
}
