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

// Create викликає створення нового чату
func (c *ChatService) Create(chat models.Chat) (int, error) {
	return c.repository.Create(chat)
}

// Get викликає отримання даних чату
func (c *ChatService) Get(chatId int) (models.Chat, error) {
	return c.repository.Get(chatId)
}

// Update викликає оновлення даних чату
func (c *ChatService) Update(chat models.Chat) error {
	return c.repository.Update(chat)
}

// Delete викликає видалення чату
func (c *ChatService) Delete(chatId int) error {
	return c.repository.Delete(chatId)
}

// AddUser викликає додання користувача до чату
func (c *ChatService) AddUser(users models.ChatUsers) (int, error) {
	return c.repository.AddUser(users)
}

// GetUsers викликає отримання масиву користувачів чатом
func (c *ChatService) GetUsers(chatId int) ([]models.User, error) {
	return c.repository.GetUsers(chatId)
}

// DeleteUser викликає видалення користувача із чату
func (c *ChatService) DeleteUser(userId, chatId int) error {
	return c.repository.DeleteUser(userId, chatId)
}

// GetPrivates отримує два ID користувачів, повертає : при помилці - -1;
// якщо чат вже існує - його ID; якщо чату немає - 0
func (c *ChatService) GetPrivates(firstUser, secondUser int) (int, error) {

	// Отримуємо список ПРИВАТНИХ чатів, в яких присутні перший чи другий користувачі
	first, second, err := c.repository.GetPrivates(firstUser, secondUser)
	if err != nil {
		return -1, err
	}

	// якщо обидва ID однакові, отже це особистий чат
	//повертаємо ID особистого чату
	if firstUser == secondUser {
		for _, chat := range first {
			list, err := c.repository.GetUsers(chat.Id)
			if err != nil {
				return -1, err
			}
			if len(list) == 1 {
				return chat.Id, err
			}
		}
		return 0, err
	}

	//якщо ID різні, шукаємо спільний чат

	// Збираємо усі ID чатів до єдиного масиву
	var check []int
	for i := range first {
		check = append(check, first[i].Id)
	}
	for i := range second {
		check = append(check, second[i].Id)
	}

	// Перевіряємо наявність дублікатів
	// Необхідний чат згадується двічі
	result := make(map[int]int)
	for _, v := range check {
		result[v]++
	}
	for k, v := range result {
		if v == 2 {
			return k, err
		}
	}
	return 0, nil
}

// GetPrivateChats викликає отримання масиву публічних чатів користувача
func (c *ChatService) GetPrivateChats(userId int) ([]models.Chat, error) {
	return c.repository.GetPrivateChats(userId)
}

// GetPublicChats викликає отримання масиву приватних чатів користувача
func (c *ChatService) GetPublicChats(userId int) ([]models.Chat, error) {
	return c.repository.GetPublicChats(userId)
}

// SearchChat викликає отримання масиву чатів, назви яких повністю чи
// частково збігаються з аргументом
func (c *ChatService) SearchChat(name string) ([]models.Chat, error) {
	return c.repository.SearchChat(name)
}

// DeleteAllMessages викликає видалення усіх повідомлень чата за його ID
func (c *ChatService) DeleteAllMessages(chatId int) error {
	return c.repository.DeleteAllMessages(chatId)
}

// GetUserById викликає отримання даних користувача за його ID
func (c *ChatService) GetUserById(userId int) (models.User, error) {
	return c.repository.GetUserById(userId)
}
