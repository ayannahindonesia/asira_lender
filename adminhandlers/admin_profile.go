package adminhandlers

import (
	"asira_lender/models"
	"fmt"
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

	userID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)
	err := userModel.FindbyID(userID)
	if err != nil {
		NLog("warning", "AdminProfile", map[string]interface{}{"message": fmt.Sprintf("user id %v profile not found", userID), "error": err}, c.Get("user").(*jwt.Token), "", true)

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
	origin := models.User{}

	userID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)
	err := userModel.FindbyID(userID)
	if err != nil {
		NLog("warning", "UserFirstLoginChangePassword", map[string]interface{}{"message": fmt.Sprintf("user id %v profile not found", userID), "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusForbidden, err, "Tidak memiliki akses.")
	}
	origin = userModel

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
			NLog("warning", "UserFirstLoginChangePassword", map[string]interface{}{"message": "validation error", "error": validate}, c.Get("user").(*jwt.Token), "", false)

			return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
		}
		userModel.FirstLoginChangePassword(pass.Pass)
		NLog("info", "UserFirstLoginChangePassword", map[string]interface{}{"message": "changed password"}, c.Get("user").(*jwt.Token), "", false)

		NAudittrail(origin, userModel, token, "user", fmt.Sprint(userModel.ID), "user first login change password")

		return c.JSON(http.StatusOK, "Password anda telah diganti.")
	}
	NLog("warning", "UserFirstLoginChangePassword", map[string]interface{}{"message": "cant change password, not first login"}, c.Get("user").(*jwt.Token), "", false)

	return c.JSON(http.StatusUnauthorized, "Akun anda bukan akun baru.")
}
