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

type BankServicePayload struct {
	ServiceID string `json:"service_id"`
	BankID    string `json:"bank_id"`
}

func BankServiceList(c echo.Context) error {
	defer c.Request().Body.Close()

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	order := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")

	// filters
	bankID := customSplit(c.QueryParam("bank_id"), ",")
	serviceID := customSplit(c.QueryParam("service_id"), ",")

	type Filter struct {
		BankID    []string `json:"bank_id"`
		ServiceID []string `json:"service_id"`
	}

	bank_service := models.BankService{}
	result, err := bank_service.PagedFindFilter(page, rows, order, sort, &Filter{
		BankID:    bankID,
		ServiceID: serviceID,
	})
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "pencarian tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

func BankServiceNew(c echo.Context) error {
	defer c.Request().Body.Close()

	bankServicePayload := BankServicePayload{}

	payloadRules := govalidator.MapData{
		"service_id": []string{"required", "valid_id:services"},
		"bank_id":    []string{"required", "valid_id:banks"},
	}

	validate := validateRequestPayload(c, payloadRules, &bankServicePayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	uint64ServiceID, _ := strconv.ParseUint(bankServicePayload.ServiceID, 10, 64)
	uint64BankID, _ := strconv.ParseUint(bankServicePayload.BankID, 10, 64)
	bankService := models.BankService{
		ServiceID: uint64ServiceID,
		BankID:    uint64BankID,
	}

	err := bankService.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat layanan bank baru")
	}

	return c.JSON(http.StatusCreated, bankService)
}

func BankServiceDetail(c echo.Context) error {
	defer c.Request().Body.Close()

	bankServiceID, _ := strconv.Atoi(c.Param("id"))

	bankService := models.BankService{}
	err := bankService.FindbyID(bankServiceID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("layanan %v tidak ditemukan", bankServiceID))
	}

	return c.JSON(http.StatusOK, bankService)
}

func BankServicePatch(c echo.Context) error {
	defer c.Request().Body.Close()

	bankServiceID, _ := strconv.Atoi(c.Param("id"))

	bankService := models.BankService{}
	err := bankService.FindbyID(bankServiceID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("layanan %v tidak ditemukan", bankServiceID))
	}

	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("layanan %v tidak ditemukan", bankServiceID))
	}
	servicePayloadRules := govalidator.MapData{
		"service_id": []string{"valid_id:services"},
		"bank_id":    []string{"valid_id:banks"},
	}
	validate := validateRequestPayload(c, servicePayloadRules, &bankService)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err = bankService.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update layanan %v", bankServiceID))
	}

	return c.JSON(http.StatusOK, bankService)
}

func BankServiceDelete(c echo.Context) error {
	defer c.Request().Body.Close()

	bankServiceID, _ := strconv.Atoi(c.Param("id"))

	bankService := models.BankService{}
	err := bankService.FindbyID(bankServiceID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank type %v tidak ditemukan", bankServiceID))
	}

	err = bankService.Delete()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank tipe %v", bankServiceID))
	}

	return c.JSON(http.StatusOK, bankService)
}
