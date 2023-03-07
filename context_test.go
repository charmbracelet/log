package log

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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
	l := New(WithOutput(&buf), WithLevel(DebugLevel))
	ctx := WithContext(context.Background(), l, "foo", "bar")
	l = FromContext(ctx)
	require.NotNil(t, l)
	l.Debug("test")
	require.Equal(t, "DEBUG test foo=bar\n", buf.String())
}

func TestLogUpdateContext(t *testing.T) {
	var buf bytes.Buffer
	logger := New(WithOutput(&buf))
	logger.SetLevel(DebugLevel)
	t.Run("with extra parameters", func(t *testing.T) {
		ctx := WithContext(context.Background(), logger, "key", "value")

		assert.True(t, UpdateContext(ctx, func(logger Logger) Logger {
			return logger.With("key2", "value2")
		}))

		buf.Reset()
		FromContext(ctx).Info("test")
		assert.Equal(t, "INFO test key=value key2=value2\n", buf.String())

		// original logger should not be changed
		buf.Reset()
		logger.Debug("test2")
		assert.Equal(t, "DEBUG test2\n", buf.String())
	})

	t.Run("without extra parameters", func(t *testing.T) {
		ctx := WithContext(context.Background(), logger)

		assert.True(t, UpdateContext(ctx, func(logger Logger) Logger {
			return logger.With("key3", "value3")
		}))

		buf.Reset()
		FromContext(ctx).Info("test3")
		assert.Equal(t, "INFO test3 key3=value3\n", buf.String())

		// original logger should not be changed
		buf.Reset()
		logger.Debug("test4")
		assert.Equal(t, "DEBUG test4\n", buf.String())
	})

	t.Run("with consecutive calls", func(t *testing.T) {
		ctx := WithContext(context.Background(), logger)

		assert.True(t, UpdateContext(ctx, func(logger Logger) Logger {
			return logger.With("key4", "value4")
		}))
		assert.True(t, UpdateContext(ctx, func(logger Logger) Logger {
			return logger.With("key5", "value5")
		}))

		buf.Reset()
		FromContext(ctx).Info("test5")
		assert.Equal(t, "INFO test5 key4=value4 key5=value5\n", buf.String())
	})

	t.Run("without log in context", func(t *testing.T) {
		assert.False(t, UpdateContext(context.Background(), func(logger Logger) Logger {
			return logger
		}))
	})
}
