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
	var usersId []int
	tx := s.db.Begin()
	err := tx.Table(StatusesTable).Select("sender_id").Where("relationship = ? and recipient_id = ?",
		StatusFriends, userId).Find(&usersId).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	result := usersId
	errt := tx.Table(StatusesTable).Select("recipient_id").Where("relationship = ? and sender_id = ?",
		StatusFriends, userId).Find(&usersId).Error
	if errt != nil {
		tx.Rollback()
		return nil, errt
	}
	for i := range usersId {
		result = append(result, usersId[i])
	}
	return result, tx.Commit().Error
}

func (s *StatusRepository) GetBlackList(userId int) ([]int, error) {
	var usersId []int
	err := s.db.Table(StatusesTable).Select("recipient_id").Where("relationship = ? and sender_id = ?",
		StatusBL, userId).Find(&usersId).Error
	return usersId, err
}

func (s *StatusRepository) GetBlackListToUser(userId int) ([]int, error) {
	var usersId []int
	err := s.db.Table(StatusesTable).Select("sender_id").Where("relationship = ? and recipient_id = ?",
		StatusBL, userId).Find(&usersId).Error
	return usersId, err
}

func (s *StatusRepository) GetSentInvites(userId int) ([]int, error) {
	var usersId []int
	err := s.db.Table(StatusesTable).Select("recipient_id").Where("relationship = ? and sender_id = ?",
		StatusInvitation, userId).Find(&usersId).Error
	return usersId, err
}

func (s *StatusRepository) GetInvites(userId int) ([]int, error) {
	var usersId []int
	err := s.db.Table(StatusesTable).Select("sender_id").Where("relationship = ? and recipient_id = ?",
		StatusInvitation, userId).Find(&usersId).Error
	return usersId, err
}

func (s *StatusRepository) SearchUser(username string) ([]models.User, error) {
	var user []models.User
	err := s.db.Table(UsersTable).Select("id", "username").Where("username LIKE ?", fmt.Sprintf("%%%s%%", username)).Find(&user).Error
	return user, err
}
