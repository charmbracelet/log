package main

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func main() {
	logger := log.New(os.Stderr)
	logger.SetPrefix("Baking 🍪 ")
	logger.SetTimeFormat(time.Kitchen)
	logger.SetReportTimestamp(true)
	logger.SetReportCaller(true)
	logger.Info("Starting oven!", "degree", 375)
	time.Sleep(3 * time.Second)
	logger.Info("Finished baking")
}
