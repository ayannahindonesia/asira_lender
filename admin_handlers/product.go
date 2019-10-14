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

func ProductList(c echo.Context) error {
	defer c.Request().Body.Close()

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	order := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")

	// filters
	id := customSplit(c.QueryParam("id"), ",")
	name := c.QueryParam("name")
	serviceID := c.QueryParam("service_id")
	minTimespan := c.QueryParam("min_timespan")
	maxTimespan := c.QueryParam("max_timespan")
	interest := c.QueryParam("interest")
	minLoan := c.QueryParam("min_loan")
	maxLoan := c.QueryParam("max_loan")
	fee := c.QueryParam("fee")
	collaterals := c.QueryParam("collaterals")
	financingSector := c.QueryParam("financing_sector")
	assurance := c.QueryParam("assurance")
	status := c.QueryParam("status")

	type Filter struct {
		ID              []string `json:"id"`
		Name            string   `json:"name" condition:"LIKE"`
		ServiceID       string   `json:"service_id"`
		MinTimeSpan     string   `json:"min_timespan"`
		MaxTimeSpan     string   `json:"max_timespan"`
		Interest        string   `json:"interest" condition:"LIKE"`
		MinLoan         string   `json:"min_loan"`
		MaxLoan         string   `json:"max_loan"`
		Fees            string   `json:"fees" condition:"LIKE"`
		Collaterals     string   `json:"collaterals" condition:"LIKE"`
		FinancingSector string   `json:"financing_sector" condition:"LIKE"`
		Assurance       string   `json:"assurance" condition:"LIKE"`
		Status          string   `json:"status" condition:"LIKE"`
	}

	product := models.Product{}
	result, err := product.PagedFindFilter(page, rows, order, sort, &Filter{
		ID:              id,
		Name:            name,
		ServiceID:       serviceID,
		MinTimeSpan:     minTimespan,
		MaxTimeSpan:     maxTimespan,
		Interest:        interest,
		MinLoan:         minLoan,
		MaxLoan:         maxLoan,
		Fees:            fee,
		Collaterals:     collaterals,
		FinancingSector: financingSector,
		Assurance:       assurance,
		Status:          status,
	})
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "pencarian tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

func ProductNew(c echo.Context) error {
	defer c.Request().Body.Close()

	product := models.Product{}

	payloadRules := govalidator.MapData{
		"name":             []string{"required"},
		"service_id":       []string{"required", "valid_id:services"},
		"min_timespan":     []string{"required", "numeric"},
		"max_timespan":     []string{"required", "numeric"},
		"interest":         []string{"required", "numeric"},
		"min_loan":         []string{"required", "numeric"},
		"max_loan":         []string{"required", "numeric"},
		"fees":             []string{},
		"collaterals":      []string{},
		"financing_sector": []string{},
		"assurance":        []string{},
		"status":           []string{"required", "active_inactive"},
	}

	validate := validateRequestPayload(c, payloadRules, &product)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err := product.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat layanan bank baru")
	}

	return c.JSON(http.StatusCreated, product)
}

func ProductDetail(c echo.Context) error {
	defer c.Request().Body.Close()

	productId, _ := strconv.Atoi(c.Param("id"))

	product := models.Product{}
	err := product.FindbyID(productId)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("product %v tidak ditemukan", productId))
	}

	return c.JSON(http.StatusOK, product)
}

func ProductPatch(c echo.Context) error {
	defer c.Request().Body.Close()

	productId, _ := strconv.Atoi(c.Param("id"))

	product := models.Product{}
	err := product.FindbyID(productId)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("product %v tidak ditemukan", productId))
	}

	payloadRules := govalidator.MapData{
		"name":             []string{},
		"service_id":       []string{"valid_id:services"},
		"min_timespan":     []string{"numeric"},
		"max_timespan":     []string{"numeric"},
		"interest":         []string{"numeric"},
		"min_loan":         []string{"numeric"},
		"max_loan":         []string{"numeric"},
		"fees":             []string{},
		"collaterals":      []string{},
		"financing_sector": []string{},
		"assurance":        []string{},
		"status":           []string{"active_inactive"},
	}
	validate := validateRequestPayload(c, payloadRules, &product)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err = product.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update layanan %v", productId))
	}

	return c.JSON(http.StatusOK, product)
}

func ProductDelete(c echo.Context) error {
	defer c.Request().Body.Close()

	productId, _ := strconv.Atoi(c.Param("id"))

	product := models.Product{}
	err := product.FindbyID(productId)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("product %v tidak ditemukan", productId))
	}

	err = product.Delete()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank tipe %v", productId))
	}

	return c.JSON(http.StatusOK, product)
}
