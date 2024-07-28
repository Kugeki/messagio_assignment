package rest

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

func LogMiddleware(log *slog.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/log"),
		)
		log.Info("logger middleware enabled")

		handleFunc := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			entry.Info("request received")

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			before := time.Now()

			next.ServeHTTP(ww, r)

			entry.Info("request completed",
				slog.Int("status", ww.Status()),
				slog.Int("bytes", ww.BytesWritten()),
				slog.String("duration", time.Since(before).String()),
			)
		}

		return http.HandlerFunc(handleFunc)
	}

}
