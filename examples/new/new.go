package main

import (
	"os"

	"charm.land/log/v2"
)

func main() {
	logger := log.New(os.Stderr)
	logger.Warn("chewy!", "butter", true)
}
