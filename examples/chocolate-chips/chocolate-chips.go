package main

import (
	"charm.land/log/v2"
)

func main() {
	logger := log.Default().With()

	logger.SetPrefix("Baking 🍪 ")
	logger.SetReportTimestamp(false)
	logger.SetReportCaller(false)
	logger.SetLevel(log.DebugLevel)
	logger.Debug("Preparing batch 2...") // DEBUG baking 🍪: Preparing batch 2...}

	batch2 := logger.With("batch", 2, "chocolateChips", true)
	batch2.Debug("Adding chocolate chips")
}
