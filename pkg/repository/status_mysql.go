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

func (s *StatusRepository) AddStatus(status models.Status) (int, error) {
	err := s.db.Table(StatusesTable).Create(&status).Error
	return status.Id, err
}

func (s *StatusRepository) UpdateStatus(status models.Status) error {
	err := s.db.Select(StatusesTable, "relationship").Updates(&status).Error
	return err
}

func (s *StatusRepository) GetStatus(senderId, recipientId int) (models.Status, error) {
	var status models.Status
	err := s.db.Table(StatusesTable).Where("sender_id = ? and recipient_id = ?", senderId, recipientId).First(&status).Error
	return status, err
}

func (s *StatusRepository) DeleteStatus(status models.Status) error {
	//err := s.db.Table(StatusesTable).Where("sender_id = ? and recipient_id = ? and relationship = ?",
	//status.SenderId, status.RecipientId, status.Relationship).Delete(&models.Status{}).Error
	err := s.db.Table(StatusesTable).Delete(&status).Error
	return err
}

func (s *StatusRepository) GetFriends(userId int) ([]int, error) {
	var usersId []models.Status
	var result []int
	tx := s.db.Begin()
	//err := tx.Table(StatusesTable).Where("relationship = ? and recipient_id = ?",
	//	StatusFriends, userId).Find(&usersId).Error
	query := fmt.Sprintf("SELECT sender_id FROM %s WHERE relationship = ? and recipient_id = ?", StatusesTable)
	err := tx.Raw(query, StatusFriends, userId).Scan(&usersId).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for i := range usersId {
		result = append(result, usersId[i].SenderId)
	}
	//errt := tx.Table(StatusesTable).Where("relationship = ? and sender_id = ?",
	//	StatusFriends, userId).Find(&usersId).Error
	queryt := fmt.Sprintf("SELECT recipient_id FROM %s WHERE relationship = ? and sender_id = ?", StatusesTable)
	errt := tx.Raw(queryt, StatusFriends, userId).Scan(&usersId).Error
	if errt != nil {
		tx.Rollback()
		return nil, errt
	}
	for i := range usersId {
		result = append(result, usersId[i].RecipientId)
	}
	return result, tx.Commit().Error
}

func (s *StatusRepository) GetBlackList(userId int) ([]int, error) {
	var usersId []models.Status
	var result []int
	//err := s.db.Table(StatusesTable).Select("recipient_id").Where("relationship = ? and sender_id = ?",
	//	StatusBL, userId).Find(&usersId).Error
	query := fmt.Sprintf("SELECT recipient_id FROM %s WHERE relationship = ? and sender_id = ?", StatusesTable)
	err := s.db.Raw(query, StatusBL, userId).Scan(&usersId).Error
	for i := range usersId {
		result = append(result, usersId[i].RecipientId)
	}
	return result, err
}

func (s *StatusRepository) GetBlackListToUser(userId int) ([]int, error) {
	var usersId []models.Status
	var result []int
	//err := s.db.Table(StatusesTable).Select("sender_id").Where("relationship = ? and recipient_id = ?",
	//	StatusBL, userId).Find(&usersId).Error
	query := fmt.Sprintf("SELECT sender_id FROM %s WHERE relationship = ? and recipient_id = ?", StatusesTable)
	err := s.db.Raw(query, StatusBL, userId).Scan(&usersId).Error
	for i := range usersId {
		result = append(result, usersId[i].SenderId)
	}
	return result, err
}

func (s *StatusRepository) GetSentInvites(userId int) ([]int, error) {
	var usersId []models.Status
	var result []int
	//err := s.db.Table(StatusesTable).Select("recipient_id").Where("relationship = ? and sender_id = ?",
	//	StatusInvitation, userId).Find(&usersId).Error
	query := fmt.Sprintf("SELECT recipient_id FROM %s WHERE relationship = ? and sender_id = ?", StatusesTable)
	err := s.db.Raw(query, StatusInvitation, userId).Scan(&usersId).Error
	for i := range usersId {
		result = append(result, usersId[i].RecipientId)
	}
	return result, err
}

func (s *StatusRepository) GetInvites(userId int) ([]int, error) {
	var usersId []models.Status
	var result []int
	//err := s.db.Table(StatusesTable).Select("sender_id").Where("relationship = ? and recipient_id = ?",
	//	StatusInvitation, userId).Find(&usersId).Error
	query := fmt.Sprintf("SELECT sender_id FROM %s WHERE relationship = ? and recipient_id = ?", StatusesTable)
	err := s.db.Raw(query, StatusInvitation, userId).Scan(&usersId).Error
	for i := range usersId {
		result = append(result, usersId[i].SenderId)
	}
	return result, err
}

func (s *StatusRepository) SearchUser(username string) ([]models.User, error) {
	var users []models.User
	//err := s.db.Table(UsersTable).Where("username LIKE ?", fmt.Sprintf("%%%s%%", username)).Find(&users).Error
	query := fmt.Sprintf("SELECT id, username, icon FROM %s WHERE username LIKE ?", UsersTable)
	err := s.db.Raw(query, fmt.Sprintf("%%%s%%", username)).Scan(&users).Error
	return users, err
}
