package handler

import (
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

	//Посилання на зображення
	router.Static("/image/", "./uploads")
	//SWAGGER
	//router.GET("/swagger/server/*", echoSwagger.WrapHandler)

	//WebSocket
	router.GET("/ws/:roomId", func(c echo.Context) error {
		roomId := c.Param("roomId")
		websocket.ServeWs(c.Response(), c.Request(), roomId)
		return nil
	})

	api := router.Group("/api")

	auth := api.Group("/auth")
	{
		//Реєстрація
		auth.POST("/sign-up", h.SignUp)
		//Авторизація
		auth.POST("/sign-in", h.SignIn)
		//Отримати ID активного користувача
		auth.GET("/get-me", h.GetMe, h.userIdentify)
		//Змінити пароль
		auth.PUT("change/password", h.ChangePassword, h.userIdentify)
		//Змінити нікнейм
		auth.PUT("change/username", h.ChangeUsername, h.userIdentify)
		//Змінити аватар
		auth.PUT("change/icon", h.ChangeIcon, h.userIdentify)
	}

	users := api.Group("/users/:id", h.userIdentify)
	//Пошук користувачів за нікнеймом
	api.GET("/users/search/:username", h.SearchUser)
	{

		// МОЖНА ОБ'ЄДНАТИ ПУБЛІЧНІ ТА ПРИВАТНІ ЧАТИ
		//Отримати усі ПУБЛІЧНІ чати користувача
		users.GET("/public", h.GetUserPublicChats)
		//Отримати усі ОСОБИСТІ чати користувача
		users.GET("/private", h.GetUserPrivateChats)
		//Отримати дані користувача за його ID
		users.GET("", h.GetUserById)
		//Отримати список усіх користувачів, пов'язаних з вами
		users.GET("/all", h.GetUserLists)
		//Запит на дружбу
		users.POST("/invite", h.InvitedToFriends)
		//Скасувати запит на дружбу
		users.DELETE("/cancel", h.CancelInvite)
		//Прийняти запит на дружбу
		users.PUT("/accept", h.AcceptInvitation)
		//Відмовити запиту на дружбу
		users.DELETE("/refuse", h.RefuseInvitation)
		//Заблокувати користувача
		users.POST("/addToBL", h.AddToBlackList)
		//Розблокувати
		users.DELETE("/deleteFromBlacklist", h.DeleteFromBlacklist)
		//Видалити з друзів
		users.DELETE("/deleteFriend", h.DeleteFriend)
	}

	chat := api.Group("/chats", h.userIdentify)
	{
		//Створити ПУБЛІЧНИЙ чат
		chat.POST("/create", h.CreatePublicChat)
		//Створити ОСОБИСТИЙ чат
		chat.GET("/:userId/private", h.PrivateChat)
		//Отримати дані чату за його ID
		chat.GET("/:id", h.GetChat)
		//Отримати дані чату та користувача (тільки у
		// приватному чаті) за ID чату
		chat.GET("/:id/link", h.GetById)
		//Отримати список користувачів чату
		chat.GET("/:id/users", h.GetUsers)
		//Додати користувачів до чату
		chat.PUT("/:id/add", h.AddUserToChat)
		//Видалити користувачів із чату
		chat.PUT("/:id/delete", h.DeleteUserFromChat)
		//Видалити чат
		chat.DELETE("/:id", h.DeleteChat)
		//Пошук чатів за назвою
		chat.GET("/search/:name", h.SearchChat)

	}

	message := chat.Group("/:chatId/messages")
	{
		//Створити повідомлення
		message.POST("", h.CreateMessage)
		//Отримати повідомлення
		message.GET("", h.GetAllMessages)
		//
		message.GET("/:id", h.GetMessage)
	}
	return router
}
