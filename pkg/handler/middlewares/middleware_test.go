package middlewares

import (
	"cmd/pkg/service"
	mockService "cmd/pkg/service/mocks"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_userIdentify(t *testing.T) {
	type mockBehavior func(s *mockService.MockAuthorization, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "OK",
			headerName: "Authorization",
			token:      "token",
			mockBehavior: func(s *mockService.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1" + "\n",
		},
		{
			name:       "Empty auth header",
			headerName: "Authorization",
			token:      "",
			mockBehavior: func(s *mockService.MockAuthorization, token string) {
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"empty auth header"}` + "\n",
		},
		{
			name:       "Empty auth header",
			headerName: "Authorization",
			token:      "token",
			mockBehavior: func(s *mockService.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(0, errors.New("some error")).AnyTimes()
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"create token error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.token)

			services := &service.Service{Authorization: auth}
			handler := NewMiddlewareHandler(services)

			e := echo.New()
			e.GET("/protected", nil, handler.UserIdentify, func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					id := c.Get(UserCtx).(int)
					c.JSON(200, id)
					return nil
				}
			})
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(testCase.headerName, testCase.token)

			e.ServeHTTP(rec, req)

			assert.Equal(t, testCase.expectedStatusCode, rec.Code)
			assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
		})
	}
}

func Test_GetParam(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	ctx := e.NewContext(nil, rec)
	ctx.SetPath("/example/:value")
	ctx.SetParamNames("value")

	want := 4
	ctx.SetParamValues(fmt.Sprintf("%d", want))

	ok := struct {
		param int
		err   error
	}{}

	ok.param, ok.err = GetParam(ctx, "value")
	if ok.err != nil {
		t.Error("FAILED. Value of param dont have only numbers or null")
	} else if ok.param != want {
		t.Errorf("FAILED. Exepted %d, got %d", want, ok.param)
	} else {
		t.Logf("PASSED. Exepted %d, got %d", want, ok.param)
	}
}

func Test_GetUserId(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	ctx := e.NewContext(nil, rec)
	want := 4
	ctx.Set(UserCtx, want)

	ok := struct {
		param int
		err   error
	}{}

	ok.param, ok.err = GetUserId(ctx)
	if ok.err != nil {
		t.Error("FAILED. User not found or value is not a numbers")
	} else if ok.param != want {
		t.Errorf("FAILED. Exepted %d, got %d", want, ok.param)
	} else {
		t.Logf("PASSED. Exepted %d, got %d", want, ok.param)
	}
}
