package main

import "github.com/charmbracelet/log"

func startOven(degree int) {
	log.Helper()
	log.Info("Starting oven", "degree", degree)
}

func main() {
	log.SetReportCaller(true)
	startOven(400)
}
