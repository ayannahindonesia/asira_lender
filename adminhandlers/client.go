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
		NLog("warning", "CreateClient", map[string]interface{}{"message": "validation create client error", "error": validate}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
	}

	err = client.Create()
	if err != nil {
		NLog("warning", "CreateClient", map[string]interface{}{"message": "error create client", "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat Client Config")
	}

	NAudittrail(models.Client{}, client, c.Get("user").(*jwt.Token), "client", fmt.Sprint(client.ID), "create")

	return c.JSON(http.StatusCreated, client)
}
