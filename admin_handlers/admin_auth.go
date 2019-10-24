package admin_handlers

import (
	"asira_lender/asira"
	"asira_lender/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
	"golang.org/x/crypto/bcrypt"
)

type (
	AdminLoginCreds struct {
		Key      string `json:"key"`
		Password string `json:"password"`
	}
)

// admin login
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
		return returnInvalidResponse(http.StatusBadRequest, validate, "invalid login")
	}

	// check if theres record
	validKey = asira.App.DB.Where("username = ?", credentials.Key).Find(&user).RecordNotFound()

	if !validKey { // check the password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
		if err != nil {
			return returnInvalidResponse(http.StatusUnauthorized, err, "invalid login")
		}

		if user.Status == "inactive" {
			return returnInvalidResponse(http.StatusUnauthorized, err, "invalid login")
		}

		tokenrole := "admin"
		token, err = createJwtToken(strconv.FormatUint(user.ID, 10), tokenrole, "1")
		if err != nil {
			return returnInvalidResponse(http.StatusInternalServerError, err, "error creating token")
		}
	} else {
		return returnInvalidResponse(http.StatusUnauthorized, "", "invalid login")
	}

	jwtConf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", asira.App.ENV))
	expiration := time.Duration(jwtConf["duration"].(int)) * time.Minute
	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":      token,
		"expires_in": expiration.Seconds(),
	})
}
