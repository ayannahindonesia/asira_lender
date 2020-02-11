package adminhandlers

import (
	"asira_lender/asira"
	"asira_lender/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ayannahindonesia/northstar/lib/northstarlib"
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

	userID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)
	err := userModel.FindbyID(userID)
	if err != nil {
		asira.App.Northstar.SubmitKafkaLog(northstarlib.Log{
			Level:    "error",
			Tag:      "AdminProfile",
			Messages: fmt.Sprintf("user id %v profile. error : %v", userID, err),
			UID:      fmt.Sprint(userID),
		}, "log")
		return returnInvalidResponse(http.StatusForbidden, err, "Tidak memiliki akses.")
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

	userID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)
	err := userModel.FindbyID(userID)
	if err != nil {
		asira.App.Northstar.SubmitKafkaLog(northstarlib.Log{
			Level:    "error",
			Tag:      "UserFirstLoginChangePassword",
			Messages: fmt.Sprintf("user id %v profile. error : %v", userID, err),
			UID:      fmt.Sprint(userID),
		}, "log")

		return returnInvalidResponse(http.StatusForbidden, err, "Tidak memiliki akses.")
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
			asira.App.Northstar.SubmitKafkaLog(northstarlib.Log{
				Level:    "error",
				Tag:      "UserFirstLoginChangePassword",
				Messages: fmt.Sprintf("validation error : %v", validate),
				UID:      fmt.Sprint(userID),
				Username: userModel.Username,
			}, "log")

			return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
		}
		userModel.FirstLoginChangePassword(pass.Pass)
		asira.App.Northstar.SubmitKafkaLog(northstarlib.Log{
			Level:    "info",
			Tag:      "UserFirstLoginChangePassword",
			Messages: fmt.Sprint("changed password"),
			UID:      fmt.Sprint(userID),
			Username: userModel.Username,
		}, "log")

		return c.JSON(http.StatusOK, "Password anda telah diganti.")
	}
	asira.App.Northstar.SubmitKafkaLog(northstarlib.Log{
		Level:    "error",
		Tag:      "UserFirstLoginChangePassword",
		Messages: fmt.Sprint("not first login"),
		UID:      fmt.Sprint(userID),
		Username: userModel.Username,
	}, "log")

	return c.JSON(http.StatusUnauthorized, "Akun anda bukan akun baru.")
}
