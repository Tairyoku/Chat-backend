package chat

import "cmd/pkg/repository/models"

type IdResponse struct {
	Id string `json:"id"`
}

type ChatIdResponse struct {
	ChatId string `json:"chat_id"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Chat models.Chat `json:"chat"`
}

type ChatAndUserResponse struct {
	Chat models.Chat `json:"chat"`
	User models.User `json:"user"`
}

type NameInput struct {
	Name string `json:"name"`
}

type ChatListResponse struct {
	List []models.Chat `json:"list"`
}

type UserIdInput struct {
	UserId int `json:"user_id"`
}
