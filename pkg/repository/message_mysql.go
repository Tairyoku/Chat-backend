package repository

import (
	"cmd/pkg/repository/models"
	"github.com/jinzhu/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// Create отримує дані повідомлення ТА повертає його ID
func (m *MessageRepository) Create(msg models.Message) (int, error) {
	err := m.db.Table(MessagesTable).Create(&msg).Error
	return msg.Id, err
}

// Get отримує ID повідомлення ТА повертає його дані
func (m *MessageRepository) Get(msgId int) (models.Message, error) {
	var msg models.Message
	err := m.db.Table(MessagesTable).First(&msg, msgId).Error
	return msg, err
}

// GetAll отримує ID чату ТА повертає його повідомлення
func (m *MessageRepository) GetAll(chatId int) ([]models.Message, error) {
	var msg []models.Message
	err := m.db.Table(MessagesTable).Where("chat_id = ?", chatId).Find(&msg).Error
	return msg, err
}

// DeleteAll отримує ID чату ТА видаляє його повідомлення
func (m *MessageRepository) DeleteAll(chatId int) error {
	err := m.db.Table(MessagesTable).Where("chat_id = ?", chatId).Delete(&models.Message{}).Error
	return err
}
