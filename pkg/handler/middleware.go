package handler

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	ParamId             = "id"
	ChatId              = "chatId"
	Username            = "username"
	ChatName            = "name"
)

func (h *Handler) userIdentify(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get(authorizationHeader)
		if header == "" {
			NewErrorResponse(c, http.StatusUnauthorized, "empty auth header")
			return nil
		}

		//headerParts := strings.Split(header, " ")
		//if len(headerParts) != 2 {
		//	NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		//	return nil
		//}

		userId, err := h.services.Authorization.ParseToken(header)
		//userId, err := h.services.Authorization.ParseToken(headerParts[1])

		if err != nil {
			NewErrorResponse(c, http.StatusUnauthorized, err.Error())
			return nil
		}
		c.Set(userCtx, userId)
		return next(c)
	}
}

func GetUserId(c echo.Context) (int, error) {
	id := c.Get(userCtx)
	if id == 0 {
		NewErrorResponse(c, http.StatusNotFound, "user id not found")
		return 0, errors.New("user id not found")
	}
	idInt, ok := id.(int)
	if !ok {
		NewErrorResponse(c, http.StatusBadRequest, "user id is of valid type")
		return 0, errors.New("user id is of valid type")
	}
	return idInt, nil
}

func GetParam(c echo.Context, name string) (int, error) {
	param, errReq := strconv.Atoi(c.Param(name))
	if errReq != nil {
		NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%s is not integer", name))
		return 0, errReq
	}
	return param, nil
}

func GetRequest(c echo.Context, i interface{}) error {
	if err := c.Bind(&i); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return err
	}
	return nil
}

func UploadImage(c echo.Context) (string, error) {
	fmt.Println("working...")
	defer fmt.Println("finish")

	//Обмежуємо розмір завантажуваних файлів
	c.Request().ParseMultipartForm(10 << 20)

	//Отримуємо файл зображення
	file, err := c.FormFile("image")
	if err != nil {
		fmt.Println(1)
		NewErrorResponse(c, http.StatusBadRequest, "incorrect file error")
		return "", err
	}

	resizeFile, resizeHand, err := c.Request().FormFile("image")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resizeFile.Close()

	//Відкриваємо дані файлу
	handler, err := file.Open()
	if err != nil {
		fmt.Println(2)
		NewErrorResponse(c, http.StatusConflict, "open file error")
		return "", err
	}

	defer handler.Close()

	//Створюємо порожні файли за необхідних розташуванням
	tempFile, err := os.CreateTemp("uploads", "upload-*.jpeg")
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "create file error")
		fmt.Println(3)
		return "", err
	}
	resFile, err := os.Create(fmt.Sprintf("uploads\\resize-%s", strings.TrimPrefix(tempFile.Name(), "uploads\\")))
	if err != nil {
		fmt.Println(4)
		NewErrorResponse(c, http.StatusInternalServerError, "create file error")
		return "", err
	}
	defer tempFile.Close()
	defer resFile.Close()

	//Розкодування зображення за типом
	var img image.Image
	imgFmt := strings.Split(resizeHand.Filename, ".")

	switch imgFmt[len(imgFmt)-1] {
	case "jpeg":
		img, err = jpeg.Decode(resizeFile)
		break
	case "jpg":
		img, err = jpeg.Decode(resizeFile)
		break
	case "png":
		img, err = png.Decode(resizeFile)
		break
	case "gif":
		img, err = gif.Decode(resizeFile)
		break
	default:
		fmt.Println(5)
		NewErrorResponse(c, http.StatusBadRequest, "incorrect file type error")
		return "", err
	}

	//Приведення зображень до необхідних форми й розмірів
	var crop = []int{10, 10}
	if crop != nil && len(crop) == 2 {
		analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
		topCrop, _ := analyzer.FindBestCrop(img, crop[0], crop[1])
		type SubImager interface {
			SubImage(r image.Rectangle) image.Image
		}
		img = img.(SubImager).SubImage(topCrop)
	}
	imgWidth := uint(math.Min(float64(100), float64(img.Bounds().Max.X)))
	resizedImg := resize.Resize(imgWidth, 0, img, resize.Lanczos3)

	//Збереження зображень у новосотворених файлах
	fileBytes, err := io.ReadAll(handler)
	if err != nil {
		fmt.Println(6)
		return "", err
	}
	tempFile.Write(fileBytes)

	err = jpeg.Encode(resFile, resizedImg, nil)
	return tempFile.Name(), nil
}
