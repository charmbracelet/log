package log

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogContext_empty(t *testing.T) {
	require.Equal(t, Default(), FromContext(context.TODO()))
}

func TestLogContext_simple(t *testing.T) {
	l := New(io.Discard)
	ctx := WithContext(context.Background(), l)
	require.Equal(t, l, FromContext(ctx))
}

func TestLogContext_fields(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf)
	l.SetLevel(DebugLevel)
	ctx := WithContext(context.Background(), l.With("foo", "bar"))
	l = FromContext(ctx)
	require.NotNil(t, l)
	l.Debug("test")
	require.Equal(t, "DEBU test foo=bar\n", buf.String())
}
