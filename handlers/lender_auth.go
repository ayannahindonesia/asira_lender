package handlers

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
	LenderLoginCreds struct {
		Key      string `json:"key"`
		Password string `json:"password"`
	}
)

// lender login, lender can choose either login with email / phone
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
		return returnInvalidResponse(http.StatusBadRequest, validate, "invalid login")
	}

	// check if theres record
	validKey = asira.App.DB.Where("username = ?", credentials.Key).Find(&lender).RecordNotFound()

	if !validKey { // check the password
		err = bcrypt.CompareHashAndPassword([]byte(lender.Password), []byte(credentials.Password))
		if err != nil {
			return returnInvalidResponse(http.StatusUnauthorized, err, "invalid login")
		}

		token, err = createJwtToken(strconv.FormatUint(lender.ID, 10), "lender")
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
