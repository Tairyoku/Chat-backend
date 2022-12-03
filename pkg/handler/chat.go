package handler

import (
	"cmd/pkg/repository/models"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) CreateChat(c echo.Context) error {
	var chat models.Chat
	errReq := GetRequest(c, chat)
	if errReq != nil {
		return errReq
	}
	if chat.Name == "" {
		NewErrorResponse(c, http.StatusBadRequest, "name is empty")
		return nil
	}
	id, err := h.services.Chat.Create(chat)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
	}
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) GetChat(c echo.Context) error {
	id, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}

	chat, err := h.services.Chat.Get(id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return nil
	}
	_, errEnCd := json.Marshal(chat)
	if errEnCd != nil {
		return errEnCd
	}
	errRes := c.JSON(http.StatusOK, chat)
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) DeleteChat(c echo.Context) error {
	id, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}

	err := h.services.Chat.Delete(id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return nil
	}

	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("chat with id %d deleted", id),
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) AddUserToChat(c echo.Context) error {
	chatId, errParamC := GetParam(c, ParamId)
	if errParamC != nil {
		return errParamC
	}

	var list models.ChatUsers
	errReq := GetRequest(c, list)
	if errReq != nil {
		return errReq
	}
	list.ChatId = chatId
	id, err := h.services.Chat.AddUser(list)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return nil
	}

	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) DeleteUserFromChat(c echo.Context) error {
	chatId, errParamC := GetParam(c, ParamId)
	if errParamC != nil {
		return errParamC
	}

	userId, errParamU := GetParam(c, userCtx)
	if errParamU != nil {
		return errParamU
	}

	err := h.services.Chat.DeleteUser(chatId, userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return nil
	}

	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("user with id %d deleted from chat with id %d", chatId, userId),
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) GetUsers(c echo.Context) error {
	chatId, errParamC := GetParam(c, ParamId)
	if errParamC != nil {
		return errParamC
	}

	users, err := h.services.Chat.GetUsers(chatId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return nil
	}

	_, errEnCd := json.Marshal(users)
	if errEnCd != nil {
		return errEnCd
	}
	errRes := c.JSON(http.StatusOK, users)
	if errRes != nil {
		return errRes
	}
	return nil
}
