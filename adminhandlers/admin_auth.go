package adminhandlers

import (
	"asira_lender/asira"
	"asira_lender/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
	"golang.org/x/crypto/bcrypt"
)

type (
	// AdminLoginCreds admin credentials container
	AdminLoginCreds struct {
		Key      string `json:"key"`
		Password string `json:"password"`
	}
)

// AdminLogin func
func AdminLogin(c echo.Context) error {
	defer c.Request().Body.Close()

	var (
		credentials AdminLoginCreds
		user        models.User
		validKey    bool
		token       string
		err         error
	)

	rules := govalidator.MapData{
		"key":      []string{"required"},
		"password": []string{"required"},
	}

	validate := validateRequestPayload(c, rules, &credentials)
	if validate != nil {
		NLog("warning", "AdminLogin", map[string]interface{}{"message": "error validation", "detail": validate}, c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusBadRequest, validate, "Login tidak valid")
	}

	// check if theres record
	validKey = asira.App.DB.Where("username = ?", credentials.Key).Find(&user).RecordNotFound()

	if !validKey { // check the password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
		if err != nil {
			NLog("warning", "AdminLogin", map[string]interface{}{"message": "password error", "detail": err}, c.Get("user").(*jwt.Token), "", true)

			return returnInvalidResponse(http.StatusUnauthorized, err, "Login tidak valid")
		}

		if user.Status == "inactive" {
			NLog("warning", "AdminLogin", map[string]interface{}{"message": "inactive username", "detail": user}, c.Get("user").(*jwt.Token), "", true)

			return returnInvalidResponse(http.StatusUnauthorized, err, "Login tidak valid")
		}

		token, err = createJwtToken(strconv.FormatUint(user.ID, 10), "users")
		if err != nil {
			NLog("error", "AdminLogin", map[string]interface{}{"message": "error generating token", "detail": err}, c.Get("user").(*jwt.Token), "", true)

			return returnInvalidResponse(http.StatusInternalServerError, err, "error creating token")
		}
	} else {
		NLog("error", "AdminLogin", map[string]interface{}{"message": "error generating token", "detail": err}, c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusUnauthorized, "", "Login tidak valid")
	}

	jwtConf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", asira.App.ENV))
	expiration := time.Duration(jwtConf["duration"].(int)) * time.Minute

	NLog("info", "AdminLogin", map[string]interface{}{"message": fmt.Sprintf("%v login", user.Username)}, c.Get("user").(*jwt.Token), "", true)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":      token,
		"expires_in": expiration.Seconds(),
	})
}
