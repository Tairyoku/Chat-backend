package handler

import (
	"cmd/pkg/repository"
	"cmd/pkg/repository/models"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) GetUserById(c echo.Context) error {
	userId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}
	user, err := h.services.Authorization.GetUserById(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
	}

	_, errEnCd := json.Marshal(user)
	if errEnCd != nil {
		return errEnCd
	}
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) GetUserPublicChats(c echo.Context) error {
	userId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}
	chats, err := h.services.Chat.GetPublicChats(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
	}

	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"list": chats,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}
func (h *Handler) GetUserPrivateChats(c echo.Context) error {
	userId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}
	chats, err := h.services.Chat.GetPrivateChats(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
	}

	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"list": chats,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) InvitedToFriends(c echo.Context) error {
	recipientId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}
	senderId, errId := GetUserId(c)
	if errId != nil {
		return errId
	}
	var status = models.Status{
		SenderId:     senderId,
		RecipientId:  recipientId,
		Relationship: repository.StatusInvitation,
	}
	id, err := h.services.Status.AddStatus(status)
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

func (h *Handler) DeleteFriend(c echo.Context) error {
	recipientId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}
	senderId, errId := GetUserId(c)
	if errId != nil {
		return errId
	}
	var status = models.Status{
		SenderId:     senderId,
		RecipientId:  recipientId,
		Relationship: repository.StatusFriends,
	}
	err := h.services.Status.DeleteStatus(status)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
	}

	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": "friend deleted",
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) CancelInvite(c echo.Context) error {
	recipientId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}
	senderId, errId := GetUserId(c)
	if errId != nil {
		return errId
	}
	var status = models.Status{
		SenderId:     senderId,
		RecipientId:  recipientId,
		Relationship: repository.StatusInvitation,
	}
	err := h.services.Status.DeleteStatus(status)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
	}

	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": "invite deleted",
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) RefuseInvitation(c echo.Context) error {
	senderId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}
	recipientId, errId := GetUserId(c)
	if errId != nil {
		return errId
	}
	var status = models.Status{
		SenderId:     senderId,
		RecipientId:  recipientId,
		Relationship: repository.StatusInvitation,
	}
	err := h.services.Status.DeleteStatus(status)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
	}

	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Invitation refused",
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) AcceptInvitation(c echo.Context) error {
	senderId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}
	recipientId, errId := GetUserId(c)
	if errId != nil {
		return errId
	}

	var status = models.Status{
		SenderId:     senderId,
		RecipientId:  recipientId,
		Relationship: repository.StatusFriends,
	}
	err := h.services.Status.UpdateStatus(status)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
	}

	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Invitation accepted",
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) AddToBlackList(c echo.Context) error {
	recipientId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}
	senderId, errId := GetUserId(c)
	if errId != nil {
		return errId
	}
	var status = models.Status{
		SenderId:     senderId,
		RecipientId:  recipientId,
		Relationship: repository.StatusBL,
	}

	id, err := h.services.Status.AddStatus(status)
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

func (h *Handler) DeleteFromBlacklist(c echo.Context) error {
	recipientId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}
	senderId, errId := GetUserId(c)
	if errId != nil {
		return errId
	}
	var status = models.Status{
		SenderId:     senderId,
		RecipientId:  recipientId,
		Relationship: repository.StatusBL,
	}
	err := h.services.Status.DeleteStatus(status)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
	}

	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User deleted from black list",
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) GetFriends(c echo.Context) error {
	userId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}
	friends, err := h.services.Status.GetFriends(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return err
	}
	var list []models.User
	for _, number := range friends {
		user, errUser := h.services.Authorization.GetUserById(number)
		if errUser != nil {
			return errParam
		}
		list = append(list, user)
	}
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"list": list,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) GetBlackList(c echo.Context) error {
	userId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}
	bl, err := h.services.Status.GetBlackList(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return err
	}
	var list []models.User
	for _, number := range bl {
		user, errUser := h.services.Authorization.GetUserById(number)
		if errUser != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "error get user")
			return errParam
		}
		list = append(list, user)
	}
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"list": list,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) GetBlackListToUser(c echo.Context) error {
	userId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}
	ubl, err := h.services.Status.GetBlackListToUser(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return err
	}

	var list []models.User
	for _, number := range ubl {
		user, errUser := h.services.Authorization.GetUserById(number)
		if errUser != nil {
			return errParam
		}
		list = append(list, user)
	}
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"list": list,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) GetSentInvites(c echo.Context) error {
	userId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}
	invited, err := h.services.Status.GetSentInvites(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return err
	}

	var list []models.User
	for _, number := range invited {
		user, errUser := h.services.Authorization.GetUserById(number)
		if errUser != nil {
			return errParam
		}
		list = append(list, user)
	}
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"list": list,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) GetInvites(c echo.Context) error {
	userId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}
	invitations, err := h.services.Status.GetInvites(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return err
	}
	var list []models.User
	for _, number := range invitations {
		user, errUser := h.services.Authorization.GetUserById(number)
		if errUser != nil {
			return errParam
		}
		list = append(list, user)
	}
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"list": list,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) SearchUser(c echo.Context) error {
	username := c.Param(Username)
	if len(username) == 0 {
		return nil
	}
	users, err := h.services.Status.SearchUser(username)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return err
	}
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"list": users,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}
