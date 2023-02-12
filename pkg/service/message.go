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

// Create викликає створення нового повідомлення та повертає його ID
func (m *MessageService) Create(msg models.Message) (int, error) {
	return m.repository.Create(msg)
}

// Get викликає повернення повідомлення за його ID
func (m *MessageService) Get(msgId int) (models.Message, error) {
	return m.repository.Get(msgId)
}

// GetAll викликає повернення усіх повідомлень чату за його ID
func (m *MessageService) GetAll(chatId int) ([]models.Message, error) {
	return m.repository.GetAll(chatId)
}

// GetLimit викликає повернення певної кількості повідомлень чату за його ID
func (m *MessageService) GetLimit(chatId, limit int) ([]models.Message, error) {
	return m.repository.GetLimit(chatId, limit)
}

// DeleteAll викликає видалення усіх повідомлень чата за його ID
func (m *MessageService) DeleteAll(chatId int) error {
	return m.repository.DeleteAll(chatId)
}
