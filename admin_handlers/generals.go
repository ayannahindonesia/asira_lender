package admin_handlers

import (
	"asira_lender/asira"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

type (
	JWTclaims struct {
		Username string `json:"username"`
		Role     string `json:"role"`
		RoleID   string `json:"role_id"`
		jwt.StandardClaims
	}
)

// general function to validate all kind of api request payload / body
func validateRequestPayload(c echo.Context, rules govalidator.MapData, data interface{}) (i interface{}) {
	opts := govalidator.Options{
		Request: c.Request(),
		Data:    data,
		Rules:   rules,
	}

	v := govalidator.New(opts)

	mappedError := v.ValidateJSON()

	if len(mappedError) > 0 {
		i = mappedError
	}

	return i
}

// general function to validate all kind of api request url query
func validateRequestQuery(c echo.Context, rules govalidator.MapData) (i interface{}) {
	opts := govalidator.Options{
		Request: c.Request(),
		Rules:   rules,
	}

	v := govalidator.New(opts)

	mappedError := v.Validate()

	if len(mappedError) > 0 {
		i = mappedError
	}

	return i
}

func returnInvalidResponse(httpcode int, details interface{}, message string) error {
	responseBody := map[string]interface{}{
		"message": message,
		"details": details,
	}

	return echo.NewHTTPError(httpcode, responseBody)
}

// self explanation
func createJwtToken(id string, role string, roleID string) (string, error) {
	jwtConf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", asira.App.ENV))

	claim := JWTclaims{
		id,
		role,
		roleID,
		jwt.StandardClaims{
			Id:        id,
			ExpiresAt: time.Now().Add(time.Duration(jwtConf["duration"].(int)) * time.Minute).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)
	token, err := rawToken.SignedString([]byte(jwtConf["jwt_secret"].(string)))
	if err != nil {
		return "", err
	}

	return token, nil
}

func customSplit(str string, separator string) []string {
	split := strings.Split(str, separator)
	log.Println(str)
	log.Println(split)
	if len(split) == 1 {
		if split[0] == "" {
			split = nil
		}
	}

	return split
}
