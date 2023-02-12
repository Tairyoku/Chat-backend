package repository

import (
	"cmd/pkg/repository/models"
	"fmt"
	"github.com/jinzhu/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

// Create отримує назву чату ТА створює новий чат
func (c *ChatRepository) Create(chat models.Chat) (int, error) {
	err := c.db.Table(ChatsTable).Select("name", "types").Create(&chat).Error
	return chat.Id, err
}

// Get отримує ID чату ТА повертає дані чату за його ID
func (c *ChatRepository) Get(chatId int) (models.Chat, error) {
	var chat models.Chat
	err := c.db.Table(ChatsTable).First(&chat, chatId).Error
	return chat, err
}

// Update отримує дані чату ТА оновлює їх
func (c *ChatRepository) Update(chat models.Chat) error {
	err := c.db.Table(ChatsTable).Select("name", "icon").Where("id = ?", chat.Id).Updates(&chat).Error
	return err
}

// Delete отримує ID чату ТА видаляє чат
func (c *ChatRepository) Delete(chatId int) error {
	err := c.db.Table(ChatsTable).Delete(&models.Chat{}, chatId).Error
	return err
}

// AddUser отримує ID чату ТА ID користувача, та додає користувача до чату
func (c *ChatRepository) AddUser(user models.ChatUsers) (int, error) {
	err := c.db.Select(ChatUsersList, "chat_id", "user_id").Create(&user).Error
	return user.Id, err
}

// CheckPrivates отримує ID двох можливих користувачів одного приватного
// чату ТА повертає масив ID приватних чатів, у яких двічі згадується ID
// користувачів, поданих як аргументи
func (c *ChatRepository) CheckPrivates(firstUser, secondUser int) ([]int, error) {
	var check []int
	var list []models.ChatUsers
	var chats []int
	result := make(map[int]int)
	queryFirst := fmt.Sprintf("SELECT chul.chat_id, chul.user_id FROM %s chul INNER JOIN %s chl ON chl.id = chul.chat_id WHERE (chul.user_id = ? or chul.user_id = ?) and chl.types = ?",
		ChatUsersList, ChatsTable)
	errFirst := c.db.Raw(queryFirst, firstUser, secondUser, ChatPrivate).Scan(&list).Error
	if errFirst != nil {
		return nil, errFirst
	}
	for i := range list {
		if list[i].UserId != firstUser {
			check = append(check, list[i].ChatId)
		}
	}
	for _, v := range check {
		result[v]++
	}
	for k, v := range result {
		if v >= 1 {
			chats = append(chats, k)
		}
	}

	return chats, nil
}

// GetByName отримує назву чату ТА повертає його дані
// ?????????????імена вже не унікальні????????????????
//func (c *ChatRepository) GetByName(name string) (models.Chat, error) {
//	var chat models.Chat
//	err := c.db.Table(ChatsTable).Where("name = ?", name).First(&chat).Error
//	return chat, err
//}

// GetUsers отримує ID чату ТА повертає масив користувачів, що приєднані до чату
func (c *ChatRepository) GetUsers(chatId int) ([]models.User, error) {
	var users []models.User
	query := fmt.Sprintf("SELECT u.id, u.username, u.icon FROM %s u INNER JOIN %s chl ON u.id = chl.user_id WHERE chl.chat_id = ?", UsersTable, ChatUsersList)
	err := c.db.Raw(query, chatId).Scan(&users).Error
	return users, err
}

// GetPrivates отримує ID двох користувачів, ТА повертає масиви приватних
// чатів, до яких належать кожен із користувачів
func (c *ChatRepository) GetPrivates(firstUser, secondUser int) ([]models.Chat, []models.Chat, error) {
	var first []models.Chat
	var second []models.Chat
	tx := c.db.Begin()
	queryFirst := fmt.Sprintf("SELECT chl.id FROM %s chl INNER JOIN %s chul ON chl.id = chul.chat_id WHERE chul.user_id = ? and chl.types = ?",
		ChatsTable, ChatUsersList)
	errFirst := c.db.Raw(queryFirst, firstUser, ChatPrivate).Scan(&first).Error
	if errFirst != nil {
		tx.Rollback()
		return nil, nil, errFirst
	}
	querySecond := fmt.Sprintf("SELECT chl.id FROM %s chl INNER JOIN %s chul ON chl.id = chul.chat_id WHERE chul.user_id = ? and chl.types = ?",
		ChatsTable, ChatUsersList)
	errSecond := c.db.Raw(querySecond, secondUser, ChatPrivate).Scan(&second).Error
	if errSecond != nil {
		tx.Rollback()
		return nil, nil, errSecond
	}
	return first, second, nil
}

// GetPrivateChats отримує ID користувача ТА повертає масив ПРИВАТНИХ чатів,
// до яких він належить
func (c *ChatRepository) GetPrivateChats(userId int) ([]models.Chat, error) {
	var chats []models.Chat
	query := fmt.Sprintf("SELECT * FROM %s ch INNER JOIN %s chl ON ch.id = chl.chat_id WHERE chl.user_id = ? and ch.types = ?", ChatsTable, ChatUsersList)
	err := c.db.Raw(query, userId, ChatPrivate).Scan(&chats).Error
	return chats, err
}

// GetPublicChats отримує ID користувача ТА повертає масив ПУБЛІЧНИХ чатів,
// до яких він належить
func (c *ChatRepository) GetPublicChats(userId int) ([]models.Chat, error) {
	var chats []models.Chat
	query := fmt.Sprintf("SELECT * FROM %s ch INNER JOIN %s chl ON ch.id = chl.chat_id WHERE chl.user_id = ? and ch.types = ?", ChatsTable, ChatUsersList)
	err := c.db.Raw(query, userId, ChatPublic).Scan(&chats).Error
	return chats, err
}

// DeleteUser отримує ID чату ТА ID користувача, та видаляє користувача із чату
func (c *ChatRepository) DeleteUser(userId, chatId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = ? and chat_id = ?", ChatUsersList)
	err := c.db.Raw(query, userId, chatId).Scan(&models.ChatUsers{}).Error
	return err
}

// SearchChat отримує назву чату (або його частину) ТА повертає масив чатів,
// назви яких збігаються з аргументом
func (c *ChatRepository) SearchChat(name string) ([]models.Chat, error) {
	var chats []models.Chat
	//query := fmt.Sprintf("SELECT id, username, icon FROM %s  WHERE name LIKE ? and types = ?", UsersTable)
	//err := c.db.Raw(query, fmt.Sprintf("%%%s%%", name), ChatPublic).Scan(&chats).Error
	err := c.db.Table(ChatsTable).Where("types = ? AND name LIKE ?", ChatPublic, fmt.Sprintf("%%%s%%", name)).Find(&chats).Error
	return chats, err
}
