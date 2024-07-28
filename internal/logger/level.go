package logger

import (
	"context"
	"log/slog"
)

func GetLevel(ctx context.Context, log *slog.Logger) slog.Level {
	result := slog.LevelError

	levels := []slog.Level{slog.LevelWarn, slog.LevelInfo, slog.LevelDebug}
	for _, lvl := range levels {
		if log.Enabled(ctx, lvl) {
			result = lvl
		}
	}

	return result
}
