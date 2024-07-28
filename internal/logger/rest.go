package logger

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

func ForRest(logger *slog.Logger, handler string, reqCtx context.Context) *slog.Logger {
	return logger.With(
		slog.String("handler", handler),
		slog.String("request_id", middleware.GetReqID(reqCtx)),
	)
}
