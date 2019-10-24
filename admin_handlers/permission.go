package admin_handlers

import (
	"asira_lender/asira"
	"asira_lender/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

func GetAllPermission(c echo.Context) error {
	defer c.Request().Body.Close()

	Permission := models.Permissions{}
	// pagination parameters
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	name := c.QueryParam("name")
	id := customSplit(c.QueryParam("id"), ",")
	role_id := customSplit(c.QueryParam("role_id"), ",")

	type Filter struct {
		ID         []string `json:"id"`
		RoleID     []string `json:"role_id"`
		Permission string   `json:"permissions" condition:"LIKE"`
	}

	result, err := Permission.FilterSearch(offset, limit, orderby, sort, &Filter{
		ID:         id,
		RoleID:     role_id,
		Permission: name,
	})

	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Permission tidak Ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

func PermissionGetDetails(c echo.Context) error {
	defer c.Request().Body.Close()

	Permission := models.Permissions{}

	PermissionID, _ := strconv.Atoi(c.Param("permission_id"))
	err := Permission.FindbyID(PermissionID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Permission tidak ditemukan")
	}

	return c.JSON(http.StatusOK, Permission)
}

func AddPermission(c echo.Context) error {
	defer c.Request().Body.Close()

	type (
		validatePermissions struct {
			RoleID      int      `json:"role_id"`
			Permissions []string `json:"permissions"`
		}
	)
	valPermissions := validatePermissions{}

	payloadRules := govalidator.MapData{
		"role_id":     []string{"required", "valid_id:roles", "numeric"},
		"permissions": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &valPermissions)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}
	Permissions := []models.Permissions{}
	for _, n := range valPermissions.Permissions {
		Permissions = append(Permissions, models.Permissions{
			RoleID:      valPermissions.RoleID,
			Permissions: n,
		})
	}

	for _, per := range Permissions {
		per.Create()
	}
	return c.JSON(http.StatusCreated, valPermissions)
}

func UpdatePermission(c echo.Context) error {
	defer c.Request().Body.Close()
	type (
		validatePermissions struct {
			RoleID      int      `json:"role_id"`
			Permissions []string `json:"permissions"`
		}
	)
	valPermissions := validatePermissions{}
	Permissions := []models.Permissions{}

	payloadRules := govalidator.MapData{
		"role_id":     []string{"required", "valid_id:roles", "numeric"},
		"permissions": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &valPermissions)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	asira.App.DB.Where("role_id = ?", valPermissions.RoleID).Delete(&Permissions)

	for _, n := range valPermissions.Permissions {
		Permissions = append(Permissions, models.Permissions{
			RoleID:      valPermissions.RoleID,
			Permissions: n,
		})
	}

	for _, per := range Permissions {
		per.Create()
	}

	return c.JSON(http.StatusOK, valPermissions)
}
