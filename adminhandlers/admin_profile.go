package adminhandlers

import (
	"asira_lender/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

// AdminProfile check admin profile
func AdminProfile(c echo.Context) error {
	defer c.Request().Body.Close()

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	userModel := models.User{}

	userID, _ := strconv.Atoi(claims["jti"].(string))
	err := userModel.FindbyID(userID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "unauthorized")
	}

	return c.JSON(http.StatusOK, userModel)
}

// UserFirstLoginChangePassword check if user is first login and change the password
func UserFirstLoginChangePassword(c echo.Context) error {
	defer c.Request().Body.Close()

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	userModel := models.User{}

	userID, _ := strconv.Atoi(claims["jti"].(string))
	err := userModel.FindbyID(userID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "unauthorized")
	}

	if userModel.FirstLogin {
		type Password struct {
			Pass string `json:"password"`
		}
		var pass Password
		payloadRules := govalidator.MapData{
			"password": []string{"required"},
		}

		validate := validateRequestPayload(c, payloadRules, &pass)
		if validate != nil {
			return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
		}
		userModel.FirstLoginChangePassword(pass.Pass)

		return c.JSON(http.StatusOK, "Password anda telah diganti.")
	}

	return c.JSON(http.StatusUnauthorized, "Akun anda bukan akun baru.")
}
