package service

import (
	"cmd/pkg/repository"
	"cmd/pkg/repository/models"
)
//go:generate mockgen -source=service.go -destination=mocks/mock.go
type Authorization interface {
	// CreateUser кодує пароль викликає створення нового користувача
	CreateUser(user models.User) (int, error)
	// GetByName викликає повернення даних користувача за ім'ям
	GetByName(name string) (models.User, error)
	// GetUserById викликає отримання даних користувача за його ID
	GetUserById(userId int) (models.User, error)
	// GenerateToken отримує за ім'ям та паролем користувача його ID,
	// далі цей ID зашифровується у токен та повертається токен
	GenerateToken(username, password string) (string, error)
	// ParseToken отримує зашифрований токен, розшифровує його та
	// повертає ID користувача
	ParseToken(token string) (int, error)
	// UpdateData оновлює ім'я або зображення
	UpdateData(user models.User) error
	// UpdatePassword кодує пароль та оновлює його
	UpdatePassword(user models.User) error
}

type Chat interface {
	// Create викликає створення нового чату
	Create(chat models.Chat) (int, error)
	// Get викликає отримання даних чату
	Get(chatId int) (models.Chat, error)
	// Update викликає оновлення чату
	Update(chat models.Chat) error
	// Delete викликає видалення чату
	Delete(chatId int) error
	// AddUser викликає додання користувача до чату
	AddUser(users models.ChatUsers) (int, error)
	// GetUsers викликає отримання масиву користувачів чатом
	GetUsers(chatId int) ([]models.User, error)
	// DeleteUser викликає видалення користувача із чату
	DeleteUser(userId, chatId int) error
	// GetPrivates отримує два ID користувачів, повертає : при помилці - -1;
	// якщо чат вже існує - його ID; якщо чату немає - 0
	GetPrivates(firstUser, secondUser int) (int, error)
	// GetPrivateChats викликає отримання масиву публічних чатів користувача
	GetPrivateChats(userId int) ([]models.Chat, error)
	// GetPublicChats викликає отримання масиву приватних чатів користувача
	GetPublicChats(userId int) ([]models.Chat, error)
	// SearchChat викликає отримання масиву чатів, назви яких повністю чи
	// частково збігаються з аргументом
	SearchChat(name string) ([]models.Chat, error)
	// DeleteAllMessages викликає видалення усіх повідомлень чата за його ID
	DeleteAllMessages(chatId int) error
	// GetUserById викликає отримання даних користувача за його ID
	GetUserById(userId int) (models.User, error)
}

type Status interface {
	// AddStatus викликає створення нового статусу та повернення його ID
	AddStatus(status models.Status) (int, error)
	// GetStatuses викликає повернення даних щодо відносин між двома користувачами
	GetStatuses(senderId, recipientId int) ([]models.Status, error)
	// UpdateStatus викликає оновлення даних статусу
	UpdateStatus(status models.Status) error
	// DeleteStatus викликає видалення відносин між двома користувачами
	DeleteStatus(status models.Status) error
	// GetFriends викликає отримання списку користувачів, що мають статус друзів
	GetFriends(userId int) ([]models.User, error)
	// GetBlackList викликає отримання списку користувачів,
	// що для вас мають статус заблокованих
	GetBlackList(userId int) ([]models.User, error)
	// GetBlackListToUser викликає отримання списку користувачів,
	// для яких ви маєте статус заблокованого
	GetBlackListToUser(userId int) ([]models.User, error)
	// GetSentInvites викликає отримання списку користувачів,
	// що для вас мають статус запрошених у друзі
	GetSentInvites(userId int) ([]models.User, error)
	// GetInvites викликає отримання списку користувачів,
	// для яких ви маєте статус запрошеного у друзі
	GetInvites(userId int) ([]models.User, error)
	// SearchUser викликає отримання списку чатів, що мають частково або
	// повністю збіг з аргументом
	SearchUser(username string) ([]models.User, error)
	// GetUserById викликає отримання даних користувача за його ID
	GetUserById(userId int) (models.User, error)
}

type Message interface {
	// Create викликає створення нового повідомлення та повертає його ID
	Create(msg models.Message) (int, error)
	// Get викликає повернення повідомлення за його ID
	Get(msgId int) (models.Message, error)
	// GetAll викликає повернення усіх повідомлень чату за його ID
	GetAll(chatId int) ([]models.Message, error)
	// GetLimit викликає повернення певної кількості повідомлень чату за його ID
	GetLimit(chatId, limit int) ([]models.Message, error)
	// DeleteAll викликає видалення усіх повідомлень чата за його ID
	DeleteAll(chatId int) error
}

type Service struct {
	Authorization
	Chat
	Status
	Message
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Chat:          NewChatService(repos.Chat),
		Status:        NewStatusService(repos.Status),
		Message:       NewMessageService(repos.Message),
	}
}
