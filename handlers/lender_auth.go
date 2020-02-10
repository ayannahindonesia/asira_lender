package handlers

import (
	"asira_lender/asira"
	"asira_lender/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ayannahindonesia/northstar/lib/northstarlib"

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
		asira.App.Northstar.SubmitKafkaLog(northstarlib.Log{Level: "error", Tag: "LenderLogin", Messages: fmt.Sprintf("error validation : %v", validate)}, "log")
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
			asira.App.Northstar.SubmitKafkaLog(northstarlib.Log{Level: "error", Tag: "LenderLogin", Messages: fmt.Sprintf("password error : %v username : %v", err, credentials.Key)}, "log")
			return returnInvalidResponse(http.StatusUnauthorized, err, "Login tidak valid")
		}

		token, err = createJwtToken(strconv.FormatUint(lender.ID, 10), "users")
		if err != nil {
			asira.App.Northstar.SubmitKafkaLog(northstarlib.Log{Level: "error", Tag: "LenderLogin", Messages: fmt.Sprintf("error generating token : %v", err)}, "log")
			return returnInvalidResponse(http.StatusInternalServerError, err, "Terjadi kesalahan")
		}
	} else {
		asira.App.Northstar.SubmitKafkaLog(northstarlib.Log{Level: "error", Tag: "LenderLogin", Messages: fmt.Sprintf("not found username : %v", credentials.Key)}, "log")
		return returnInvalidResponse(http.StatusUnauthorized, "username not found", "Login tidak valid")
	}

	jwtConf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", asira.App.ENV))
	expiration := time.Duration(jwtConf["duration"].(int)) * time.Minute

	asira.App.Northstar.SubmitKafkaLog(northstarlib.Log{Level: "event", Tag: "LenderLogin", Messages: fmt.Sprintf("%v login", credentials.Key)}, "log")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":      token,
		"expires_in": expiration.Seconds(),
	})
}
