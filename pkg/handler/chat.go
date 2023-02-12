package handler

import (
	"cmd/pkg/repository"
	"cmd/pkg/repository/models"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strings"
)

func (h *Handler) CreatePublicChat(c echo.Context) error {

	// Отримуємо дані з сайту (ім'я) та перевіряємо їх
	var chat models.Chat
	if err := c.Bind(&chat); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return nil
	}
	if chat.Name == "" {
		NewErrorResponse(c, http.StatusBadRequest, "name is empty")
		return nil
	}

	// Призначаємо публічний тип чату
	chat.Types = repository.ChatPublic

	// Отримуємо ID активного користувача
	userId, errId := GetUserId(c)
	if errId != nil {
		return errId
	}

	// Створюємо чат
	chatId, err := h.services.Chat.Create(chat)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "create message error")
		return err
	}

	// Додаємо активного користувача до новоствореного чату
	newUser := models.ChatUsers{
		ChatId: chatId,
		UserId: userId,
	}
	_, errAdd := h.services.Chat.AddUser(newUser)
	if errAdd != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "add user to chat error")
		return errAdd
	}

	// Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"id": userId,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) GetChat(c echo.Context) error {

	// Отримуємо ID чату
	id, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}

	// Отримує дані чату за його ID
	chat, err := h.services.Chat.Get(id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return nil
	}
	//_, errEnCd := json.Marshal(chat)
	//if errEnCd != nil {
	//	return errEnCd
	//}

	//Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"chat": chat,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) GetById(c echo.Context) error {

	// Отримання власного ID
	creatorId := c.Get(userCtx).(int)

	// Отримання назви чату
	chatId, errParamC := GetParam(c, ParamId)
	if errParamC != nil {
		return errParamC
	}
	if chatId == 0 {
		NewErrorResponse(c, http.StatusBadRequest, "no chat id")
		return nil
	}

	// Отримання даних чату
	chat, err := h.services.Chat.Get(chatId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "no chat error")
		return err
	}

	//Отримання даних користувача (для приватного чату)
	var user models.User
	if chat.Types == repository.ChatPrivate {
		// Отримання користувачів приватного чату
		users, err := h.services.Chat.GetUsers(chatId)
		if err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "user name error")
			return nil
		}
		fmt.Println(len(users))
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
	fmt.Println(user)
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

func (h *Handler) GetUsers(c echo.Context) error {

	// отримуємо ID чату
	chatId, errParamC := GetParam(c, ParamId)
	if errParamC != nil {
		return errParamC
	}

	// Отримуємо усіх користувачів чату
	users, err := h.services.Chat.GetUsers(chatId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return nil
	}

	// Відгук сервера
	errRes := c.JSON(http.StatusOK, users)
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) GetUserPublicChats(c echo.Context) error {

	// Отримуємо ID користувача
	userId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}

	// Отримує список публічних чатів, в яких присутній користувач
	chats, err := h.services.Chat.GetPublicChats(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
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

func (h *Handler) GetUserPrivateChats(c echo.Context) error {

	// Отримання власного ID
	creatorId := c.Get(userCtx).(int)

	// Отримуємо ID користувача
	userId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}

	// Отримуємо список приватних чатів
	chats, err := h.services.Chat.GetPrivateChats(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return err
	}

	// Перевіряємо кількість користувачів у чаті
	// якщо їх не двоє, то не виводимо
	var result []models.Chat

	for _, chat := range chats {
		users, err := h.services.Chat.GetUsers(chat.Id)
		if err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "server error")
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

func (h *Handler) AddUserToChat(c echo.Context) error {

	// Отримуємо ID чату
	chatId, errParamC := GetParam(c, ParamId)
	if errParamC != nil {
		return errParamC
	}

	// Отримуємо від сайту ID користувача
	var list models.ChatUsers
	if errReq := c.Bind(&list); errReq != nil {
		NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return nil
	}
	list.ChatId = chatId

	fmt.Println(list)
	// Додаємо користувача до чату
	id, err := h.services.Chat.AddUser(list)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
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

func (h *Handler) DeleteUserFromChat(c echo.Context) error {

	// Отримуємо ID користувача від сайту
	var list models.ChatUsers
	if errReq := c.Bind(&list); errReq != nil {
		NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return nil
	}
	fmt.Println(list.UserId)

	// Отримуємо ID чату
	chatId, errParamC := GetParam(c, ParamId)
	if errParamC != nil {
		return errParamC
	}
	fmt.Println(chatId)
	// Видаляємо користувача з чату
	err := h.services.Chat.DeleteUser(list.UserId, chatId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return err
	}
	fmt.Println("finish")

	// Отримуємо усіх користувачів чату
	users, errUser := h.services.Chat.GetUsers(chatId)
	if errUser != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "get chat users server error")
		return nil
	}

	// Якщо в чаті не залишилося користувачів - видаляємо чат
	if len(users) == 0 {

		// Видаляємо чат
		err := h.services.Chat.Delete(chatId)
		if err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "chat delete server error")
			return nil
		}

		// Видаляємо усі повідомлення чату
		errMsg := h.services.Message.DeleteAll(chatId)
		if errMsg != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "messages delete server error")
			return nil
		}
	}
	// Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("user with id %d deleted from chat with id %d", chatId, list.UserId),
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) ChangeChatIcon(c echo.Context) error {
	// Отримуємо ID чату
	chatId, errParamC := GetParam(c, ParamId)
	if errParamC != nil {
		return errParamC
	}

	fileName, err := UploadImage(c)
	if err != nil {
		return err
	}
	fmt.Println("work")
	//Отримуємо дані чату
	chat, errCh := h.services.Chat.Get(chatId)
	if errCh != nil {
		fmt.Println(7)

		NewErrorResponse(c, http.StatusBadRequest, "incorrect chat data")
		return nil
	}

	//Замінюємо дані у БД
	var oldIcon = chat.Icon
	chat.Icon = strings.TrimPrefix(fileName, "uploads\\")
	errPut := h.services.Chat.Update(chat)
	if errPut != nil {
		fmt.Println(8)

		NewErrorResponse(c, http.StatusInternalServerError, "update icon error")
		return nil
	}

	//Видалення застарілих файлів
	if len(oldIcon) != 0 {
		if err := os.Remove(fmt.Sprintf("uploads/%s", oldIcon)); err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "delete icon error")
			return nil
		}
		if err := os.Remove(fmt.Sprintf("uploads/resize-%s", oldIcon)); err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "delete icon error")
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

func (h *Handler) DeleteChat(c echo.Context) error {

	// Отримуємо ID чату
	chatId, errParam := GetParam(c, ParamId)
	if errParam != nil {
		return errParam
	}

	// Отримуємо усіх користувачів чату
	users, errUser := h.services.Chat.GetUsers(chatId)
	if errUser != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "get chat users server error")
		return nil
	}

	// Видаляємо усіх користувачів з чату
	for _, user := range users {
		errDelUser := h.services.Chat.DeleteUser(user.Id, chatId)
		if errDelUser != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "delete user from chat server error")
			return nil
		}
	}

	// Видаляємо чат
	err := h.services.Chat.Delete(chatId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "chat delete server error")
		return nil
	}

	// Видаляємо усі повідомлення чату
	errMsg := h.services.Message.DeleteAll(chatId)
	if errMsg != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "messages delete server error")
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

func (h *Handler) PrivateChat(c echo.Context) error {

	// Отримання власного ID
	creatorId := c.Get(userCtx).(int)

	// Отримання ID користувача, з яким створюється чат
	userId, errParamC := GetParam(c, userCtx)
	if errParamC != nil {
		return errParamC
	}

	// Отримуємо ID чату або 0
	code, errP := h.services.Chat.GetPrivates(creatorId, userId)
	if errP != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "wrong users")
		return errP
	}

	// Якщо не існує чату, створюємо новий
	if code == 0 {
		// Отримуємо дані переданого користувача
		user, errUser := h.services.Authorization.GetUserById(userId)
		if errUser != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "wrong user")
			return errUser
		}

		// Створюємо новий чат
		chat := models.Chat{
			Name:  user.Username,
			Types: repository.ChatPrivate,
		}
		chatId, err := h.services.Chat.Create(chat)
		if err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "create chat error")
			return errUser
		}
		code = chatId

		// Додаємо активного користувача до чату
		_, errAdd := h.services.Chat.AddUser(models.ChatUsers{
			ChatId: chatId,
			UserId: creatorId,
		})
		if errAdd != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "add user to chat error")
			return errAdd
		}

		// Додаємо отриманого користувача до чату (якщо це не особистий чат)
		if creatorId != userId {
			_, errNewAdd := h.services.Chat.AddUser(models.ChatUsers{
				ChatId: chatId,
				UserId: userId,
			})
			if errNewAdd != nil {
				NewErrorResponse(c, http.StatusInternalServerError, "add user to chat error")
				return errNewAdd
			}
		}
	}

	// Отримуємо повідомлення чату
	messages, errMes := h.services.Message.GetAll(code)
	if errMes != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "messages error")
		return errMes
	}

	// Відгук чату
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"messages": messages,
		"chatId":   code,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) SearchChat(c echo.Context) error {

	// Отримуємо сегмент назви чату
	name := c.Param(ChatName)
	if len(name) == 0 {
		return nil
	}
	// Отримуємо список чатів
	chats, err := h.services.Chat.SearchChat(name)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return err
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
