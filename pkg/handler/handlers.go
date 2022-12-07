package handler

import (
	"cmd/pkg/service"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *echo.Echo {
	router := echo.New()
	//router.GET("/swagger/server/*", echoSwagger.WrapHandler)
	api := router.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/sign-in", h.SignIn)
		//ChangePassword
		//ChangeUsername
		//ChangeIcon
	}

	users := api.Group("/users/:id")
	{
		users.GET("/chats", h.GetUserChats)
		users.GET("/friends", h.GetFriends)
		users.GET("/blacklist", h.GetBlackList)
		users.GET("/inBlacklist", h.GetBlackListToUser)
		users.GET("/invited", h.GetSentInvites)
		users.GET("/invitations", h.GetInvites)
		users.POST("/invite", h.InvitedToFriends, h.userIdentify)
		users.POST("/addToBL", h.AddToBlackList, h.userIdentify)
		users.PUT("/accept", h.AcceptInvitation, h.userIdentify)
		users.DELETE("/refuseInvitation", h.RefuseInvitation, h.userIdentify)
		users.DELETE("/deleteFriend", h.DeleteFriend, h.userIdentify)
		users.DELETE("/deleteFromBlacklist", h.DeleteFromBlacklist, h.userIdentify)

	}
	chat := api.Group("/chats")
	{
		chat.POST("", h.CreateChat, h.userIdentify)
		chat.GET("/:id", h.GetChat)
		chat.DELETE("/:id", h.DeleteChat, h.userIdentify)
		chat.POST("/:id/user", h.AddUserToChat, h.userIdentify)
		chat.DELETE("/:id/user", h.DeleteUserFromChat, h.userIdentify)
		chat.GET("/:id/users", h.GetUsers, h.userIdentify)

	}

	message := chat.Group("/:chatId/messages")
	{
		message.POST("", h.CreateMessage)
		message.GET("", h.GetAllMessages)
		message.GET("/:id", h.GetMessage)
	}
	return router
}
