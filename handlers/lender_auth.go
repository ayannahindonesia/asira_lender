package handlers

import (
	"asira_lender/adminhandlers"
	"asira_lender/asira"
	"asira_lender/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
	"golang.org/x/crypto/bcrypt"
)

type (
	// LenderLoginCreds type
	LenderLoginCreds struct {
		Key      string `json:"key"`
		Password string `json:"password"`
	}
)

// LenderLogin lender can choose either login with email / phone
func LenderLogin(c echo.Context) error {
	defer c.Request().Body.Close()

	var (
		credentials LenderLoginCreds
		lender      models.User
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
		adminhandlers.NLog("warning", "LenderLogin", map[string]interface{}{"message": "validation error", "error": validate}, c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusBadRequest, validate, "Login tidak valid")
	}

	// check if theres record
	validKey = asira.App.DB.
		Where("username = ?", credentials.Key).
		Where("status = ?", "active").
		Find(&lender).RecordNotFound()

	if !validKey { // check the password
		err = bcrypt.CompareHashAndPassword([]byte(lender.Password), []byte(credentials.Password))
		if err != nil {
			adminhandlers.NLog("warning", "LenderLogin", map[string]interface{}{"message": fmt.Sprintf("password error on user %v", credentials.Key), "error": err}, c.Get("user").(*jwt.Token), "", true)

			return returnInvalidResponse(http.StatusUnauthorized, err, "Login tidak valid")
		}

		token, err = createJwtToken(strconv.FormatUint(lender.ID, 10), "users")
		if err != nil {
			adminhandlers.NLog("warning", "LenderLogin", map[string]interface{}{"message": "error generating token", "error": err}, c.Get("user").(*jwt.Token), "", true)

			return returnInvalidResponse(http.StatusInternalServerError, err, "Terjadi kesalahan")
		}
	} else {
		adminhandlers.NLog("warning", "LenderLogin", map[string]interface{}{"message": fmt.Sprintf("user not found %v", credentials.Key)}, c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusUnauthorized, "username not found", "Login tidak valid")
	}

	jwtConf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", asira.App.ENV))
	expiration := time.Duration(jwtConf["duration"].(int)) * time.Minute

	adminhandlers.NLog("info", "LenderLogin", map[string]interface{}{"message": fmt.Sprintf("%v login", credentials.Key)}, c.Get("user").(*jwt.Token), "", true)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":      token,
		"expires_in": expiration.Seconds(),
	})
}
