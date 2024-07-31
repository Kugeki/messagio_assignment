package rest

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log/slog"
	_ "messagio_assignment/docs" // for swagger
	"messagio_assignment/internal/config"
	"net/http"
	"net/url"
)

//go:generate

//	@title			Messagio Assigment
//	@version		0.1
//	@description	Test task to Messagio.

// @BasePath	/

func NewServer(httpCfg config.HTTPServer, msgUC MessageUsecase, log *slog.Logger) *http.Server {
	router := chi.NewRouter()
	msgHandler := NewMessageHandler(router, msgUC, log, MessageHandlerConfig{
		CreateMsgPerMinute: httpCfg.Handlers.Message.CreateMsgPerMinute,
		GetStatsPerMinute:  httpCfg.Handlers.Message.GetStatsPerMinute,
	})
	handler := NewHandler(router, msgHandler, log)

	swaggerURL := url.URL{
		Path: "/swagger/doc.json",
	}
	handler.Router.Get("/swagger/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		httpSwagger.Handler(
			httpSwagger.URL(swaggerURL.String()),
		)(w, r)
	})

	handler.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusPermanentRedirect)
	})

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
