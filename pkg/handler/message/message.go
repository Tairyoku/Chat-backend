package messages

import (
	"cmd/pkg/handler/middlewares"
	"cmd/pkg/handler/responses"
	"cmd/pkg/repository/models"
	"cmd/pkg/service"
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

// CreateMessage godoc
// @Summary      Create message
// @Description  Отримує текст повідомлення.
// @Description  Створює повідомлення.
// @Security ApiKeyAuth
// @Tags         message
// @Accept       json
// @Produce      json
// @Param        chatId		path     int   true  "Chat ID"
// @Param        message_text	body     TextInput   true  "Message text"
// @Success      200 	{object} IdResponse			"return message ID"
// @Failure 	 400 	{object} responses.ErrorResponse	 "body is empty"
// @Failure 	 500 	{object} responses.ErrorResponse	 "create message error"
// @Router       /chats/{chatId}/messages [post]
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
	msg.SentAt = time.Now().Round(10 * time.Millisecond)

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

// GetMessage godoc
// @Summary      Get message by ID
// @Description  Отримує ID повідомлення.
// @Description  Повертає повідомлення.
// @Security ApiKeyAuth
// @Tags         message
// @Accept       json
// @Produce      json
// @Param        chatId		path     int   true  "Chat ID"
// @Param        id		path     int   true  "Message ID"
// @Success      200 	{object} MessageResponse			"return message ID"
// @Failure 	 500 	{object} responses.ErrorResponse	 "get message error"
// @Router       /chats/{chatId}/messages/{id} [get]
func (h *MessageHandler) GetMessage(c echo.Context) error {
	msgId, errParam := middlewares.GetParam(c, middlewares.ParamId)
	if errParam != nil {
		return errParam
	}

	msg, err := h.services.Message.Get(msgId)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "get message error")
	}

	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"list": msg,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

// GetLimitMessages godoc
// @Summary      Get limited message list
// @Description  Отримує ID чату.
// @Description  Повертає певну кількість повідомлень чату.
// @Security ApiKeyAuth
// @Tags         message
// @Accept       json
// @Produce      json
// @Param        chatId		path     int   true  "Chat ID"
// @Param        id		path     int   true  "Message ID"
// @Success      200 	{object} MessageListResponse			"return message ID"
// @Failure 	 500 	{object} responses.ErrorResponse	 "get limit error"
// @Router       /chats/{chatId}/messages/limit/{id} [get]
func (h *MessageHandler) GetLimitMessages(c echo.Context) error {

	// Отримуємо ID чату
	chatId, errParam := middlewares.GetParam(c, middlewares.ChatId)
	if errParam != nil {
		return errParam
	}

	// Отримуємо кількість необхідних повідомлень
	limit, errParamId := middlewares.GetParam(c, middlewares.ParamId)
	if errParamId != nil {
		return errParamId
	}

	// Отримуємо список повідомлень зі зворотним порядком
	msg, err := h.services.Message.GetLimit(chatId, limit)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "get limit error")
		return nil
	}

	// Повертаємо правильний порядок
	var result []models.Message
	var length = len(msg) - 1
	for i := range msg {
		result = append(result, msg[length-i])
	}

	// Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"list": result,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}
