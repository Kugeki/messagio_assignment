package rest

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"messagio_assignment/internal/config"
	"messagio_assignment/internal/domain/message"
	"messagio_assignment/internal/usecases"
	"net/http"
)

func NewServer(cfg config.Config, msgRepo message.Repository, log *slog.Logger) *http.Server {
	msgUC := usecases.NewMessageUC(msgRepo)

	router := chi.NewRouter()
	msgHandler := NewMessageHandler(router, msgUC, log)
	handler := NewHandler(router, msgHandler, log)

	return &http.Server{
		Addr:              cfg.HTTPServer.Addr,
		Handler:           handler,
		ReadTimeout:       cfg.HTTPServer.Timeouts.Read,
		ReadHeaderTimeout: cfg.HTTPServer.Timeouts.ReadHeader,
		WriteTimeout:      cfg.HTTPServer.Timeouts.Write,
		IdleTimeout:       cfg.HTTPServer.Timeouts.Idle,
		MaxHeaderBytes:    cfg.HTTPServer.MaxHeaderBytes,
	}
}
