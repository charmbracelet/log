package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/log/v2"
)

func main() {
	// baseline
	fmt.Println(time.Now().UTC().Format(time.RFC3339), "foo")
	fmt.Println(time.Now().Format(time.RFC3339), "bar")

	handler := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
		TimeFunction:    log.NowUTC,
		TimeFormat:      time.RFC3339,
	})
	handler.Info("foobar")

	logger := slog.New(handler)
	logger.Info("foobar")
}
