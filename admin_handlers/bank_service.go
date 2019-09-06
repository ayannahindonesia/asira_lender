package admin_handlers

import (
	"asira_lender/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

type BankServicePayload struct {
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
}

func BankServiceList(c echo.Context) error {
	defer c.Request().Body.Close()

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	// filters
	name := c.QueryParam("name")

	type Filter struct {
		Name string `json:"name" condition:"LIKE"`
	}

	bank_service := models.BankService{}
	result, err := bank_service.PagedFilterSearch(page, rows, orderby, sort, &Filter{
		Name: name,
	})
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "pencarian tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

func BankServiceNew(c echo.Context) error {
	defer c.Request().Body.Close()

	bank_service_payload := BankServicePayload{}

	payloadRules := govalidator.MapData{
		"name":        []string{"required"},
		"image":       []string{"required"},
		"status":      []string{"required", "active_inactive"},
		"description": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &bank_service_payload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	image := models.Image{
		Image_string: bank_service_payload.Image,
	}
	image.Create()

	bankService := models.BankService{
		Name:    bank_service_payload.Name,
		ImageID: int(image.ID),
		Status:  bank_service_payload.Status,
	}

	newBankService, err := bankService.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat layanan bank baru")
	}

	return c.JSON(http.StatusCreated, newBankService)
}

func BankServiceDetail(c echo.Context) error {
	defer c.Request().Body.Close()

	bank_service_id, _ := strconv.Atoi(c.Param("bank_service_id"))

	bankService := models.BankService{}
	result, err := bankService.FindbyID(bank_service_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("layanan %v tidak ditemukan", bank_service_id))
	}

	return c.JSON(http.StatusOK, result)
}

func BankServicePatch(c echo.Context) error {
	defer c.Request().Body.Close()

	bank_service_id, _ := strconv.Atoi(c.Param("bank_service_id"))

	bankService := models.BankService{}
	bankServiceRes, err := bankService.FindbyID(bank_service_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("layanan %v tidak ditemukan", bank_service_id))
	}

	bankServiceImg := models.Image{}
	bankServiceImgRes, _ := bankServiceImg.FindbyID(bankServiceRes.ImageID)

	payloadBucket := BankServicePayload{}
	servicePayloadRules := govalidator.MapData{
		"name":        []string{},
		"image":       []string{},
		"status":      []string{"active_inactive"},
		"description": []string{},
	}
	validate := validateRequestPayload(c, servicePayloadRules, &payloadBucket)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	if len(payloadBucket.Name) > 0 {
		bankServiceRes.Name = payloadBucket.Name
	}
	if len(payloadBucket.Image) > 0 {
		bankServiceImgRes.Image_string = payloadBucket.Image
	}
	if len(payloadBucket.Status) > 0 {
		bankServiceRes.Status = payloadBucket.Status
	}

	_, err = bankServiceRes.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update layanan %v", bank_service_id))
	}
	_, err = bankServiceImgRes.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update layanan %v", bank_service_id))
	}

	return c.JSON(http.StatusOK, bankServiceRes)
}

func BankServiceDelete(c echo.Context) error {
	defer c.Request().Body.Close()

	bank_service_id, _ := strconv.Atoi(c.Param("bank_service_id"))

	bankService := models.BankService{}
	result, err := bankService.FindbyID(bank_service_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank type %v tidak ditemukan", bank_service_id))
	}

	_, err = result.Delete()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank tipe %v", bank_service_id))
	}

	return c.JSON(http.StatusOK, result)
}
