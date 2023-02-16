package handler

//
//import (
//	"encoding/json"
//	"github.com/labstack/echo/v4"
//	"io"
//	"net/http"
//	"test/pkg/repository/models"
//)
//
//type GoogleAuth struct {
//	Name   string `json:"name"`
//	Email  string `json:"email"`
//	UserId string `json:"id"`
//}
//
//func (h *Handler) GoogleLogin(c echo.Context) error {
//	googleConfig := SetupConfig()
//
//	url := googleConfig.AuthCodeURL("randomState")
//	err := c.Redirect(http.StatusSeeOther, url)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, "Error with google auth")
//		return err
//	}
//	return nil
//}
//
//func (h *Handler) GoogleCallback(c echo.Context) error {
//	state := c.Request().FormValue("state")
//	if state != "randomState" {
//		NewErrorResponse(c, http.StatusBadRequest, "wrong state")
//		return nil
//	}
//
//	code := c.Request().FormValue("code")
//
//	googleConfig := SetupConfig()
//
//	token, err := googleConfig.Exchange(c.Request().Context(), code)
//	if err != nil {
//		NewErrorResponse(c, http.StatusBadRequest, "wrong token")
//		return err
//	}
//
//	res, errRes := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
//	if errRes != nil {
//		NewErrorResponse(c, http.StatusBadRequest, "wrong response")
//		return errRes
//	}
//
//	userData, errData := io.ReadAll(res.Body)
//	if errData != nil {
//		NewErrorResponse(c, http.StatusBadRequest, "wrong body data")
//		return errData
//	}
//
//	var googleAuth GoogleAuth
//
//	errJSON := json.Unmarshal(userData, &googleAuth)
//	if errJSON != nil {
//		return errJSON
//	}
//
//	var input = models.User{
//		Id:       0,
//		Name:     googleAuth.Name,
//		Username: googleAuth.Email,
//		Password: googleAuth.UserId,
//	}
//
//	//check username: if it is not used (it will get an error), try to create new user
//	errUser := h.services.Authorization.CheckUser(input.Username)
//	if errUser != nil {
//		errSignUp := h.GoogleSignUp(c, input)
//		if errSignUp != nil {
//			return errSignUp
//		}
//		return nil
//	} else {
//		errSignIn := h.GoogleSignIn(c, input)
//		if errSignIn != nil {
//			return errSignIn
//		}
//	}
//	return nil
//}
