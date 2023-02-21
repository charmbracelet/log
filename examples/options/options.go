package main

import (
	"time"

	"github.com/charmbracelet/log"
)

func main() {
	logger := log.New(log.WithTimestamp(), log.WithTimeFormat(time.Kitchen),
		log.WithCaller(), log.WithPrefix("Baking 🍪 "))
	logger.Info("Starting oven!", "degree", 375)
	time.Sleep(3 * time.Second)
	logger.Info("Finished baking")
}
