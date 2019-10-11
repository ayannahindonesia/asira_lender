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

func BankProductList(c echo.Context) error {
	defer c.Request().Body.Close()

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	order := strings.Split(c.QueryParam("order"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")

	// filters
	productID := customSplit(c.QueryParam("product_id"), ",")
	bankID := customSplit(c.QueryParam("bank_id"), ",")

	type Filter struct {
		ProductID []string `json:"product_id"`
		BankID    []string `json:"bank_id"`
	}

	bank_product := models.BankProduct{}
	result, err := bank_product.PagedFindFilter(page, rows, order, sort, &Filter{
		ProductID: productID,
		BankID:    bankID,
	})
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "pencarian tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

func BankProductNew(c echo.Context) error {
	defer c.Request().Body.Close()

	bankProduct := models.BankProduct{}

	payloadRules := govalidator.MapData{
		"product_id": []string{"required", "valid_id:products"},
		"bank_id":    []string{"required", "valid_id:banks"},
	}

	validate := validateRequestPayload(c, payloadRules, &bankProduct)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err := bankProduct.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat layanan bank baru")
	}

	return c.JSON(http.StatusCreated, bankProduct)
}

func BankProductDetail(c echo.Context) error {
	defer c.Request().Body.Close()

	product_id, _ := strconv.Atoi(c.Param("id"))

	bankProduct := models.BankProduct{}
	err := bankProduct.FindbyID(product_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("product %v tidak ditemukan", bankProduct))
	}

	return c.JSON(http.StatusOK, bankProduct)
}

func BankProductPatch(c echo.Context) error {
	defer c.Request().Body.Close()

	product_id, _ := strconv.Atoi(c.Param("id"))

	bankProduct := models.BankProduct{}
	err := bankProduct.FindbyID(product_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("product %v tidak ditemukan", product_id))
	}

	productPayloadRules := govalidator.MapData{
		"product_id": []string{"required", "valid_id:products"},
		"bank_id":    []string{"required", "valid_id:banks"},
	}
	validate := validateRequestPayload(c, productPayloadRules, &bankProduct)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err = bankProduct.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update product %v", product_id))
	}

	return c.JSON(http.StatusOK, bankProduct)
}

func BankProductDelete(c echo.Context) error {
	defer c.Request().Body.Close()

	product_id, _ := strconv.Atoi(c.Param("id"))

	bankProduct := models.BankProduct{}
	err := bankProduct.FindbyID(product_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank type %v tidak ditemukan", product_id))
	}

	err = bankProduct.Delete()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank tipe %v", product_id))
	}

	return c.JSON(http.StatusOK, bankProduct)
}
