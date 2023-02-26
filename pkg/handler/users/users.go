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

// GetUserById godoc
// @Summary      Get user`s data by ID
// @Description  Отримує ID користувача.
// @Description  Повертає дані користувача.
// @Security ApiKeyAuth
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id		path     int   true  "User ID"
// @Success      200 	{object} UserResponse			"return user`s data"
// @Failure 	 500 	{object} responses.ErrorResponse	 "get user error"
// @Router       /users/{id} [get]
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

// GetUserLists godoc
// @Summary      Get user`s relationship lists by ID
// @Description  Отримує ID користувача.
// @Description  Повертає списки відносин між користувачем та іншими користувачами.
// @Security ApiKeyAuth
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id		path     int   true  "User ID"
// @Success      200 	{object} StatusesListResponse		"return user`s status lists"
// @Failure 	 500 	{object} responses.ErrorResponse	 "friends list error"
// @Failure 	 500 	{object} responses.ErrorResponse	 "black list error"
// @Failure 	 500 	{object} responses.ErrorResponse	 "on black list error"
// @Failure 	 500 	{object} responses.ErrorResponse	 "friend invites list error"
// @Failure 	 500 	{object} responses.ErrorResponse	 "friend requires list error"
// @Router       /users/{id}/all [get]
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

// InvitedToFriends godoc
// @Summary      Create friendship invitation
// @Description  Отримує ID користувача.
// @Description  Створює запит на дружбу поточного користувача з отриманим.
// @Security ApiKeyAuth
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id		path     int   true  "User ID"
// @Success      200 	{object} IdResponse			"require is sent"
// @Failure 	 500 	{object} responses.ErrorResponse	 "add status error"
// @Router       /users/{id}/invite [post]
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

// CancelInvite godoc
// @Summary      Cancel friendship invitation
// @Description  Отримує ID користувача.
// @Description  Видаляє запит на дружбу поточного користувача з отриманим.
// @Security ApiKeyAuth
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id		path     int   true  "User ID"
// @Success      200 	{object} MessageResponse			"invite deleted"
// @Failure 	 500 	{object} responses.ErrorResponse	 "delete status error"
// @Router       /users/{id}/cancel [delete]
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

// AcceptInvitation godoc
// @Summary      Accept friendship invitation
// @Description  Отримує ID користувача.
// @Description  Підтверджує запит на дружбу поточного користувача з отриманим.
// @Security ApiKeyAuth
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id		path     int   true  "User ID"
// @Success      200 	{object} MessageResponse			"invitation accepted"
// @Failure 	 500 	{object} responses.ErrorResponse	 "update status error"
// @Router       /users/{id}/accept [put]
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
		"message": "invitation accepted",
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

// RefuseInvitation godoc
// @Summary      Refuse friendship invitation
// @Description  Отримує ID користувача.
// @Description  Відхиляє запит на дружбу поточного користувача з отриманим.
// @Security ApiKeyAuth
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id		path     int   true  "User ID"
// @Success      200 	{object} MessageResponse			"invitation refused"
// @Failure 	 500 	{object} responses.ErrorResponse	 "delete status error"
// @Router       /users/{id}/refuse [put]
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

// DeleteFriend godoc
// @Summary      Delete friend
// @Description  Отримує ID користувача.
// @Description  Видаляє користувача зі списку друзів.
// @Security ApiKeyAuth
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id		path     int   true  "User ID"
// @Success      200 	{object} MessageResponse			"friend deleted"
// @Failure 	 500 	{object} responses.ErrorResponse	 "delete status error"
// @Router       /users/{id}/deleteFriend [delete]
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

// AddToBlackList godoc
// @Summary      Block user
// @Description  Отримує ID користувача.
// @Description  Блокує користувача.
// @Security ApiKeyAuth
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id		path     int   true  "User ID"
// @Success      200 	{object} IdResponse			"user is blocked"
// @Failure 	 500 	{object} responses.ErrorResponse	 "add status error"
// @Router       /users/{id}/addToBL [post]
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

// DeleteFromBlacklist godoc
// @Summary      Delete from blacklist
// @Description  Отримує ID користувача.
// @Description  Видаляє користувача зі списку заблокованих.
// @Security ApiKeyAuth
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id		path     int   true  "User ID"
// @Success      200 	{object} MessageResponse			"user deleted from black list"
// @Failure 	 500 	{object} responses.ErrorResponse	 "delete status error"
// @Router       /users/{id}/deleteFromBlacklist [delete]
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
		"message": "user deleted from black list",
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

// SearchUser godoc
// @Summary      Get found chats
// @Description  Отримує частину імені чату.
// @Description  Повертає список користувачів, ім'я яких повністю або частково збігається.
// @Security ApiKeyAuth
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        username		path     string   true  "Slice of users"
// @Success      200 	{object} ListResponse			"lisl of found chats"
// @Failure 	 500 	{object} responses.ErrorResponse	 "search users error"
// @Router       /users/search/{username} [get]
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
