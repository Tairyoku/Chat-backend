package messages

import (
	"cmd/pkg/handler/middlewares"
	"cmd/pkg/repository/models"
	"cmd/pkg/service"
	mockService "cmd/pkg/service/mocks"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestMessageHandler_CreateMessage(t *testing.T) {
	type mockBehavior func(s *mockService.MockMessage, message models.Message)

	testTable := []struct {
		name                 string
		inputText            string
		inputMessage         models.Message
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			inputText: `{"text":"test body"}`,
			inputMessage: models.Message{
				Author: 5,
				ChatId: 3,
				Text:   "test body",
				SentAt: time.Now().Round(10 * time.Millisecond),
			},
			mockBehavior: func(s *mockService.MockMessage, msg models.Message) {
				s.EXPECT().Create(msg).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}` + "\n",
		},
		{
			name:      "empty body",
			inputText: `{"text":""}`,
			mockBehavior: func(s *mockService.MockMessage, msg models.Message) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"body is empty"}` + "\n",
		},
		{
			name:      "server error",
			inputText: `{"text":"test body"}`,
			inputMessage: models.Message{
				Author: 5,
				ChatId: 3,
				Text:   "test body",
				SentAt: time.Now().Round(10 * time.Millisecond),
			},
			mockBehavior: func(s *mockService.MockMessage, msg models.Message) {
				s.EXPECT().Create(msg).Return(0, errors.New("create message error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"create message error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			msg := mockService.NewMockMessage(c)
			testCase.mockBehavior(msg, testCase.inputMessage)

			services := &service.Service{Message: msg}
			handler := NewMessageHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodPost, "/api/chats/:chatId/messages",
				strings.NewReader(testCase.inputText))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set(middlewares.UserCtx, 5)
			ctx.SetPath("/api/chats/:chatId/messages")
			ctx.SetParamNames("chatId")
			ctx.SetParamValues("3")

			//Перевірка результатів
			if assert.NoError(t, handler.CreateMessage(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestMessageHandler_GetLimitMessagesMessage(t *testing.T) {
	type mockBehavior func(s *mockService.MockMessage, chatId, limit int)

	testTable := []struct {
		name                 string
		inputChatId          int
		inputLimit           int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "ok",
			inputChatId: 13,
			inputLimit:  2,
			mockBehavior: func(s *mockService.MockMessage, chatId, limit int) {
				ret := []models.Message{
					{
						Id:     14,
						Author: 5,
						ChatId: 13,
						Text:   "test body",
						SentAt: time.Date(2023, 10, 10, 10, 10, 10, 10, time.UTC),
					},
					{
						Id:     15,
						Author: 5,
						ChatId: 13,
						Text:   "test body",
						SentAt: time.Date(2023, 10, 10, 10, 11, 10, 10, time.UTC),
					},
				}
				s.EXPECT().GetLimit(chatId, limit).Return(ret, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"list":[{"id":15,"chat_id":13,"author":5,"text":"test body","sent_at":"2023-10-10T10:11:10.00000001Z"},{"id":14,"chat_id":13,"author":5,"text":"test body","sent_at":"2023-10-10T10:10:10.00000001Z"}]}` + "\n",
		},
		{
			name:        "server error",
			inputChatId: 13,
			inputLimit:  2,
			mockBehavior: func(s *mockService.MockMessage, chatId, limit int) {
				s.EXPECT().GetLimit(chatId, limit).Return(nil, errors.New("get limit error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"get limit error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			msg := mockService.NewMockMessage(c)
			testCase.mockBehavior(msg, testCase.inputChatId, testCase.inputLimit)

			services := &service.Service{Message: msg}
			handler := NewMessageHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodGet, "/api/chats/:chatId/messages/limit/:id", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			//ctx.Set(middlewares.UserCtx, 5)
			ctx.SetPath("/api/chats/:chatId/messages/limit/:id")
			ctx.SetParamNames("chatId", "id")
			ctx.SetParamValues("13", "2")

			//Перевірка результатів
			if assert.NoError(t, handler.GetLimitMessages(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}
