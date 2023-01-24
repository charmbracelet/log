package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type cup int

func (c cup) String() string {
	s := fmt.Sprintf("%d cup", c)
	if c > 1 {
		s += "s"
	}
	return s
}

func startOven(degree int) {
	log.Helper()
	log.Debug("Starting oven", "temprature", degree)
}

func main() {
	log.SetTimeFormat(time.Kitchen)
	log.SetLevel(log.DebugLevel)

	var (
		butter    = cup(1)
		chocolate = cup(2)
		flour     = cup(3)
		sugar     = cup(5)
		temp      = 375
		bakeTime  = 10
	)

	startOven(temp)
	time.Sleep(time.Second)
	log.Debug("Mixing ingredients", "ingredients",
		strings.Join([]string{
			"butter " + butter.String(),
			"chocolate " + chocolate.String(),
			"flour " + flour.String(),
			"sugar " + sugar.String(),
		}, "\n"),
	)
	time.Sleep(time.Second)
	if sugar > 2 {
		log.Warn("That's a lot of sugar", "amount", sugar)
	}
	log.Info("Baking cookies", "time", fmt.Sprintf("%d minutes", bakeTime))
	time.Sleep(2 * time.Second)
	log.Info("Increasing temprature", "amount", 300)
	temp += 300
	time.Sleep(time.Second)
	if temp > 500 {
		log.Error("Oven is too hot", "temprature", temp)
		os.Exit(1)
	}
}
