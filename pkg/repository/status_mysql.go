package repository

import (
	"cmd/pkg/repository/models"
	"gorm.io/gorm"
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
	err := s.db.Select(StatusesTable, "sender_id").Where("relationship = ? and recipient_id = ?",
		StatusFriends, userId).Select(StatusesTable, "recipient_id").Where("relationship = ? and sender_id = ?",
		StatusFriends, userId).Find(&usersId).Error
	return usersId, err
}

func (s *StatusRepository) GetBlackList(userId int) ([]int, error) {
	var usersId []int
	err := s.db.Select(StatusesTable, "recipient_id").Where("relationship = ? and sender_id = ?",
		StatusBL, userId).Find(&usersId).Error
	return usersId, err
}

func (s *StatusRepository) GetBlackListToUser(userId int) ([]int, error) {
	var usersId []int
	err := s.db.Select(StatusesTable, "sender_id").Where("relationship = ? and recipient_id = ?",
		StatusBL, userId).Find(&usersId).Error
	return usersId, err
}

func (s *StatusRepository) GetSentInvites(userId int) ([]int, error) {
	var usersId []int
	err := s.db.Select(StatusesTable, "recipient_id").Where("relationship = ? and sender_id = ?",
		StatusInvitation, userId).Find(&usersId).Error
	return usersId, err
}

func (s *StatusRepository) GetInvites(userId int) ([]int, error) {
	var usersId []int
	err := s.db.Select(StatusesTable, "sender_id").Where("relationship = ? and recipient_id = ?",
		StatusInvitation, userId).Find(&usersId).Error
	return usersId, err
}
