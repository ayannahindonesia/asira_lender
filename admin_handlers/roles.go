package admin_handlers

import (
	"asira_lender/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

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

func RoleNew(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_role_new")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	Iroles := models.Roles{}

	payloadRules := govalidator.MapData{
		"name":        []string{"required"},
		"description": []string{},
		"system":      []string{"required"},
		"status":      []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &Iroles)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err = Iroles.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat Internal Roles")
	}

	return c.JSON(http.StatusCreated, Iroles)
}

func RolePatch(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_role_patch")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	Iroles_id, _ := strconv.Atoi(c.Param("role_id"))

	Iroles := models.Roles{}
	err = Iroles.FindbyID(Iroles_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Internal Role %v tidak ditemukan", Iroles_id))
	}

	payloadRules := govalidator.MapData{
		"name":        []string{"required"},
		"description": []string{},
		"system":      []string{"required"},
		"status":      []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &Iroles)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err = Iroles.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update Internal Roles %v", Iroles_id))
	}

	return c.JSON(http.StatusOK, Iroles)
}

func RoleRange(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_role_range")
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
