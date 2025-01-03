package logger

import (
	"log/slog"
	"os"
)

func New(isProd bool) *slog.Logger {
	var logger *slog.Logger

	if isProd {
		optsProd := &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
		logger = slog.New(slog.NewJSONHandler(os.Stdout, optsProd))

	}

	if !isProd {
		optsDev := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		logger = slog.New(slog.NewTextHandler(os.Stdout, optsDev))
	}

	return logger
}
