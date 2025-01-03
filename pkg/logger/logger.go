package logger

import "log/slog"

type Logger struct {
	IsProd bool
	logger *slog.Logger
}

func New() *slog.Logger {
}
