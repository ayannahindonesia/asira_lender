package adminhandlers

import (
	"asira_lender/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ayannahindonesia/basemodel"
	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo"
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
	err := validatePermission(c, "lender_service_list")
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

// ServiceDetail get service by id
func ServiceDetail(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "lender_service_list_detail")
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
