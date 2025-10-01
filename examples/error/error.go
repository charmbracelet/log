package main

import (
	"fmt"

	"github.com/charmbracelet/log/v2"
)

func main() {
	err := fmt.Errorf("too much sugar")
	log.Error("failed to bake cookies", "err", err)
}
