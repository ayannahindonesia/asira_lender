package admin_handlers

import (
	"asira_lender/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

type ServicePayload struct {
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
}

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

	// filters
	id := customSplit(c.QueryParam("id"), ",")
	name := c.QueryParam("name")
	status := c.QueryParam("status")

	type Filter struct {
		ID     []string `json:"id"`
		Name   string   `json:"name" condition:"LIKE"`
		Status string   `json:"status"`
	}

	service := models.Service{}
	result, err := service.PagedFindFilter(page, rows, order, sort, &Filter{
		ID:     id,
		Name:   name,
		Status: status,
	})
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "pencarian tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

func ServiceNew(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_service_new")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	servicePayload := ServicePayload{}

	payloadRules := govalidator.MapData{
		"name":   []string{"required"},
		"image":  []string{"required"},
		"status": []string{"required", "active_inactive"},
	}

	validate := validateRequestPayload(c, payloadRules, &servicePayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	image := models.Image{
		Image_string: servicePayload.Image,
	}
	image.Create()

	service := models.Service{
		Name:    servicePayload.Name,
		ImageID: image.ID,
		Status:  servicePayload.Status,
	}
	err = service.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat layanan bank baru")
	}

	return c.JSON(http.StatusCreated, service)
}

func ServiceDetail(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_service_detail")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	serviceId, _ := strconv.Atoi(c.Param("id"))

	service := models.Service{}
	err = service.FindbyID(serviceId)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("layanan %v tidak ditemukan", serviceId))
	}

	return c.JSON(http.StatusOK, service)
}

func ServicePatch(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_service_patch")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	serviceId, _ := strconv.Atoi(c.Param("id"))

	service := models.Service{}
	err = service.FindbyID(serviceId)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("layanan %v tidak ditemukan", serviceId))
	}

	serviceImg := models.Image{}
	err = serviceImg.FindbyID(int(service.ImageID))

	servicePayload := ServicePayload{}
	payloadRules := govalidator.MapData{
		"name":   []string{},
		"image":  []string{},
		"status": []string{"active_inactive"},
	}
	validate := validateRequestPayload(c, payloadRules, &servicePayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	if len(servicePayload.Name) > 0 {
		service.Name = servicePayload.Name
	}
	if len(servicePayload.Image) > 0 {
		serviceImg.Image_string = servicePayload.Image
	}
	if len(servicePayload.Status) > 0 {
		service.Status = servicePayload.Status
	}

	err = service.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update layanan %v", serviceId))
	}
	err = serviceImg.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update layanan %v", serviceId))
	}

	return c.JSON(http.StatusOK, service)
}

func ServiceDelete(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_service_delete")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	serviceId, _ := strconv.Atoi(c.Param("id"))

	service := models.Service{}
	err = service.FindbyID(serviceId)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank type %v tidak ditemukan", serviceId))
	}

	err = service.Delete()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank tipe %v", serviceId))
	}

	return c.JSON(http.StatusOK, service)
}
