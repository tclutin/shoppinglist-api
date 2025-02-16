package logger

import (
	"log/slog"
	"os"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Warn(msg string, args ...any)
	Debug(msg string, args ...any)
	With(args ...any) Logger
}

type SLogger struct {
	logger *slog.Logger
}

func New(isProd bool) Logger {
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

	return &SLogger{logger: logger}
}
func (l *SLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *SLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *SLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *SLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *SLogger) With(args ...any) Logger {
	return &SLogger{logger: l.logger.With(args...)}
}
