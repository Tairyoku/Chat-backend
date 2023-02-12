package repository

import (
	"cmd/pkg/repository/models"
	"github.com/jinzhu/gorm"
)

type Authorization interface {
	// CreateUser отримує ім'я та пароль ТА створює нового користувача
	CreateUser(user models.User) (int, error)
	// GetUser отримує ім'я та пароль ТА повертає його дані
	GetUser(username, password string) (models.User, error)
	// GetUserById отримує ID користувача ТА повертає його дані
	GetUserById(userId int) (models.User, error)
	// GetByName отримує ім'я користувача ТА повертає його дані
	GetByName(username string) (models.User, error)
	// UpdateUser отримує дані користувача ТА оновлює їх
	UpdateUser(user models.User) error
}

type Chat interface {
	// Create отримує назву чату ТА створює новий чат
	Create(chat models.Chat) (int, error)
	// Get отримує ID чату ТА повертає дані чату за його ID
	Get(chatId int) (models.Chat, error)
	// Delete отримує ID чату ТА видаляє чат
	Delete(chatId int) error
	// Update отримує ID чату ТА оновлює дані чату
	Update(chat models.Chat) error
	// AddUser отримує ID чату ТА ID користувача, та додає користувача до чату
	AddUser(users models.ChatUsers) (int, error)
	// CheckPrivates отримує ID двох можливих користувачів одного приватного
	// чату ТА повертає масив ID приватних чатів, у яких двічі згадується ID
	// користувачів, поданих як аргументи
	CheckPrivates(firstUser, secondUser int) ([]int, error)
	// GetUsers отримує ID чату ТА повертає масив користувачів, що приєднані до чату
	GetUsers(chatId int) ([]models.User, error)
	// GetPrivates отримує ID двох користувачів, ТА повертає масив приватних
	// чатів, до яких належать кожен із користувачів
	GetPrivates(firstUser, secondUser int) ([]models.Chat, []models.Chat, error)
	// GetPrivateChats отримує ID користувача ТА повертає масив ПРИВАТНИХ чатів,
	// до яких він належить
	GetPrivateChats(userId int) ([]models.Chat, error)
	// GetPublicChats отримує ID користувача ТА повертає масив ПУБЛІЧНИХ чатів,
	// до яких він належить
	GetPublicChats(userId int) ([]models.Chat, error)
	// DeleteUser отримує ID чату ТА ID користувача, та видаляє користувача із чату
	DeleteUser(userId, chatId int) error
	// SearchChat отримує назву чату (або його частину) ТА повертає масив чатів,
	// назви яких збігаються з аргументом
	SearchChat(name string) ([]models.Chat, error)
}

type Status interface {
	// AddStatus отримує ID двох користувачів та їх тип відносин ТА повертає ID створеного статусу
	AddStatus(status models.Status) (int, error)
	// GetStatuses отримує ID двох користувачів ТА повертає дані їх відносин
	GetStatuses(senderId, recipientId int) ([]models.Status, error)
	// UpdateStatus отримує ID двох користувачів та їх тип відносин ТА оновлює дані
	UpdateStatus(status models.Status) error
	// DeleteStatus отримує ID двох користувачів та їх тип відносин ТА видаляє ці відносини
	DeleteStatus(status models.Status) error
	// GetFriends отримує ID користувача ТА повертає масив користувачів, що є ДРУЗЯМИ
	GetFriends(userId int) ([]models.User, error)
	// GetBlackList отримує ID користувача ТА повертає масив ЗАБЛОКОВАНИХ користувачів
	GetBlackList(userId int) ([]models.User, error)
	// GetBlackListToUser отримує ID користувача ТА повертає масив користувачів, що
	// ЗАБЛОКУВАЛИ його
	GetBlackListToUser(userId int) ([]models.User, error)
	// GetSentInvites отримує ID користувача ТА повертає масив користувачів, що
	// ОТРИМАЛИ його запрошення у друзі
	GetSentInvites(userId int) ([]models.User, error)
	// GetInvites отримує ID користувача ТА повертає масив користувачів, що
	// НАДІСЛАЛИ йому запрошення в друзі
	GetInvites(userId int) ([]models.User, error)
	// SearchUser отримує ім'я (або його частину) ТА повертає масив користувачів, що
	// мають збіг з аргументом
	SearchUser(username string) ([]models.User, error)
}

type Message interface {
	// Create отримує дані повідомлення ТА повертає його ID
	Create(msg models.Message) (int, error)
	// Get отримує ID повідомлення ТА повертає його дані
	Get(msgId int) (models.Message, error)
	// GetAll отримує ID чату ТА повертає його повідомлення
	GetAll(chatId int) ([]models.Message, error)
	// GetLimit отримує ID чату ліміт кількості повідомлень ТА повертає їх
	GetLimit(chatId, limit int) ([]models.Message, error)
	// DeleteAll отримує ID чату ТА видаляє його повідомлення
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
