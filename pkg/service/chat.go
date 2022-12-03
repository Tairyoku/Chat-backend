package service

import (
	"cmd/pkg/repository"
	"cmd/pkg/repository/models"
)

type ChatService struct {
	repository repository.Chat
}

func NewChatService(repository repository.Chat) *ChatService {
	return &ChatService{repository: repository}
}

func (c *ChatService) Create(chat models.Chat) (int, error) {
	return c.repository.Create(chat)
}
func (c *ChatService) Delete(chatId int) error {
	return c.repository.Delete(chatId)
}
func (c *ChatService) Get(chatId int) (models.Chat, error) {
	return c.repository.Get(chatId)
}
func (c *ChatService) AddUser(users models.ChatUsers) (int, error) {
	return c.repository.AddUser(users)
}

func (c *ChatService) DeleteUser(userId, chatId int) error {
	return c.repository.DeleteUser(userId, chatId)
}

func (c *ChatService) GetUsers(chatId int) ([]models.User, error) {
	return c.repository.GetUsers(chatId)
}

func (c *ChatService) GetUserChats(userId int) ([]models.Chat, error) {
	return c.repository.GetUserChats(userId)
}
