package auth

import (
	"cmd/pkg/handler/middlewares"
	"cmd/pkg/repository/models"
	"cmd/pkg/service"
	mockService "cmd/pkg/service/mocks"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthHandler_SignUp(t *testing.T) {
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
				token := "token"
				s.EXPECT().CreateUser(user).Return(1, nil)
				s.EXPECT().GenerateToken(user.Username, user.Password).Return(token, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"token":"token"}` + "\n",
		},
		{
			name:      "Error request data",
			inputBody: "error",
			inputUser: models.User{
				Username: "test username",
				Password: "password",
			},
			mockBehavior: func(s *mockService.MockAuthorization, user models.User) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"incorrect request data"}` + "\n",
		},
		{
			name:      "Wrong Input UserName",
			inputBody: `{"username": "", "password": "password"}`,
			mockBehavior: func(r *mockService.MockAuthorization, user models.User) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"You must enter a username"}` + "\n",
		},
		{
			name:      "Wrong Input Password",
			inputBody: `{"username": "username", "password": "fifth"}`,
			mockBehavior: func(r *mockService.MockAuthorization, user models.User) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Password must be at least 6 symbols"}` + "\n",
		},
		{
			name:      "Create User Error",
			inputBody: `{"username":"test username","password":"password"}`,
			inputUser: models.User{
				Username: "test username",
				Password: "password",
			},
			mockBehavior: func(r *mockService.MockAuthorization, user models.User) {
				r.EXPECT().CreateUser(user).Return(0, errors.New("username is already used"))
			},
			expectedStatusCode:   409,
			expectedResponseBody: `{"message":"username is already used"}` + "\n",
		},
		{
			name:      "Generate Token Error",
			inputBody: `{"username":"test username","password":"password"}`,
			inputUser: models.User{
				Username: "test username",
				Password: "password",
			},
			mockBehavior: func(s *mockService.MockAuthorization, user models.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
				s.EXPECT().GenerateToken(user.Username, user.Password).Return("", errors.New("generate token error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"generate token error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewAuthHandler(services)

			e := echo.New()

			req := httptest.NewRequest(http.MethodPost, "/auth/sign-up",
				strings.NewReader(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			if assert.NoError(t, handler.SignUp(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestAuthHandler_SignIn(t *testing.T) {
	type mockBehavior func(s *mockService.MockAuthorization, user SignInInput)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            SignInInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			inputBody: `{"username":"test username","password":"password"}`,
			inputUser: SignInInput{
				Username: "test username",
				Password: "password",
			},
			mockBehavior: func(s *mockService.MockAuthorization, user SignInInput) {
				res := models.User{
					Id:       2,
					Username: "test username",
					Icon:     "",
					Password: "",
				}
				s.EXPECT().GetByName(user.Username).Return(res, nil)
				s.EXPECT().GenerateToken(user.Username, user.Password).Return("token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"token"}` + "\n",
		},
		{
			name:      "Error request data",
			inputBody: "error",
			inputUser: SignInInput{
				Username: "test username",
				Password: "password",
			},
			mockBehavior: func(s *mockService.MockAuthorization, user SignInInput) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"incorrect request data"}` + "\n",
		},
		{
			name:      "Check user error",
			inputBody: `{"username":"test username","password":"password"}`,
			inputUser: SignInInput{
				Username: "test username",
				Password: "password",
			},
			mockBehavior: func(s *mockService.MockAuthorization, user SignInInput) {
				res := models.User{
					Id:       0,
					Username: "",
					Icon:     "",
					Password: "",
				}
				s.EXPECT().GetByName(user.Username).Return(res, errors.New("check user error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"check user error"}` + "\n",
		},
		{
			name:      "User not found",
			inputBody: `{"username":"test username","password":"password"}`,
			inputUser: SignInInput{
				Username: "test username",
				Password: "password",
			},
			mockBehavior: func(s *mockService.MockAuthorization, user SignInInput) {
				res := models.User{
					Id:       0,
					Username: "",
					Icon:     "",
					Password: "",
				}
				s.EXPECT().GetByName(user.Username).Return(res, nil)
			},
			expectedStatusCode:   404,
			expectedResponseBody: `{"message":"user not found"}` + "\n",
		},
		{
			name:      "Incorrect password",
			inputBody: `{"username":"test username","password":"password"}`,
			inputUser: SignInInput{
				Username: "test username",
				Password: "password",
			},
			mockBehavior: func(s *mockService.MockAuthorization, user SignInInput) {
				res := models.User{
					Id:       2,
					Username: "test username",
					Icon:     "",
					Password: "",
				}
				s.EXPECT().GetByName(user.Username).Return(res, nil)
				s.EXPECT().GenerateToken(user.Username, user.Password).Return("", errors.New("incorrect password"))
			},
			expectedStatusCode:   409,
			expectedResponseBody: `{"message":"incorrect password"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewAuthHandler(services)

			e := echo.New()

			req := httptest.NewRequest(http.MethodPost, "/auth/sign-in",
				strings.NewReader(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			if assert.NoError(t, handler.SignIn(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestAuthHandler_GetMe(t *testing.T) {
	type mockBehavior func(s *mockService.MockAuthorization, userId int)

	testTable := []struct {
		name                 string
		inputUserId          int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "ok",
			inputUserId: 5,
			mockBehavior: func(s *mockService.MockAuthorization, userId int) {
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":5}` + "\n",
		},
		{
			name:        "User not found",
			inputUserId: 0,
			mockBehavior: func(s *mockService.MockAuthorization, userId int) {
			},
			expectedStatusCode:   404,
			expectedResponseBody: `{"message":"user not found"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUserId)

			services := &service.Service{Authorization: auth}
			handler := NewAuthHandler(services)

			e := echo.New()

			req := httptest.NewRequest(http.MethodPost, "/auth/get-me", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set(middlewares.UserCtx, testCase.inputUserId)

			if assert.NoError(t, handler.GetMe(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestAuthHandler_ChangePassword(t *testing.T) {
	type mockBehavior func(s *mockService.MockAuthorization, userId int, passwords ChangePassword)

	testTable := []struct {
		name                 string
		inputUserId          int
		inputBody            string
		inputPasswords       ChangePassword
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "ok",
			inputUserId: 4,
			inputBody:   `{"old_password":"old password","new_password":"new password"}`,
			inputPasswords: ChangePassword{
				OldPassword: "old password",
				NewPassword: "new password",
			},
			mockBehavior: func(s *mockService.MockAuthorization, userId int, passwords ChangePassword) {
				res := models.User{
					Id:       4,
					Username: "test username",
					Icon:     "",
					Password: "",
				}
				s.EXPECT().GetUserById(userId).Return(res, nil)
				s.EXPECT().GenerateToken(res.Username, passwords.OldPassword).Return("token", nil)
				res.Password = passwords.NewPassword
				s.EXPECT().UpdatePassword(res).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"password changed"}` + "\n",
		},
		{
			name:        "Incorrect request data",
			inputUserId: 4,
			inputBody:   `{"error"}`,

			mockBehavior: func(s *mockService.MockAuthorization, userId int, passwords ChangePassword) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"incorrect request data"}` + "\n",
		},
		{
			name:        "Password length error",
			inputUserId: 4,
			inputBody:   `{"old_password":"old password","new_password":"fifth"}`,
			inputPasswords: ChangePassword{
				OldPassword: "old password",
				NewPassword: "fifth",
			},
			mockBehavior: func(s *mockService.MockAuthorization, userId int, passwords ChangePassword) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"password must be at least 6 symbols"}` + "\n",
		},
		{
			name:        "Incorrect user data",
			inputUserId: 4,
			inputBody:   `{"old_password":"old password","new_password":"new password"}`,
			inputPasswords: ChangePassword{
				OldPassword: "old password",
				NewPassword: "new password",
			},
			mockBehavior: func(s *mockService.MockAuthorization, userId int, passwords ChangePassword) {
				res := models.User{
					Id:       0,
					Username: "",
					Icon:     "",
					Password: "",
				}
				s.EXPECT().GetUserById(userId).Return(res, errors.New("incorrect user data"))
			},
			expectedStatusCode:   404,
			expectedResponseBody: `{"message":"incorrect user data"}` + "\n",
		},
		{
			name:        "Incorrect old password",
			inputUserId: 4,
			inputBody:   `{"old_password":"old password","new_password":"new password"}`,
			inputPasswords: ChangePassword{
				OldPassword: "old password",
				NewPassword: "new password",
			},
			mockBehavior: func(s *mockService.MockAuthorization, userId int, passwords ChangePassword) {
				res := models.User{
					Id:       4,
					Username: "test username",
					Icon:     "",
					Password: "",
				}
				s.EXPECT().GetUserById(userId).Return(res, nil)
				s.EXPECT().GenerateToken(res.Username, passwords.OldPassword).Return("", errors.New("incorrect password"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"incorrect password"}` + "\n",
		},
		{
			name:        "Update password error",
			inputUserId: 4,
			inputBody:   `{"old_password":"old password","new_password":"new password"}`,
			inputPasswords: ChangePassword{
				OldPassword: "old password",
				NewPassword: "new password",
			},
			mockBehavior: func(s *mockService.MockAuthorization, userId int, passwords ChangePassword) {
				res := models.User{
					Id:       4,
					Username: "test username",
					Icon:     "",
					Password: "",
				}
				s.EXPECT().GetUserById(userId).Return(res, nil)
				s.EXPECT().GenerateToken(res.Username, passwords.OldPassword).Return("token", nil)
				res.Password = passwords.NewPassword
				s.EXPECT().UpdatePassword(res).Return(errors.New("update password error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"update password error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUserId, testCase.inputPasswords)

			services := &service.Service{Authorization: auth}
			handler := NewAuthHandler(services)

			e := echo.New()

			req := httptest.NewRequest(http.MethodPost, "/auth/change/password",
				strings.NewReader(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set(middlewares.UserCtx, testCase.inputUserId)

			if assert.NoError(t, handler.ChangePassword(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestAuthHandler_ChangeUsername(t *testing.T) {
	type mockBehavior func(s *mockService.MockAuthorization, userId int, user models.User)

	testTable := []struct {
		name                 string
		inputUserId          int
		inputBody            string
		inputUserName        models.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "ok",
			inputUserId: 4,
			inputBody:   `{"username":"new username"}`,
			inputUserName: models.User{
				Username: "new username",
			},
			mockBehavior: func(s *mockService.MockAuthorization, userId int, user models.User) {
				res := models.User{
					Id:       4,
					Username: "test username",
					Icon:     "",
					Password: "",
				}
				check := models.User{}
				s.EXPECT().GetUserById(userId).Return(res, nil)
				s.EXPECT().GetByName(user.Username).Return(check, errors.New("record not found"))
				res.Username = user.Username
				s.EXPECT().UpdateData(res).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"username changed"}` + "\n",
		},
		{
			name:        "Incorrect request data",
			inputUserId: 4,
			inputBody:   `{"error"}`,
			mockBehavior: func(s *mockService.MockAuthorization, userId int, user models.User) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"incorrect request data"}` + "\n",
		},
		{
			name:        "Incorrect user data",
			inputUserId: 4,
			inputBody:   `{"username":"new username"}`,
			inputUserName: models.User{
				Username: "new username",
			},
			mockBehavior: func(s *mockService.MockAuthorization, userId int, user models.User) {
				res := models.User{}
				s.EXPECT().GetUserById(userId).Return(res, errors.New("incorrect user data"))

			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"incorrect user data"}` + "\n",
		},
		{
			name:        "Username is used",
			inputUserId: 4,
			inputBody:   `{"username":"new username"}`,
			inputUserName: models.User{
				Username: "new username",
			},
			mockBehavior: func(s *mockService.MockAuthorization, userId int, user models.User) {
				res := models.User{
					Id:       4,
					Username: "test username",
					Icon:     "",
					Password: "",
				}
				check := models.User{
					Id:       2,
					Username: "new username",
				}
				s.EXPECT().GetUserById(userId).Return(res, nil)
				s.EXPECT().GetByName(user.Username).Return(check, nil)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"username is used"}` + "\n",
		},
		{
			name:        "Update username error",
			inputUserId: 4,
			inputBody:   `{"username":"new username"}`,
			inputUserName: models.User{
				Username: "new username",
			},
			mockBehavior: func(s *mockService.MockAuthorization, userId int, user models.User) {
				res := models.User{
					Id:       4,
					Username: "test username",
					Icon:     "",
					Password: "",
				}
				check := models.User{}
				s.EXPECT().GetUserById(userId).Return(res, nil)
				s.EXPECT().GetByName(user.Username).Return(check, errors.New("record is not found"))
				res.Username = user.Username
				s.EXPECT().UpdateData(res).Return(errors.New("update username error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"update username error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUserId, testCase.inputUserName)

			services := &service.Service{Authorization: auth}
			handler := NewAuthHandler(services)

			e := echo.New()

			req := httptest.NewRequest(http.MethodPost, "/auth/change/username",
				strings.NewReader(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set(middlewares.UserCtx, testCase.inputUserId)

			if assert.NoError(t, handler.ChangeUsername(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestAuthHandler_ChangeIcon(t *testing.T) {
	type mockBehavior func(s *mockService.MockAuthorization, userId int, filename string)

	testTable := []struct {
		name                 string
		inputUserId          int
		inputFilename        string
		inputFile            multipart.File
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:          "ok",
			inputUserId:   4,
			inputFilename: "filename",
			mockBehavior: func(s *mockService.MockAuthorization, userId int, filename string) {
				res := models.User{
					Id:       4,
					Username: "test username",
					Icon:     "",
					Password: "",
				}
				s.EXPECT().GetUserById(userId).Return(res, nil)
				res.Icon = filename
				s.EXPECT().UpdateData(res).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"icon changed"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUserId, testCase.inputFilename)

			services := &service.Service{Authorization: auth}
			handler := NewAuthHandler(services)

			e := echo.New()

			req := httptest.NewRequest(http.MethodPut, "/auth/change/username", testCase.inputFile)
			req.Header.Set(echo.HeaderContentType, echo.MIMEMultipartForm)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set(middlewares.UserCtx, testCase.inputUserId)

			if assert.NoError(t, handler.ChangeIcon(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}
