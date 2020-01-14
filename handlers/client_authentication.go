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

// ClientLogin func
func ClientLogin(c echo.Context) error {
	defer c.Request().Body.Close()

	data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Basic "))
	if err != nil {
		return returnInvalidResponse(http.StatusUnauthorized, "", "Login tidak valid")
	}

	auth := strings.Split(string(data), ":")
	if len(auth) < 2 {
		return returnInvalidResponse(http.StatusUnauthorized, "", "Login tidak valid")
	}
	type Login struct {
		Key    string `json:"key"`
		Secret string `json:"secret"`
	}

	client := models.Client{}
	err = client.FilterSearchSingle(&Login{
		Key:    auth[0],
		Secret: auth[1],
	})
	if err != nil {
		return returnInvalidResponse(http.StatusUnauthorized, "", "Login tidak valid")
	}

	token, err := createJwtToken(strconv.FormatUint(client.ID, 10), "client")
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Terjadi kesalahan")
	}

	jwtConf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", asira.App.ENV))
	expiration := time.Duration(jwtConf["duration"].(int)) * time.Minute

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":      token,
		"expires_in": expiration.Seconds(),
	})
}
