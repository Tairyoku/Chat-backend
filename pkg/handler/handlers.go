package handler

import (
	auth2 "cmd/pkg/handler/auth"
	chat2 "cmd/pkg/handler/chat"
	message2 "cmd/pkg/handler/message"
	"cmd/pkg/handler/middlewares"
	users2 "cmd/pkg/handler/users"
	"cmd/pkg/handler/websocket"
	"cmd/pkg/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *echo.Echo {
	router := echo.New()
	router.Use(middleware.CORS())
	UsersHandler := middlewares.NewMiddlewareHandler(h.services)
	messageHandler := message2.NewMessageHandler(h.services)
	chatHandler := chat2.NewChatHandler(h.services)
	authHandler := auth2.NewAuthHandler(h.services)
	usersHandler := users2.NewUsersHandler(h.services)
	//SWAGGER
	//router.GET("/swagger/server/*", echoSwagger.WrapHandler)

	//WebSocket
	router.GET("/ws/:roomId", func(c echo.Context) error {
		roomId := c.Param("roomId")
		websocket.ServeWs(c.Response(), c.Request(), roomId)
		return nil
	})

	api := router.Group("/api")

	//Посилання на зображення
	api.Static("/image/", "./uploads")

	auth := api.Group("/auth")
	{
		//Реєстрація
		auth.POST("/sign-up", authHandler.SignUp)
		//Авторизація
		auth.POST("/sign-in", authHandler.SignIn)
		//Отримати ID активного користувача
		auth.GET("/get-me", authHandler.GetMe, UsersHandler.UserIdentify)
		//Змінити пароль
		auth.PUT("/change/password", authHandler.ChangePassword, UsersHandler.UserIdentify)
		//Змінити нікнейм
		auth.PUT("/change/username", authHandler.ChangeUsername, UsersHandler.UserIdentify)
		//Змінити аватар
		auth.PUT("/change/icon", authHandler.ChangeIcon, UsersHandler.UserIdentify)
	}

	users := api.Group("/users/:id", UsersHandler.UserIdentify)
	//Пошук користувачів за нікнеймом
	api.GET("/users/search/:username", usersHandler.SearchUser)
	{

		// МОЖНА ОБ'ЄДНАТИ ПУБЛІЧНІ ТА ПРИВАТНІ ЧАТИ
		//Отримати усі ПУБЛІЧНІ чати користувача
		users.GET("/public", chatHandler.GetUserPublicChats)
		//Отримати усі ОСОБИСТІ чати користувача
		users.GET("/private", chatHandler.GetUserPrivateChats)
		//Отримати дані користувача за його ID
		users.GET("", usersHandler.GetUserById)
		//Отримати список усіх користувачів, пов'язаних з вами
		users.GET("/all", usersHandler.GetUserLists)
		//Запит на дружбу
		users.POST("/invite", usersHandler.InvitedToFriends)
		//Скасувати запит на дружбу
		users.DELETE("/cancel", usersHandler.CancelInvite)
		//Прийняти запит на дружбу
		users.PUT("/accept", usersHandler.AcceptInvitation)
		//Відмовити запиту на дружбу
		users.DELETE("/refuse", usersHandler.RefuseInvitation)
		//Заблокувати користувача
		users.POST("/addToBL", usersHandler.AddToBlackList)
		//Розблокувати
		users.DELETE("/deleteFromBlacklist", usersHandler.DeleteFromBlacklist)
		//Видалити з друзів
		users.DELETE("/deleteFriend", usersHandler.DeleteFriend)
	}

	chat := api.Group("/chats", UsersHandler.UserIdentify)
	{
		//Створити ПУБЛІЧНИЙ чат
		chat.POST("/create", chatHandler.CreatePublicChat)
		//Створити ОСОБИСТИЙ чат
		chat.GET("/:userId/private", chatHandler.PrivateChat)
		//Отримати дані чату за його ID
		chat.GET("/:id", chatHandler.GetChat)
		//Отримати дані чату та користувача (тільки у
		// приватному чаті) за ID чату
		chat.GET("/:id/link", chatHandler.GetById)
		//Отримати список користувачів чату
		chat.GET("/:id/users", chatHandler.GetUsers)
		//Додати користувачів до чату
		chat.POST("/:id/add", chatHandler.AddUserToChat)
		//Видалити користувачів із чату
		chat.PUT("/:id/delete", chatHandler.DeleteUserFromChat)
		//Оновити зображення чату
		chat.PUT("/:id/icon", chatHandler.ChangeChatIcon)
		//Видалити чат
		chat.DELETE("/:id", chatHandler.DeleteChat)
		//Пошук чатів за назвою
		chat.GET("/search/:name", chatHandler.SearchChat)

	}

	message := chat.Group("/:chatId/messages")
	{
		//Створити повідомлення
		message.POST("", messageHandler.CreateMessage)
		//Отримати повідомлення
		message.GET("", messageHandler.GetAllMessages)
		//Отримати певну кількість повідомлень
		message.GET("/limit/:id", messageHandler.GetLimitMessages)
		//
		message.GET("/:id", messageHandler.GetMessage)
	}
	return router
}
