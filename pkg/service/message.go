package service

import (
	"cmd/pkg/repository"
	"cmd/pkg/repository/models"
)

type MessageService struct {
	repository repository.Message
}

func NewMessageService(repository repository.Message) *MessageService {
	return &MessageService{repository: repository}
}

func (m *MessageService) Create(msg models.Message) (int, error) {
	return m.repository.Create(msg)
}

func (m *MessageService) GetAll(chatId int) ([]models.Message, error) {
	return m.repository.GetAll(chatId)
}

func (m *MessageService) Get(msgId int) (models.Message, error) {
	return m.repository.Get(msgId)
}
