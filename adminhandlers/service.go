package adminhandlers

import (
	"asira_lender/middlewares"
	"asira_lender/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
		NLog("warning", "ServiceList", map[string]interface{}{"message": "error listing services", "error": err}, c.Get("user").(*jwt.Token), "", false)

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

	url, err := UploadCloudImage("svc", servicePayload.Image)
	if err != nil {
		NLog("error", "ServiceNew", map[string]interface{}{"message": "upload image error", "error": err}, c.Get("user").(*jwt.Token), "", false)

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
		NLog("error", "ServiceNew", map[string]interface{}{"message": "service create error", "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat layanan baru")
	}

	middlewares.SubmitKafkaPayload(service, "service_create")
	if err != nil {
		NLog("error", "ServiceNew", map[string]interface{}{"message": "kafka submit error", "error": err, "service": service}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat layanan baru")
	}

	NAudittrail(models.Service{}, service, c.Get("user").(*jwt.Token), "service", fmt.Sprint(service.ID), "create")

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
		NLog("warning", "ServiceDetail", map[string]interface{}{"message": fmt.Sprintf("error finding service %v", serviceID), "error": err}, c.Get("user").(*jwt.Token), "", false)

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
		NLog("warning", "ServicePatch", map[string]interface{}{"message": fmt.Sprintf("error finding service %v", serviceID), "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Layanan %v tidak ditemukan", serviceID))
	}
	origin := service

	servicePayload := ServicePayload{}
	payloadRules := govalidator.MapData{
		"name":        []string{},
		"image":       []string{},
		"status":      []string{"active_inactive"},
		"description": []string{},
	}
	validate := validateRequestPayload(c, payloadRules, &servicePayload)
	if validate != nil {
		NLog("warning", "ServicePatch", map[string]interface{}{"message": "validation error", "error": validate}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
	}

	if len(servicePayload.Name) > 0 {
		service.Name = servicePayload.Name
	}
	if len(servicePayload.Image) > 0 {
		url, err := UploadCloudImage("svc", servicePayload.Image)
		if err != nil {
			NLog("error", "ServicePatch", map[string]interface{}{"message": "error upload image", "error": err}, c.Get("user").(*jwt.Token), "", false)

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
		NLog("error", "ServicePatch", map[string]interface{}{"message": "kafka submit error", "error": err, "service": service}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update layanan %v", serviceID))
	}

	NAudittrail(origin, service, c.Get("user").(*jwt.Token), "service", fmt.Sprint(service.ID), "update")

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
		NLog("warning", "ServiceDelete", map[string]interface{}{"message": fmt.Sprintf("error finding service %v", serviceID), "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Layanan %v tidak ditemukan", serviceID))
	}

	err = middlewares.SubmitKafkaPayload(service, "service_delete")
	if err != nil {
		NLog("error", "ServiceDelete", map[string]interface{}{"message": "error submit kafka", "error": err, "service": service}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal delete layanan %v", serviceID))
	}

	NAudittrail(service, models.Service{}, c.Get("user").(*jwt.Token), "service", fmt.Sprint(service.ID), "delete")

	return c.JSON(http.StatusOK, service)
}
