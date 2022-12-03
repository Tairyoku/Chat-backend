package handler

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	ParamId             = "id"
	ChatId              = "chatId"
)

func (h *Handler) userIdentify(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get(authorizationHeader)
		if header == "" {
			NewErrorResponse(c, http.StatusUnauthorized, "empty auth header")
			return nil
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
			return nil
		}

		userId, err := h.services.Authorization.ParseToken(headerParts[1])

		if err != nil {
			NewErrorResponse(c, http.StatusUnauthorized, err.Error())
			return nil
		}
		c.Set(userCtx, userId)
		return next(c)
	}
}

func GetUserId(c echo.Context) (int, error) {
	id := c.Get(userCtx)
	if id == 0 {
		NewErrorResponse(c, http.StatusNotFound, "user id not found")
		return 0, errors.New("user id not found")
	}
	idInt, ok := id.(int)
	if !ok {
		NewErrorResponse(c, http.StatusBadRequest, "user id is of valid type")
		return 0, errors.New("user id is of valid type")
	}
	return idInt, nil
}

func GetParam(c echo.Context, name string) (int, error) {
	param, errReq := strconv.Atoi(c.Param(name))
	if errReq != nil {
		NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%s is not integer", name))
		return 0, errReq
	}
	return param, nil
}

func GetRequest(c echo.Context, i interface{}) error {
	if err := c.Bind(&i); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return err
	}
	return nil
}

//
//func SetupConfig() *oauth2.Config {
//	googleOauthConfig := &oauth2.Config{
//		ClientID:     "64460222459-j6mfme4oj54o4vasuifrlcip3l5ohdkk.apps.googleusercontent.com",
//		ClientSecret: "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",
//		Endpoint:     google.Endpoint,
//		RedirectURL:  "http://localhost:8080/google/callback",
//		Scopes: []string{
//			"https://www.googleapis.com/auth/userinfo.email",   //See your primary Google Account email address
//			"https://www.googleapis.com/auth/userinfo.profile", //See your personal info, including any personal info you've made publicly available
//		},
//	}
//	return googleOauthConfig
//}
//
//func (h *Handler) GoogleSignUp(c echo.Context, input models.User) error {
//	id, err := h.services.Authorization.CreateUser(input)
//
//	if err != nil {
//		NewErrorResponse(c, http.StatusInternalServerError, "something went wrong")
//		return nil
//	}
//	errRes := c.JSON(http.StatusOK, map[string]interface{}{
//		"id": id,
//	})
//	if errRes != nil {
//		return errRes
//	}
//	return nil
//}
//
//func (h *Handler) GoogleSignIn(c echo.Context, input models.User) error {
//	// if username is required, to login and to generate token
//	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
//	if err != nil {
//		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
//		return nil
//	}
//	errRes := c.JSON(http.StatusOK, map[string]interface{}{
//		"token": token,
//	})
//	if errRes != nil {
//		return errRes
//	}
//	return nil
//}
