package handler

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	ParamId             = "id"
	ChatId              = "chatId"
	Username            = "username"
	ChatName            = "name"
)

func (h *Handler) userIdentify(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get(authorizationHeader)
		if header == "" {
			NewErrorResponse(c, http.StatusUnauthorized, "empty auth header")
			return nil
		}

		//headerParts := strings.Split(header, " ")
		//if len(headerParts) != 2 {
		//	NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		//	return nil
		//}

		userId, err := h.services.Authorization.ParseToken(header)
		//userId, err := h.services.Authorization.ParseToken(headerParts[1])

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
