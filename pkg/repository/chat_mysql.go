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

func (c *ChatRepository) Create(chat models.Chat) (int, error) {
	err := c.db.Select(ChatsTable, "name", "types").Create(&chat).Error
	return chat.Id, err
}

func (c *ChatRepository) Delete(chatId int) error {
	err := c.db.Table(ChatsTable).Delete(&models.Chat{}, chatId).Error
	return err
}

func (c *ChatRepository) Get(chatId int) (models.Chat, error) {
	var chat models.Chat
	err := c.db.Table(ChatsTable).First(&chat, chatId).Error
	return chat, err
}

func (c *ChatRepository) CheckPrivates(firstUser, secondUser int) ([]int, error) {
	var check []int
	var chats []int
	result := make(map[int]int)
	queryFirst := fmt.Sprintf("SELECT chul.chat_id FROM %s chul INNER JOIN %s chl ON chl.id = chul.chat_id WHERE (chul.user_id = ? or chul.user_id = ?) and chl.types = ?",
		ChatUsersList, ChatsTable)
	errFirst := c.db.Raw(queryFirst, firstUser, secondUser, ChatPrivate).Scan(&check).Error
	if errFirst != nil {
		return nil, errFirst
	}
	for _, v := range check {
		result[v]++
	}
	for k, v := range result {
		if v >= 2 {
			chats = append(chats, k)
		}
	}

	return chats, nil
}

func (c *ChatRepository) AddUser(user models.ChatUsers) (int, error) {
	err := c.db.Select(ChatUsersList, "chat_id", "user_id").Create(&user).Error
	return user.Id, err
}

func (c *ChatRepository) DeleteUser(userId, chatId int) error {
	err := c.db.Table(ChatUsersList).Where("chat_id = ?", chatId).Delete(&models.ChatUsers{}, userId).Error
	return err
}

func (c *ChatRepository) GetUsers(chatId int) ([]models.User, error) {
	var users []models.User
	query := fmt.Sprintf("SELECT u.id, u.username, u.icon FROM %s u INNER JOIN %s chl ON u.id = chl.user_id WHERE chl.chat_id = ?", UsersTable, ChatUsersList)
	err := c.db.Raw(query, chatId).Scan(&users).Error
	return users, err
}

func (c *ChatRepository) GetPrivateChats(userId int) ([]models.Chat, error) {
	var chats []models.Chat
	query := fmt.Sprintf("SELECT * FROM %s ch INNER JOIN %s chl ON ch.id = chl.chat_id WHERE chl.user_id = ? and ch.types = ?", ChatsTable, ChatUsersList)
	err := c.db.Raw(query, userId, ChatPrivate).Scan(&chats).Error
	return chats, err
}

func (c *ChatRepository) GetPublicChats(userId int) ([]models.Chat, error) {
	var chats []models.Chat
	query := fmt.Sprintf("SELECT * FROM %s ch INNER JOIN %s chl ON ch.id = chl.chat_id WHERE chl.user_id = ? and ch.types = ?", ChatsTable, ChatUsersList)
	err := c.db.Raw(query, userId, ChatPublic).Scan(&chats).Error
	return chats, err
}

func (c *ChatRepository) SearchChat(name string) ([]models.Chat, error) {
	var chats []models.Chat
	err := c.db.Table(UsersTable).Select("id", "name").Where("name LIKE ? and types = ?", fmt.Sprintf("%%%s%%", name), ChatPublic).Find(&chats).Error
	return chats, err
}

//tx := c.db.Begin()
//query := fmt.Sprintf("INSERT INTO %s 'name' values $1", ChatsTable)
//errFirstTx := tx.Raw(query, chat.Name).Error
//if errFirstTx != nil {
//	tx.Rollback()
//	return 0, errFirstTx
//}
//
//row := fmt.Sprintf("INSERT INTO %s 'name' values $1", ChatUsersList)
//errSecondTx := tx.Raw(row, chat.Name).Error
//if errSecondTx != nil {
//	tx.Rollback()
//	return 0, errFirstTx
//}
