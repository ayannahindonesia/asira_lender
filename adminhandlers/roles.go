package adminhandlers

import (
	"asira_lender/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ayannahindonesia/basemodel"
	"github.com/dgrijalva/jwt-go"
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

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")

	var (
		roles  models.Roles
		result basemodel.PagedFindResult
	)

	if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
		type Filter struct {
			ID     int64  `json:"id" condition:"optional"`
			Name   string `json:"name" condition:"LIKE,optional"`
			Status string `json:"status" condition:"optional"`
		}
		id, _ := strconv.ParseInt(searchAll, 10, 64)
		result, err = roles.PagedFindFilter(page, rows, orderby, sort, &Filter{
			ID:     id,
			Name:   searchAll,
			Status: searchAll,
		})
	} else {
		type Filter struct {
			ID     []string `json:"id"`
			Name   string   `json:"name" condition:"LIKE"`
			Status string   `json:"status"`
		}
		result, err = roles.PagedFindFilter(page, rows, orderby, sort, &Filter{
			ID:     customSplit(c.QueryParam("id"), ","),
			Name:   c.QueryParam("name"),
			Status: c.QueryParam("status"),
		})
	}

	if err != nil {
		NLog("warning", "RoleList", map[string]interface{}{"message": "role listing error", "error": err}, c.Get("user").(*jwt.Token), "", false)

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

	IrolesID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	err = Iroles.FindbyID(IrolesID)
	if err != nil {
		NLog("warning", "RoleDetails", map[string]interface{}{"message": "error finding role", "error": err}, c.Get("user").(*jwt.Token), "", false)

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
		NLog("warning", "RoleNew", map[string]interface{}{"message": "validation error", "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
	}

	marshal, _ := json.Marshal(rolePayload)
	json.Unmarshal(marshal, &Iroles)

	err = Iroles.Create()
	if err != nil {
		NLog("error", "RoleNew", map[string]interface{}{"message": "error creating role", "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat Internal Roles")
	}

	NAudittrail(models.Roles{}, Iroles, c.Get("user").(*jwt.Token), "role", fmt.Sprint(Iroles.ID), "create")

	return c.JSON(http.StatusCreated, Iroles)
}

// RolePatch edit role by id
func RolePatch(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_role_patch")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	IrolesID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	Iroles := models.Roles{}
	rolePayload := RolePayload{}
	err = Iroles.FindbyID(IrolesID)
	if err != nil {
		NLog("warning", "RolePatch", map[string]interface{}{"message": fmt.Sprintf("error finding role %v", IrolesID), "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Internal Role %v tidak ditemukan", IrolesID))
	}
	origin := Iroles

	payloadRules := govalidator.MapData{
		"name":        []string{},
		"description": []string{},
		"system":      []string{},
		"status":      []string{"active_inactive"},
		"permissions": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &rolePayload)
	if validate != nil {
		NLog("warning", "RolePatch", map[string]interface{}{"message": "validation error", "error": validate}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
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
		NLog("error", "RolePatch", map[string]interface{}{"message": "error saving role", "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update Internal Roles %v", IrolesID))
	}

	NAudittrail(origin, Iroles, c.Get("user").(*jwt.Token), "role", fmt.Sprint(Iroles.ID), "update")

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
	orderby := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")

	name := c.QueryParam("name")
	id := customSplit(c.QueryParam("id"), ",")
	status := c.QueryParam("status")

	type Filter struct {
		ID     []string `json:"id"`
		Name   string   `json:"name" condition:"LIKE"`
		Status string   `json:"status"`
	}

	result, err := Iroles.FindFilter(offset, limit, orderby, sort, &Filter{
		ID:     id,
		Name:   name,
		Status: status,
	})

	if err != nil {
		NLog("warning", "RoleRange", map[string]interface{}{"message": "error listing roles", "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusNotFound, err, "Role tidak Ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}
