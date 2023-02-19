package main

import "github.com/charmbracelet/log"

func main() {
	logger := log.New()
	logger.Warn("chewy!", "butter", true)
}
