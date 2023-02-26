package middlewares

import (
	"cmd/pkg/handler/responses"
	"cmd/pkg/service"
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

type MiddlewareHandler struct {
	services *service.Service
}

func NewMiddlewareHandler(services *service.Service) *MiddlewareHandler {
	return &MiddlewareHandler{services: services}
}

const (
	authorizationHeader = "Authorization"
	UserCtx             = "userId"
	ParamId             = "id"
	ChatId              = "chatId"
	Username            = "username"
	ChatName            = "name"
	imageSize           = 100
	imgWidth            = 10
	imgHeight           = 10
)

func (h *MiddlewareHandler) UserIdentify(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get(authorizationHeader)
		if header == "" {
			responses.NewErrorResponse(c, http.StatusUnauthorized, "empty auth header")
			return nil
		}

		userId, err := h.services.Authorization.ParseToken(header)
		if err != nil {
			responses.NewErrorResponse(c, http.StatusUnauthorized, "create token error")
			return nil
		}
		c.Set(UserCtx, userId)
		return next(c)
	}
}

func GetUserId(c echo.Context) (int, error) {
	id := c.Get(UserCtx)
	if id == 0 {
		responses.NewErrorResponse(c, http.StatusNotFound, "user id not found")
		return 0, errors.New("user id not found")
	}
	idInt, ok := id.(int)
	if !ok {
		responses.NewErrorResponse(c, http.StatusBadRequest, "user id is of valid type")
		return 0, errors.New("user id is of valid type")
	}
	return idInt, nil
}

func GetParam(c echo.Context, name string) (int, error) {
	return strconv.Atoi(c.Param(name))
}

func UploadImage(c echo.Context) (string, error) {

	//Обмежуємо розмір завантажуваних файлів
	c.Request().ParseMultipartForm(10 << 20)

	//Отримуємо файл зображення
	file, err := c.FormFile("image")
	if err != nil {
		responses.NewErrorResponse(c, http.StatusBadRequest, "incorrect file error")
		return "", err
	}

	resizeFile, resizeHand, err := c.Request().FormFile("image")
	if err != nil {
		return "", err
	}
	defer resizeFile.Close()

	//Відкриваємо дані файлу
	handler, err := file.Open()
	if err != nil {
		responses.NewErrorResponse(c, http.StatusConflict, "open file error")
		return "", err
	}

	defer handler.Close()

	//Створюємо порожні файли за необхідних розташуванням
	tempFile, err := os.CreateTemp("uploads", "upload-*.jpeg")
	if err != nil {
		errDel := os.Remove(tempFile.Name())
		if errDel != nil {
			return "", errDel
		}
		defer tempFile.Close()
		responses.NewErrorResponse(c, http.StatusInternalServerError, "create file error")
		return "", err
	}

	resFile, err := os.Create(fmt.Sprintf("uploads/resize-%s", strings.TrimPrefix(tempFile.Name(), "uploads/")))
	if err != nil {
		errDel := os.Remove(tempFile.Name())
		if errDel != nil {
			return "", errDel
		}
		defer tempFile.Close()
		defer resFile.Close()
		responses.NewErrorResponse(c, http.StatusInternalServerError, "create file error")
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
		responses.NewErrorResponse(c, http.StatusBadRequest, "incorrect file type error")
		return "", err
	}

	//Приведення зображень до необхідних форми й розмірів
	var crop = []int{imgWidth, imgHeight}
	if crop != nil && len(crop) == 2 {
		analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
		topCrop, _ := analyzer.FindBestCrop(img, crop[0], crop[1])
		type SubImager interface {
			SubImage(r image.Rectangle) image.Image
		}
		img = img.(SubImager).SubImage(topCrop)
	}
	imgWidth := uint(math.Min(float64(imageSize), float64(img.Bounds().Max.X)))
	resizedImg := resize.Resize(imgWidth, 0, img, resize.Lanczos3)

	//Збереження зображень у новосотворених файлах
	fileBytes, err := io.ReadAll(handler)
	if err != nil {
		return "", err
	}

	tempFile.Write(fileBytes)

	err = jpeg.Encode(resFile, resizedImg, nil)
	return strings.TrimPrefix(tempFile.Name(), "uploads/"), nil
}
