package handler

import (
	"cmd/pkg/repository"
	"cmd/pkg/repository/models"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) CreatePublicChat(c echo.Context) error {
	var chat models.Chat
	if err := c.Bind(&chat); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return nil
	}
	if chat.Name == "" {
		NewErrorResponse(c, http.StatusBadRequest, "name is empty")
		return nil
	}
	chat.Types = repository.ChatPublic

	userId, errId := GetUserId(c) //получаем "свой" id
	if errId != nil {
		return errId
	}
	chatId, err := h.services.Chat.Create(chat) //создаем публичний чат, получаем id чата
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return err
	}

	newUser := models.ChatUsers{
		ChatId: chatId,
		UserId: userId,
	}
	_, errAdd := h.services.Chat.AddUser(newUser) //присвоение чату пользователя (себя)
	if errAdd != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "add user to chat error")
		return errAdd
	}
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"id": chatId,
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
	chatId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}
	//получить пользователей из чата
	users, errUser := h.services.Chat.GetUsers(chatId)
	if errUser != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "get chat users server error")
		return nil
	}
	//удалить всех пользователей из чата
	for _, user := range users {
		errDelUser := h.services.Chat.DeleteUser(user.Id, chatId)
		if errDelUser != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "delete user from chat server error")
			return nil
		}
	}
	err := h.services.Chat.Delete(chatId) //удалить сам чат
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "chat delete server error")
		return nil
	}
	errMsg := h.services.Message.DeleteAll(chatId) //удалить все сообщения чата
	if errMsg != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "messages delete server error")
		return nil
	}
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("chat with id %d deleted", chatId),
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) GetCorrespondence(c echo.Context) error {
	chatId, errParam := GetParam(c, ChatId)
	if errParam != nil {
		return errParam
	}

	msg, err := h.services.Message.GetAll(chatId)
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
	var list models.ChatUsers
	errReq := GetRequest(c, list)
	if errReq != nil {
		return errReq
	}
	chatId, errParamC := GetParam(c, ParamId)
	if errParamC != nil {
		return errParamC
	}

	err := h.services.Chat.DeleteUser(chatId, list.UserId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return nil
	}

	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("user with id %d deleted from chat with id %d", chatId, list.UserId),
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

	errRes := c.JSON(http.StatusOK, users)
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) PrivateChat(c echo.Context) error {
	userId, errParamC := GetParam(c, userCtx)
	if errParamC != nil {
		return errParamC
	}
	creatorId := c.Get(userCtx).(int)
	chats, err := h.services.Chat.CheckPrivates(creatorId, userId)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "error check chats")
	}
	switch len(chats) {
	// chat is not created
	case 0:
		{
			// users data
			user, errUser := h.services.Authorization.GetUserById(userId)
			if errUser != nil {
				NewErrorResponse(c, http.StatusInternalServerError, "wrong user")
				return errUser
			}
			chat := models.Chat{
				Name:  user.Username,
				Types: repository.ChatPrivate,
			}
			//create chat
			chatId, err := h.services.Chat.Create(chat)
			if err != nil {
				NewErrorResponse(c, http.StatusInternalServerError, "server error")
				return err
			}
			chats = append(chats, chatId)
			//add creator
			_, errAdd := h.services.Chat.AddUser(models.ChatUsers{
				ChatId: chatId,
				UserId: creatorId,
			})
			if errAdd != nil {
				NewErrorResponse(c, http.StatusInternalServerError, "add user to chat error")
				return errAdd
			}
			//add user
			_, errNewAdd := h.services.Chat.AddUser(models.ChatUsers{
				ChatId: chatId,
				UserId: userId,
			})
			if errNewAdd != nil {
				NewErrorResponse(c, http.StatusInternalServerError, "add user to chat error")
				return errNewAdd
			}
			break
		}
	case 1:
		{
			break
		}
	default:
		{
			NewErrorResponse(c, http.StatusInternalServerError, "private chat had a duplicate")
			return nil
		}
	}
	messages, errMes := h.services.Message.GetAll(chats[0])
	if errMes != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "messages error")
		return errMes
	}

	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"messages": messages,
		"chatId":   chats[0],
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) GetChatByName(c echo.Context) error {
	name := c.Param(ChatName)
	firstId, errParamC := GetUserId(c)
	if errParamC != nil {
		return errParamC
	}
	second, err := h.services.Authorization.GetByName(name)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "user name error")
		return nil
	}
	chats, err := h.services.Chat.CheckPrivates(firstId, second.Id)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "error check chats")
	}
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"chat": chats[0],
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) SearchChat(c echo.Context) error {
	name := c.Param(ChatName)
	if len(name) == 0 {
		return nil
	}
	chats, err := h.services.Chat.SearchChat(name)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return err
	}
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"list": chats,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}
