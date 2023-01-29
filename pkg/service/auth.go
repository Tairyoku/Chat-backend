package service

import (
	"cmd/pkg/repository"
	"cmd/pkg/repository/models"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

const (
	tokenTTL = 48 * time.Hour
)

type AuthService struct {
	repository repository.Authorization
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func NewAuthService(repository repository.Authorization) *AuthService {
	return &AuthService{repository: repository}
}

func (a *AuthService) CreateUser(user models.User) (int, error) {
	user.Password = CreatePasswordHash(user.Password)
	return a.repository.CreateUser(user)
}

func (a *AuthService) GetUserById(userId int) (models.User, error) {
	return a.repository.GetUserById(userId)
}

func (a *AuthService) CheckUsername(username string) error {
	return a.repository.CheckUsername(username)
}

func (a *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := a.repository.GetUser(username, CreatePasswordHash(password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.Id,
	})

	return token.SignedString([]byte(os.Getenv("signInKey")))
}

func (a *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("signInKey")), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}
	return claims.UserId, nil
}

func (a *AuthService) UpdatePassword(user models.User) error {
	user.Password = CreatePasswordHash(user.Password)
	err := a.repository.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthService) UpdateData(user models.User) error {
	err := a.repository.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthService) GetByName(name string) (models.User, error) {
	return a.repository.GetByName(name)
}

func CreatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("salt"))))
}
