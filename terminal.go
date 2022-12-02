package log

import (
	"io"

	"github.com/mattn/go-isatty"
)

// IsTerminal returns true if w writes to a terminal.
func IsTerminal(w io.Writer) bool {
	fw, ok := w.(interface {
		Fd() uintptr
	})
	if !ok {
		return false
	}
	return isatty.IsTerminal(fw.Fd())
}
