package logging

import (
	"context"
	"log/slog"
)

var defaultLogger = &Logger{logger: slog.Default()}

func GetDefault() *Logger {
	return defaultLogger
}

type Logger struct {
	logger *slog.Logger
}

func (l *Logger) Info(msg string, attrs ...slog.Attr) {
	logAttrs(context.Background(), l.logger, slog.LevelInfo, msg, attrs)
}

func (l *Logger) Warn(msg string, attrs ...slog.Attr) {
	logAttrs(context.Background(), l.logger, slog.LevelWarn, msg, attrs)
}

func (l *Logger) Error(msg string, attrs ...slog.Attr) {
	logAttrs(context.Background(), l.logger, slog.LevelError, msg, attrs)
}

func (l *Logger) ErrorErr(msg string, err error, attrs ...slog.Attr) {
	attrs = append(attrs, slog.String("error", err.Error()))
	logAttrs(context.Background(), l.logger, slog.LevelError, msg, attrs)
}

func (l *Logger) Debug(msg string, attrs ...slog.Attr) {
	logAttrs(context.Background(), l.logger, slog.LevelDebug, msg, attrs)
}

func (l *Logger) With(attr slog.Attr) *Logger {
	return &Logger{logger: l.logger.With(attr)}
}

func NewWith(attr slog.Attr) *Logger {
	return &Logger{logger: defaultLogger.logger.With(attr)}
}
