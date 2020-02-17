package adminhandlers

import (
	"asira_lender/asira"
	"asira_lender/middlewares"
	"asira_lender/models"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ayannahindonesia/basemodel"
	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

// ServicePayload handles request body
type ServicePayload struct {
	Name        string `json:"name"`
	Image       string `json:"image"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// ServiceList gets all services
func ServiceList(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_service_list")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	order := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")

	var (
		service models.Service
		result  basemodel.PagedFindResult
	)

	if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
		type Filter struct {
			ID     int64  `json:"id" condition:"optional"`
			Name   string `json:"name" condition:"LIKE,optional"`
			Status string `json:"status" condition:"optional"`
		}
		id, _ := strconv.ParseInt(searchAll, 10, 64)
		result, err = service.PagedFindFilter(page, rows, order, sort, &Filter{
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
		result, err = service.PagedFindFilter(page, rows, order, sort, &Filter{
			ID:     customSplit(c.QueryParam("id"), ","),
			Name:   c.QueryParam("name"),
			Status: c.QueryParam("status"),
		})
	}

	if err != nil {
		NLog("warning", "ServiceList", fmt.Sprintf("error : %v", err), c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Pencarian tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

// ServiceNew add new service
func ServiceNew(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_service_new")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	servicePayload := ServicePayload{}

	payloadRules := govalidator.MapData{
		"name":        []string{"required"},
		"image":       []string{"required"},
		"status":      []string{"required", "active_inactive"},
		"description": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &servicePayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
	}

	unbased, _ := base64.StdEncoding.DecodeString(servicePayload.Image)
	filename := "svc" + strconv.FormatInt(time.Now().Unix(), 10)
	url, err := asira.App.S3.UploadJPEG(unbased, filename)
	if err != nil {
		NLog("error", "ServiceNew", fmt.Sprintf("upload image error : %v", err), c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat layanan bank baru")
	}

	service := models.Service{
		Name:        servicePayload.Name,
		Image:       url,
		Status:      servicePayload.Status,
		Description: servicePayload.Description,
	}

	err = service.Create()
	if err != nil {
		NLog("error", "ServiceNew", fmt.Sprintf("service create error : %v", err), c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat layanan baru")
	}

	middlewares.SubmitKafkaPayload(service, "service_create")
	if err != nil {
		NLog("error", "ServiceNew", fmt.Sprintf("kafka submit error : %v", err), c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat layanan baru")
	}

	return c.JSON(http.StatusCreated, service)
}

// ServiceDetail get service by id
func ServiceDetail(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_service_detail")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	serviceID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	service := models.Service{}
	err = service.FindbyID(serviceID)
	if err != nil {
		NLog("warning", "ServiceDetail", fmt.Sprintf("error finding service %v : %v", serviceID, err), c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Layanan %v tidak ditemukan", serviceID))
	}

	return c.JSON(http.StatusOK, service)
}

// ServicePatch update service by id
func ServicePatch(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_service_patch")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	serviceID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	service := models.Service{}
	err = service.FindbyID(serviceID)
	if err != nil {
		NLog("warning", "ServicePatch", fmt.Sprintf("error finding service %v : %v", serviceID, err), c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Layanan %v tidak ditemukan", serviceID))
	}

	servicePayload := ServicePayload{}
	payloadRules := govalidator.MapData{
		"name":        []string{},
		"image":       []string{},
		"status":      []string{"active_inactive"},
		"description": []string{},
	}
	validate := validateRequestPayload(c, payloadRules, &servicePayload)
	if validate != nil {
		NLog("warning", "ServicePatch", fmt.Sprintf("validation error : %v", validate), c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
	}

	if len(servicePayload.Name) > 0 {
		service.Name = servicePayload.Name
	}
	if len(servicePayload.Image) > 0 {
		unbased, _ := base64.StdEncoding.DecodeString(servicePayload.Image)
		filename := "svc" + strconv.FormatInt(time.Now().Unix(), 10)
		url, err := asira.App.S3.UploadJPEG(unbased, filename)
		if err != nil {
			NLog("error", "ServicePatch", fmt.Sprintf("error upload image : %v", err), c.Get("user").(*jwt.Token), "", false)

			return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat layanan bank baru")
		}

		service.Image = url
	}
	if len(servicePayload.Status) > 0 {
		service.Status = servicePayload.Status
	}

	if len(servicePayload.Description) > 0 {
		service.Description = servicePayload.Description
	}

	err = middlewares.SubmitKafkaPayload(service, "service_update")
	if err != nil {
		NLog("error", "ServicePatch", fmt.Sprintf("kafka submit error : %v", err), c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update layanan %v", serviceID))
	}

	return c.JSON(http.StatusOK, service)
}

// ServiceDelete delete service by id
func ServiceDelete(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_service_delete")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	serviceID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	service := models.Service{}
	err = service.FindbyID(serviceID)
	if err != nil {
		NLog("warning", "ServiceDelete", fmt.Sprintf("error finding service %v : %v", serviceID, err), c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Layanan %v tidak ditemukan", serviceID))
	}

	err = middlewares.SubmitKafkaPayload(service, "service_delete")
	if err != nil {
		NLog("error", "ServiceDelete", fmt.Sprintf("error submit kafka : %v", err), c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal delete layanan %v", serviceID))
	}

	return c.JSON(http.StatusOK, service)
}
