package rest

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"messagio_assignment/internal/domain"
	"messagio_assignment/internal/domain/message"
	"messagio_assignment/internal/logger"
	"messagio_assignment/internal/ports/rest/dto"
	"net/http"
)

//go:generate mockery --name MessageUsecase
type MessageUsecase interface {
	CreateMessage(ctx context.Context, msg *message.Message) error
	GetStats(ctx context.Context) (*message.Stats, error)
	UpdateProcessedMessage(ctx context.Context, msg *message.Message) error
}

type MessageHandler struct {
	router chi.Router
	uc     MessageUsecase

	Log *slog.Logger
}

func NewMessageHandler(router chi.Router, uc MessageUsecase, log *slog.Logger) *MessageHandler {
	if log == nil {
		log = logger.NewEraseLogger()
	}

	log = log.With(
		slog.String("component", "ports/rest/message_handler"),
	)
	return &MessageHandler{router: router, uc: uc, Log: log}
}

func (h *MessageHandler) SetupRoutes(r chi.Router) {
	r.Route("/messages", func(r chi.Router) {
		r.Post("/", h.CreateMessage())
		r.Get("/stats", h.GetStats())
	})
}

func (h *MessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *MessageHandler) CreateMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.ForRest(h.Log, "create message", r.Context())

		msgReq := dto.CreateMessageReq{}

		if err := json.NewDecoder(r.Body).Decode(&msgReq); err != nil {
			log.Warn("failed to decode request body", logger.Err(err))
			h.error(w, http.StatusBadRequest, err)
			return
		}

		log.Info("request body is decoded", slog.Any("msgReq", msgReq))
		msg := msgReq.ToDomain()

		err := h.uc.CreateMessage(r.Context(), msg)
		if err != nil {
			log.Error("failed to create message", logger.Err(err))
			switch {
			case errors.Is(err, domain.ErrAlreadyExists):
				h.error(w, http.StatusConflict, err)
			default:
				h.error(w, http.StatusUnprocessableEntity, err) // or InternalError?
			}
			return
		}

		log.Info("message is created", slog.Any("msg", msg))

		var msgResp dto.CreateMessageResp
		msgResp.FromDomain(msg)

		h.respond(w, http.StatusCreated, msgResp)
	}
}

func (h *MessageHandler) GetStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.ForRest(h.Log, "get stats", r.Context())

		stats, err := h.uc.GetStats(r.Context())
		if err != nil {
			log.Error("failed to get stats", logger.Err(err))
			h.error(w, http.StatusInternalServerError, err)
			return
		}

		log.Info("stats is gotten", slog.Any("stats", stats))

		var statsResp dto.GetStatsResp
		statsResp.FromDomain(stats)

		h.respond(w, http.StatusOK, statsResp)
	}
}

func (h *MessageHandler) error(w http.ResponseWriter, code int, err error) {
	h.respond(w, code, map[string]string{"error": err.Error()})
}

func (h *MessageHandler) respond(w http.ResponseWriter, code int, data interface{}) {
	var jsonData []byte
	var err error

	if data != nil {
		jsonData, err = json.Marshal(data)

		if err != nil {
			h.Log.Error("json marshal error", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(code)
	_, err = w.Write(jsonData)
	if err != nil {
		h.Log.Error("response write json data error", slog.String("error", err.Error()))
	}
}
