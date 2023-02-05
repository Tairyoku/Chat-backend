package repository

import (
	"cmd/pkg/repository/models"
	"fmt"
	"github.com/jinzhu/gorm"
)

type StatusRepository struct {
	db *gorm.DB
}

func NewStatusRepository(db *gorm.DB) *StatusRepository {
	return &StatusRepository{db: db}
}

// AddStatus отримує ID двох користувачів та їх тип відносин ТА повертає ID статусу
func (s *StatusRepository) AddStatus(status models.Status) (int, error) {
	err := s.db.Table(StatusesTable).Create(&status).Error
	return status.Id, err
}

// GetStatuses отримує ID двох користувачів ТА повертає дані їх відносин
func (s *StatusRepository) GetStatuses(senderId, recipientId int) ([]models.Status, error) {
	var status []models.Status
	err := s.db.Table(StatusesTable).Where("(sender_id = ? and recipient_id = ?) or (sender_id = ? and recipient_id = ? and relationship = ?)",
		senderId, recipientId, recipientId, senderId, StatusFriends).First(&status).Error
	return status, err
}

// UpdateStatus отримує ID двох користувачів та їх тип відносин ТА оновлює дані
func (s *StatusRepository) UpdateStatus(status models.Status) error {
	err := s.db.Table(StatusesTable).Updates(&status).Error
	return err
}

// DeleteStatus отримує ID двох користувачів та їх тип відносин ТА видаляє ці відносини
func (s *StatusRepository) DeleteStatus(status models.Status) error {
	query := fmt.Sprintf("DELETE FROM %s stl WHERE relationship = ? and (recipient_id = ? and sender_id = ?) or (sender_id = ? and recipient_id = ?)", StatusesTable)
	err := s.db.Raw(query, status.Relationship, status.SenderId, status.RecipientId, status.SenderId, status.RecipientId).Scan(&status).Error
	//err := s.db.Table(StatusesTable).Where("sender_id = ? and recipient_id = ? and relationship = ?",
	//status.SenderId, status.RecipientId, status.Relationship).Delete(&models.Status{}).Error
	//err := s.db.Table(StatusesTable).Delete(&status).Error
	return err
}

// GetFriends отримує ID користувача ТА повертає масив користувачів, що є ДРУЗЯМИ
func (s *StatusRepository) GetFriends(userId int) ([]models.User, error) {
	var usersId []models.User
	var result []models.User
	tx := s.db.Begin()

	query := fmt.Sprintf("SELECT u.id, u.username, u.icon FROM %s u INNER JOIN %s chul ON chul.sender_id = u.id WHERE relationship = ? and recipient_id = ?", UsersTable, StatusesTable)
	err := tx.Raw(query, StatusFriends, userId).Scan(&usersId).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	result = usersId

	queryt := fmt.Sprintf("SELECT u.id, u.username, u.icon FROM %s u INNER JOIN %s chul ON chul.recipient_id = u.id WHERE relationship = ? and sender_id = ?", UsersTable, StatusesTable)
	errt := tx.Raw(queryt, StatusFriends, userId).Scan(&usersId).Error
	if errt != nil {
		tx.Rollback()
		return nil, errt
	}
	for i := range usersId {
		result = append(result, usersId[i])
	}
	return result, tx.Commit().Error
}

// GetBlackList отримує ID користувача ТА повертає масив ЗАБЛОКОВАНИХ користувачів
func (s *StatusRepository) GetBlackList(userId int) ([]models.User, error) {
	var users []models.User
	//query := fmt.Sprintf("SELECT recipient_id FROM %s WHERE relationship = ? and sender_id = ?", StatusesTable)
	//err := s.db.Raw(query, StatusBL, userId).Scan(&usersId).Error
	query := fmt.Sprintf("SELECT u.id, u.username, u.icon FROM %s u INNER JOIN %s chul ON chul.recipient_id = u.id WHERE relationship = ? and sender_id = ?", UsersTable, StatusesTable)
	err := s.db.Raw(query, StatusBL, userId).Scan(&users).Error

	return users, err
}

// GetBlackListToUser отримує ID користувача ТА повертає масив користувачів, що
// ЗАБЛОКУВАЛИ його
func (s *StatusRepository) GetBlackListToUser(userId int) ([]models.User, error) {
	var users []models.User
	//query := fmt.Sprintf("SELECT sender_id FROM %s WHERE relationship = ? and recipient_id = ?", StatusesTable)
	//err := s.db.Raw(query, StatusBL, userId).Scan(&usersId).Error
	query := fmt.Sprintf("SELECT u.id, u.username, u.icon FROM %s u INNER JOIN %s chul ON chul.sender_id = u.id WHERE relationship = ? and recipient_id = ?", UsersTable, StatusesTable)
	err := s.db.Raw(query, StatusBL, userId).Scan(&users).Error

	return users, err
}

// GetSentInvites отримує ID користувача ТА повертає масив користувачів, що
// ОТРИМАЛИ його запрошення у друзі
func (s *StatusRepository) GetSentInvites(userId int) ([]models.User, error) {
	var users []models.User
	//query := fmt.Sprintf("SELECT recipient_id FROM %s WHERE relationship = ? and sender_id = ?", StatusesTable)
	//err := s.db.Raw(query, StatusInvitation, userId).Scan(&usersId).Error
	query := fmt.Sprintf("SELECT u.id, u.username, u.icon FROM %s u INNER JOIN %s chul ON chul.recipient_id = u.id WHERE relationship = ? and sender_id = ?", UsersTable, StatusesTable)
	err := s.db.Raw(query, StatusInvitation, userId).Scan(&users).Error

	return users, err
}

// GetInvites отримує ID користувача ТА повертає масив користувачів, що
// НАДІСЛАЛИ йому запрошення в друзі
func (s *StatusRepository) GetInvites(userId int) ([]models.User, error) {
	var users []models.User
	//query := fmt.Sprintf("SELECT sender_id FROM %s WHERE relationship = ? and recipient_id = ?", StatusesTable)
	//err := s.db.Raw(query, StatusInvitation, userId).Scan(&usersId).Error
	query := fmt.Sprintf("SELECT u.id, u.username, u.icon FROM %s u INNER JOIN %s chul ON chul.sender_id = u.id WHERE relationship = ? and recipient_id = ?", UsersTable, StatusesTable)
	err := s.db.Raw(query, StatusInvitation, userId).Scan(&users).Error

	return users, err
}

// SearchUser отримує ім'я (або його частину) ТА повертає масив користувачів, що
// мають збіг з аргументом
func (s *StatusRepository) SearchUser(username string) ([]models.User, error) {
	var users []models.User
	//err := s.db.Table(UsersTable).Where("username LIKE ?", fmt.Sprintf("%%%s%%", username)).Find(&users).Error
	query := fmt.Sprintf("SELECT id, username, icon FROM %s WHERE username LIKE ?", UsersTable)
	err := s.db.Raw(query, fmt.Sprintf("%%%s%%", username)).Scan(&users).Error
	return users, err
}

//func (s *StatusRepository) GetFriends(userId int) ([]int, error) {
//	var usersId []models.Status
//	var result []int
//	tx := s.db.Begin()
//	//err := tx.Table(StatusesTable).Where("relationship = ? and recipient_id = ?",
//	//	StatusFriends, userId).Find(&usersId).Error
//	query := fmt.Sprintf("SELECT sender_id FROM %s WHERE relationship = ? and recipient_id = ?", StatusesTable)
//	err := tx.Raw(query, StatusFriends, userId).Scan(&usersId).Error
//	if err != nil {
//		tx.Rollback()
//		return nil, err
//	}
//	for i := range usersId {
//		result = append(result, usersId[i].SenderId)
//	}
//	//errt := tx.Table(StatusesTable).Where("relationship = ? and sender_id = ?",
//	//	StatusFriends, userId).Find(&usersId).Error
//	queryt := fmt.Sprintf("SELECT recipient_id FROM %s WHERE relationship = ? and sender_id = ?", StatusesTable)
//	errt := tx.Raw(queryt, StatusFriends, userId).Scan(&usersId).Error
//	if errt != nil {
//		tx.Rollback()
//		return nil, errt
//	}
//	for i := range usersId {
//		result = append(result, usersId[i].RecipientId)
//	}
//	return result, tx.Commit().Error
//}
