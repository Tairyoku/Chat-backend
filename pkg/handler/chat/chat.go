package chat

import (
	"cmd/pkg/handler/middlewares"
	"cmd/pkg/handler/responses"
	"cmd/pkg/repository"
	"cmd/pkg/repository/models"
	"cmd/pkg/service"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strings"
)

type ChatHandler struct {
	services *service.Service
}

func NewChatHandler(services *service.Service) *ChatHandler {
	return &ChatHandler{services: services}
}

// CreatePublicChat godoc
// @Summary      Create a new public chat
// @Description  Отримує ім'я чату. Створює новий чат. Повертає ID чату.
// @Security ApiKeyAuth
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        chat_name	body     NameInput   true  "Chat name"
// @Success      200 	{object} IdResponse   "result is chat ID"
// @Failure 	 400 	{object} responses.ErrorResponse	 "incorrect request data"
// @Failure 	 400 	{object} responses.ErrorResponse	 "name is empty"
// @Failure 	 500 	{object} responses.ErrorResponse	 "create message error"
// @Failure 	 500 	{object} responses.ErrorResponse	 "add user to chat error"
// @Router       /chats/create [post]
func (h *ChatHandler) CreatePublicChat(c echo.Context) error {

	// Отримуємо дані з сайту (ім'я) та перевіряємо їх
	var chat models.Chat
	if err := c.Bind(&chat); err != nil {
		responses.NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return nil
	}
	if chat.Name == "" {
		responses.NewErrorResponse(c, http.StatusBadRequest, "name is empty")
		return nil
	}

	// Призначаємо публічний тип чату
	chat.Types = repository.ChatPublic

	// Отримуємо ID активного користувача
	userId := c.Get(middlewares.UserCtx).(int)

	// Створюємо чат
	chatId, err := h.services.Chat.Create(chat)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "create message error")
		return nil
	}

	// Додаємо активного користувача до новоствореного чату
	newUser := models.ChatUsers{
		ChatId: chatId,
		UserId: userId,
	}
	_, errAdd := h.services.Chat.AddUser(newUser)
	if errAdd != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "add user to chat error")
		return nil
	}

	// Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"id": chatId,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

// GetChat godoc
// @Summary      Get chat info
// @Description  Отримує ID чату. Повертає дані чату.
// @Security ApiKeyAuth
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        id		path     int   true  "Chat ID"
// @Success      200 	{object} ChatResponse   "result is chat data"
// @Failure 	 500 	{object} responses.ErrorResponse	 "get chat error"
// @Router       /chats/{id} [get]
func (h *ChatHandler) GetChat(c echo.Context) error {

	// Отримуємо ID чату
	id, errParam := middlewares.GetParam(c, middlewares.ParamId)
	if errParam != nil {
		return errParam
	}

	// Отримує дані чату за його ID
	chat, err := h.services.Chat.Get(id)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "get chat error")
		return nil
	}

	//Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"chat": chat,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

// GetById godoc
// @Summary      Get chat (and if chat is private - user) info
// @Description  Отримує ID чату.
// @Description  Повертає дані чату та, якщо чат приватний, користувача.
// @Security ApiKeyAuth
// @Accept       json
// @Tags         chat
// @Produce      json
// @Param        id		path     int   true  "Chat ID"
// @Success      200 	{object} ChatAndUserResponse   "result is chat data (and user data)"
// @Failure 	 500 	{object} responses.ErrorResponse	 "no chat error"
// @Failure 	 500 	{object} responses.ErrorResponse	 "get users error"
// @Router       /chats/{id}/link [get]
func (h *ChatHandler) GetById(c echo.Context) error {

	// Отримання власного ID
	creatorId := c.Get(middlewares.UserCtx).(int)

	// Отримання назви чату
	chatId, errParamC := middlewares.GetParam(c, middlewares.ParamId)
	if errParamC != nil {
		return errParamC
	}

	// Отримання даних чату
	chat, err := h.services.Chat.Get(chatId)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "no chat error")
		return nil
	}

	//Отримання даних користувача (для приватного чату)
	var user models.User
	if chat.Types == repository.ChatPrivate {
		// Отримання користувачів приватного чату
		users, err := h.services.Chat.GetUsers(chatId)
		if err != nil {
			responses.NewErrorResponse(c, http.StatusInternalServerError, "get users error")
			return nil
		}

		// Фільтрація користувачів
		if len(users) == 1 {
			user = users[0]
		} else {
			for _, u := range users {
				if u.Id != creatorId {
					user = u
				}
			}
		}
	}

	// Відгук чату
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
		"chat": chat,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

// GetUsers godoc
// @Summary      Get chat`s users
// @Description  Отримує ID чату. Повертає список користувачів чату.
// @Security ApiKeyAuth
// @Accept       json
// @Tags         chat
// @Produce      json
// @Param        id		path     int   true  "Chat ID"
// @Success      200 	{array} []models.User  "result is list of chat`s users"
// @Failure 	 500 	{object} responses.ErrorResponse	 "get users error"
// @Router       /chats/{id}/users [get]
func (h *ChatHandler) GetUsers(c echo.Context) error {

	// отримуємо ID чату
	chatId, errParamC := middlewares.GetParam(c, middlewares.ParamId)
	if errParamC != nil {
		return errParamC
	}

	// Отримуємо усіх користувачів чату
	users, err := h.services.Chat.GetUsers(chatId)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "get users error")
		return nil
	}

	// Відгук сервера
	errRes := c.JSON(http.StatusOK, users)
	if errRes != nil {
		return errRes
	}
	return nil
}

// GetUserPublicChats godoc
// @Summary      Get user`s public chats
// @Description  Отримує ID користувача.
// @Description  Повертає список публічних чатів, в яких є користувач.
// @Security ApiKeyAuth
// @Accept       json
// @Tags         chat
// @Produce      json
// @Param        id		path     int   true  "User ID"
// @Success      200 	{object} ChatListResponse  "result is list of chats"
// @Failure 	 500 	{object} responses.ErrorResponse	 "get chats error"
// @Router       /chats/users/{id}/public [get]
func (h *ChatHandler) GetUserPublicChats(c echo.Context) error {

	// Отримуємо ID користувача
	userId, errParam := middlewares.GetParam(c, middlewares.ParamId)
	if errParam != nil {
		return errParam
	}

	// Отримує список публічних чатів, в яких присутній користувач
	chats, err := h.services.Chat.GetPublicChats(userId)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "get chats error")
		return nil
	}

	// Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"list": chats,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

// GetUserPrivateChats godoc
// @Summary      Get user`s private chats
// @Description  Отримує ID користувача.
// @Description  Повертає список приватних чатів користувача.
// @Security ApiKeyAuth
// @Accept       json
// @Tags         chat
// @Produce      json
// @Param        id		path     int   true  "User ID"
// @Success      200 	{object} ChatListResponse  "result is list of chats"
// @Failure 	 500 	{object} responses.ErrorResponse	 "get chats error"
// @Failure 	 500 	{object} responses.ErrorResponse	 "get users error"
// @Router       /chats/users/{id}/private [get]
func (h *ChatHandler) GetUserPrivateChats(c echo.Context) error {

	// Отримання власного ID
	creatorId := c.Get(middlewares.UserCtx).(int)

	// Отримуємо ID користувача
	userId, errParam := middlewares.GetParam(c, middlewares.ParamId)
	if errParam != nil {
		return errParam
	}

	// Отримуємо список приватних чатів
	chats, err := h.services.Chat.GetPrivateChats(userId)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "get chats error")
		return nil
	}

	// Перевіряємо кількість користувачів у чаті
	// якщо їх не двоє, то не виводимо
	var result []models.Chat

	for _, chat := range chats {
		users, err := h.services.Chat.GetUsers(chat.Id)
		if err != nil {
			responses.NewErrorResponse(c, http.StatusInternalServerError, "get users error")
			return nil
		}
		if len(users) == 2 {
			for _, u := range users {
				if u.Id != creatorId {
					chat.Name = u.Username
					result = append(result, chat)
				}
			}
		}
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

// AddUserToChat godoc
// @Summary      Add user to chat
// @Description  Отримує ID чату та користувача.
// @Description  Додає користувача до чату.
// @Description  Повертає ID зв'язку між чатами та користувачами.
// @Security ApiKeyAuth
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        id		path     int   true  "Chat ID"
// @Param        user_id	body     UserIdInput   true  "User ID"
// @Success      200 	{object} IdResponse   "result is ID of chats and users relations"
// @Failure 	 400 	{object} responses.ErrorResponse	 "incorrect request data"
// @Failure 	 500 	{object} responses.ErrorResponse	 "add user to chat error"
// @Router       /chats/{id}/add [post]
func (h *ChatHandler) AddUserToChat(c echo.Context) error {

	// Отримуємо ID чату
	chatId, errParamC := middlewares.GetParam(c, middlewares.ParamId)
	if errParamC != nil {
		return errParamC
	}

	// Отримуємо від сайту ID користувача
	var list models.ChatUsers
	if errReq := c.Bind(&list); errReq != nil {
		responses.NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return nil
	}
	list.ChatId = chatId

	// Додаємо користувача до чату
	id, err := h.services.Chat.AddUser(list)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "add user to chat error")
		return nil
	}

	// Відгук сайту
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

// DeleteUserFromChat godoc
// @Summary      Delete user from chat
// @Description  Отримує ID чату та користувача.
// @Description  Видаляє користувача з чату.
// @Security ApiKeyAuth
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        id		path     int   true  "Chat ID"
// @Param        user_id	body     UserIdInput   true  "User ID"
// @Success      200 	{object} MessageResponse			"user deleted from chat"
// @Failure 	 400 	{object} responses.ErrorResponse	 "incorrect request data"
// @Failure 	 500 	{object} responses.ErrorResponse	 "delete user error"
// @Failure 	 500 	{object} responses.ErrorResponse	 "get chat users error"
// @Failure 	 500 	{object} responses.ErrorResponse	 "chat delete error"
// @Failure 	 500 	{object} responses.ErrorResponse	 "messages delete error"
// @Router       /chats/{id}/delete [put]
func (h *ChatHandler) DeleteUserFromChat(c echo.Context) error {

	// Отримуємо ID користувача від сайту
	var list models.ChatUsers
	if errReq := c.Bind(&list); errReq != nil {
		responses.NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return nil
	}

	// Отримуємо ID чату
	chatId, errParamC := middlewares.GetParam(c, middlewares.ParamId)
	if errParamC != nil {
		return errParamC
	}

	// Видаляємо користувача з чату
	err := h.services.Chat.DeleteUser(list.UserId, chatId)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "delete user error")
		return nil
	}

	// Отримуємо усіх користувачів чату
	users, errUser := h.services.Chat.GetUsers(chatId)
	if errUser != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "get chat users error")
		return nil
	}

	// Якщо в чаті не залишилося користувачів - видаляємо чат
	if len(users) == 0 {

		// Видаляємо чат
		err := h.services.Chat.Delete(chatId)
		if err != nil {
			responses.NewErrorResponse(c, http.StatusInternalServerError, "chat delete error")
			return nil
		}

		// Видаляємо усі повідомлення чату
		errMsg := h.services.Chat.DeleteAllMessages(chatId)
		if errMsg != nil {
			responses.NewErrorResponse(c, http.StatusInternalServerError, "messages delete error")
			return nil
		}
	}
	// Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("user with id %d deleted from chat with id %d", list.UserId, chatId),
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

// ChangeChatIcon godoc
// @Summary      Change chat icon
// @Description  Отримує ID чату та файл зображення.
// @Description  Оновлює зображення користувача.
// @Security ApiKeyAuth
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        id		path     int   true  "Chat ID"
// @Success      200 	{object} MessageResponse			"icon changed"
// @Failure 	 400 	{object} responses.ErrorResponse	 "incorrect chat data"
// @Failure 	 500 	{object} responses.ErrorResponse	 "update icon error"
// @Failure 	 500 	{object} responses.ErrorResponse	 "delete icon error"
// @Router       /chats/{id}/icon [put]
func (h *ChatHandler) ChangeChatIcon(c echo.Context) error {
	// Отримуємо ID чату
	chatId, errParamC := middlewares.GetParam(c, middlewares.ParamId)
	if errParamC != nil {
		return errParamC
	}

	fileName, err := middlewares.UploadImage(c)
	if err != nil {
		return err
	}

	//Отримуємо дані чату
	chat, errCh := h.services.Chat.Get(chatId)
	if errCh != nil {
		responses.NewErrorResponse(c, http.StatusBadRequest, "incorrect chat data")
		return nil
	}
	//Замінюємо дані у БД
	var oldIcon = chat.Icon
	chat.Icon = strings.TrimPrefix(fileName, "uploads/")

	errPut := h.services.Chat.Update(chat)
	if errPut != nil {

		responses.NewErrorResponse(c, http.StatusInternalServerError, "update icon error")
		return nil
	}

	//Видалення застарілих файлів
	if len(oldIcon) != 0 {
		if err := os.Remove(fmt.Sprintf("uploads\\%s", oldIcon)); err != nil {
			responses.NewErrorResponse(c, http.StatusInternalServerError, "delete icon error")
			return nil
		}
		if err := os.Remove(fmt.Sprintf("uploads\\resize-%s", oldIcon)); err != nil {
			responses.NewErrorResponse(c, http.StatusInternalServerError, "delete icon error")
			return nil
		}
	}

	//Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": "icon changed",
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

// DeleteChat godoc
// @Summary      Delete chat
// @Description  Отримує ID чату. Видаляє чат.
// @Security ApiKeyAuth
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        id		path     int   true  "Chat ID"
// @Success      200 	{object} MessageResponse			"chat  deleted"
// @Failure 	 500 	{object} responses.ErrorResponse	 "get chat users error"
// @Failure 	 500 	{object} responses.ErrorResponse	 "delete user from chat error"
// @Failure 	 500 	{object} responses.ErrorResponse	 "chat delete error"
// @Failure 	 500 	{object} responses.ErrorResponse	 "messages delete error"
// @Router       /chats/{id} [delete]
func (h *ChatHandler) DeleteChat(c echo.Context) error {

	// Отримуємо ID чату
	chatId, errParam := middlewares.GetParam(c, middlewares.ParamId)
	if errParam != nil {
		return errParam
	}

	// Отримуємо усіх користувачів чату
	users, errUser := h.services.Chat.GetUsers(chatId)
	if errUser != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "get chat users error")
		return nil
	}
	// Видаляємо усіх користувачів з чату
	for _, user := range users {
		errDelUser := h.services.Chat.DeleteUser(user.Id, chatId)
		if errDelUser != nil {
			responses.NewErrorResponse(c, http.StatusInternalServerError, "delete user from chat error")
			return nil
		}
	}

	// Видаляємо чат
	err := h.services.Chat.Delete(chatId)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "chat delete error")
		return nil
	}

	// Видаляємо усі повідомлення чату
	errMsg := h.services.Chat.DeleteAllMessages(chatId)
	if errMsg != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "messages delete error")
		return nil
	}

	// Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("chat with id %d deleted", chatId),
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

// PrivateChat godoc
// @Summary      Create private chat and/or get chat ID
// @Description  Отримує ID другого користувача приватного чату.
// @Description  Якщо чат вже існує, повертає його ID.
// @Description  Якщо чату не існує, створює його, та повертає ID чату.
// @Description  Якщо другий користувач є поточним користувачем,
// @Description  створює (за необхідністю) персональний чат та повертає його ID.
// @Security ApiKeyAuth
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        userId		path     int   true  "User ID"
// @Success      200 	{object} ChatIdResponse			"return chat ID"
// @Failure 	 500 	{object} responses.ErrorResponse	 "wrong users"
// @Failure 	 500 	{object} responses.ErrorResponse	 "wrong user"
// @Failure 	 500 	{object} responses.ErrorResponse	 "create chat error"
// @Failure 	 500 	{object} responses.ErrorResponse	 "add active user to chat error"
// @Failure 	 500 	{object} responses.ErrorResponse	 "add second user to chat error"
// @Router       /chats/{userId}/private [get]
func (h *ChatHandler) PrivateChat(c echo.Context) error {

	// Отримання власного ID
	creatorId := c.Get(middlewares.UserCtx).(int)

	// Отримання ID користувача, з яким створюється чат
	userId, errParamC := middlewares.GetParam(c, middlewares.UserCtx)
	if errParamC != nil {
		return errParamC
	}

	// Отримуємо ID чату або 0
	code, errP := h.services.Chat.GetPrivates(creatorId, userId)
	if errP != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "wrong users")
		return nil
	}

	// Якщо не існує чату, створюємо новий
	if code == 0 {
		// Отримуємо дані переданого користувача
		user, errUser := h.services.Chat.GetUserById(userId)
		if errUser != nil {
			responses.NewErrorResponse(c, http.StatusInternalServerError, "wrong user")
			return nil
		}

		// Створюємо новий чат
		chat := models.Chat{
			Name:  user.Username,
			Types: repository.ChatPrivate,
		}
		chatId, err := h.services.Chat.Create(chat)
		if err != nil {
			responses.NewErrorResponse(c, http.StatusInternalServerError, "create chat error")
			return nil
		}
		code = chatId

		// Додаємо активного користувача до чату
		_, errAdd := h.services.Chat.AddUser(models.ChatUsers{
			ChatId: chatId,
			UserId: creatorId,
		})
		if errAdd != nil {
			responses.NewErrorResponse(c, http.StatusInternalServerError, "add active user to chat error")
			return nil
		}

		// Додаємо отриманого користувача до чату (якщо це не особистий чат)
		if creatorId != userId {
			_, errNewAdd := h.services.Chat.AddUser(models.ChatUsers{
				ChatId: chatId,
				UserId: userId,
			})
			if errNewAdd != nil {
				responses.NewErrorResponse(c, http.StatusInternalServerError, "add second user to chat error")
				return nil
			}
		}
	}

	// Відгук чату
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"chatId": code,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

// SearchChat godoc
// @Summary      Get found users
// @Description  Отримує частину імені користувача.
// @Description  Повертає список користувачів, ім'я яких повністю або частково збігається.
// @Security ApiKeyAuth
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        name		path     string   true  "Slice of chats"
// @Success      200 	{object} ChatListResponse			"return found users"
// @Failure 	 500 	{object} responses.ErrorResponse	 "found chats error"
// @Router       /chats/search/{name} [get]
func (h *ChatHandler) SearchChat(c echo.Context) error {

	// Отримуємо сегмент назви чату
	name := c.Param(middlewares.ChatName)
	if len(name) == 0 {
		return nil
	}
	// Отримуємо список чатів
	chats, err := h.services.Chat.SearchChat(name)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "found chats error")
		return nil
	}
	// Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"list": chats,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}
