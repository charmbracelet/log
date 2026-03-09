package main

import (
	"fmt"

	"charm.land/log/v2"
)

func main() {
	err := fmt.Errorf("too much sugar")
	log.Error("failed to bake cookies", "err", err)
}
