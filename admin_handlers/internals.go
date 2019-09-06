package admin_handlers

import (
	"asira_lender/models"
	"net/http"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

func CreateInternal(c echo.Context) error {
	defer c.Request().Body.Close()

	internals := models.Internals{}

	payloadRules := govalidator.MapData{
		"name": []string{"required"},
		"key":  []string{"required"},
		"role": []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &internals)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err := internals.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat Client Config")
	}

	return c.JSON(http.StatusCreated, internals)
}
