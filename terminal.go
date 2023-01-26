package log

import (
	"io"
	"os"

	"github.com/mattn/go-isatty"
)

func isTerminal(w io.Writer) bool {
	if len(os.Getenv("CI")) > 0 {
		return false
	}

	if f, ok := w.(interface {
		Fd() uintptr
	}); ok {
		return isatty.IsTerminal(f.Fd())
	}

	return false
}
