package loggers

import (
	"log/slog"
	"os"
	"sync"
)

// Declare a global logger instance and a mutex to ensure thread-safety
var (
	loggerInstance *slog.Logger
	once           sync.Once
)

// InitializeLogger initializes and returns a singleton logger with a JSON handler
func InitializeLogger() *slog.Logger {
	// Use sync.Once to ensure the logger is created only once
	once.Do(func() {
		// Create the logger only if it hasn't been created yet
		loggerInstance = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	})
	return loggerInstance
}
