package auth

import (
	"cmd/pkg/handler/middlewares"
	"cmd/pkg/handler/responses"
	"cmd/pkg/repository/models"
	"cmd/pkg/service"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strings"
)

type AuthHandler struct {
	services *service.Service
}

func NewAuthHandler(services *service.Service) *AuthHandler {
	return &AuthHandler{services: services}
}

// SignUp godoc
// @Summary      Create a new user
// @Security ApiKeyAuth
// @Description  Користувач відправляє ім'я та пароль. За отриманими даними буде створено нового користувача.
// Сервер поверне ID нового користувача
// @ID  add-new-user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user	body     UserResponse   true  "Add user"
// @Success      200 	{object} IdResponse		 "result is id of user"
// @Failure 	 400 	{object} responses.ErrorResponse	 "incorrect request data"
// @Failure 	 404 	{object} responses.ErrorResponse	 "user id not found"
// @Failure 	 500 	{object} responses.ErrorResponse	 "something went wrong"
// @Router       /auth/sign-up [post]
func (h *AuthHandler) SignUp(c echo.Context) error {

	// Отримуємо дані з сайту (ім'я та пароль)
	var input models.User
	if errReq := c.Bind(&input); errReq != nil {
		responses.NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return nil
	}

	// Перевіряємо отримані дані
	{
		//username is not empty
		if len(input.Username) == 0 {
			responses.NewErrorResponse(c, http.StatusBadRequest, "You must enter a username")
			return nil
		}

		// password length
		if len(input.Password) < 6 {
			responses.NewErrorResponse(c, http.StatusBadRequest, "Password must be at least 6 symbols")
			return nil
		}
	}

	// Створюємо нового користувача
	id, errUser := h.services.Authorization.CreateUser(input)
	// При спробі створення користувача з однаковим ім'ям викличеться помилка
	if errUser != nil {
		responses.NewErrorResponse(c, http.StatusConflict, "username is already used")
		return nil
	}

	// Генеруємо токен та шифруємо в ньому ID користувача
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "generate token error")
		return nil
	}

	// Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"id":    id,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

type SignInInput struct {
	Username string `json:"username" form:"username"  binding:"required"`
	Password string `json:"password" form:"password"  binding:"required"`
}

// SignIn godoc
// @Summary      Generate a new user token
// @Description  get user token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user	body     SignInInput   true  "Get user token"
// @Success      200 	{object} TokenResponse   "result is user token"
// @Failure 	 400 	{object} ErrorResponse	 "incorrect request data"
// @Failure 	 400 	{object} ErrorResponse	 "incorrect password"
// @Failure 	 404 	{object} ErrorResponse	 "user not found"
// @Failure 	 500 	{object} ErrorResponse	 "something went wrong"
// @Router       /auth/sign-in [post]
func (h *AuthHandler) SignIn(c echo.Context) error {

	// Отримуємо дані з сайту (ім'я та пароль)
	var input SignInInput
	if err := c.Bind(&input); err != nil {
		responses.NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return nil
	}

	//Перевіряємо чи існує користувач за його іменем
	user, errCheck := h.services.Authorization.GetByName(input.Username)
	if errCheck != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "check user error")
		return nil
	}
	if user.Username == "" {
		responses.NewErrorResponse(c, http.StatusNotFound, "user not found")
		return nil
	}

	// Генеруємо токен (якщо ім'я та пароль правильні)
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusConflict, "incorrect password")
		return nil
	}

	// Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *AuthHandler) GetMe(c echo.Context) error {

	// Отримуємо ID активного користувача
	userId := c.Get(middlewares.UserCtx)

	if userId == 0 {
		responses.NewErrorResponse(c, http.StatusNotFound, "user not found")
		return nil
	}

	// Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"id": userId,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

type ChangePassword struct {
	OldPassword string `json:"old_password" form:"old_password"  binding:"required"`
	NewPassword string `json:"new_password" form:"new_password"  binding:"required"`
}

func (h *AuthHandler) ChangePassword(c echo.Context) error {
	//Отримуємо власний ID з контексту
	userId := c.Get(middlewares.UserCtx).(int)

	//Отримуємо актуальний та новий паролі
	var passwords ChangePassword
	if errReq := c.Bind(&passwords); errReq != nil {
		responses.NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return nil
	}
	if len(passwords.NewPassword) < 6 {
		responses.NewErrorResponse(c, http.StatusBadRequest, "password must be at least 6 symbols")
		return nil
	}

	//Отримуємо дані активного користувача
	user, errU := h.services.Authorization.GetUserById(userId)
	if errU != nil {
		responses.NewErrorResponse(c, http.StatusNotFound, "incorrect user data")
		return nil
	}

	//Перевіряємо вірність введеного паролю
	_, errCheck := h.services.Authorization.GenerateToken(user.Username, passwords.OldPassword)
	if errCheck != nil {
		responses.NewErrorResponse(c, http.StatusBadRequest, "incorrect password")
		return nil
	}

	//Оновлюємо пароль у БД
	user.Password = passwords.NewPassword
	err := h.services.Authorization.UpdatePassword(user)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "update password error")
		return nil
	}

	//Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": "password changed",
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *AuthHandler) ChangeUsername(c echo.Context) error {
	//Отримуємо власний ID з контексту
	userId := c.Get(middlewares.UserCtx).(int)

	//Отримуємо новий нікнейм
	var username models.User
	if errReq := c.Bind(&username); errReq != nil {
		responses.NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return nil
	}

	//Отримуємо дані активного користувача
	user, errU := h.services.Authorization.GetUserById(userId)
	if errU != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "incorrect user data")
		return nil
	}

	//Перевіряємо чи існує користувач за його іменем
	check, errCheck := h.services.Authorization.GetByName(username.Username)
	if errCheck != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "check user error")
		return nil
	}
	if check.Id != 0 {
		responses.NewErrorResponse(c, http.StatusConflict, "username is used")
		return nil
	}

	//Оновлюємо нікнейм у БД
	user.Username = username.Username
	errPut := h.services.Authorization.UpdateData(user)
	if errPut != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "update username error")
		return nil
	}

	//Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": "username changed",
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *AuthHandler) ChangeIcon(c echo.Context) error {
	//Отримуємо власний ID з контексту
	userId := c.Get(middlewares.UserCtx).(int)

	//Отримуємо ім'я файлу зображення
	fileName, err := middlewares.UploadImage(c)
	if err != nil {
		return err
	}

	//Отримуємо дані активного користувача
	user, errU := h.services.Authorization.GetUserById(userId)
	if errU != nil {
		responses.NewErrorResponse(c, http.StatusBadRequest, "incorrect user data")
		return nil
	}

	//Замінюємо дані у БД
	var oldIcon = user.Icon
	user.Icon = strings.TrimPrefix(fileName, "uploads\\")
	errPut := h.services.Authorization.UpdateData(user)
	if errPut != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "update icon error")
		return nil
	}

	//Видалення застарілих файлів
	if len(oldIcon) != 0 {
		if err := os.Remove(fmt.Sprintf("uploads/%s", oldIcon)); err != nil {
			responses.NewErrorResponse(c, http.StatusInternalServerError, "delete icon error")
			return nil
		}
		if err := os.Remove(fmt.Sprintf("uploads/resize-%s", oldIcon)); err != nil {
			responses.NewErrorResponse(c, http.StatusInternalServerError, "delete icon error")
			return nil
		}
	}

	//Відгук сервера
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"message": "icon changed",
	})
	if errRes != nil {
		return errRes
	}
	return nil
}
