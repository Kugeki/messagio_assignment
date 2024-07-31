package rest

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"log/slog"
	"messagio_assignment/internal/domain"
	"messagio_assignment/internal/domain/message"
	"messagio_assignment/internal/logger"
	"messagio_assignment/internal/ports/rest/dto"
	"net/http"
	"time"
)

//go:generate mockery --name MessageUsecase
type MessageUsecase interface {
	CreateMessage(ctx context.Context, msg *message.Message) error
	GetStats(ctx context.Context) (*message.Stats, error)
}

type MessageHandlerConfig struct {
	CreateMsgPerMinute int
	GetStatsPerMinute  int
}

type MessageHandler struct {
	router chi.Router
	uc     MessageUsecase
	cfg    MessageHandlerConfig

	Log *slog.Logger
}

func NewMessageHandler(router chi.Router, uc MessageUsecase,
	log *slog.Logger, cfg MessageHandlerConfig) *MessageHandler {
	if log == nil {
		log = logger.NewEraseLogger()
	}

	log = log.With(
		slog.String("component", "ports/rest/message_handler"),
	)
	return &MessageHandler{router: router, uc: uc, Log: log, cfg: cfg}
}

func (h *MessageHandler) SetupRoutes(r chi.Router) {
	r.Route("/messages", func(r chi.Router) {
		if h.cfg.CreateMsgPerMinute != 0 {
			r.Use(httprate.Limit(h.cfg.CreateMsgPerMinute, time.Minute,
				httprate.WithLimitHandler(h.Limit()),
			))
		}
		r.Post("/", h.CreateMessage())
	})
	r.Route("/messages/stats", func(r chi.Router) {
		if h.cfg.GetStatsPerMinute != 0 {
			r.Use(httprate.Limit(h.cfg.GetStatsPerMinute, time.Minute,
				httprate.WithLimitHandler(h.Limit()),
			))
		}
		r.Get("/", h.GetStats())
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

func (h *MessageHandler) Limit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.error(w, http.StatusTooManyRequests, errors.New("too many requests"))
	}
}

func (h *MessageHandler) error(w http.ResponseWriter, code int, err error) {
	h.respond(w, code, dto.HTTPError{Error: err.Error()})
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
