package handler

import (
	"cmd/pkg/repository"
	"cmd/pkg/repository/models"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) GetUserChats(c echo.Context) error {
	userId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}
	chats, err := h.services.Chat.GetUserChats(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
	}

	_, errEnCd := json.Marshal(chats)
	if errEnCd != nil {
		return errEnCd
	}
	errRes := c.JSON(http.StatusOK, chats)
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
	}

	_, errEnCd := json.Marshal(friends)
	if errEnCd != nil {
		return errEnCd
	}
	errRes := c.JSON(http.StatusOK, friends)
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
	}

	_, errEnCd := json.Marshal(bl)
	if errEnCd != nil {
		return errEnCd
	}
	errRes := c.JSON(http.StatusOK, bl)
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
	}

	_, errEnCd := json.Marshal(ubl)
	if errEnCd != nil {
		return errEnCd
	}
	errRes := c.JSON(http.StatusOK, ubl)
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
	}

	_, errEnCd := json.Marshal(invited)
	if errEnCd != nil {
		return errEnCd
	}
	errRes := c.JSON(http.StatusOK, invited)
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
	}

	_, errEnCd := json.Marshal(invitations)
	if errEnCd != nil {
		return errEnCd
	}
	errRes := c.JSON(http.StatusOK, invitations)
	if errRes != nil {
		return errRes
	}
	return nil
}
