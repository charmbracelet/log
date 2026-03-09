package main

import (
	"fmt"
	"time"

	"charm.land/log/v2"
)

func main() {
	for temp := 375; temp <= 400; temp++ {
		log.Info("Increasing temperature", "degree", fmt.Sprintf("%d°F", temp))
		time.Sleep(100 * time.Millisecond)
	}
}
