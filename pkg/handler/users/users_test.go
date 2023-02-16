package users

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
	"testing"
)

func TestUsersHandler_GetUserById(t *testing.T) {
	type mockBehavior func(s *mockService.MockStatus, userId int)

	testTable := []struct {
		name                 string
		inputUserId          int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "ok",
			inputUserId: 13,
			mockBehavior: func(s *mockService.MockStatus, userId int) {
				ret := models.User{
					Id:       13,
					Username: "user",
				}
				s.EXPECT().GetUserById(userId).Return(ret, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"user":{"id":13,"username":"user","password":"","icon":""}}` + "\n",
		},
		{
			name:        "Get user error",
			inputUserId: 13,
			mockBehavior: func(s *mockService.MockStatus, userId int) {
				var ret models.User
				s.EXPECT().GetUserById(userId).Return(ret, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"get user error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			status := mockService.NewMockStatus(c)
			testCase.mockBehavior(status, testCase.inputUserId)

			services := &service.Service{Status: status}
			handler := NewUsersHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodGet, "/api/users/:id", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetPath("/api/users/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues(strconv.Itoa(testCase.inputUserId))

			//Перевірка результатів
			if assert.NoError(t, handler.GetUserById(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestUsersHandler_GetUserLists(t *testing.T) {
	type mockBehavior func(s *mockService.MockStatus, userId int)

	testTable := []struct {
		name                 string
		inputUserId          int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			inputUserId: 13,
			mockBehavior: func(s *mockService.MockStatus, userId int) {
				friends := []models.User{
					{
						Id:       13,
						Username: "user",
					},
					{
						Id:       2,
						Username: "friend",
					},
				}
				s.EXPECT().GetFriends(userId).Return(friends, nil)
				bl := []models.User{
					{
						Id:       7,
						Username: "blocked",
					},
				}
				s.EXPECT().GetBlackList(userId).Return(bl, nil)
				onBL := []models.User{
					{
						Id:       19,
						Username: "on block",
					},
				}
				s.EXPECT().GetBlackListToUser(userId).Return(onBL, nil)
				invites := []models.User{
					{
						Id:       21,
						Username: "invites",
					},
				}
				s.EXPECT().GetSentInvites(userId).Return(invites, nil)
				requires := []models.User{
					{
						Id:       4,
						Username: "requires",
					},
				}
				s.EXPECT().GetInvites(userId).Return(requires, nil)

			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"blacklist":[{"id":7,"username":"blocked","password":"","icon":""}],"friends":[{"id":13,"username":"user","password":"","icon":""},{"id":2,"username":"friend","password":"","icon":""}],"invites":[{"id":21,"username":"invites","password":"","icon":""}],"onBlacklist":[{"id":19,"username":"on block","password":"","icon":""}],"requires":[{"id":4,"username":"requires","password":"","icon":""}]}` + "\n",
		},
		{
			name:        "Friends list error",
			inputUserId: 13,
			mockBehavior: func(s *mockService.MockStatus, userId int) {
				var friends []models.User
				s.EXPECT().GetFriends(userId).Return(friends, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"friends list error"}` + "\n",
		},
		{
			name:        "Black list error",
			inputUserId: 13,
			mockBehavior: func(s *mockService.MockStatus, userId int) {
				friends := []models.User{
					{
						Id:       13,
						Username: "user",
					},
					{
						Id:       2,
						Username: "friend",
					},
				}
				s.EXPECT().GetFriends(userId).Return(friends, nil)
				var bl []models.User
				s.EXPECT().GetBlackList(userId).Return(bl, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"black list error"}` + "\n",
		},
		{
			name:        "On black list error",
			inputUserId: 13,
			mockBehavior: func(s *mockService.MockStatus, userId int) {
				friends := []models.User{
					{
						Id:       13,
						Username: "user",
					},
					{
						Id:       2,
						Username: "friend",
					},
				}
				s.EXPECT().GetFriends(userId).Return(friends, nil)
				bl := []models.User{
					{
						Id:       7,
						Username: "blocked",
					},
				}
				s.EXPECT().GetBlackList(userId).Return(bl, nil)
				var onBL []models.User
				s.EXPECT().GetBlackListToUser(userId).Return(onBL, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"on black list error"}` + "\n",
		},
		{
			name:        "Friend invites list error",
			inputUserId: 13,
			mockBehavior: func(s *mockService.MockStatus, userId int) {
				friends := []models.User{
					{
						Id:       13,
						Username: "user",
					},
					{
						Id:       2,
						Username: "friend",
					},
				}
				s.EXPECT().GetFriends(userId).Return(friends, nil)
				bl := []models.User{
					{
						Id:       7,
						Username: "blocked",
					},
				}
				s.EXPECT().GetBlackList(userId).Return(bl, nil)
				onBL := []models.User{
					{
						Id:       19,
						Username: "on block",
					},
				}
				s.EXPECT().GetBlackListToUser(userId).Return(onBL, nil)
				var invites []models.User
				s.EXPECT().GetSentInvites(userId).Return(invites, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"friend invites list error"}` + "\n",
		},
		{
			name:        "Friend requires list error",
			inputUserId: 13,
			mockBehavior: func(s *mockService.MockStatus, userId int) {
				friends := []models.User{
					{
						Id:       13,
						Username: "user",
					},
					{
						Id:       2,
						Username: "friend",
					},
				}
				s.EXPECT().GetFriends(userId).Return(friends, nil)
				bl := []models.User{
					{
						Id:       7,
						Username: "blocked",
					},
				}
				s.EXPECT().GetBlackList(userId).Return(bl, nil)
				onBL := []models.User{
					{
						Id:       19,
						Username: "on block",
					},
				}
				s.EXPECT().GetBlackListToUser(userId).Return(onBL, nil)
				invites := []models.User{
					{
						Id:       21,
						Username: "invites",
					},
				}
				s.EXPECT().GetSentInvites(userId).Return(invites, nil)
				var requires []models.User
				s.EXPECT().GetInvites(userId).Return(requires, errors.New("some error"))

			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"friend requires list error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			status := mockService.NewMockStatus(c)
			testCase.mockBehavior(status, testCase.inputUserId)

			services := &service.Service{Status: status}
			handler := NewUsersHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodGet, "/api/users/:id/all", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetPath("/api/users/:id/all")
			ctx.SetParamNames("id")
			ctx.SetParamValues(strconv.Itoa(testCase.inputUserId))

			//Перевірка результатів
			if assert.NoError(t, handler.GetUserLists(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestUsersHandler_InvitedToFriends(t *testing.T) {
	type mockBehavior func(s *mockService.MockStatus, senderId, recipientId int)

	testTable := []struct {
		name                 string
		inputSenderId        int
		inputRecipientId     int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:             "Ok",
			inputSenderId:    13,
			inputRecipientId: 2,
			mockBehavior: func(s *mockService.MockStatus, senderId, recipientId int) {
				status := models.Status{
					SenderId:     senderId,
					RecipientId:  recipientId,
					Relationship: "invitation",
				}
				res := 2
				s.EXPECT().AddStatus(status).Return(res, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":2}` + "\n",
		},
		{
			name:             "Add status error",
			inputSenderId:    13,
			inputRecipientId: 2,
			mockBehavior: func(s *mockService.MockStatus, senderId, recipientId int) {
				status := models.Status{
					SenderId:     senderId,
					RecipientId:  recipientId,
					Relationship: "invitation",
				}
				res := 0
				s.EXPECT().AddStatus(status).Return(res, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"add status error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			status := mockService.NewMockStatus(c)
			testCase.mockBehavior(status, testCase.inputSenderId, testCase.inputRecipientId)

			services := &service.Service{Status: status}
			handler := NewUsersHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodPost, "/api/users/:id/invite", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set(middlewares.UserCtx, testCase.inputSenderId)
			ctx.SetPath("/api/users/:id/invite")
			ctx.SetParamNames("id")
			ctx.SetParamValues(strconv.Itoa(testCase.inputRecipientId))

			//Перевірка результатів
			if assert.NoError(t, handler.InvitedToFriends(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestUsersHandler_CancelInvite(t *testing.T) {
	type mockBehavior func(s *mockService.MockStatus, senderId, recipientId int)

	testTable := []struct {
		name                 string
		inputSenderId        int
		inputRecipientId     int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:             "Ok",
			inputSenderId:    13,
			inputRecipientId: 2,
			mockBehavior: func(s *mockService.MockStatus, senderId, recipientId int) {
				status := models.Status{
					SenderId:     senderId,
					RecipientId:  recipientId,
					Relationship: "invitation",
				}
				s.EXPECT().DeleteStatus(status).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"invite deleted"}` + "\n",
		},
		{
			name:             "Delete status error",
			inputSenderId:    13,
			inputRecipientId: 2,
			mockBehavior: func(s *mockService.MockStatus, senderId, recipientId int) {
				status := models.Status{
					SenderId:     senderId,
					RecipientId:  recipientId,
					Relationship: "invitation",
				}
				s.EXPECT().DeleteStatus(status).Return(errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"delete status error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			status := mockService.NewMockStatus(c)
			testCase.mockBehavior(status, testCase.inputSenderId, testCase.inputRecipientId)

			services := &service.Service{Status: status}
			handler := NewUsersHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodDelete, "/api/users/:id/cancel", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set(middlewares.UserCtx, testCase.inputSenderId)
			ctx.SetPath("/api/users/:id/cancel")
			ctx.SetParamNames("id")
			ctx.SetParamValues(strconv.Itoa(testCase.inputRecipientId))

			//Перевірка результатів
			if assert.NoError(t, handler.CancelInvite(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestUsersHandler_AcceptInvitation(t *testing.T) {
	type mockBehavior func(s *mockService.MockStatus, senderId, recipientId int)

	testTable := []struct {
		name                 string
		inputSenderId        int
		inputRecipientId     int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:             "Ok",
			inputSenderId:    13,
			inputRecipientId: 2,
			mockBehavior: func(s *mockService.MockStatus, senderId, recipientId int) {
				status := models.Status{
					SenderId:     senderId,
					RecipientId:  recipientId,
					Relationship: "friends",
				}
				s.EXPECT().UpdateStatus(status).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"Invitation accepted"}` + "\n",
		},
		{
			name:             "Update status error",
			inputSenderId:    13,
			inputRecipientId: 2,
			mockBehavior: func(s *mockService.MockStatus, senderId, recipientId int) {
				status := models.Status{
					SenderId:     senderId,
					RecipientId:  recipientId,
					Relationship: "friends",
				}
				s.EXPECT().UpdateStatus(status).Return(errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"update status error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			status := mockService.NewMockStatus(c)
			testCase.mockBehavior(status, testCase.inputSenderId, testCase.inputRecipientId)

			services := &service.Service{Status: status}
			handler := NewUsersHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodPut, "/api/users/:id/accept", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set(middlewares.UserCtx, testCase.inputRecipientId)
			ctx.SetPath("/api/users/:id/accept")
			ctx.SetParamNames("id")
			ctx.SetParamValues(strconv.Itoa(testCase.inputSenderId))

			//Перевірка результатів
			if assert.NoError(t, handler.AcceptInvitation(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestUsersHandler_RefuseInvitation(t *testing.T) {
	type mockBehavior func(s *mockService.MockStatus, senderId, recipientId int)

	testTable := []struct {
		name                 string
		inputSenderId        int
		inputRecipientId     int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:             "Ok",
			inputSenderId:    13,
			inputRecipientId: 2,
			mockBehavior: func(s *mockService.MockStatus, senderId, recipientId int) {
				status := models.Status{
					SenderId:     senderId,
					RecipientId:  recipientId,
					Relationship: "invitation",
				}
				s.EXPECT().DeleteStatus(status).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"invitation refused"}` + "\n",
		},
		{
			name:             "Delete status error",
			inputSenderId:    13,
			inputRecipientId: 2,
			mockBehavior: func(s *mockService.MockStatus, senderId, recipientId int) {
				status := models.Status{
					SenderId:     senderId,
					RecipientId:  recipientId,
					Relationship: "invitation",
				}
				s.EXPECT().DeleteStatus(status).Return(errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"delete status error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			status := mockService.NewMockStatus(c)
			testCase.mockBehavior(status, testCase.inputSenderId, testCase.inputRecipientId)

			services := &service.Service{Status: status}
			handler := NewUsersHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodDelete, "/api/users/:id/refuse", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set(middlewares.UserCtx, testCase.inputRecipientId)
			ctx.SetPath("/api/users/:id/refuse")
			ctx.SetParamNames("id")
			ctx.SetParamValues(strconv.Itoa(testCase.inputSenderId))

			//Перевірка результатів
			if assert.NoError(t, handler.RefuseInvitation(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestUsersHandler_DeleteFriend(t *testing.T) {
	type mockBehavior func(s *mockService.MockStatus, senderId, recipientId int)

	testTable := []struct {
		name                 string
		inputSenderId        int
		inputRecipientId     int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:             "Ok",
			inputSenderId:    13,
			inputRecipientId: 2,
			mockBehavior: func(s *mockService.MockStatus, senderId, recipientId int) {
				status := models.Status{
					SenderId:     senderId,
					RecipientId:  recipientId,
					Relationship: "friends",
				}
				s.EXPECT().DeleteStatus(status).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"friend deleted"}` + "\n",
		},
		{
			name:             "Delete status error",
			inputSenderId:    13,
			inputRecipientId: 2,
			mockBehavior: func(s *mockService.MockStatus, senderId, recipientId int) {
				status := models.Status{
					SenderId:     senderId,
					RecipientId:  recipientId,
					Relationship: "friends",
				}
				s.EXPECT().DeleteStatus(status).Return(errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"delete status error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			status := mockService.NewMockStatus(c)
			testCase.mockBehavior(status, testCase.inputSenderId, testCase.inputRecipientId)

			services := &service.Service{Status: status}
			handler := NewUsersHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodDelete, "/api/users/:id/refuse", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set(middlewares.UserCtx, testCase.inputSenderId)
			ctx.SetPath("/api/users/:id/refuse")
			ctx.SetParamNames("id")
			ctx.SetParamValues(strconv.Itoa(testCase.inputRecipientId))

			//Перевірка результатів
			if assert.NoError(t, handler.DeleteFriend(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestUsersHandler_AddToBlackList(t *testing.T) {
	type mockBehavior func(s *mockService.MockStatus, senderId, recipientId int)

	testTable := []struct {
		name                 string
		inputSenderId        int
		inputRecipientId     int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:             "Ok",
			inputSenderId:    13,
			inputRecipientId: 2,
			mockBehavior: func(s *mockService.MockStatus, senderId, recipientId int) {
				status := models.Status{
					SenderId:     senderId,
					RecipientId:  recipientId,
					Relationship: "black_list",
				}
				res := 2
				s.EXPECT().AddStatus(status).Return(res, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":2}` + "\n",
		},
		{
			name:             "Add status error",
			inputSenderId:    13,
			inputRecipientId: 2,
			mockBehavior: func(s *mockService.MockStatus, senderId, recipientId int) {
				status := models.Status{
					SenderId:     senderId,
					RecipientId:  recipientId,
					Relationship: "black_list",
				}
				res := 0
				s.EXPECT().AddStatus(status).Return(res, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"add status error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			status := mockService.NewMockStatus(c)
			testCase.mockBehavior(status, testCase.inputSenderId, testCase.inputRecipientId)

			services := &service.Service{Status: status}
			handler := NewUsersHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodPost, "/api/users/:id/addToBL", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set(middlewares.UserCtx, testCase.inputSenderId)
			ctx.SetPath("/api/users/:id/addToBL")
			ctx.SetParamNames("id")
			ctx.SetParamValues(strconv.Itoa(testCase.inputRecipientId))

			//Перевірка результатів
			if assert.NoError(t, handler.AddToBlackList(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestUsersHandler_DeleteFromBlacklist(t *testing.T) {
	type mockBehavior func(s *mockService.MockStatus, senderId, recipientId int)

	testTable := []struct {
		name                 string
		inputSenderId        int
		inputRecipientId     int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:             "Ok",
			inputSenderId:    13,
			inputRecipientId: 2,
			mockBehavior: func(s *mockService.MockStatus, senderId, recipientId int) {
				status := models.Status{
					SenderId:     senderId,
					RecipientId:  recipientId,
					Relationship: "black_list",
				}
				s.EXPECT().DeleteStatus(status).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"User deleted from black list"}` + "\n",
		},
		{
			name:             "Delete status error",
			inputSenderId:    13,
			inputRecipientId: 2,
			mockBehavior: func(s *mockService.MockStatus, senderId, recipientId int) {
				status := models.Status{
					SenderId:     senderId,
					RecipientId:  recipientId,
					Relationship: "black_list",
				}
				s.EXPECT().DeleteStatus(status).Return(errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"delete status error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			status := mockService.NewMockStatus(c)
			testCase.mockBehavior(status, testCase.inputSenderId, testCase.inputRecipientId)

			services := &service.Service{Status: status}
			handler := NewUsersHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodDelete, "/api/users/:id/deleteFromBlacklist", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set(middlewares.UserCtx, testCase.inputSenderId)
			ctx.SetPath("/api/users/:id/deleteFromBlacklist")
			ctx.SetParamNames("id")
			ctx.SetParamValues(strconv.Itoa(testCase.inputRecipientId))

			//Перевірка результатів
			if assert.NoError(t, handler.DeleteFromBlacklist(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestUsersHandler_SearchUser(t *testing.T) {
	type mockBehavior func(s *mockService.MockStatus, name string)

	testTable := []struct {
		name                 string
		inputName            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			inputName: "fi",
			mockBehavior: func(s *mockService.MockStatus, name string) {
				ret := []models.User{
					{
						Id:       3,
						Username: "first",
					},
					{
						Id:       6,
						Username: "fifth",
					},
				}
				s.EXPECT().SearchUser(name).Return(ret, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"list":[{"id":3,"username":"first","password":"","icon":""},{"id":6,"username":"fifth","password":"","icon":""}]}` + "\n",
		},
		{
			name:      "Search users error",
			inputName: "fi",
			mockBehavior: func(s *mockService.MockStatus, name string) {
				var ret []models.User
				s.EXPECT().SearchUser(name).Return(ret, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"search users error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			// Початкові значення
			// Налаштовуємо логіку оболонок (підключаємо усі рівні)
			c := gomock.NewController(t)
			defer c.Finish()

			status := mockService.NewMockStatus(c)
			testCase.mockBehavior(status, testCase.inputName)

			services := &service.Service{Status: status}
			handler := NewUsersHandler(services)

			//Тестовий сервер
			e := echo.New()

			//Тестовий запит
			req := httptest.NewRequest(http.MethodGet, "/api/users/search/:username", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetPath("/api/users/search/:username")
			ctx.SetParamNames("username")
			ctx.SetParamValues(testCase.inputName)

			//Перевірка результатів
			if assert.NoError(t, handler.SearchUser(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}
