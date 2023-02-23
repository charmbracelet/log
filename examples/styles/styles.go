package main

import (
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func main() {
	log.ErrorLevelStyle = lipgloss.NewStyle().
		SetString("ERROR!!").
		Padding(0, 1, 0, 1).
		Background(lipgloss.AdaptiveColor{
			Light: "203",
			Dark:  "204",
		}).
		Foreground(lipgloss.Color("0"))
	log.KeyStyles["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	logger := log.New()
	logger.Error("Whoops!", "err", "kitchen on fire")
	time.Sleep(3 * time.Second)
}
