package logger

import (
	"log/slog"
	"os"
)

// Logger is the application-wide structured logger.
type Logger struct {
	*slog.Logger
}

// New creates a logger for the application.
//
// JSON logs are useful for production because logging systems
// can easily search and filter fields from JSON output.
func New(environment string) *Logger {
	handlerOptions := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	var handler slog.Handler

	// Development logs are easier for humans to read.
	// Production logs use JSON so log aggregation tools can
	// process them more easily.
	if environment == "development" {
		handler = slog.NewTextHandler(os.Stdout,handlerOptions,)
	} else {
		handler = slog.NewJSONHandler(os.Stdout,handlerOptions,)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}