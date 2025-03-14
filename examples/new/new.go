package main

import (
	"os"

	"github.com/charmbracelet/log/v2"
)

func main() {
	logger := log.New(os.Stderr)
	logger.Warn("chewy!", "butter", true)
}
