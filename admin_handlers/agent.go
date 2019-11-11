package admin_handlers

import (
	"asira_lender/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

type AgentPayload struct {
	Name          string  `json:"name"`
	Username      string  `json:"username"`
	Email         string  `json:"email"`
	Phone         string  `json:"phone"`
	Category      string  `json:"category"`
	AgentProvider int64   `json:"agent_provider"`
	Banks         []int64 `json:"banks"`
	Status        string  `json:"status"`
}

// AgentList get all agents
func AgentList(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_agent_list")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	agent := models.Agent{}

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	order := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")
	// filters
	name := c.QueryParam("name")
	username := c.QueryParam("username")
	id := customSplit(c.QueryParam("id"), ",")
	email := c.QueryParam("email")
	phone := c.QueryParam("phone")
	category := c.QueryParam("category")
	agentProvider := customSplit(c.QueryParam("agent_provider"), ",")
	status := c.QueryParam("status")

	type Filter struct {
		Name          string   `json:"name" condition:"LIKE"`
		Username      string   `json:"username" condition:"LIKE"`
		ID            []string `json:"id"`
		Email         string   `json:"email"`
		Phone         string   `json:"phone"`
		Category      string   `json:"category"`
		AgentProvider []string `json:"agent_provider"`
		Status        string   `json:"status"`
	}
	result, err := agent.PagedFilterSearch(page, rows, order, sort, &Filter{
		Name:          name,
		Username:      username,
		ID:            id,
		Email:         email,
		Phone:         phone,
		Category:      category,
		AgentProvider: agentProvider,
		Status:        status,
	})

	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Agent provider tidak Ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

// AgentDetails find agent by id
func AgentDetails(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_agent_details")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	id, _ := strconv.Atoi(c.Param("id"))

	agent := models.Agent{}
	err = agent.FindbyID(id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("agent %v tidak ditemukan", id))
	}

	return c.JSON(http.StatusOK, agent)
}

// AgentNew create agent
func AgentNew(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_agent_new")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	agentPayload := AgentPayload{}

	payloadRules := govalidator.MapData{
		"name":           []string{"required"},
		"username":       []string{"required"},
		"email":          []string{"required", "unique:agents,email"},
		"phone":          []string{"required", "unique:agents,phone"},
		"category":       []string{"required", "agent_categories"},
		"agent_provider": []string{"valid_id:agent_providers"},
		"banks":          []string{"required", "valid_id:banks"},
		"status":         []string{"active_inactive"},
	}

	validate := validateRequestPayload(c, payloadRules, &agentPayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	agent := models.Agent{}

	marshal, _ := json.Marshal(agentPayload)
	json.Unmarshal(marshal, &agent)

	if agentPayload.AgentProvider != 0 {
		agent.AgentProvider = sql.NullInt64{
			Int64: agentPayload.AgentProvider,
			Valid: true,
		}
	}

	err = agent.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat agent baru")
	}

	return c.JSON(http.StatusCreated, agent)
}

// AgentPatch edit agent
func AgentPatch(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_agent_patch")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	id, _ := strconv.Atoi(c.Param("id"))

	agent := models.Agent{}
	err = agent.FindbyID(id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("agent %v tidak ditemukan", id))
	}

	payloadRules := govalidator.MapData{
		"name":           []string{},
		"email":          []string{"unique:agents,email,1"},
		"phone":          []string{"unique:agents,phone,1"},
		"category":       []string{"agent_categories"},
		"agent_provider": []string{"valid_id:agent_providers"},
		"banks":          []string{},
		"status":         []string{"active_inactive"},
	}

	validate := validateRequestPayload(c, payloadRules, &agent)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err = agent.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal mengubah agent baru")
	}

	return c.JSON(http.StatusOK, agent)
}

// AgentDelete edit agent
func AgentDelete(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_agent_delete")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	id, _ := strconv.Atoi(c.Param("id"))

	agent := models.Agent{}
	err = agent.FindbyID(id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("agent %v tidak ditemukan", id))
	}

	err = agent.Delete()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal mengubah agent baru")
	}

	return c.JSON(http.StatusOK, agent)
}
