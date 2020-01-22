package adminhandlers

import (
	"asira_lender/asira"
	"asira_lender/middlewares"
	"asira_lender/models"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ayannahindonesia/basemodel"

	"github.com/lib/pq"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

type (
	// AgentPayload request body container
	AgentPayload struct {
		Name          string  `json:"name"`
		Username      string  `json:"username"`
		Image         string  `json:"image"`
		Email         string  `json:"email"`
		Phone         string  `json:"phone"`
		Category      string  `json:"category"`
		AgentProvider int64   `json:"agent_provider"`
		Banks         []int64 `json:"banks"`
		Status        string  `json:"status"`
	}
	// AgentSelect query result container
	AgentSelect struct {
		models.Agent
		AgentProviderName string         `json:"agent_provider_name"`
		BankNames         pq.StringArray `json:"bank_names"`
	}
)

// AgentList get all agents
func AgentList(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_agent_list")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	db := asira.App.DB
	var (
		totalRows int
		offset    int
		rows      int
		page      int
		lastPage  int
		agents    []AgentSelect
	)

	// pagination parameters
	rows, _ = strconv.Atoi(c.QueryParam("rows"))
	if rows > 0 {
		page, _ = strconv.Atoi(c.QueryParam("page"))
		if page <= 0 {
			page = 1
		}
		offset = (page * rows) - rows
	}

	db = db.Table("agents").
		Select("agents.*, ap.name as agent_provider_name, (SELECT ARRAY_AGG(name) FROM banks WHERE id IN (SELECT UNNEST(agents.banks))) as bank_names").
		Joins("LEFT JOIN agent_providers ap ON agents.agent_provider = ap.id")

	if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
		searchLike := "%" + strings.ToLower(searchAll) + "%"
		extraquery := fmt.Sprintf("CAST(agents.id as varchar(255)) = ?") + // use searchAll #1
			fmt.Sprintf(" OR LOWER(agents.name) LIKE ?") + // use searchLike #2
			fmt.Sprintf(" OR LOWER(agents.category) LIKE ?") + // use searchLike #3
			fmt.Sprintf(" OR LOWER(ap.name) LIKE ?") + // use searchLike #4
			fmt.Sprintf(" OR LOWER(agents.status) LIKE ?") // use searchLike #5

		db = db.Where(extraquery, searchAll, searchLike, searchLike, searchLike, searchLike)
	} else {
		if id := customSplit(c.QueryParam("id"), ","); len(id) > 0 {
			db = db.Where("agents.id IN (?)", id)
		}
		if name := c.QueryParam("name"); len(name) > 0 {
			db = db.Where("LOWER(agents.name) LIKE ?", "%"+strings.ToLower(name)+"%")
		}
		if username := c.QueryParam("username"); len(username) > 0 {
			db = db.Where("LOWER(agents.username) = ?", "%"+strings.ToLower(username)+"%")
		}
		if email := c.QueryParam("email"); len(email) > 0 {
			db = db.Where("LOWER(agents.email) LIKE ?", "%"+strings.ToLower(email)+"%")
		}
		if phone := c.QueryParam("phone"); len(phone) > 0 {
			db = db.Where("agents.phone LIKE ?", "%"+phone+"%")
		}
		if category := c.QueryParam("category"); len(category) > 0 {
			db = db.Where("LOWER(agents.category) LIKE ?", "%"+strings.ToLower(category)+"%")
		}
		if agentProvider := customSplit(c.QueryParam("agent_provider"), ","); len(agentProvider) > 0 {
			db = db.Where("ap.id IN (?)", agentProvider)
		}
		if status := c.QueryParam("status"); len(status) > 0 {
			db = db.Where("LOWER(agents.status) LIKE ?", "%"+strings.ToLower(status)+"%")
		}
		if agentProviderName := c.QueryParam("agent_provider_name"); len(agentProviderName) > 0 {
			db = db.Where("LOWER(ap.name) LIKE ?", "%"+strings.ToLower(agentProviderName)+"%")
		}
		if bankID := c.QueryParam("bank_id"); len(bankID) > 0 {
			db = db.Where("agents.banks LIKE ?", "%"+bankID+"%")
		}
	}

	if order := strings.Split(c.QueryParam("orderby"), ","); len(order) > 0 {
		if sort := strings.Split(c.QueryParam("sort"), ","); len(sort) > 0 {
			for k, v := range order {
				q := v
				if len(sort) > k {
					value := sort[k]
					if strings.ToUpper(value) == "ASC" || strings.ToUpper(value) == "DESC" {
						q = v + " " + strings.ToUpper(value)
					}
				}
				db = db.Order(q)
			}
		}
	}

	tempDB := db
	tempDB.Where("agents.deleted_at IS NULL").Count(&totalRows)

	if rows > 0 {
		db = db.Limit(rows).Offset(offset)
		lastPage = int(math.Ceil(float64(totalRows) / float64(rows)))
	}
	err = db.Find(&agents).Error
	if err != nil {
		log.Println(err)
	}

	result := basemodel.PagedFindResult{
		TotalData:   totalRows,
		Rows:        rows,
		CurrentPage: page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + rows,
		Data:        agents,
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

	agent := AgentSelect{}

	db := asira.App.DB

	err = db.Table("agents").
		Select("agents.*, ap.name as agent_provider_name, (SELECT ARRAY_AGG(name) FROM banks WHERE id IN (SELECT UNNEST(agents.banks))) as bank_names").
		Joins("LEFT JOIN agent_providers ap ON agents.agent_provider = ap.id").
		Where("agents.id = ?", id).
		Find(&agent).Error

	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Agen %v tidak ditemukan", id))
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
		"image":          []string{},
		"email":          []string{"required", "unique:agents,email"},
		"phone":          []string{"required", "unique:agents,phone"},
		"category":       []string{"required", "agent_categories"},
		"agent_provider": []string{"valid_id:agent_providers", "status:agent_providers,active"},
		"banks":          []string{"required", "valid_id:banks"},
		"status":         []string{"active_inactive"},
	}

	validate := validateRequestPayload(c, payloadRules, &agentPayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
	}

	if agentPayload.Category == "account_executive" {
		if agentPayload.AgentProvider > 0 || len(agentPayload.Banks) > 1 {
			return returnInvalidResponse(http.StatusUnprocessableEntity, "Account executive tidak dapat memiliki penyedia agen dan hanya memiliki 1 bank.", "Hambatan validasi")
		}
	}
	if agentPayload.Category == "agent" {
		if agentPayload.AgentProvider <= 0 {
			return returnInvalidResponse(http.StatusUnprocessableEntity, "Agen harus memilih penyedia agent.", "validation error")
		}
	}

	agent := models.Agent{}

	marshal, _ := json.Marshal(agentPayload)
	json.Unmarshal(marshal, &agent)

	if len(agentPayload.Image) > 0 {
		unbased, _ := base64.StdEncoding.DecodeString(agentPayload.Image)
		filename := "agt" + strconv.FormatInt(time.Now().Unix(), 10)
		url, err := asira.App.S3.UploadJPEG(unbased, filename)
		if err != nil {
			return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat agent baru")
		}

		agent.Image = url
	}

	if agentPayload.AgentProvider != 0 {
		agent.AgentProvider = sql.NullInt64{
			Int64: agentPayload.AgentProvider,
			Valid: true,
		}
	}

	agent.Create()
	err = middlewares.SubmitKafkaPayload(agent, "agent_create")
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

	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	agent := models.Agent{}
	agentPayload := AgentPayload{}
	err = agent.FindbyID(id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Agen %v tidak ditemukan", id))
	}

	payloadRules := govalidator.MapData{
		"name":           []string{},
		"image":          []string{},
		"email":          []string{},
		"phone":          []string{},
		"category":       []string{"agent_categories"},
		"agent_provider": []string{"valid_id:agent_providers", "status:agent_providers,active"},
		"banks":          []string{"valid_id:banks"},
		"status":         []string{"active_inactive"},
	}

	validate := validateRequestPayload(c, payloadRules, &agentPayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
	}

	if agentPayload.Category == "account_executive" {
		if agentPayload.AgentProvider > 0 || len(agentPayload.Banks) > 1 {
			return returnInvalidResponse(http.StatusUnprocessableEntity, "Account executive tidak dapat memiliki penyedia agen dan hanya boleh memiliki 1 bank.", "Hambatan validasi")
		}
	}
	if agentPayload.Category == "agent" {
		if agentPayload.AgentProvider <= 0 {
			return returnInvalidResponse(http.StatusUnprocessableEntity, "Agen harus memilih penyedia agen.", "Hambatan validasi")
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
	if len(agentPayload.Image) > 0 {
		unbased, _ := base64.StdEncoding.DecodeString(agentPayload.Image)
		filename := "agt" + strconv.FormatInt(time.Now().Unix(), 10)
		url, err := asira.App.S3.UploadJPEG(unbased, filename)
		if err != nil {
			return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat agent baru")
		}

		i := strings.Split(agent.Image, "/")
		delImage := i[len(i)-1]
		err = asira.App.S3.DeleteObject(delImage)
		if err != nil {
			log.Printf("failed to delete image %v from s3 bucket", delImage)
		}

		agent.Image = url
	}

	err = middlewares.SubmitKafkaPayload(agent, "agent_update")
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

	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	agent := models.Agent{}
	err = agent.FindbyID(id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Agen %v tidak ditemukan", id))
	}

	err = middlewares.SubmitKafkaPayload(agent, "agent_delete")
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal mengubah agen baru")
	}

	return c.JSON(http.StatusOK, agent)
}
