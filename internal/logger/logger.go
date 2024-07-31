package logger

import (
	"log/slog"
	"os"
)

type Environment interface {
	IsDev() bool
	IsProd() bool
	String() string
}

var DefaultWriter = os.Stdout

func New(environment Environment, level slog.Level) *slog.Logger {
	return slog.New(
		NewHandler(environment, &slog.HandlerOptions{Level: level}),
	).With(slog.String("env", environment.String()))
}

func NewHandler(env Environment, opts *slog.HandlerOptions) slog.Handler {
	if env.IsDev() {
		return slog.NewTextHandler(DefaultWriter, opts)
	}
	if env.IsProd() {
		return slog.NewJSONHandler(DefaultWriter, opts)
	}

	return slog.NewJSONHandler(DefaultWriter, opts)
}
