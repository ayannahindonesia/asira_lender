package handlers

import (
	"asira_lender/adminhandlers"
	"asira_lender/asira"
	"asira_lender/models"
	"fmt"
	"math"
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

// LenderServiceList gets all services
func LenderServiceList(c echo.Context) error {
	defer c.Request().Body.Close()

	var (
		// service   models.Service
		result    basemodel.PagedFindResult
		totalRows int
		offset    int
		rows      int
		page      int
		lastPage  int
		services  []models.Service
	)

	err := validatePermission(c, "lender_service_list")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	jti := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["jti"].(string)
	lenderID, _ := strconv.ParseUint(jti, 10, 64)
	bankRep := models.BankRepresentatives{}

	//get bank representatives
	err = bankRep.FindbyUserID(int(lenderID))
	if err != nil {
		adminhandlers.NLog("warning", "LenderServiceList", map[string]interface{}{"message": "error listing services", "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	// pagination parameters
	rows, err = strconv.Atoi(c.QueryParam("rows"))
	page, err = strconv.Atoi(c.QueryParam("page"))
	order := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")

	db := asira.App.DB

	// pagination parameters
	if rows > 0 {
		if page <= 0 {
			page = 1
		}
		offset = (page * rows) - rows
	}

	db = db.Table("services").
		Select("services.*").
		Joins("INNER JOIN banks b ON services.id IN (SELECT UNNEST(b.services)) ").
		Where("b.id = ?", bankRep.BankID)

	if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
		//all value for all params
		extraquery := fmt.Sprintf("CAST(services.id as varchar(255)) = ?") + // use searchAll #1
			fmt.Sprintf(" OR LOWER(services.name) LIKE ?") + // use searchLike #2
			fmt.Sprintf(" OR LOWER(services.status) LIKE ?") // use searchLike #3

		db = db.Where(extraquery, searchAll, "%"+searchAll+"%", "%"+searchAll+"%")

	} else {
		if id := c.QueryParam("id"); len(id) > 0 {
			db = db.Where("services.id IN (?)", id)
		}
		if name := c.QueryParam("name"); len(name) > 0 {
			db = db.Where("LOWER(services.name) LIKE ?", "%"+strings.ToLower(name)+"%")
		}
		if status := c.QueryParam("status"); len(status) > 0 {
			db = db.Where("LOWER(services.status) LIKE ?", "%"+strings.ToLower(status)+"%")
		}
	}

	if len(order) > 0 {
		if len(sort) > 0 {
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
	tempDB.Count(&totalRows)

	if rows > 0 {
		db = db.Limit(rows).Offset(offset)
		lastPage = int(math.Ceil(float64(totalRows) / float64(rows)))
	}

	err = db.Find(&services).Error
	if err != nil {
		adminhandlers.NLog("warning", "LenderServiceList", map[string]interface{}{"message": "error listing services", "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Pencarian tidak ditemukan")
	}

	result = basemodel.PagedFindResult{
		TotalData:   totalRows,
		Rows:        rows,
		CurrentPage: page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + rows,
		Data:        services,
	}

	return c.JSON(http.StatusOK, result)
}

// LenderServiceLListDetail get service by id
func LenderServiceLListDetail(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "lender_service_list_detail")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	serviceID, _ := strconv.ParseUint(c.Param("service_id"), 10, 64)

	jti := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["jti"].(string)
	lenderID, _ := strconv.ParseUint(jti, 10, 64)
	bankRep := models.BankRepresentatives{}

	//get bank representatives
	err = bankRep.FindbyUserID(int(lenderID))
	if err != nil {
		adminhandlers.NLog("warning", "LenderServiceList", map[string]interface{}{"message": "error listing services", "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	db := asira.App.DB

	db = db.Table("services").
		Select("*").
		Joins("INNER JOIN banks b ON services.id IN (SELECT UNNEST(b.services)) ").
		Where("b.id = ?", bankRep.BankID).
		Where("services.id = ?", serviceID)

	var service models.Service

	err = db.Find(&service).Error
	if err != nil {
		adminhandlers.NLog("warning", "ServiceDetail", map[string]interface{}{"message": fmt.Sprintf("error finding service %v", serviceID), "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Layanan %v tidak ditemukan", serviceID))
	}

	return c.JSON(http.StatusOK, service)
}
