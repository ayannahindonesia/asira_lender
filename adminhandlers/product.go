package adminhandlers

import (
	"asira_lender/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ayannahindonesia/basemodel"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/lib/pq"
	"github.com/thedevsaddam/govalidator"
)

// ProductPayload handles product request body
type ProductPayload struct {
	Name            string         `json:"name"`
	ServiceID       uint64         `json:"service_id"`
	MinTimeSpan     int            `json:"min_timespan"`
	MaxTimeSpan     int            `json:"max_timespan"`
	Interest        float64        `json:"interest"`
	MinLoan         int            `json:"min_loan"`
	MaxLoan         int            `json:"max_loan"`
	Fees            postgres.Jsonb `json:"fees"`
	Collaterals     pq.StringArray `json:"collaterals"`
	FinancingSector pq.StringArray `json:"financing_sector"`
	Assurance       string         `json:"assurance"`
	Status          string         `json:"status"`
}

// ProductList get all product list
func ProductList(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_product_list")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	order := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")

	var (
		product models.Product
		result  basemodel.PagedFindResult
	)

	if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
		type Filter struct {
			ID              int64  `json:"id" condition:"optional"`
			Name            string `json:"name" condition:"LIKE,optional"`
			Interest        string `json:"interest" condition:"LIKE,optional"`
			Fees            string `json:"fees" condition:"LIKE,optional"`
			Collaterals     string `json:"collaterals" condition:"LIKE,optional"`
			FinancingSector string `json:"financing_sector" condition:"LIKE,optional"`
			Assurance       string `json:"assurance" condition:"LIKE,optional"`
			Status          string `json:"status" condition:"LIKE,optional"`
		}
		id, _ := strconv.ParseInt(searchAll, 10, 64)
		result, err = product.PagedFindFilter(page, rows, order, sort, &Filter{
			ID:              id,
			Name:            searchAll,
			Interest:        searchAll,
			Fees:            searchAll,
			Collaterals:     searchAll,
			FinancingSector: searchAll,
			Assurance:       searchAll,
			Status:          searchAll,
		})
	} else {
		type Filter struct {
			ID              []string `json:"id"`
			Name            string   `json:"name" condition:"LIKE"`
			ServiceID       []string `json:"service_id"`
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
		result, err = product.PagedFindFilter(page, rows, order, sort, &Filter{
			ID:              customSplit(c.QueryParam("id"), ","),
			Name:            c.QueryParam("name"),
			ServiceID:       customSplit(c.QueryParam("service_id"), ","),
			MinTimeSpan:     c.QueryParam("min_timespan"),
			MaxTimeSpan:     c.QueryParam("max_timespan"),
			Interest:        c.QueryParam("interest"),
			MinLoan:         c.QueryParam("min_loan"),
			MaxLoan:         c.QueryParam("max_loan"),
			Fees:            c.QueryParam("fee"),
			Collaterals:     c.QueryParam("collaterals"),
			FinancingSector: c.QueryParam("financing_sector"),
			Assurance:       c.QueryParam("assurance"),
			Status:          c.QueryParam("status"),
		})
	}

	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "pencarian tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

// ProductNew add new product
func ProductNew(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_product_new")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	product := models.Product{}
	productPayload := ProductPayload{}

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

	validate := validateRequestPayload(c, payloadRules, &productPayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	marshal, _ := json.Marshal(productPayload)
	json.Unmarshal(marshal, &product)

	err = product.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat layanan bank baru")
	}

	return c.JSON(http.StatusCreated, product)
}

// ProductDetail get product detail by id
func ProductDetail(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_product_detail")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	productID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	product := models.Product{}
	err = product.FindbyID(productID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("product %v tidak ditemukan", productID))
	}

	return c.JSON(http.StatusOK, product)
}

// ProductPatch edit product by id
func ProductPatch(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_product_patch")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	productID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	product := models.Product{}
	productPayload := ProductPayload{}
	err = product.FindbyID(productID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("product %v tidak ditemukan", productID))
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
	validate := validateRequestPayload(c, payloadRules, &productPayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	if len(productPayload.Name) > 0 {
		product.Name = productPayload.Name
	}
	if productPayload.ServiceID > 0 {
		product.ServiceID = productPayload.ServiceID
	}
	if productPayload.MinTimeSpan > 0 {
		product.MinTimeSpan = productPayload.MinTimeSpan
	}
	if productPayload.MaxTimeSpan > 0 {
		product.MaxTimeSpan = productPayload.MaxTimeSpan
	}
	if productPayload.Interest > 0 {
		product.Interest = productPayload.Interest
	}
	if productPayload.MinLoan > 0 {
		product.MinLoan = productPayload.MinLoan
	}
	if productPayload.MaxLoan > 0 {
		product.MaxLoan = productPayload.MaxLoan
	}
	if len(string(productPayload.Fees.RawMessage)) > 2 {
		product.Fees = productPayload.Fees
	}
	if len(productPayload.Collaterals) > 0 {
		product.Collaterals = pq.StringArray(productPayload.Collaterals)
	}
	if len(productPayload.FinancingSector) > 0 {
		product.FinancingSector = pq.StringArray(productPayload.FinancingSector)
	}
	if len(productPayload.Assurance) > 0 {
		product.Assurance = productPayload.Assurance
	}
	if len(productPayload.Status) > 0 {
		product.Status = productPayload.Status
	}

	err = product.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update layanan %v", productID))
	}

	return c.JSON(http.StatusOK, product)
}

// ProductDelete delete product
func ProductDelete(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_product_delete")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	productID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	product := models.Product{}
	err = product.FindbyID(productID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("product %v tidak ditemukan", productID))
	}

	err = product.Delete()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank tipe %v", productID))
	}

	return c.JSON(http.StatusOK, product)
}
