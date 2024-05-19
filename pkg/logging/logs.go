package logging

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

func logAttrs(ctx context.Context, logger *slog.Logger, level slog.Level, msg string, attrs []slog.Attr) {
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip: [Callers, this func, log wrapper func]
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.AddAttrs(attrs...)
	_ = logger.Handler().Handle(ctx, r)
}

func Info(msg string, attrs ...slog.Attr) {
	logAttrs(context.Background(), slog.Default(), slog.LevelInfo, msg, attrs)
}

func Warn(msg string, attrs ...slog.Attr) {
	logAttrs(context.Background(), slog.Default(), slog.LevelWarn, msg, attrs)
}

func Error(msg string, attrs ...slog.Attr) {
	logAttrs(context.Background(), slog.Default(), slog.LevelError, msg, attrs)
}

func ErrorErr(msg string, err error, attrs ...slog.Attr) {
	attrs = append(attrs, slog.String("error", err.Error()))
	logAttrs(context.Background(), slog.Default(), slog.LevelError, msg, attrs)
}

func Debug(msg string, attrs ...slog.Attr) {
	logAttrs(context.Background(), slog.Default(), slog.LevelDebug, msg, attrs)
}
