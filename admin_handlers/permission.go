package admin_handlers

import (
	"asira_lender/asira"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func GetAllPermission(c echo.Context) error {
	defer c.Request().Body.Close()

	permissions := asira.App.Permission.GetStringMap(fmt.Sprintf("%s", "permissions"))

	return c.JSON(http.StatusOK, permissions)
}
