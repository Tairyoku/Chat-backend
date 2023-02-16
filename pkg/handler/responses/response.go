package responses

import "github.com/labstack/echo/v4"

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(c echo.Context, statusCode int, message string) {
	errRes := c.JSON(statusCode, ErrorResponse{Message: message})
	if errRes != nil {
		return
	}
}
