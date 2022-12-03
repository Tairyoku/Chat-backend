package handler

import (
	"cmd/pkg/repository/models"
	"github.com/labstack/echo/v4"
	"net/http"
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
	var input models.User

	if errReq := c.Bind(&input); errReq != nil {
		NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return nil
	}
	{
		//username is not empty
		if len(input.Username) == 0 {
			NewErrorResponse(c, http.StatusBadRequest, "You must enter a username")
			return nil
		}

		//name is not empty
		if len(input.Name) == 0 {
			NewErrorResponse(c, http.StatusBadRequest, "You must enter a name")
			return nil
		}

		// password length
		if len(input.Password) < 6 {
			NewErrorResponse(c, http.StatusBadRequest, "Password must be at least 6 symbols")
			return nil
		}
	}
	id, errCreate := h.services.Authorization.CreateUser(input)
	if errCreate != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "something went wrong")
		return nil
	}
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"id": id})
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
	var input SignInInput
	if err := c.Bind(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return nil
	}
	errCheck := h.services.Authorization.CheckUsername(input.Username)
	if errCheck != nil {
		NewErrorResponse(c, http.StatusNotFound, "user not found")
		return nil
	}
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "incorrect password")
		return nil
	}
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"token": token})
	if errRes != nil {
		return errRes
	}
	return nil
}
