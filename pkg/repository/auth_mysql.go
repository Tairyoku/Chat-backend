package repository

import (
	"cmd/pkg/repository/models"
	"github.com/jinzhu/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (a *AuthRepository) CreateUser(user models.User) (int, error) {
	err := a.db.Table(UsersTable).Select("username", "password_hash").Create(&user).Error
	return user.Id, err
}

func (a *AuthRepository) GetUserById(userId int) (models.User, error) {
	var user models.User
	err := a.db.Table(UsersTable).Select("id", "username", "icon").First(&user, userId).Error
	return user, err
}

func (a *AuthRepository) GetUser(username, password string) (models.User, error) {
	var user models.User
	err := a.db.Table(UsersTable).Where("username = ? and password_hash = ?", username, password).First(&user).Error
	return user, err
}

func (a *AuthRepository) UpdateUser(user models.User) error {
	err := a.db.Table(UsersTable).Where("id = ?", user.Id).Updates(&user).Error
	return err
}

func (a *AuthRepository) CheckUsername(username string) error {
	var user models.User
	err := a.db.Table(UsersTable).Where("username = ?", username).First(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthRepository) GetByName(name string) (models.User, error) {
	var user models.User
	err := a.db.Table(UsersTable).Select("id", "username", "icon").Where("username = ?", name).First(&user).Error
	return user, err
}
