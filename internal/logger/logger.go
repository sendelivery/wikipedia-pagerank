package logger

import (
	"log"
	"log/slog"
	"os"
)

// MakeLogger creates a new logger that writes to a `log` file in the current directory.
func MakeLogger(logLevel slog.Level, logfile string) *slog.Logger {
	file, err := os.Create(logfile)
	if err != nil {
		log.Fatalf("Failed to create log file: %s", err)
	}
	options := &slog.HandlerOptions{Level: logLevel}
	handler := slog.NewJSONHandler(file, options)
	return slog.New(handler)
}
