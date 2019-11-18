package adminhandlers

import (
	"asira_lender/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lib/pq"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

// RolePayload handles role request body
type RolePayload struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	System      string   `json:"system"`
	Status      string   `json:"status"`
	Permissions []string `json:"permissions"`
}

// RoleList get all roles
func RoleList(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_role_list")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	Iroles := models.Roles{}
	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	name := c.QueryParam("name")
	id := customSplit(c.QueryParam("id"), ",")
	status := c.QueryParam("status")

	type Filter struct {
		ID     []string `json:"id"`
		Name   string   `json:"name" condition:"LIKE"`
		Status string   `json:"status"`
	}

	result, err := Iroles.PagedFilterSearch(page, rows, orderby, sort, &Filter{
		ID:     id,
		Name:   name,
		Status: status,
	})

	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Role tidak Ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

// RoleDetails get role detail by id
func RoleDetails(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_role_details")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	Iroles := models.Roles{}

	IrolesID, _ := strconv.Atoi(c.Param("role_id"))
	err = Iroles.FindbyID(IrolesID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Role ID tidak ditemukan")
	}

	return c.JSON(http.StatusOK, Iroles)
}

// RoleNew create new role
func RoleNew(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_role_new")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	Iroles := models.Roles{}
	rolePayload := RolePayload{}

	payloadRules := govalidator.MapData{
		"name":        []string{"required"},
		"description": []string{},
		"system":      []string{"required"},
		"status":      []string{"active_inactive"},
		"permissions": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &rolePayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	marshal, _ := json.Marshal(rolePayload)
	json.Unmarshal(marshal, &Iroles)

	err = Iroles.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat Internal Roles")
	}

	return c.JSON(http.StatusCreated, Iroles)
}

// RolePatch edit role by id
func RolePatch(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_role_patch")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	IrolesID, _ := strconv.Atoi(c.Param("role_id"))

	Iroles := models.Roles{}
	rolePayload := RolePayload{}
	err = Iroles.FindbyID(IrolesID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Internal Role %v tidak ditemukan", IrolesID))
	}

	payloadRules := govalidator.MapData{
		"name":        []string{"required"},
		"description": []string{},
		"system":      []string{"required"},
		"status":      []string{"active_inactive"},
		"permissions": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &rolePayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	if len(rolePayload.Name) > 0 {
		Iroles.Name = rolePayload.Name
	}
	if len(rolePayload.Description) > 0 {
		Iroles.Description = rolePayload.Description
	}
	if len(rolePayload.System) > 0 {
		Iroles.System = rolePayload.System
	}
	if len(rolePayload.Status) > 0 {
		Iroles.Status = rolePayload.Status
	}
	if len(rolePayload.Permissions) > 0 {
		Iroles.Permissions = pq.StringArray(rolePayload.Permissions)
	}

	err = Iroles.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update Internal Roles %v", IrolesID))
	}

	return c.JSON(http.StatusOK, Iroles)
}

// RoleRange get all role without pagination
func RoleRange(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_role_list")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	Iroles := models.Roles{}
	// pagination parameters
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	name := c.QueryParam("name")
	id := customSplit(c.QueryParam("id"), ",")
	status := c.QueryParam("status")

	type Filter struct {
		ID     []string `json:"id"`
		Name   string   `json:"name" condition:"LIKE"`
		Status string   `json:"status"`
	}

	result, err := Iroles.FilterSearch(offset, limit, orderby, sort, &Filter{
		ID:     id,
		Name:   name,
		Status: status,
	})

	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Role tidak Ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}
