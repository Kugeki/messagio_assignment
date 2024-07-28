package logger

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

func TestGetLevel(t *testing.T) {
	levels := []slog.Level{
		slog.LevelDebug, slog.LevelInfo,
		slog.LevelWarn, slog.LevelError,
	}
	for _, lvl := range levels {
		slogger := NewSlogLogger(lvl)

		gotLvl := GetLevel(context.Background(), slogger)
		assert.Equal(t, lvl, gotLvl)
	}
}

func NewSlogLogger(level slog.Level) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}
