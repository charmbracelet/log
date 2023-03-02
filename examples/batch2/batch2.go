package main

import (
	"github.com/charmbracelet/log"
)

func main() {
	logger := log.Default().With("batch", 2, "chocolateChips", true)
	logger.SetPrefix("baking 🍪 ")
	logger.SetReportTimestamp(false)
	logger.SetReportCaller(false)
	logger.SetLevel(log.DebugLevel)
	logger.Debug("Preparing batch 2...")
	logger.Debug("Adding chocolate chips")
}
