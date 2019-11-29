package adminhandlers

import (
	"asira_lender/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gitlab.com/asira-ayannah/basemodel"

	"github.com/lib/pq"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

// AgentPayload request body container
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

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	order := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")

	var (
		agent  models.Agent
		result basemodel.PagedFindResult
	)

	if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
		type Filter struct {
			Name          string `json:"name" condition:"LIKE,optional"`
			Username      string `json:"username" condition:"LIKE,optional"`
			ID            int64  `json:"id" condition:"optional"`
			Email         string `json:"email" condition:"optional"`
			Phone         string `json:"phone" condition:"optional"`
			Category      string `json:"category" condition:"optional"`
			AgentProvider int64  `json:"agent_provider" condition:"optional"`
			Status        string `json:"status" condition:"optional"`
		}
		id, _ := strconv.ParseInt(searchAll, 10, 64)
		result, err = agent.PagedFilterSearch(page, rows, order, sort, &Filter{
			Name:          searchAll,
			Username:      searchAll,
			ID:            id,
			Email:         searchAll,
			Phone:         searchAll,
			Category:      searchAll,
			AgentProvider: id,
			Status:        searchAll,
		})
	} else {
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
		result, err = agent.PagedFilterSearch(page, rows, order, sort, &Filter{
			Name:          c.QueryParam("name"),
			Username:      c.QueryParam("username"),
			ID:            customSplit(c.QueryParam("id"), ","),
			Email:         c.QueryParam("email"),
			Phone:         c.QueryParam("phone"),
			Category:      c.QueryParam("category"),
			AgentProvider: customSplit(c.QueryParam("agent_provider"), ","),
			Status:        c.QueryParam("status"),
		})
	}

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
		"username":       []string{"required", "unique:agents,username"},
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

	if agentPayload.Category == "account_executive" {
		if agentPayload.AgentProvider > 0 || len(agentPayload.Banks) > 1 {
			return returnInvalidResponse(http.StatusUnprocessableEntity, "account executive cannot have any agent_providers and only allowed 1 bank at a time.", "validation error")
		}
	}
	if agentPayload.Category == "agent" {
		if agentPayload.AgentProvider <= 0 {
			return returnInvalidResponse(http.StatusUnprocessableEntity, "agent must choose an agent provider.", "validation error")
		}
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
	agentPayload := AgentPayload{}
	err = agent.FindbyID(id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("agent %v tidak ditemukan", id))
	}

	payloadRules := govalidator.MapData{
		"name":           []string{},
		"email":          []string{},
		"phone":          []string{},
		"category":       []string{"agent_categories"},
		"agent_provider": []string{"valid_id:agent_providers"},
		"banks":          []string{"valid_id:banks"},
		"status":         []string{"active_inactive"},
	}

	validate := validateRequestPayload(c, payloadRules, &agentPayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	if agentPayload.Category == "account_executive" {
		if agentPayload.AgentProvider > 0 || len(agentPayload.Banks) > 1 {
			return returnInvalidResponse(http.StatusUnprocessableEntity, "account executive cannot have any agent_providers and only allowed 1 bank at a time.", "validation error")
		}
	}
	if agentPayload.Category == "agent" {
		if agentPayload.AgentProvider <= 0 {
			return returnInvalidResponse(http.StatusUnprocessableEntity, "agent must choose an agent provider.", "validation error")
		}
	}

	if len(agentPayload.Name) > 0 {
		agent.Name = agentPayload.Name
	}
	if len(agentPayload.Email) > 0 {
		agent.Email = agentPayload.Email
	}
	if len(agentPayload.Phone) > 0 {
		agent.Phone = agentPayload.Phone
	}
	if len(agentPayload.Category) > 0 {
		agent.Category = agentPayload.Category
	}
	if agentPayload.AgentProvider > 0 {
		agent.AgentProvider = sql.NullInt64{
			Int64: agentPayload.AgentProvider,
			Valid: true,
		}
	}
	if len(agentPayload.Banks) > 0 {
		agent.Banks = pq.Int64Array(agentPayload.Banks)
	}
	if len(agentPayload.Status) > 0 {
		agent.Status = agentPayload.Status
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
