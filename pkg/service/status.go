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

// AddStatus викликає створення нового статусу та повернення його ID
func (s *StatusService) AddStatus(status models.Status) (int, error) {
	return s.repository.AddStatus(status)
}

// GetStatuses викликає повернення даних щодо відносин між двома користувачами
func (s *StatusService) GetStatuses(senderId, recipientId int) ([]models.Status, error) {
	return s.repository.GetStatuses(senderId, recipientId)
}

// UpdateStatus викликає оновлення даних статусу
func (s *StatusService) UpdateStatus(status models.Status) error {
	return s.repository.UpdateStatus(status)
}

// DeleteStatus викликає видалення відносин між двома користувачами
func (s *StatusService) DeleteStatus(status models.Status) error {
	return s.repository.DeleteStatus(status)
}

// GetFriends викликає отримання списку користувачів, що мають статус друзів
func (s *StatusService) GetFriends(userId int) ([]models.User, error) {
	return s.repository.GetFriends(userId)
}

// GetBlackList викликає отримання списку користувачів,
// що для вас мають статус заблокованих
func (s *StatusService) GetBlackList(userId int) ([]models.User, error) {
	return s.repository.GetBlackList(userId)
}

// GetBlackListToUser викликає отримання списку користувачів,
// для яких ви маєте статус заблокованого
func (s *StatusService) GetBlackListToUser(userId int) ([]models.User, error) {
	return s.repository.GetBlackListToUser(userId)
}

// GetSentInvites викликає отримання списку користувачів,
// що для вас мають статус запрошених у друзі
func (s *StatusService) GetSentInvites(userId int) ([]models.User, error) {
	return s.repository.GetSentInvites(userId)
}

// GetInvites викликає отримання списку користувачів,
// для яких ви маєте статус запрошеного у друзі
func (s *StatusService) GetInvites(userId int) ([]models.User, error) {
	return s.repository.GetInvites(userId)
}

// SearchUser викликає отримання списку чатів, що мають частково або
// повністю збіг з аргументом
func (s *StatusService) SearchUser(username string) ([]models.User, error) {
	return s.repository.SearchUser(username)
}
