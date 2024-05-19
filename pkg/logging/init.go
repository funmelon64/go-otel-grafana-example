package logging

import (
	"io"
	"log/slog"
	"os"
)

func InitLogging(cfg Config) {
	var logOut io.Writer = os.Stdout

	if cfg.FileOut != "" {
		file, err := os.OpenFile(cfg.FileOut, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			logOut = file
		} else {
			slog.Error("Failed to log to file, using default stdout")
		}
	}

	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(logOut, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}))

	slog.SetDefault(logger)
}
