package main

import "charm.land/log/v2"

func startOven(degree int) {
	log.Helper()
	log.Info("Starting oven", "degree", degree)
}

func main() {
	log.SetReportCaller(true)
	startOven(400)
}
