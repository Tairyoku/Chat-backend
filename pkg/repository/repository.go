package repository

import (
	"cmd/pkg/repository/models"
	"gorm.io/gorm"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
	CheckUsername(username string) error
}

type Chat interface {
	Create(chat models.Chat) (int, error)
	Delete(chatId int) error
	Get(chatId int) (models.Chat, error)
	AddUser(users models.ChatUsers) (int, error)
	DeleteUser(userId, chatId int) error
	GetUsers(chatId int) ([]models.User, error)
	GetUserChats(userId int) ([]models.Chat, error)
}

type Status interface {
	AddStatus(status models.Status) (int, error)
	GetStatus(senderId, recipientId int) (models.Status, error)
	DeleteStatus(status models.Status) error
	UpdateStatus(status models.Status) error
	GetFriends(userId int) ([]int, error)
	GetBlackList(userId int) ([]int, error)
	GetBlackListToUser(userId int) ([]int, error)
	GetSentInvites(userId int) ([]int, error)
	GetInvites(userId int) ([]int, error)
}

type Message interface {
	Create(msg models.Message) (int, error)
	Get(msgId int) (models.Message, error)
	GetAll(chatId int) ([]models.Message, error)
}

type Repository struct {
	Authorization
	Chat
	Status
	Message
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db),
		Chat:          NewChatRepository(db),
		Status:        NewStatusRepository(db),
		Message:       NewMessageRepository(db),
	}
}
