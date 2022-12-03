package repository

import (
	"cmd/pkg/repository/models"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (a *AuthRepository) CreateUser(user models.User) (int, error) {
	err := a.db.Select(UsersTable, "username", "password_hash").Create(&user).Error
	return user.Id, err
}

func (a *AuthRepository) GetUser(username, password string) (models.User, error) {
	var user models.User
	err := a.db.Where("username = ? and password_hash = ?", username, password).Find(&user).Error
	return user, err
}

func (a *AuthRepository) CheckUsername(username string) error {
	var user models.User
	err := a.db.Table(UsersTable).Where("username = ?", username).First(&user).Error
	return err
}
