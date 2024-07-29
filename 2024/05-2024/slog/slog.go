package main

import (
	"log"
	"log/slog"
	"os"
)

func main() {
	log.Print("Info message")
	slog.Info("Info message")
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	log.Println("Hello from old logger")
}
