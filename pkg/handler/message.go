package handler

import (
	"cmd/pkg/repository/models"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func (h *Handler) CreateMessage(c echo.Context) error {

	// Отримуємо дані з сайту (текст повідомлення)
	var msg models.Message
	if err := c.Bind(&msg); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return nil
	}
	if msg.Text == "" {
		NewErrorResponse(c, http.StatusBadRequest, "body is empty")
		return nil
	}

	// Отримуємо ID чату
	chatId, errParam := GetParam(c, ChatId)
	if errParam != nil {
		return errParam
	}

	// Отримуємо ID активного користувача
	userId, errId := GetUserId(c)
	if errId != nil {
		return errId
	}

	// Заповнюємо форму повідомлення
	msg.ChatId = chatId
	msg.Author = userId
	msg.SentAt = time.Now()

	// Створюємо нове повідомлення
	id, err := h.services.Message.Create(msg)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return nil
	}

	// Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) GetMessage(c echo.Context) error {
	msgId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}

	msg, err := h.services.Message.Get(msgId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
	}
	_, errEnCd := json.Marshal(msg)
	if errEnCd != nil {
		return errEnCd
	}
	errRes := c.JSON(http.StatusOK, msg)
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) GetAllMessages(c echo.Context) error {
	chatId, errParam := GetParam(c, ChatId)
	if errParam != nil {
		return errParam
	}

	msg, err := h.services.Message.GetAll(chatId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
	}

	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"list": msg,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) GetLimitMessages(c echo.Context) error {
	chatId, errParam := GetParam(c, ChatId)
	if errParam != nil {
		return errParam
	}

	limit, errParamId := GetParam(c, ParamId)
	if errParamId != nil {
		return errParamId
	}

	msg, err := h.services.Message.GetLimit(chatId, limit)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
	}

	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"list": msg,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}
