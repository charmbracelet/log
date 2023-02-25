package main

import (
	"time"

	"github.com/charmbracelet/log"
)

func main() {
	logger := log.New(log.WithTimestamp(), log.WithTimeFormat(time.Kitchen),
		log.WithCaller(), log.WithPrefix("baking 🍪 ")).With("batch", 2, "chocolateChips", true)
	logger.SetReportTimestamp(false)
	logger.SetReportCaller(false)
	logger.SetLevel(log.DebugLevel)
	logger.Debug("Preparing batch 2...")
	logger.Debug("Adding chocolate chips")
}
