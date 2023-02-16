package users

import (
	"cmd/pkg/handler/middlewares"
	"cmd/pkg/handler/responses"
	"cmd/pkg/repository"
	"cmd/pkg/repository/models"
	"cmd/pkg/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UsersHandler struct {
	services *service.Service
}

func NewUsersHandler(services *service.Service) *UsersHandler {
	return &UsersHandler{services: services}
}

func (h *UsersHandler) GetUserById(c echo.Context) error {

	// Отримуємо ID користувача
	userId, errParam := middlewares.GetParam(c, middlewares.ParamId)
	if errParam != nil {
		return errParam
	}

	// Отримуємо дані користувача
	user, err := h.services.Status.GetUserById(userId)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "get user error")
		return nil
	}

	// Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *UsersHandler) GetUserLists(c echo.Context) error {

	// Отримуємо ID користувача
	userId, errParam := middlewares.GetParam(c, middlewares.ParamId)
	if errParam != nil {
		return errParam
	}

	// Отримуємо список друзів
	friends, errFr := h.services.Status.GetFriends(userId)
	if errFr != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "friends list error")
		return nil
	}

	// Отримуємо список заблокованих користувачів
	bl, errBL := h.services.Status.GetBlackList(userId)
	if errBL != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "black list error")
		return nil
	}

	// Отримуємо список користувачів, що заблокували користувача
	onBL, errOnBL := h.services.Status.GetBlackListToUser(userId)
	if errOnBL != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "on black list error")
		return nil
	}

	// Отримуємо список користувачів, яким відправлено запрошення в друзі
	invites, errInv := h.services.Status.GetSentInvites(userId)
	if errInv != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "friend invites list error")
		return nil
	}

	// Отримуємо список користувачів, які отримали від користувача запрошення у друзі
	requires, errReq := h.services.Status.GetInvites(userId)
	if errReq != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "friend requires list error")
		return nil
	}

	//Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"friends":     friends,
		"blacklist":   bl,
		"onBlacklist": onBL,
		"invites":     invites,
		"requires":    requires,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *UsersHandler) InvitedToFriends(c echo.Context) error {

	// Отримуємо ID активного користувача
	senderId := c.Get(middlewares.UserCtx).(int)

	// Отримуємо ID запрошуваного користувача
	recipientId, errParam := middlewares.GetParam(c, middlewares.ParamId)
	if errParam != nil {
		return errParam
	}

	// Заповнюємо модель відносин
	var status = models.Status{
		SenderId:     senderId,
		RecipientId:  recipientId,
		Relationship: repository.StatusInvitation,
	}

	//Створюємо нові відносини
	id, err := h.services.Status.AddStatus(status)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "add status error")
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

func (h *UsersHandler) CancelInvite(c echo.Context) error {

	// Отримуємо ID активного користувача
	senderId := c.Get(middlewares.UserCtx).(int)

	// Отримуємо ID запрошуваного користувача
	recipientId, errParam := middlewares.GetParam(c, middlewares.ParamId)
	if errParam != nil {
		return errParam
	}

	// Заповнюємо модель відносин
	var status = models.Status{
		SenderId:     senderId,
		RecipientId:  recipientId,
		Relationship: repository.StatusInvitation,
	}

	// Видаляємо відносини за моделлю
	err := h.services.Status.DeleteStatus(status)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "delete status error")
		return nil
	}

	// Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": "invite deleted",
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *UsersHandler) AcceptInvitation(c echo.Context) error {

	// Отримуємо ID активного користувача
	recipientId := c.Get(middlewares.UserCtx).(int)

	// Отримуємо ID запрошуваного користувача
	senderId, errParam := middlewares.GetParam(c, middlewares.ParamId)
	if errParam != nil {
		return errParam
	}

	// Заповнюємо модель відносин
	var status = models.Status{
		SenderId:     senderId,
		RecipientId:  recipientId,
		Relationship: repository.StatusFriends,
	}

	// Оновлюємо відносини за моделлю
	err := h.services.Status.UpdateStatus(status)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "update status error")
		return nil
	}

	// Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Invitation accepted",
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *UsersHandler) RefuseInvitation(c echo.Context) error {

	// Отримуємо ID активного користувача
	recipientId := c.Get(middlewares.UserCtx).(int)

	// Отримуємо ID запрошуваного користувача
	senderId, errParam := middlewares.GetParam(c, middlewares.ParamId)
	if errParam != nil {
		return errParam
	}

	// Заповнюємо модель відносин
	var status = models.Status{
		SenderId:     senderId,
		RecipientId:  recipientId,
		Relationship: repository.StatusInvitation,
	}

	// Видаляємо відносини за моделлю
	err := h.services.Status.DeleteStatus(status)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "delete status error")
		return nil
	}

	// Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": "invitation refused",
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *UsersHandler) DeleteFriend(c echo.Context) error {

	// Отримуємо ID активного користувача
	senderId := c.Get(middlewares.UserCtx).(int)

	// Отримуємо ID запрошуваного користувача
	recipientId, errParam := middlewares.GetParam(c, middlewares.ParamId)
	if errParam != nil {
		return errParam
	}

	// Заповнюємо модель відносин
	var status = models.Status{
		SenderId:     senderId,
		RecipientId:  recipientId,
		Relationship: repository.StatusFriends,
	}

	// Видаляємо відносини за моделлю
	err := h.services.Status.DeleteStatus(status)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "delete status error")
		return nil
	}

	// Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": "friend deleted",
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *UsersHandler) AddToBlackList(c echo.Context) error {

	// Отримуємо ID активного користувача
	senderId := c.Get(middlewares.UserCtx).(int)

	// Отримуємо ID запрошуваного користувача
	recipientId, errParam := middlewares.GetParam(c, middlewares.ParamId)
	if errParam != nil {
		return errParam
	}

	// Заповнюємо модель відносин
	var status = models.Status{
		SenderId:     senderId,
		RecipientId:  recipientId,
		Relationship: repository.StatusBL,
	}

	//Створюємо нові відносини
	id, err := h.services.Status.AddStatus(status)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "add status error")
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

func (h *UsersHandler) DeleteFromBlacklist(c echo.Context) error {

	// Отримуємо ID активного користувача
	senderId := c.Get(middlewares.UserCtx).(int)

	// Отримуємо ID запрошуваного користувача
	recipientId, errParam := middlewares.GetParam(c, middlewares.ParamId)
	if errParam != nil {
		return errParam
	}

	// Заповнюємо модель відносин
	var status = models.Status{
		SenderId:     senderId,
		RecipientId:  recipientId,
		Relationship: repository.StatusBL,
	}

	// Видаляємо відносини за моделлю
	err := h.services.Status.DeleteStatus(status)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "delete status error")
		return nil
	}

	// Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User deleted from black list",
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *UsersHandler) SearchUser(c echo.Context) error {

	// Отримуємо фрагмент імені користувача
	username := c.Param(middlewares.Username)
	if len(username) == 0 {
		return nil
	}

	// Отримуємо список користувачів, що мають в імені отриманий фрагмент
	users, err := h.services.Status.SearchUser(username)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "search users error")
		return nil
	}

	// Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"list": users,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}
