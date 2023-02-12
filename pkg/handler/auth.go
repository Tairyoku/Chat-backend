package handler

import (
	"cmd/pkg/repository/models"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strings"
)

// SignUp godoc
// @Summary      Create a new user
// @Description  add new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user	body     UserResponse   true  "Add user"
// @Success      200 	{object} IdResponse		 "result is id of user"
// @Failure 	 400 	{object} ErrorResponse	 "incorrect request data"
// @Failure 	 404 	{object} ErrorResponse	 "user id not found"
// @Failure 	 500 	{object} ErrorResponse	 "something went wrong"
// @Router       /auth/sign-up [post]
func (h *Handler) SignUp(c echo.Context) error {

	// Отримуємо дані з сайту (ім'я та пароль)
	var input models.User
	if errReq := c.Bind(&input); errReq != nil {
		NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return nil
	}

	// Перевіряємо отримані дані
	{
		//username is not empty
		if len(input.Username) == 0 {
			NewErrorResponse(c, http.StatusBadRequest, "You must enter a username")
			return nil
		}

		// password length
		if len(input.Password) < 6 {
			NewErrorResponse(c, http.StatusBadRequest, "Password must be at least 6 symbols")
			return nil
		}
	}

	// Створюємо нового користувача
	id, errUser := h.services.Authorization.CreateUser(input)
	if errUser != nil {
		NewErrorResponse(c, http.StatusBadRequest, "username is already used")
		return nil
	}

	// Генеруємо токен та шифруємо в ньому ID користувача
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "generate token error")
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
func (h *Handler) SignIn(c echo.Context) error {

	// Отримуємо дані з сайту (ім'я та пароль)
	var input SignInInput
	if err := c.Bind(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return nil
	}

	//Перевіряємо чи існує користувач за його іменем
	user, errCheck := h.services.Authorization.GetByName(input.Username)
	if errCheck != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return nil
	}
	if user.Username == "" {
		NewErrorResponse(c, http.StatusNotFound, "user not found")
		return nil
	}

	// Генеруємо токен (якщо ім'я та пароль правильні)
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "incorrect password")
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

func (h *Handler) GetMe(c echo.Context) error {

	// Отримуємо ID активного користувача
	userId := c.Get(userCtx)

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

func (h *Handler) ChangePassword(c echo.Context) error {
	//Отримуємо власний ID з контексту
	userId := c.Get(userCtx).(int)

	//Отримуємо актуальний та новий паролі
	var passwords ChangePassword
	if errReq := c.Bind(&passwords); errReq != nil {
		NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return nil
	}

	//Отримуємо дані активного користувача
	user, errU := h.services.Authorization.GetUserById(userId)
	if errU != nil {
		NewErrorResponse(c, http.StatusBadRequest, "incorrect user data")
		return nil
	}

	//Перевіряємо вірність введеного паролю
	_, errCheck := h.services.Authorization.GenerateToken(user.Username, passwords.OldPassword)
	if errCheck != nil {
		NewErrorResponse(c, http.StatusBadRequest, "wrong password error")
		return nil
	}

	//Оновлюємо пароль у БД
	user.Password = passwords.NewPassword
	err := h.services.Authorization.UpdatePassword(user)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
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

func (h *Handler) ChangeUsername(c echo.Context) error {
	//Отримуємо власний ID з контексту
	userId := c.Get(userCtx).(int)

	//Отримуємо новий нікнейм
	var username models.User
	if errReq := c.Bind(&username); errReq != nil {
		NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return nil
	}

	//Отримуємо дані активного користувача
	user, errU := h.services.Authorization.GetUserById(userId)
	if errU != nil {
		NewErrorResponse(c, http.StatusBadRequest, "incorrect user data")
		return nil
	}

	//Перевіряємо чи існує користувач за його іменем
	check, errCheck := h.services.Authorization.GetByName(username.Username)
	if errCheck != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return nil
	}
	if check.Id != 0 {
		NewErrorResponse(c, http.StatusNotFound, "username is used")
		return nil
	}

	//Оновлюємо нікнейм у БД
	user.Username = username.Username
	errPut := h.services.Authorization.UpdateData(user)
	if errPut != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
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

func (h *Handler) ChangeIcon(c echo.Context) error {
	//Отримуємо власний ID з контексту
	userId := c.Get(userCtx).(int)
	//
	////Обмежуємо розмір завантажуваних файлів
	//c.Request().ParseMultipartForm(10 << 20)
	//
	////Отримуємо файл зображення
	//file, err := c.FormFile("image")
	//if err != nil {
	//	fmt.Println(1)
	//	NewErrorResponse(c, http.StatusBadRequest, "incorrect file error")
	//	return err
	//}
	//
	//resizeFile, resizeHand, err := c.Request().FormFile("image")
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}
	//defer resizeFile.Close()
	//
	////Відкриваємо дані файлу
	//handler, err := file.Open()
	//if err != nil {
	//	fmt.Println(2)
	//	NewErrorResponse(c, http.StatusConflict, "open file error")
	//	return err
	//}
	//
	//defer handler.Close()
	//
	////Створюємо порожні файли за необхідних розташуванням
	//tempFile, err := os.CreateTemp("uploads", "upload-*.jpeg")
	//if err != nil {
	//	NewErrorResponse(c, http.StatusInternalServerError, "create file error")
	//	fmt.Println(3)
	//	return err
	//}
	//resFile, err := os.Create(fmt.Sprintf("uploads\\resize-%s", strings.TrimPrefix(tempFile.Name(), "uploads\\")))
	//if err != nil {
	//	fmt.Println(4)
	//	NewErrorResponse(c, http.StatusInternalServerError, "create file error")
	//	return err
	//}
	//defer tempFile.Close()
	//defer resFile.Close()
	//
	////Розкодування зображення за типом
	//var img image.Image
	//imgFmt := strings.Split(resizeHand.Filename, ".")
	//
	//switch imgFmt[len(imgFmt)-1] {
	//case "jpeg":
	//	img, err = jpeg.Decode(resizeFile)
	//	break
	//case "jpg":
	//	img, err = jpeg.Decode(resizeFile)
	//	break
	//case "png":
	//	img, err = png.Decode(resizeFile)
	//	break
	//case "gif":
	//	img, err = gif.Decode(resizeFile)
	//	break
	//default:
	//	fmt.Println(5)
	//	NewErrorResponse(c, http.StatusBadRequest, "incorrect file type error")
	//	return nil
	//}
	//
	////Приведення зображень до необхідних форми й розмірів
	//var crop = []int{10, 10}
	//if crop != nil && len(crop) == 2 {
	//	analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
	//	topCrop, _ := analyzer.FindBestCrop(img, crop[0], crop[1])
	//	type SubImager interface {
	//		SubImage(r image.Rectangle) image.Image
	//	}
	//	img = img.(SubImager).SubImage(topCrop)
	//}
	//imgWidth := uint(math.Min(float64(100), float64(img.Bounds().Max.X)))
	//resizedImg := resize.Resize(imgWidth, 0, img, resize.Lanczos3)
	//
	////Збереження зображень у новосотворених файлах
	//fileBytes, err := io.ReadAll(handler)
	//if err != nil {
	//	fmt.Println(6)
	//	return err
	//}
	//tempFile.Write(fileBytes)
	//
	//err = jpeg.Encode(resFile, resizedImg, nil)
	//Отримуємо ім'я файлу зображення
	fileName, err := UploadImage(c)
	if err != nil {
		return err
	}

	//Отримуємо дані активного користувача
	user, errU := h.services.Authorization.GetUserById(userId)
	if errU != nil {
		fmt.Println(7)

		NewErrorResponse(c, http.StatusBadRequest, "incorrect user data")
		return nil
	}

	//Замінюємо дані у БД
	var oldIcon = user.Icon
	user.Icon = strings.TrimPrefix(fileName, "uploads\\")
	errPut := h.services.Authorization.UpdateData(user)
	if errPut != nil {
		fmt.Println(8)

		NewErrorResponse(c, http.StatusInternalServerError, "update icon error")
		return nil
	}

	//Видалення застарілих файлів
	if len(oldIcon) != 0 {
		if err := os.Remove(fmt.Sprintf("uploads/%s", oldIcon)); err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "delete icon error")
			return nil
		}
		if err := os.Remove(fmt.Sprintf("uploads/resize-%s", oldIcon)); err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "delete icon error")
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
