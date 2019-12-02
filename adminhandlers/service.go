package adminhandlers

import (
	"asira_lender/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gitlab.com/asira-ayannah/basemodel"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

// ServicePayload handles request body
type ServicePayload struct {
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
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
		return returnInvalidResponse(http.StatusInternalServerError, err, "pencarian tidak ditemukan")
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

// ServiceDetail get service by id
func ServiceDetail(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_service_detail")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	serviceID, _ := strconv.Atoi(c.Param("id"))

	service := models.Service{}
	err = service.FindbyID(serviceID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("layanan %v tidak ditemukan", serviceID))
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

	serviceID, _ := strconv.Atoi(c.Param("id"))

	service := models.Service{}
	err = service.FindbyID(serviceID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("layanan %v tidak ditemukan", serviceID))
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
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update layanan %v", serviceID))
	}
	err = serviceImg.Save()
	if err != nil {
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

	serviceID, _ := strconv.Atoi(c.Param("id"))

	service := models.Service{}
	err = service.FindbyID(serviceID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank type %v tidak ditemukan", serviceID))
	}

	err = service.Delete()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank tipe %v", serviceID))
	}

	return c.JSON(http.StatusOK, service)
}
