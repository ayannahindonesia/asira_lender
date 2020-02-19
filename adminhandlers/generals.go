package adminhandlers

import (
	"asira_lender/asira"
	"asira_lender/models"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/ayannahindonesia/northstar/lib/northstarlib"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

// JWTclaims type
type JWTclaims struct {
	Username    string   `json:"username"`
	Group       string   `json:"group"`
	Permissions []string `json:"permissions"`
	jwt.StandardClaims
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

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
func createJwtToken(id string, group string) (string, error) {
	jwtConf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", asira.App.ENV))

	type PermModel struct {
		Permission string `json:"permissions" gorm:"column:permissions"`
	}
	var permissions []string
	var permModel []PermModel
	var db = asira.App.DB
	switch group {
	case "users":
		err := db.Table("roles").
			Select("DISTINCT TRIM(UNNEST(roles.permissions)) as permissions").
			Joins("INNER JOIN users u ON roles.id IN (SELECT UNNEST(u.roles))").
			Where("roles.status = ?", "active").
			Where("u.id = ?", id).Scan(&permModel).Error
		if err != nil {
			return "", err
		}
		for _, v := range permModel {
			permissions = append(permissions, v.Permission)
		}
		break
	}

	claim := JWTclaims{
		id,
		group,
		permissions,
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

// RandString random string alphanumeric. parameter length
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func customSplit(str string, separator string) []string {
	split := strings.Split(str, separator)
	if len(split) == 1 {
		if split[0] == "" {
			split = nil
		}
	}

	return split
}

func validatePermission(c echo.Context, permission string) error {
	user := c.Get("user")
	token := user.(*jwt.Token)

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if claimPermissions, ok := claims["permissions"]; ok {
			s := strings.Split(strings.Trim(fmt.Sprintf("%v", claimPermissions), "[]"), " ")
			for _, v := range s {
				if strings.ToLower(v) == strings.ToLower(permission) || strings.ToLower(v) == "all" {
					return nil
				}
			}
		}

		NLog("warning", "validatePermission", fmt.Sprintf("user dont have permission %v", permission), c.Get("user").(*jwt.Token), "", false)

		return fmt.Errorf("Tidak memiliki hak akses")
	}

	NLog("warning", "validatePermission", "invalid token. error claims", c.Get("user").(*jwt.Token), "", true)

	return fmt.Errorf("Tidak memiliki hak akses")
}

// NLog send log to northstar service
func NLog(level string, tag string, message string, jwttoken *jwt.Token, note string, nouser bool) {
	var (
		uid      string
		username string
		err      error
	)

	if !nouser {
		jti, _ := strconv.ParseUint(jwttoken.Claims.(jwt.MapClaims)["jti"].(string), 10, 64)
		user := models.User{}
		err = user.FindbyID(jti)
		if err == nil {
			uid = fmt.Sprint(user.ID)
			username = user.Username
		}
	}

	err = asira.App.Northstar.SubmitKafkaLog(northstarlib.Log{
		Level:    level,
		Tag:      tag,
		Messages: message,
		UID:      uid,
		Username: username,
		Note:     note,
	}, "log")

	if err != nil {
		log.Printf("error northstar log : %v", err)
	}
}

// NAudittrail send audit trail log to northstar service
func NAudittrail(ori interface{}, new interface{}, jwttoken *jwt.Token, entity string, entityID string, action string) {
	var (
		uid      string
		username string
		err      error
	)

	jti, _ := strconv.ParseUint(jwttoken.Claims.(jwt.MapClaims)["jti"].(string), 10, 64)
	user := models.User{}
	err = user.FindbyID(jti)
	if err == nil {
		uid = fmt.Sprint(user.ID)
		username = user.Username
	} else {
		uid = "0"
		username = "not found"
	}

	err = asira.App.Northstar.SubmitKafkaLog(northstarlib.Audittrail{
		Client:   asira.App.Northstar.Secret,
		UserID:   uid,
		Username: username,
		Roles:    fmt.Sprint(user.Roles),
		Entity:   entity,
		EntityID: entityID,
		Action:   action,
		Original: fmt.Sprintf("%+v\n", ori),
		New:      fmt.Sprintf("%+v\n", new),
	}, "audittrail")

	if err != nil {
		log.Printf("error northstar log : %v", err)
	}
}
