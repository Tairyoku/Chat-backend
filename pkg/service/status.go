package service

import (
	"cmd/pkg/repository"
	"cmd/pkg/repository/models"
)

type StatusService struct {
	repository repository.Status
}

func NewStatusService(repository repository.Status) *StatusService {
	return &StatusService{repository: repository}
}

func (s *StatusService) AddStatus(status models.Status) (int, error) {
	return s.repository.AddStatus(status)
}

func (s *StatusService) DeleteStatus(status models.Status) error {
	return s.repository.DeleteStatus(status)
}

func (s *StatusService) UpdateStatus(status models.Status) error {
	return s.repository.UpdateStatus(status)
}

func (s *StatusService) GetStatus(senderId, recipientId int) (models.Status, error) {
	return s.repository.GetStatus(senderId, recipientId)
}

func (s *StatusService) GetFriends(userId int) ([]int, error) {
	return s.repository.GetFriends(userId)
}

func (s *StatusService) GetBlackList(userId int) ([]int, error) {
	return s.repository.GetBlackList(userId)
}

func (s *StatusService) GetBlackListToUser(userId int) ([]int, error) {
	return s.repository.GetBlackListToUser(userId)
}

func (s *StatusService) GetSentInvites(userId int) ([]int, error) {
	return s.repository.GetSentInvites(userId)
}

func (s *StatusService) GetInvites(userId int) ([]int, error) {
	return s.repository.GetInvites(userId)
}

func (s *StatusService) SearchUser(username string) ([]models.User, error) {
	return s.repository.SearchUser(username)
}
