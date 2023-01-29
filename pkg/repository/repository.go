package repository

import (
	"cmd/pkg/repository/models"
	"github.com/jinzhu/gorm"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
	CheckUsername(username string) error
	UpdateUser(user models.User) error
	GetUserById(userId int) (models.User, error)
	GetByName(name string) (models.User, error)
}

type Chat interface {
	Create(chat models.Chat) (int, error)
	Delete(chatId int) error
	Get(chatId int) (models.Chat, error)
	AddUser(users models.ChatUsers) (int, error)
	DeleteUser(userId, chatId int) error
	GetUsers(chatId int) ([]models.User, error)
	GetPrivateChats(userId int) ([]models.Chat, error)
	GetPublicChats(userId int) ([]models.Chat, error)
	CheckPrivates(firstUser, secondUser int) ([]int, error)
	SearchChat(name string) ([]models.Chat, error)
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
	SearchUser(username string) ([]models.User, error)
}

type Message interface {
	Create(msg models.Message) (int, error)
	Get(msgId int) (models.Message, error)
	GetAll(chatId int) ([]models.Message, error)
	DeleteAll(chatId int) error
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
