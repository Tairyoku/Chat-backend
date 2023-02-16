package handler

import (
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
)

func TestHandler_SignUp(t *testing.T) {
	type mockBehavior func(s *mockService.MockAuthorization, user models.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            models.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			inputBody: `{"username":"test username","password":"password"}`,
			inputUser: models.User{
				Username: "test username",
				Password: "password",
			},
			mockBehavior: func(s *mockService.MockAuthorization, user models.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}` + "\n",
		},
		{
			name:      "Error request data",
			inputBody: "error",
			inputUser: models.User{
				Username: "test username",
				Password: "password",
			},
			mockBehavior: func(r *mockService.MockAuthorization, user models.User) {
				//r.EXPECT().CreateUser(user).Return(0, errors.New("you must enter a username"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"incorrect request data"}` + "\n",
		},
		{
			name:      "Wrong Input UserName",
			inputBody: `{"name": "Test Name", "password": "qwerty"}`,
			inputUser: models.User{
				Password: "password",
			},
			mockBehavior: func(r *mockService.MockAuthorization, user models.User) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"You must enter a username"}` + "\n",
		},
		{
			name:      "Wrong Input Name",
			inputBody: `{"username": "username", "password": "qwerty"}`,
			mockBehavior: func(r *mockService.MockAuthorization, user models.User) {
				//r.EXPECT().CreateUser(user).Return(0, errors.New("invalid input body"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"You must enter a name"}` + "\n",
		},
		{
			name:      "Wrong Input Password",
			inputBody: `{"username": "username", "name": "Test Name"}`,
			mockBehavior: func(r *mockService.MockAuthorization, user models.User) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Password must be at least 6 symbols"}` + "\n",
		},
		{
			name:      "Service Error",
			inputBody: `{"name": "Test","username":"test username","password":"password"}`,
			inputUser: models.User{
				Username: "test username",
				Password: "password",
			},
			mockBehavior: func(r *mockService.MockAuthorization, user models.User) {
				r.EXPECT().CreateUser(user).Return(0, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Начальные значения
			//настраиваем логику оболочек (подключаем все уровни)
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			//Тестовый сервер
			e := echo.New()

			//Тестовый запрос
			req := httptest.NewRequest(http.MethodPost, "/sign-up",
				strings.NewReader(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			//Проверка результатов
			if assert.NoError(t, handler.SignUp(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

//
//func TestHandler_SignIn(t *testing.T) {
//	type mockBehavior func(s *mockService.MockAuthorization, user SignInInput)
//
//	testTable := []struct {
//		name                 string
//		inputBody            string
//		inputUser            SignInInput
//		mockBehavior         mockBehavior
//		expectedStatusCode   int
//		expectedResponseBody string
//	}{
//		{
//			name:      "ok",
//			inputBody: `{"username":"test username","password":"password"}`,
//			inputUser: SignInInput{
//				Username: "test username",
//				Password: "password",
//			},
//			mockBehavior: func(s *mockService.MockAuthorization, user SignInInput) {
//				s.EXPECT().CheckUser(user.Username).Return(nil)
//				s.EXPECT().GenerateToken(user.Username, user.Password).Return("token", nil)
//			},
//			expectedStatusCode:   200,
//			expectedResponseBody: `{"token":"token"}` + "\n",
//		},
//		{
//			name:      "Error request data",
//			inputBody: "error",
//			inputUser: SignInInput{
//				Username: "test username",
//				Password: "password",
//			},
//			mockBehavior: func(s *mockService.MockAuthorization, user SignInInput) {
//			},
//			expectedStatusCode:   400,
//			expectedResponseBody: `{"message":"incorrect request data"}` + "\n",
//		},
//		{
//			name:      "Incorrect username",
//			inputBody: `{"username":"test username","password":"password"}`,
//			inputUser: SignInInput{
//				Username: "test username",
//				Password: "password",
//			},
//			mockBehavior: func(s *mockService.MockAuthorization, user SignInInput) {
//				s.EXPECT().CheckUser(user.Username).Return(errors.New("user not found"))
//			},
//			expectedStatusCode:   404,
//			expectedResponseBody: `{"message":"user not found"}` + "\n",
//		},
//		{
//			name:      "incorrect password",
//			inputBody: `{"username":"test username","password":"password"}`,
//			inputUser: SignInInput{
//				Username: "test username",
//				Password: "password",
//			},
//			mockBehavior: func(s *mockService.MockAuthorization, user SignInInput) {
//				s.EXPECT().CheckUser(user.Username).Return(nil)
//				s.EXPECT().GenerateToken(user.Username, user.Password).Return("", errors.New("incorrect password"))
//			},
//			expectedStatusCode:   400,
//			expectedResponseBody: `{"message":"incorrect password"}` + "\n",
//		},
//	}
//
//	for _, testCase := range testTable {
//		t.Run(testCase.name, func(t *testing.T) {
//
//			//Начальные значения
//			//настраиваем логику оболочек (подключаем все уровни)
//			c := gomock.NewController(t)
//			defer c.Finish()
//
//			auth := mockService.NewMockAuthorization(c)
//			testCase.mockBehavior(auth, testCase.inputUser)
//
//			services := &service.Service{Authorization: auth}
//			handler := NewHandler(services)
//
//			//Тестовый сервер
//			e := echo.New()
//
//			//Тестовый запрос
//			req := httptest.NewRequest(http.MethodPost, "/sign-in",
//				strings.NewReader(testCase.inputBody))
//			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//			rec := httptest.NewRecorder()
//			ctx := e.NewContext(req, rec)
//
//			//Проверка результатов
//			if assert.NoError(t, handler.SignIn(ctx)) {
//				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
//				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
//			}
//		})
//	}
//
//}
