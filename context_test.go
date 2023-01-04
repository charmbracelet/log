package log

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogContext_empty(t *testing.T) {
	require.Equal(t, Default(), FromContext(context.TODO()))
}

func TestLogContext_simple(t *testing.T) {
	l := New()
	ctx := WithContext(context.Background(), l)
	require.Equal(t, l, FromContext(ctx))
}

func TestLogContext_fields(t *testing.T) {
	var buf bytes.Buffer
	l := New(WithOutput(&buf), WithLevel(LevelDebug), WithNoStyles())
	ctx := WithContext(context.Background(), l, "foo", "bar")
	l = FromContext(ctx)
	require.NotNil(t, l)
	l.Debug("test")
	require.Equal(t, "DEBUG test foo=bar\n", buf.String())
}
