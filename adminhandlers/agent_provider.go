package adminhandlers

import (
	"asira_lender/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ayannahindonesia/basemodel"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

// AgentProviderList get all agent providers
func AgentProviderList(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_agent_provider_list")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	order := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")

	var (
		agentProvider models.AgentProvider
		result        basemodel.PagedFindResult
	)

	if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
		type Filter struct {
			Name string `json:"name" condition:"LIKE,optional"`
			ID   int64  `json:"id" condition:"optional"`
			// PIC    string `json:"pic" condition:"LIKE,optional"`
			// Phone  string `json:"phone" condition:"LIKE,optional"`
			Status string `json:"status" condition:"LIKE,optional"`
		}
		id, _ := strconv.ParseInt(searchAll, 10, 64)
		result, err = agentProvider.PagedFilterSearch(page, rows, order, sort, &Filter{
			Name: searchAll,
			ID:   id,
			// PIC:    searchAll,
			// Phone:  searchAll,
			Status: searchAll,
		})
	} else {
		type Filter struct {
			Name   string   `json:"name" condition:"LIKE"`
			ID     []string `json:"id"`
			PIC    string   `json:"pic" condition:"LIKE"`
			Phone  string   `json:"phone" condition:"LIKE"`
			Status string   `json:"status" condition:"LIKE"`
		}
		result, err = agentProvider.PagedFilterSearch(page, rows, order, sort, &Filter{
			Name:   c.QueryParam("name"),
			ID:     customSplit(c.QueryParam("id"), ","),
			PIC:    c.QueryParam("pic"),
			Phone:  c.QueryParam("phone"),
			Status: c.QueryParam("status"),
		})
	}

	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Agent provider tidak Ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

// AgentProviderDetails find agent provider by id
func AgentProviderDetails(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_agent_provider_details")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	id, _ := strconv.Atoi(c.Param("id"))

	agentProvider := models.AgentProvider{}
	err = agentProvider.FindbyID(id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("agent provider %v tidak ditemukan", id))
	}

	return c.JSON(http.StatusOK, agentProvider)
}

// AgentProviderNew create agent providers
func AgentProviderNew(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_agent_provider_new")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	agentProvider := models.AgentProvider{}

	payloadRules := govalidator.MapData{
		"name":    []string{"required"},
		"pic":     []string{"required"},
		"phone":   []string{"required", "unique:agent_providers,phone"},
		"address": []string{"required"},
		"status":  []string{"active_inactive"},
	}

	validate := validateRequestPayload(c, payloadRules, &agentProvider)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err = agentProvider.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat tipe bank baru")
	}

	return c.JSON(http.StatusCreated, agentProvider)
}

// AgentProviderPatch edit agent providers
func AgentProviderPatch(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_agent_provider_patch")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	id, _ := strconv.Atoi(c.Param("id"))

	agentProvider := models.AgentProvider{}
	err = agentProvider.FindbyID(id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("agent provider %v tidak ditemukan", id))
	}

	payloadRules := govalidator.MapData{
		"name":    []string{},
		"pic":     []string{},
		"phone":   []string{},
		"address": []string{},
		"status":  []string{"active_inactive"},
	}

	validate := validateRequestPayload(c, payloadRules, &agentProvider)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err = agentProvider.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat tipe bank baru")
	}

	return c.JSON(http.StatusOK, agentProvider)
}
