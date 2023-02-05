package repository

import (
	"cmd/pkg/repository/models"
	"fmt"
	"github.com/jinzhu/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateUser отримує ім'я та пароль ТА створює нового користувача
func (a *AuthRepository) CreateUser(user models.User) (int, error) {
	err := a.db.Table(UsersTable).Select("username", "password_hash").Create(&user).Error
	return user.Id, err
}

// GetUser отримує ім'я та пароль ТА повертає його дані
func (a *AuthRepository) GetUser(username, password string) (models.User, error) {
	var user models.User
	err := a.db.Table(UsersTable).Where("username = ? and password_hash = ?", username, password).First(&user).Error
	return user, err
}

// GetUserById отримує ID користувача ТА повертає його дані
func (a *AuthRepository) GetUserById(userId int) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id, username, icon FROM %s WHERE id = ?", UsersTable)
	err := a.db.Raw(query, userId).Scan(&user).Error
	//err := a.db.Table(UsersTable).Select("id", "username", "icon").First(&user, userId).Error
	return user, err
}

// GetByName отримує ім'я користувача ТА повертає його дані
func (a *AuthRepository) GetByName(username string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id, username, icon FROM %s WHERE username = ?", UsersTable)
	err := a.db.Raw(query, username).Scan(&user).Error
	return user, err
}

// UpdateUser отримує дані користувача ТА оновлює їх
func (a *AuthRepository) UpdateUser(user models.User) error {
	err := a.db.Table(UsersTable).Select("username", "icon").Where("id = ?", user.Id).Updates(&user).Error
	return err
}
