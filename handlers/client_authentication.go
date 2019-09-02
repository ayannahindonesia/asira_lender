package handlers

import (
	"asira_lender/asira"
	"asira_lender/models"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

func ClientLogin(c echo.Context) error {
	defer c.Request().Body.Close()

	data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Basic "))
	if err != nil {
		return returnInvalidResponse(http.StatusUnauthorized, "", "Invalid Credentials")
	}

	auth := strings.Split(string(data), ":")
	if len(auth) < 2 {
		return returnInvalidResponse(http.StatusUnauthorized, "", "Invalid Credentials")
	}
	type Login struct {
		Key    string `json:"key"`
		Secret string `json:"secret"`
	}

	internals := models.Internals{}
	err = internals.FilterSearchSingle(&Login{
		Key:    auth[0],
		Secret: auth[1],
	})
	if err != nil {
		return returnInvalidResponse(http.StatusUnauthorized, "", "Invalid Credentials")
	}

	token, err := createJwtToken(strconv.FormatUint(internals.ID, 10), internals.Role)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, "", fmt.Sprint(err))
	}

	jwtConf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", asira.App.ENV))
	expiration := time.Duration(jwtConf["duration"].(int)) * time.Minute

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":      token,
		"expires_in": expiration.Seconds(),
	})
}
