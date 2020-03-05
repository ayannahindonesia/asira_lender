package adminhandlers

import (
	"asira_lender/middlewares"
	"asira_lender/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ayannahindonesia/basemodel"
	"github.com/dgrijalva/jwt-go"
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
	InterestType    string         `json:"interest_type"`
	MinLoan         int            `json:"min_loan"`
	MaxLoan         int            `json:"max_loan"`
	Fees            postgres.Jsonb `json:"fees"`
	Collaterals     pq.StringArray `json:"collaterals"`
	FinancingSector pq.StringArray `json:"financing_sector"`
	Assurance       string         `json:"assurance"`
	Status          string         `json:"status"`
	Form            postgres.Jsonb `json:"form"`
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
			InterestType    string `json:"interest_type" condition:"LIKE,optional"`
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
			InterestType:    searchAll,
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
			InterestType    string   `json:"interest_type"`
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
			InterestType:    c.QueryParam("interest_type"),
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
		NLog("warning", "ProductList", map[string]interface{}{"message": "error listing products", "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Pencarian tidak ditemukan")
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
		"interest_type":    []string{"required"},
		"min_loan":         []string{"required", "numeric"},
		"max_loan":         []string{"required", "numeric"},
		"fees":             []string{},
		"collaterals":      []string{},
		"financing_sector": []string{},
		"assurance":        []string{},
		"status":           []string{"required", "active_inactive"},
		"form":             []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &productPayload)
	if validate != nil {
		NLog("warning", "ProductNew", map[string]interface{}{"message": "validation error", "error": validate}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi.")
	}

	marshal, _ := json.Marshal(productPayload)
	json.Unmarshal(marshal, &product)

	err = product.Create()
	if err != nil {
		NLog("error", "ProductNew", map[string]interface{}{"message": "create product error", "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat produk baru")
	}

	err = middlewares.SubmitKafkaPayload(product, "product_create")
	if err != nil {
		NLog("error", "ProductNew", map[string]interface{}{"message": "kafka submit error", "error": err, "product": product}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat produk baru")
	}

	NAudittrail(models.Product{}, product, c.Get("user").(*jwt.Token), "product", fmt.Sprint(product.ID), "create")

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
		NLog("warning", "ProductDetail", map[string]interface{}{"message": fmt.Sprintf("find product %v error", productID), "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Product %v tidak ditemukan", productID))
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
		NLog("warning", "ProductPatch", map[string]interface{}{"message": fmt.Sprintf("patch product %v", productID), "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Product %v tidak ditemukan", productID))
	}
	origin := product

	payloadRules := govalidator.MapData{
		"name":             []string{},
		"service_id":       []string{"valid_id:services"},
		"min_timespan":     []string{"numeric"},
		"max_timespan":     []string{"numeric"},
		"interest":         []string{"numeric"},
		"interest_type":    []string{},
		"min_loan":         []string{"numeric"},
		"max_loan":         []string{"numeric"},
		"fees":             []string{},
		"collaterals":      []string{},
		"financing_sector": []string{},
		"assurance":        []string{},
		"status":           []string{"active_inactive"},
		"form":             []string{},
	}
	validate := validateRequestPayload(c, payloadRules, &productPayload)
	if validate != nil {
		NLog("warning", "ProductPatch", map[string]interface{}{"message": "validation error", "error": validate}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi.")
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
	if len(productPayload.InterestType) > 0 {
		product.InterestType = productPayload.InterestType
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
	if len(string(productPayload.Form.RawMessage)) > 2 {
		product.Form = productPayload.Form
	}

	err = middlewares.SubmitKafkaPayload(product, "product_update")
	if err != nil {
		NLog("error", "ProductPatch", map[string]interface{}{"message": fmt.Sprintf("kafka submit error on product %v", productID), "error": err, "product": product}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update produk %v", productID))
	}

	NAudittrail(origin, product, c.Get("user").(*jwt.Token), "product", fmt.Sprint(product.ID), "update")

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
		NLog("warning", "ProductDelete", map[string]interface{}{"message": fmt.Sprintf("error finding product %v", productID), "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Product %v tidak ditemukan", productID))
	}

	err = middlewares.SubmitKafkaPayload(product, "product_delete")
	if err != nil {
		NLog("error", "ProductDelete", map[string]interface{}{"message": fmt.Sprintf("kafka submit error product %v", productID), "error": err, "product": product}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal delete produk %v", productID))
	}

	NAudittrail(product, models.Product{}, c.Get("user").(*jwt.Token), "product", fmt.Sprint(product.ID), "delete")

	return c.JSON(http.StatusOK, product)
}
