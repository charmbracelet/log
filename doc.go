// Package log provides a simple and flexible logger with support for
// structured logging, log levels, and customizable output formats.
//
// # Concurrency
//
// A single Logger instance is safe for concurrent use from multiple
// goroutines. However, when multiple Logger instances share the same
// io.Writer, the writer itself must be safe for concurrent use.
// Use [NewSyncWriter] to wrap a non-thread-safe writer.
package log

import (
	"io"
	"sync"
)

// SyncWriter wraps an io.Writer with a mutex to make it safe for concurrent use.
type SyncWriter struct {
	mu sync.Mutex
	w  io.Writer
}

// NewSyncWriter returns a new SyncWriter that serializes writes to w.
// Use this when sharing a single writer between multiple Logger instances.
func NewSyncWriter(w io.Writer) *SyncWriter {
	return &SyncWriter{w: w}
}

// Write implements io.Writer.
func (sw *SyncWriter) Write(p []byte) (n int, err error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.w.Write(p)
}
