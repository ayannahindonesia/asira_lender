package adminhandlers

import (
	"asira_lender/models"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

// CreateClient func
func CreateClient(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_create_client")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	client := models.Client{}

	payloadRules := govalidator.MapData{
		"name":   []string{"required"},
		"key":    []string{"required"},
		"secret": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &client)
	if validate != nil {
		NLog("warning", "CreateClient", fmt.Sprintf("validation create client error : %v", validate), c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
	}

	err = client.Create()
	if err != nil {
		NLog("warning", "CreateClient", fmt.Sprintf("error create client : %v", err), c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat Client Config")
	}

	return c.JSON(http.StatusCreated, client)
}
