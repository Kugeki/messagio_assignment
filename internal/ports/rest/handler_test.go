package rest

import (
	"github.com/gavv/httpexpect/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ChiHandlerMock struct {
	IsRoutesSetup    bool
	IsEndpointCalled bool

	RequestID string
}

func (h *ChiHandlerMock) SetupRoutes(r chi.Router) {
	h.IsRoutesSetup = true

	r.Get("/example", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		h.IsEndpointCalled = true
	})

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("this need to be recovered")
	})

	r.Get("/request_id", func(w http.ResponseWriter, r *http.Request) {
		h.RequestID = middleware.GetReqID(r.Context())
		w.WriteHeader(http.StatusOK)
	})
}

func TestNewHandler(t *testing.T) {
	handlerMock := &ChiHandlerMock{}

	router := chi.NewRouter()
	handler := NewHandler(router, handlerMock, nil)

	require.NotNil(t, handler.Log)
	assert.True(t, handlerMock.IsRoutesSetup)

	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	t.Run("handler endpoint call", func(t *testing.T) {
		e.GET("/example").
			Expect().
			Status(http.StatusTeapot)

		assert.True(t, handlerMock.IsEndpointCalled)
	})

	t.Run("heartbeat", func(t *testing.T) {
		e.GET("/health").
			Expect().
			StatusRange(httpexpect.Status2xx)
	})

	t.Run("recover", func(t *testing.T) {
		e.GET("/panic").
			Expect().
			StatusRange(httpexpect.Status5xx)
	})

	t.Run("request id", func(t *testing.T) {
		e.GET("/request_id").
			Expect().
			Status(http.StatusOK)

		assert.NotEmpty(t, handlerMock.RequestID)
	})
}
