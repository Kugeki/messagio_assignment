package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"messagioassignment/internal/logger"
	"net/http"
)

type ChiHandler interface {
	SetupRoutes(router chi.Router)
}

type Handler struct {
	Router     chi.Router
	msgHandler ChiHandler
	Log        *slog.Logger
}

func NewHandler(router chi.Router, mh ChiHandler, log *slog.Logger) *Handler {
	if log == nil {
		log = logger.NewEraseLogger()
	}

	h := &Handler{
		Router:     router,
		msgHandler: mh,
		Log:        log,
	}
	h.Configure()
	return h
}

func (h *Handler) Configure() {
	h.Middlewares()
	h.Routes()
}

func (h *Handler) Routes() {
	h.Router.Route("/", h.setupOtherRoutes)
}

func (h *Handler) setupOtherRoutes(r chi.Router) {
	h.msgHandler.SetupRoutes(r)
}

func (h *Handler) Middlewares() {
	h.Router.Use(middleware.RequestID)
	h.Router.Use(LogMiddleware(h.Log))
	h.Router.Use(middleware.Recoverer)
	h.Router.Use(middleware.Heartbeat("/health"))
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Router.ServeHTTP(w, r)
}
