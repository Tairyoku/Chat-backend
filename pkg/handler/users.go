package handler

import (
	"cmd/pkg/repository"
	"cmd/pkg/repository/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) GetUserById(c echo.Context) error {

	// Отримуємо ID користувача
	userId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}

	// Отримуємо дані користувача
	user, err := h.services.Authorization.GetUserById(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
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

func (h *Handler) GetUserLists(c echo.Context) error {

	// Отримуємо ID користувача
	userId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}

	// Отримуємо список друзів
	friends, errFr := h.services.Status.GetFriends(userId)
	if errFr != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "friends list error")
		return errFr
	}

	// Отримуємо список заблокованих користувачів
	bl, errBL := h.services.Status.GetBlackList(userId)
	if errBL != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "black list error")
		return errBL
	}

	// Отримуємо список користувачів, що заблокували користувача
	onBL, errOnBL := h.services.Status.GetBlackListToUser(userId)
	if errOnBL != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "on black list error")
		return errOnBL
	}

	// Отримуємо список користувачів, яким відправлено запрошення в друзі
	invites, errInv := h.services.Status.GetSentInvites(userId)
	if errInv != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "friend invites list error")
		return errInv
	}

	// Отримуємо список користувачів, які отримали від користувача запрошення у друзі
	requires, errReq := h.services.Status.GetInvites(userId)
	if errReq != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "friend requires list error")
		return errReq
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

func (h *Handler) InvitedToFriends(c echo.Context) error {

	// Отримуємо ID активного користувача
	senderId, errId := GetUserId(c)
	if errId != nil {
		return errId
	}

	// Отримуємо ID запрошуваного користувача
	recipientId, errParam := GetParam(c, ParamId)
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
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
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

func (h *Handler) CancelInvite(c echo.Context) error {

	// Отримуємо ID активного користувача
	senderId, errId := GetUserId(c)
	if errId != nil {
		return errId
	}

	// Отримуємо ID запрошуваного користувача
	recipientId, errParam := GetParam(c, ParamId)
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
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
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

func (h *Handler) AcceptInvitation(c echo.Context) error {

	// Отримуємо ID активного користувача
	recipientId, errId := GetUserId(c)
	if errId != nil {
		return errId
	}

	// Отримуємо ID запрошуваного користувача
	senderId, errParam := GetParam(c, ParamId)
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
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
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

func (h *Handler) RefuseInvitation(c echo.Context) error {

	// Отримуємо ID активного користувача
	recipientId, errId := GetUserId(c)
	if errId != nil {
		return errId
	}

	// Отримуємо ID запрошуваного користувача
	senderId, errParam := GetParam(c, ParamId)
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
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
	}

	// Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Invitation refused",
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) DeleteFriend(c echo.Context) error {

	// Отримуємо ID активного користувача
	senderId, errId := GetUserId(c)
	if errId != nil {
		return errId
	}

	// Отримуємо ID запрошуваного користувача
	recipientId, errParam := GetParam(c, ParamId)
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
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
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

func (h *Handler) AddToBlackList(c echo.Context) error {

	// Отримуємо ID активного користувача
	senderId, errId := GetUserId(c)
	if errId != nil {
		return errId
	}

	// Отримуємо ID запрошуваного користувача
	recipientId, errParam := GetParam(c, ParamId)
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
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
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

func (h *Handler) DeleteFromBlacklist(c echo.Context) error {

	// Отримуємо ID активного користувача
	senderId, errId := GetUserId(c)
	if errId != nil {
		return errId
	}

	// Отримуємо ID запрошуваного користувача
	recipientId, errParam := GetParam(c, ParamId)
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
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
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

func (h *Handler) SearchUser(c echo.Context) error {

	// Отримуємо фрагмент імені користувача
	username := c.Param(Username)
	if len(username) == 0 {
		return nil
	}

	// Отримуємо список користувачів, що мають в імені отриманий фрагмент
	users, err := h.services.Status.SearchUser(username)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return err
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
