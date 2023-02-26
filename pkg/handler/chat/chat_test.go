package chat

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
	"strconv"
	"strings"
	"testing"
)

func TestChatHandler_CreatePublicChat(t *testing.T) {
	type mockBehavior func(s *mockService.MockChat, userId int, chat models.Chat)

	testTable := []struct {
		name                 string
		inputUserId          int
		inputBody            string
		inputChat            models.Chat
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "ok",
			inputUserId: 4,
			inputBody:   `{"name":"test"}`,
			inputChat: models.Chat{
				Name:  "test",
				Types: "public",
			},
			mockBehavior: func(s *mockService.MockChat, userId int, chat models.Chat) {
				res := 5
				s.EXPECT().Create(chat).Return(res, nil)
				addUser := models.ChatUsers{
					ChatId: res,
					UserId: userId,
				}
				s.EXPECT().AddUser(addUser).Return(3, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":5}` + "\n",
		},
		{
			name:        "Incorrect request data",
			inputUserId: 4,
			inputBody:   `{"error"}`,
			mockBehavior: func(s *mockService.MockChat, userId int, chat models.Chat) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"incorrect request data"}` + "\n",
		},
		{
			name:        "No name",
			inputUserId: 4,
			inputBody:   `{"name":""}`,
			mockBehavior: func(s *mockService.MockChat, userId int, chat models.Chat) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"name is empty"}` + "\n",
		},
		{
			name:        "Create message error",
			inputUserId: 4,
			inputBody:   `{"name":"test"}`,
			inputChat: models.Chat{
				Name:  "test",
				Types: "public",
			},
			mockBehavior: func(s *mockService.MockChat, userId int, chat models.Chat) {
				s.EXPECT().Create(chat).Return(0, errors.New("create message error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"create message error"}` + "\n",
		},
		{
			name:        "Add user to chat error",
			inputUserId: 4,
			inputBody:   `{"name":"test"}`,
			inputChat: models.Chat{
				Name:  "test",
				Types: "public",
			},
			mockBehavior: func(s *mockService.MockChat, userId int, chat models.Chat) {
				res := 5
				s.EXPECT().Create(chat).Return(res, nil)
				addUser := models.ChatUsers{
					ChatId: res,
					UserId: userId,
				}
				s.EXPECT().AddUser(addUser).Return(0, errors.New("add user to chat error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"add user to chat error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			chat := mockService.NewMockChat(c)
			testCase.mockBehavior(chat, testCase.inputUserId, testCase.inputChat)

			services := &service.Service{Chat: chat}
			handler := NewChatHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodPost, "/api/chats/create",
				strings.NewReader(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set(middlewares.UserCtx, testCase.inputUserId)

			//Перевірка результатів
			if assert.NoError(t, handler.CreatePublicChat(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestChatHandler_GetChat(t *testing.T) {
	type mockBehavior func(s *mockService.MockChat, chatId int)

	testTable := []struct {
		name                 string
		inputChatId          int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "ok",
			inputChatId: 4,
			mockBehavior: func(s *mockService.MockChat, chatId int) {
				res := models.Chat{
					Id:    4,
					Name:  "name",
					Types: "public",
					Icon:  "",
				}
				s.EXPECT().Get(chatId).Return(res, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"chat":{"id":4,"name":"name","types":"public","icon":""}}` + "\n",
		},
		{
			name:        "Get chat error",
			inputChatId: 4,
			mockBehavior: func(s *mockService.MockChat, chatId int) {
				res := models.Chat{}
				s.EXPECT().Get(chatId).Return(res, errors.New("get chat error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"get chat error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			chat := mockService.NewMockChat(c)
			testCase.mockBehavior(chat, testCase.inputChatId)

			services := &service.Service{Chat: chat}
			handler := NewChatHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodGet, "/api/chats/:id", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetPath("/api/chats/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues(strconv.Itoa(testCase.inputChatId))

			//Перевірка результатів
			if assert.NoError(t, handler.GetChat(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestChatHandler_GetById(t *testing.T) {
	type mockBehavior func(s *mockService.MockChat, userId, chatId int)

	testTable := []struct {
		name                 string
		inputUserId          int
		inputChatId          int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok for private chats",
			inputUserId: 4,
			inputChatId: 6,
			mockBehavior: func(s *mockService.MockChat, userId, chatId int) {
				res := models.Chat{
					Id:    6,
					Name:  "name",
					Types: "private",
					Icon:  "",
				}
				s.EXPECT().Get(chatId).Return(res, nil)
				users := []models.User{
					{
						Id:       4,
						Username: "first",
					},
					{
						Id:       6,
						Username: "second",
					},
				}
				s.EXPECT().GetUsers(chatId).Return(users, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"chat":{"id":6,"name":"name","types":"private","icon":""},"user":{"id":6,"username":"second","password":"","icon":""}}` + "\n",
		},
		{
			name:        "Ok for personal chats",
			inputUserId: 4,
			inputChatId: 6,
			mockBehavior: func(s *mockService.MockChat, userId, chatId int) {
				res := models.Chat{
					Id:    6,
					Name:  "name",
					Types: "private",
					Icon:  "",
				}
				s.EXPECT().Get(chatId).Return(res, nil)
				users := []models.User{
					{
						Id:       6,
						Username: "second",
					},
				}
				s.EXPECT().GetUsers(chatId).Return(users, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"chat":{"id":6,"name":"name","types":"private","icon":""},"user":{"id":6,"username":"second","password":"","icon":""}}` + "\n",
		},
		{
			name:        "Ok for public chats",
			inputUserId: 4,
			inputChatId: 6,
			mockBehavior: func(s *mockService.MockChat, userId, chatId int) {
				res := models.Chat{
					Id:    6,
					Name:  "name",
					Types: "public",
					Icon:  "",
				}
				s.EXPECT().Get(chatId).Return(res, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"chat":{"id":6,"name":"name","types":"public","icon":""},"user":{"id":0,"username":"","password":"","icon":""}}` + "\n",
		},
		{
			name:        "Get chat info error",
			inputUserId: 4,
			inputChatId: 6,
			mockBehavior: func(s *mockService.MockChat, userId, chatId int) {
				res := models.Chat{}
				s.EXPECT().Get(chatId).Return(res, errors.New("no chat error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"no chat error"}` + "\n",
		},
		{
			name:        "Get users info error",
			inputUserId: 4,
			inputChatId: 6,
			mockBehavior: func(s *mockService.MockChat, userId, chatId int) {
				res := models.Chat{
					Id:    6,
					Name:  "name",
					Types: "private",
					Icon:  "",
				}
				s.EXPECT().Get(chatId).Return(res, nil)
				users := []models.User{{}}
				s.EXPECT().GetUsers(chatId).Return(users, errors.New("get users error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"get users error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			chat := mockService.NewMockChat(c)
			testCase.mockBehavior(chat, testCase.inputUserId, testCase.inputChatId)

			services := &service.Service{Chat: chat}
			handler := NewChatHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodGet, "/api/chats/:id/link", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set(middlewares.UserCtx, testCase.inputUserId)
			ctx.SetPath("/api/chats/:id/link")
			ctx.SetParamNames("id")
			ctx.SetParamValues(strconv.Itoa(testCase.inputChatId))

			//Перевірка результатів
			if assert.NoError(t, handler.GetById(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestChatHandler_GetUsers(t *testing.T) {
	type mockBehavior func(s *mockService.MockChat, chatId int)

	testTable := []struct {
		name                 string
		inputChatId          int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			inputChatId: 6,
			mockBehavior: func(s *mockService.MockChat, chatId int) {
				users := []models.User{
					{
						Id:       4,
						Username: "first",
					},
					{
						Id:       6,
						Username: "second",
					},
				}
				s.EXPECT().GetUsers(chatId).Return(users, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":4,"username":"first","password":"","icon":""},{"id":6,"username":"second","password":"","icon":""}]` + "\n",
		},
		{
			name:        "Get users error",
			inputChatId: 6,
			mockBehavior: func(s *mockService.MockChat, chatId int) {
				users := []models.User{{}}
				s.EXPECT().GetUsers(chatId).Return(users, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"get users error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			chat := mockService.NewMockChat(c)
			testCase.mockBehavior(chat, testCase.inputChatId)

			services := &service.Service{Chat: chat}
			handler := NewChatHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodGet, "/api/chats/:id/users", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetPath("/api/chats/:id/users")
			ctx.SetParamNames("id")
			ctx.SetParamValues(strconv.Itoa(testCase.inputChatId))

			//Перевірка результатів
			if assert.NoError(t, handler.GetUsers(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestChatHandler_GetUserPublicChats(t *testing.T) {
	type mockBehavior func(s *mockService.MockChat, userId int)

	testTable := []struct {
		name                 string
		inputUserId          int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			inputUserId: 6,
			mockBehavior: func(s *mockService.MockChat, userId int) {
				chats := []models.Chat{
					{
						Id:    4,
						Name:  "first",
						Types: "public",
						Icon:  "some image name",
					},
					{
						Id:    6,
						Name:  "second",
						Types: "public",
					},
				}
				s.EXPECT().GetPublicChats(userId).Return(chats, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"list":[{"id":4,"name":"first","types":"public","icon":"some image name"},{"id":6,"name":"second","types":"public","icon":""}]}` + "\n",
		},
		{
			name:        "Get chats error",
			inputUserId: 6,
			mockBehavior: func(s *mockService.MockChat, userId int) {
				chats := []models.Chat{{}}
				s.EXPECT().GetPublicChats(userId).Return(chats, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"get chats error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			chat := mockService.NewMockChat(c)
			testCase.mockBehavior(chat, testCase.inputUserId)

			services := &service.Service{Chat: chat}
			handler := NewChatHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodGet, "/api/chats/:id/public", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetPath("/api/chats/:id/public")
			ctx.SetParamNames("id")
			ctx.SetParamValues(strconv.Itoa(testCase.inputUserId))

			//Перевірка результатів
			if assert.NoError(t, handler.GetUserPublicChats(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestChatHandler_GetUserPrivateChats(t *testing.T) {
	type mockBehavior func(s *mockService.MockChat, userId, personalId int)

	testTable := []struct {
		name                 string
		inputUserId          int
		inputPersonalId      int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:            "Ok",
			inputPersonalId: 3,
			inputUserId:     6,
			mockBehavior: func(s *mockService.MockChat, userId, personalId int) {
				chats := []models.Chat{
					{
						Id:    3,
						Name:  "first",
						Types: "private",
						Icon:  "some image name",
					},
					{
						Id:    6,
						Name:  "second",
						Types: "private",
					},
				}
				s.EXPECT().GetPrivateChats(userId).Return(chats, nil)
				users := [][]models.User{
					{
						{
							Id:       3,
							Username: "first",
						},
						{
							Id:       6,
							Username: "second",
						},
					},
					{
						{
							Id:       3,
							Username: "first",
						},
					},
				}
				for i, v := range chats {
					s.EXPECT().GetUsers(v.Id).Return(users[i], nil)
				}
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"list":[{"id":3,"name":"first","types":"private","icon":"some image name"}]}` + "\n",
		},
		{
			name:            "Get chats error",
			inputPersonalId: 3,
			inputUserId:     6,
			mockBehavior: func(s *mockService.MockChat, userId, personalId int) {
				chats := []models.Chat{{}}
				s.EXPECT().GetPrivateChats(userId).Return(chats, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"get chats error"}` + "\n",
		},
		{
			name:            "Get users error",
			inputPersonalId: 3,
			inputUserId:     6,
			mockBehavior: func(s *mockService.MockChat, userId, personalId int) {
				chats := []models.Chat{
					{
						Id:    3,
						Name:  "first",
						Types: "private",
						Icon:  "some image name",
					},
					{
						Id:    6,
						Name:  "second",
						Types: "private",
					},
				}
				s.EXPECT().GetPrivateChats(userId).Return(chats, nil)
				users := [][]models.User{
					{
						{
							Id:       3,
							Username: "first",
						},
						{
							Id:       6,
							Username: "second",
						},
					},
					{
						{
							Id:       0,
							Username: "",
						},
					},
				}
				errors := []interface{}{
					nil,
					errors.New("some error"),
				}
				for i, v := range chats {
					s.EXPECT().GetUsers(v.Id).Return(users[i], errors[i])
				}
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"get users error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			chat := mockService.NewMockChat(c)
			testCase.mockBehavior(chat, testCase.inputUserId, testCase.inputPersonalId)

			services := &service.Service{Chat: chat}
			handler := NewChatHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodGet, "/api/chats/:id/private", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set(middlewares.UserCtx, testCase.inputUserId)
			ctx.SetPath("/api/chats/:id/private")
			ctx.SetParamNames("id")
			ctx.SetParamValues(strconv.Itoa(testCase.inputUserId))

			//Перевірка результатів
			if assert.NoError(t, handler.GetUserPrivateChats(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestChatHandler_AddUserToChat(t *testing.T) {
	type mockBehavior func(s *mockService.MockChat, chatId int, list models.ChatUsers)

	testTable := []struct {
		name                 string
		inputChatId          int
		inputBody            string
		inputChatUsers       models.ChatUsers
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "ok",
			inputChatId: 4,
			inputBody:   `{"user_id":8}`,
			inputChatUsers: models.ChatUsers{
				UserId: 8,
			},
			mockBehavior: func(s *mockService.MockChat, chatId int, list models.ChatUsers) {
				res := 5
				list.ChatId = chatId
				s.EXPECT().AddUser(list).Return(res, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":5}` + "\n",
		},
		{
			name:        "Incorrect request data",
			inputChatId: 4,
			inputBody:   `{"error"}`,
			inputChatUsers: models.ChatUsers{
				UserId: 8,
			},
			mockBehavior: func(s *mockService.MockChat, chatId int, list models.ChatUsers) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"incorrect request data"}` + "\n",
		},
		{
			name:        "Add user to chat error",
			inputChatId: 4,
			inputBody:   `{"user_id":8}`,
			inputChatUsers: models.ChatUsers{
				UserId: 8,
			},
			mockBehavior: func(s *mockService.MockChat, chatId int, list models.ChatUsers) {
				list.ChatId = chatId
				s.EXPECT().AddUser(list).Return(0, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"add user to chat error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			chat := mockService.NewMockChat(c)
			testCase.mockBehavior(chat, testCase.inputChatId, testCase.inputChatUsers)

			services := &service.Service{Chat: chat}
			handler := NewChatHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodPost, "/api/chats/:id/add",
				strings.NewReader(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetPath("/api/chats/:id/add")
			ctx.SetParamNames("id")
			ctx.SetParamValues(strconv.Itoa(testCase.inputChatId))

			//Перевірка результатів
			if assert.NoError(t, handler.AddUserToChat(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestChatHandler_DeleteUserFromChat(t *testing.T) {
	type mockBehavior func(s *mockService.MockChat, chatId int, list models.ChatUsers)

	testTable := []struct {
		name                 string
		inputChatId          int
		inputBody            string
		inputChatUsers       models.ChatUsers
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			inputChatId: 4,
			inputBody:   `{"user_id":8}`,
			inputChatUsers: models.ChatUsers{
				UserId: 8,
			},
			mockBehavior: func(s *mockService.MockChat, chatId int, list models.ChatUsers) {
				s.EXPECT().DeleteUser(list.UserId, chatId).Return(nil)
				users := []models.User{
					{
						Id:       4,
						Username: "user",
					},
				}
				s.EXPECT().GetUsers(chatId).Return(users, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"user with id 8 deleted from chat with id 4"}` + "\n",
		},
		{
			name:        "Ok and delete chat",
			inputChatId: 4,
			inputBody:   `{"user_id":8}`,
			inputChatUsers: models.ChatUsers{
				UserId: 8,
			},
			mockBehavior: func(s *mockService.MockChat, chatId int, list models.ChatUsers) {
				s.EXPECT().DeleteUser(list.UserId, chatId).Return(nil)
				var users []models.User
				s.EXPECT().GetUsers(chatId).Return(users, nil)
				s.EXPECT().Delete(chatId).Return(nil)
				s.EXPECT().DeleteAllMessages(chatId).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"user with id 8 deleted from chat with id 4"}` + "\n",
		},
		{
			name:        "Incorrect request data",
			inputChatId: 4,
			inputBody:   `{"error"}`,
			mockBehavior: func(s *mockService.MockChat, chatId int, list models.ChatUsers) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"incorrect request data"}` + "\n",
		},
		{
			name:        "Delete user error",
			inputChatId: 4,
			inputBody:   `{"user_id":8}`,
			inputChatUsers: models.ChatUsers{
				UserId: 8,
			},
			mockBehavior: func(s *mockService.MockChat, chatId int, list models.ChatUsers) {
				s.EXPECT().DeleteUser(list.UserId, chatId).Return(errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"delete user error"}` + "\n",
		},
		{
			name:        "Get chat users error",
			inputChatId: 4,
			inputBody:   `{"user_id":8}`,
			inputChatUsers: models.ChatUsers{
				UserId: 8,
			},
			mockBehavior: func(s *mockService.MockChat, chatId int, list models.ChatUsers) {
				s.EXPECT().DeleteUser(list.UserId, chatId).Return(nil)
				var users []models.User
				s.EXPECT().GetUsers(chatId).Return(users, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"get chat users error"}` + "\n",
		},
		{
			name:        "Chat delete error",
			inputChatId: 4,
			inputBody:   `{"user_id":8}`,
			inputChatUsers: models.ChatUsers{
				UserId: 8,
			},
			mockBehavior: func(s *mockService.MockChat, chatId int, list models.ChatUsers) {
				s.EXPECT().DeleteUser(list.UserId, chatId).Return(nil)
				var users []models.User
				s.EXPECT().GetUsers(chatId).Return(users, nil)
				s.EXPECT().Delete(chatId).Return(errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"chat delete error"}` + "\n",
		},
		{
			name:        "Messages delete error",
			inputChatId: 4,
			inputBody:   `{"user_id":8}`,
			inputChatUsers: models.ChatUsers{
				UserId: 8,
			},
			mockBehavior: func(s *mockService.MockChat, chatId int, list models.ChatUsers) {
				s.EXPECT().DeleteUser(list.UserId, chatId).Return(nil)
				var users []models.User
				s.EXPECT().GetUsers(chatId).Return(users, nil)
				s.EXPECT().Delete(chatId).Return(nil)
				s.EXPECT().DeleteAllMessages(chatId).Return(errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"messages delete error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			chat := mockService.NewMockChat(c)
			testCase.mockBehavior(chat, testCase.inputChatId, testCase.inputChatUsers)

			services := &service.Service{Chat: chat}
			handler := NewChatHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodDelete, "/api/chats/:id/delete",
				strings.NewReader(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetPath("/api/chats/:id/delete")
			ctx.SetParamNames("id")
			ctx.SetParamValues(strconv.Itoa(testCase.inputChatId))

			//Перевірка результатів
			if assert.NoError(t, handler.DeleteUserFromChat(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestChatHandler_DeleteChat(t *testing.T) {
	type mockBehavior func(s *mockService.MockChat, chatId int)

	testTable := []struct {
		name                 string
		inputChatId          int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			inputChatId: 4,
			mockBehavior: func(s *mockService.MockChat, chatId int) {
				users := []models.User{
					{
						Id:       3,
						Username: "first",
					},
					{
						Id:       5,
						Username: "second",
					},
				}
				s.EXPECT().GetUsers(chatId).Return(users, nil)
				errors := []interface{}{
					nil,
					nil,
				}
				for i, v := range users {
					s.EXPECT().DeleteUser(v.Id, chatId).Return(errors[i])
				}
				s.EXPECT().Delete(chatId).Return(nil)
				s.EXPECT().DeleteAllMessages(chatId).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"chat with id 4 deleted"}` + "\n",
		},
		{
			name:        "Get chat users error",
			inputChatId: 4,
			mockBehavior: func(s *mockService.MockChat, chatId int) {
				var users []models.User
				s.EXPECT().GetUsers(chatId).Return(users, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"get chat users error"}` + "\n",
		},
		{
			name:        "Delete user from chat error",
			inputChatId: 4,
			mockBehavior: func(s *mockService.MockChat, chatId int) {
				users := []models.User{
					{
						Id:       3,
						Username: "first",
					},
					{
						Id:       5,
						Username: "second",
					},
				}
				s.EXPECT().GetUsers(chatId).Return(users, nil)
				errors := []interface{}{
					nil,
					errors.New("some error"),
				}
				for i, v := range users {
					s.EXPECT().DeleteUser(v.Id, chatId).Return(errors[i])
				}
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"delete user from chat error"}` + "\n",
		},
		{
			name:        "Chat delete error",
			inputChatId: 4,
			mockBehavior: func(s *mockService.MockChat, chatId int) {
				users := []models.User{
					{
						Id:       3,
						Username: "first",
					},
					{
						Id:       5,
						Username: "second",
					},
				}
				s.EXPECT().GetUsers(chatId).Return(users, nil)
				error := []interface{}{
					nil,
					nil,
				}
				for i, v := range users {
					s.EXPECT().DeleteUser(v.Id, chatId).Return(error[i])
				}
				s.EXPECT().Delete(chatId).Return(errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"chat delete error"}` + "\n",
		},
		{
			name:        "Messages delete error",
			inputChatId: 4,
			mockBehavior: func(s *mockService.MockChat, chatId int) {
				users := []models.User{
					{
						Id:       3,
						Username: "first",
					},
					{
						Id:       5,
						Username: "second",
					},
				}
				s.EXPECT().GetUsers(chatId).Return(users, nil)
				error := []interface{}{
					nil,
					nil,
				}
				for i, v := range users {
					s.EXPECT().DeleteUser(v.Id, chatId).Return(error[i])
				}
				s.EXPECT().Delete(chatId).Return(nil)
				s.EXPECT().DeleteAllMessages(chatId).Return(errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"messages delete error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			chat := mockService.NewMockChat(c)
			testCase.mockBehavior(chat, testCase.inputChatId)

			services := &service.Service{Chat: chat}
			handler := NewChatHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodDelete, "/api/chats/:id", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetPath("/api/chats/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues(strconv.Itoa(testCase.inputChatId))

			//Перевірка результатів
			if assert.NoError(t, handler.DeleteChat(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestChatHandler_PrivateChat(t *testing.T) {
	type mockBehavior func(s *mockService.MockChat, userId, activeUserId int)

	testTable := []struct {
		name                 string
		inputUserId          int
		inputActiveUserId    int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:              "Ok",
			inputUserId:       4,
			inputActiveUserId: 1,
			mockBehavior: func(s *mockService.MockChat, userId, activeUserId int) {
				code := 12
				s.EXPECT().GetPrivates(activeUserId, userId).Return(code, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"chatId":12}` + "\n",
		},
		{
			name:              "Ok and create new private chat",
			inputUserId:       4,
			inputActiveUserId: 1,
			mockBehavior: func(s *mockService.MockChat, userId, activeUserId int) {
				code := 0
				s.EXPECT().GetPrivates(activeUserId, userId).Return(code, nil)
				user := models.User{
					Id:       4,
					Username: "first",
				}
				s.EXPECT().GetUserById(userId).Return(user, nil)
				chat := models.Chat{
					Name:  user.Username,
					Types: "private",
				}
				newChatId := 12
				s.EXPECT().Create(chat).Return(12, nil)
				s.EXPECT().AddUser(models.ChatUsers{
					ChatId: newChatId,
					UserId: activeUserId,
				}).Return(3, nil)
				s.EXPECT().AddUser(models.ChatUsers{
					ChatId: newChatId,
					UserId: userId,
				}).Return(4, nil)

			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"chatId":12}` + "\n",
		},
		{
			name:              "Ok and create new personal chat",
			inputUserId:       1,
			inputActiveUserId: 1,
			mockBehavior: func(s *mockService.MockChat, userId, activeUserId int) {
				code := 0
				s.EXPECT().GetPrivates(activeUserId, userId).Return(code, nil)
				user := models.User{
					Id:       1,
					Username: "first",
				}
				s.EXPECT().GetUserById(userId).Return(user, nil)
				chat := models.Chat{
					Name:  user.Username,
					Types: "private",
				}
				newChatId := 12
				s.EXPECT().Create(chat).Return(12, nil)
				s.EXPECT().AddUser(models.ChatUsers{
					ChatId: newChatId,
					UserId: activeUserId,
				}).Return(3, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"chatId":12}` + "\n",
		},
		{
			name:              "Wrong users",
			inputUserId:       4,
			inputActiveUserId: 1,
			mockBehavior: func(s *mockService.MockChat, userId, activeUserId int) {
				code := 0
				s.EXPECT().GetPrivates(activeUserId, userId).Return(code, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"wrong users"}` + "\n",
		},
		{
			name:              "Wrong user",
			inputUserId:       4,
			inputActiveUserId: 1,
			mockBehavior: func(s *mockService.MockChat, userId, activeUserId int) {
				code := 0
				s.EXPECT().GetPrivates(activeUserId, userId).Return(code, nil)
				user := models.User{}
				s.EXPECT().GetUserById(userId).Return(user, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"wrong user"}` + "\n",
		},
		{
			name:              "Create chat error",
			inputUserId:       4,
			inputActiveUserId: 1,
			mockBehavior: func(s *mockService.MockChat, userId, activeUserId int) {
				code := 0
				s.EXPECT().GetPrivates(activeUserId, userId).Return(code, nil)
				user := models.User{
					Id:       4,
					Username: "first",
				}
				s.EXPECT().GetUserById(userId).Return(user, nil)
				chat := models.Chat{
					Name:  user.Username,
					Types: "private",
				}
				s.EXPECT().Create(chat).Return(0, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"create chat error"}` + "\n",
		},
		{
			name:              "Add active user to chat error",
			inputUserId:       4,
			inputActiveUserId: 1,
			mockBehavior: func(s *mockService.MockChat, userId, activeUserId int) {
				code := 0
				s.EXPECT().GetPrivates(activeUserId, userId).Return(code, nil)
				user := models.User{
					Id:       4,
					Username: "first",
				}
				s.EXPECT().GetUserById(userId).Return(user, nil)
				chat := models.Chat{
					Name:  user.Username,
					Types: "private",
				}
				newChatId := 12
				s.EXPECT().Create(chat).Return(12, nil)
				s.EXPECT().AddUser(models.ChatUsers{
					ChatId: newChatId,
					UserId: activeUserId,
				}).Return(0, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"add active user to chat error"}` + "\n",
		},
		{
			name:              "Add user to chat error",
			inputUserId:       4,
			inputActiveUserId: 1,
			mockBehavior: func(s *mockService.MockChat, userId, activeUserId int) {
				code := 0
				s.EXPECT().GetPrivates(activeUserId, userId).Return(code, nil)
				user := models.User{
					Id:       4,
					Username: "first",
				}
				s.EXPECT().GetUserById(userId).Return(user, nil)
				chat := models.Chat{
					Name:  user.Username,
					Types: "private",
				}
				newChatId := 12
				s.EXPECT().Create(chat).Return(12, nil)
				s.EXPECT().AddUser(models.ChatUsers{
					ChatId: newChatId,
					UserId: activeUserId,
				}).Return(3, nil)
				s.EXPECT().AddUser(models.ChatUsers{
					ChatId: newChatId,
					UserId: userId,
				}).Return(0, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"add second user to chat error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			chat := mockService.NewMockChat(c)
			testCase.mockBehavior(chat, testCase.inputUserId, testCase.inputActiveUserId)

			services := &service.Service{Chat: chat}
			handler := NewChatHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodGet, "/api/chats/:userId/private", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set(middlewares.UserCtx, testCase.inputActiveUserId)
			ctx.SetPath("/api/chats/:userId/private")
			ctx.SetParamNames("userId")
			ctx.SetParamValues(strconv.Itoa(testCase.inputUserId))

			//Перевірка результатів
			if assert.NoError(t, handler.PrivateChat(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestChatHandler_SearchChat(t *testing.T) {
	type mockBehavior func(s *mockService.MockChat, name string)

	testTable := []struct {
		name                 string
		inputName            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputName: "f",
			mockBehavior: func(s *mockService.MockChat, name string) {
				chats := []models.Chat{
					{
						Id:    4,
						Name:  "first",
						Types: "public",
						Icon:  "",
					},
					{
						Id:    6,
						Name:  "stuff",
						Types: "public",
						Icon:  "",
					},
				}
				s.EXPECT().SearchChat(name).Return(chats, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"list":[{"id":4,"name":"first","types":"public","icon":""},{"id":6,"name":"stuff","types":"public","icon":""}]}` + "\n",
		},
		{
			name:      "Found chats error",
			inputName: "f",
			mockBehavior: func(s *mockService.MockChat, name string) {
				var chats []models.Chat
				s.EXPECT().SearchChat(name).Return(chats, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"found chats error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			chat := mockService.NewMockChat(c)
			testCase.mockBehavior(chat, testCase.inputName)

			services := &service.Service{Chat: chat}
			handler := NewChatHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodGet, "/api/chats/search/:name", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetPath("/api/chats/search/:name")
			ctx.SetParamNames("name")
			ctx.SetParamValues(testCase.inputName)

			//Перевірка результатів
			if assert.NoError(t, handler.SearchChat(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}
