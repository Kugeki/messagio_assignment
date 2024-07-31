package rest

import (
	"errors"
	"github.com/gavv/httpexpect/v2"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"messagio_assignment/internal/domain"
	"messagio_assignment/internal/domain/message"
	"messagio_assignment/internal/ports/rest/dto"
	"messagio_assignment/internal/ports/rest/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMessageHandler(t *testing.T) {
	t.Run("logger is not nil", func(t *testing.T) {
		mh := NewMessageHandler(nil, nil, nil, MessageHandlerConfig{})
		require.NotNil(t, mh.Log)
	})
}

func TestMessageHandler_CreateMessage(t *testing.T) {
	t.Run("json requests", func(t *testing.T) {
		tcases := []struct {
			Name            string
			Message         message.Message
			ExpectedStatus  int
			IsErrorExpected bool
			UcMockInit      func(uc *mocks.MessageUsecase)
		}{
			{
				Name: "successful",
				Message: message.Message{
					ID:        0,
					Content:   "some content",
					Processed: true,
				},
				ExpectedStatus:  http.StatusCreated,
				IsErrorExpected: false,
				UcMockInit: func(uc *mocks.MessageUsecase) {
					uc.On("CreateMessage", mock.Anything,
						&message.Message{
							ID:        0,
							Content:   "some content",
							Processed: true,
						}).Return(nil).Once()
				},
			},
			{
				Name:            "already exists",
				Message:         message.Message{},
				ExpectedStatus:  http.StatusConflict,
				IsErrorExpected: true,
				UcMockInit: func(uc *mocks.MessageUsecase) {
					uc.On("CreateMessage", mock.Anything, &message.Message{}).
						Return(domain.ErrAlreadyExists).Once()
				},
			},
			{
				Name:            "some db error",
				Message:         message.Message{},
				ExpectedStatus:  http.StatusUnprocessableEntity,
				IsErrorExpected: true,
				UcMockInit: func(uc *mocks.MessageUsecase) {
					uc.On("CreateMessage", mock.Anything, &message.Message{}).
						Return(errors.New("db error")).Once()
				},
			},
		}

		uc := mocks.NewMessageUsecase(t)
		router := chi.NewRouter()
		mh := NewMessageHandler(router, uc, nil, MessageHandlerConfig{})
		mh.SetupRoutes(router)

		server := httptest.NewServer(mh)
		defer server.Close()

		e := httpexpect.Default(t, server.URL)

		for _, tc := range tcases {
			t.Run(tc.Name, func(t *testing.T) {
				tc.UcMockInit(uc)

				msgReq := dto.CreateMessageReq{
					Content:   tc.Message.Content,
					Processed: tc.Message.Processed,
				}

				obj := e.POST("/messages").WithJSON(msgReq).
					Expect().
					Status(tc.ExpectedStatus).
					HasContentType("application/json").
					JSON().Object()

				if tc.IsErrorExpected {
					obj.Keys().ContainsOnly("error")
					obj.Value("error").String().NotEmpty()
					return
				}

				obj.Keys().NotContainsAny("error")

				wantResp := dto.CreateMessageResp{
					ID:        tc.Message.ID,
					Content:   tc.Message.Content,
					Processed: tc.Message.Processed,
				}

				var gotResp dto.CreateMessageResp
				obj.Decode(&gotResp)

				assert.Equal(t, wantResp, gotResp)
			})
		}
	})

	t.Run("bad request", func(t *testing.T) {
		uc := mocks.NewMessageUsecase(t)
		router := chi.NewRouter()
		mh := NewMessageHandler(router, uc, nil, MessageHandlerConfig{})
		mh.SetupRoutes(router)

		server := httptest.NewServer(mh)
		defer server.Close()

		e := httpexpect.Default(t, server.URL)

		obj := e.POST("/messages").WithBytes([]byte("some random")).
			Expect().
			HasContentType("application/json").
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Keys().ContainsOnly("error")
		obj.Value("error").String().NotEmpty()
	})
}

func TestMessageHandler_GetStats(t *testing.T) {
	tcases := []struct {
		Name            string
		ExpectedStats   message.Stats
		ExpectedStatus  int
		IsErrorExpected bool
		UcMockInit      func(uc *mocks.MessageUsecase)
	}{
		{
			Name: "successful",
			ExpectedStats: message.Stats{
				All:       42,
				Processed: 21,
			},
			ExpectedStatus:  http.StatusOK,
			IsErrorExpected: false,
			UcMockInit: func(uc *mocks.MessageUsecase) {
				uc.On("GetStats", mock.Anything).
					Return(&message.Stats{
						All:       42,
						Processed: 21,
					}, nil).Once()
			},
		},
		{
			Name:            "some db error",
			ExpectedStats:   message.Stats{},
			ExpectedStatus:  http.StatusInternalServerError,
			IsErrorExpected: true,
			UcMockInit: func(uc *mocks.MessageUsecase) {
				uc.On("GetStats", mock.Anything).
					Return(nil, errors.New("db error")).Once()
			},
		},
	}

	uc := mocks.NewMessageUsecase(t)
	router := chi.NewRouter()
	mh := NewMessageHandler(router, uc, nil, MessageHandlerConfig{})
	mh.SetupRoutes(router)

	server := httptest.NewServer(mh)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	for _, tc := range tcases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.UcMockInit(uc)

			obj := e.GET("/messages/stats").
				Expect().
				Status(tc.ExpectedStatus).
				HasContentType("application/json").
				JSON().Object()

			if tc.IsErrorExpected {
				obj.Keys().ContainsOnly("error")
				obj.Value("error").String().NotEmpty()
				return
			}

			obj.Keys().NotContainsAny("error")

			wantResp := dto.GetStatsResp{
				All:       tc.ExpectedStats.All,
				Processed: tc.ExpectedStats.Processed,
			}

			var gotResp dto.GetStatsResp
			obj.Decode(&gotResp)

			assert.Equal(t, wantResp, gotResp)
		})
	}
}

func TestMessageHandler_Limit(t *testing.T) {
	var (
		successfulRequests = 6
		limitedRequests    = 4
	)

	uc := mocks.NewMessageUsecase(t)
	router := chi.NewRouter()
	mh := NewMessageHandler(router, uc, nil, MessageHandlerConfig{
		CreateMsgPerMinute: successfulRequests,
		GetStatsPerMinute:  successfulRequests,
	})
	mh.SetupRoutes(router)

	server := httptest.NewServer(mh)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	t.Run("create message", func(t *testing.T) {
		var (
			msg           = &message.Message{Content: "some content"}
			createdStatus = http.StatusCreated
			limitStatus   = http.StatusTooManyRequests
		)

		for range successfulRequests {
			uc.On("CreateMessage", mock.Anything,
				msg).Return(nil).Once()
		}

		msgReq := dto.CreateMessageReq{
			Content:   msg.Content,
			Processed: msg.Processed,
		}

		for range successfulRequests {
			obj := e.POST("/messages").WithJSON(msgReq).
				Expect().
				Status(createdStatus).
				HasContentType("application/json").
				JSON().Object()

			obj.Keys().NotContainsAny("error")
		}

		for range limitedRequests {
			obj := e.POST("/messages").WithJSON(msgReq).
				Expect().
				Status(limitStatus).
				HasContentType("application/json").
				JSON().Object()

			obj.Keys().ContainsOnly("error")
			obj.Value("error").String().NotEmpty()
		}
	})

	t.Run("get stats", func(t *testing.T) {
		var (
			stats       = &message.Stats{}
			getStatus   = http.StatusOK
			limitStatus = http.StatusTooManyRequests
		)

		for range successfulRequests {
			uc.On("GetStats", mock.Anything).
				Return(stats, nil).Once()
		}

		for range successfulRequests {
			obj := e.GET("/messages/stats").
				Expect().
				Status(getStatus).
				HasContentType("application/json").
				JSON().Object()

			obj.Keys().NotContainsAny("error")
		}

		for range limitedRequests {
			obj := e.GET("/messages/stats").
				Expect().
				Status(limitStatus).
				HasContentType("application/json").
				JSON().Object()

			obj.Keys().ContainsOnly("error")
			obj.Value("error").String().NotEmpty()
		}
	})
}
