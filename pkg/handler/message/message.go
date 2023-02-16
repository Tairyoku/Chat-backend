package messages

import (
	"cmd/pkg/handler/middlewares"
	"cmd/pkg/handler/responses"
	"cmd/pkg/repository/models"
	"cmd/pkg/service"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type MessageHandler struct {
	services *service.Service
}

func NewMessageHandler(services *service.Service) *MessageHandler {
	return &MessageHandler{services: services}
}

func (h *MessageHandler) CreateMessage(c echo.Context) error {

	// Отримуємо дані з сайту (текст повідомлення)
	var msg models.Message
	if err := c.Bind(&msg); err != nil {
		return err
	}
	if msg.Text == "" {
		responses.NewErrorResponse(c, http.StatusBadRequest, "body is empty")
		return nil
	}

	// Отримуємо ID чату
	chatId, errParam := middlewares.GetParam(c, middlewares.ChatId)
	if errParam != nil {
		return errParam
	}

	// Отримуємо ID активного користувача
	userId, errId := middlewares.GetUserId(c)
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
		responses.NewErrorResponse(c, http.StatusInternalServerError, "create message error")
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

func (h *MessageHandler) GetMessage(c echo.Context) error {
	msgId, errParam := middlewares.GetParam(c, middlewares.ParamId)
	if errParam != nil {
		return errParam
	}

	msg, err := h.services.Message.Get(msgId)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "server error")
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

func (h *MessageHandler) GetAllMessages(c echo.Context) error {
	chatId, errParam := middlewares.GetParam(c, middlewares.ChatId)
	if errParam != nil {
		return errParam
	}

	msg, err := h.services.Message.GetAll(chatId)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "server error")
	}

	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"list": msg,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *MessageHandler) GetLimitMessages(c echo.Context) error {
	chatId, errParam := middlewares.GetParam(c, middlewares.ChatId)
	if errParam != nil {
		return errParam
	}

	limit, errParamId := middlewares.GetParam(c, middlewares.ParamId)
	if errParamId != nil {
		return errParamId
	}

	msg, err := h.services.Message.GetLimit(chatId, limit)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "get limit error")
		return nil
	}
	var result []models.Message
	var length = len(msg) - 1
	for i := range msg {
		result = append(result, msg[length-i])
	}
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"list": result,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}
