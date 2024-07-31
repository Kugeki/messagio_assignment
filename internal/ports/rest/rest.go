package rest

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"messagio_assignment/internal/config"
	"net/http"
)

func NewServer(httpCfg config.HTTPServer, msgUC MessageUsecase, log *slog.Logger) *http.Server {
	router := chi.NewRouter()

	msgHandler := NewMessageHandler(router, msgUC, log, MessageHandlerConfig{
		CreateMsgPerMinute: httpCfg.Handlers.Message.CreateMsgPerMinute,
		GetStatsPerMinute:  httpCfg.Handlers.Message.GetStatsPerMinute,
	})
	handler := NewHandler(router, msgHandler, log)

	return &http.Server{
		Addr:              httpCfg.Addr,
		Handler:           handler,
		ReadTimeout:       httpCfg.Timeouts.Read,
		ReadHeaderTimeout: httpCfg.Timeouts.ReadHeader,
		WriteTimeout:      httpCfg.Timeouts.Write,
		IdleTimeout:       httpCfg.Timeouts.Idle,
		MaxHeaderBytes:    httpCfg.MaxHeaderBytes,
	}
}
