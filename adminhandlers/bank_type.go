package adminhandlers

import (
	"asira_lender/middlewares"
	"asira_lender/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

// BankTypePayload to handle post and patch
type BankTypePayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// BankTypeList lists all bank type
func BankTypeList(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_bank_type_list")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")

	// filters
	name := c.QueryParam("name")

	type Filter struct {
		Name string `json:"name" condition:"LIKE"`
	}

	bankType := models.BankType{}
	result, err := bankType.PagedFindFilter(page, rows, orderby, sort, &Filter{
		Name: name,
	})
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Pencarian tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

// BankTypeNew add new bank type
func BankTypeNew(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_bank_type_new")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	bankType := models.BankType{}
	bankTypePayload := BankTypePayload{}

	payloadRules := govalidator.MapData{
		"name":        []string{"required"},
		"description": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &bankTypePayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
	}

	marshal, _ := json.Marshal(bankTypePayload)
	json.Unmarshal(marshal, &bankType)

	err = middlewares.SubmitKafkaPayload(bankType, "bank_type_create")
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat tipe bank baru")
	}

	return c.JSON(http.StatusCreated, bankType)
}

// BankTypeDetail get bank type detail by id
func BankTypeDetail(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_bank_type_detail")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	bankID, _ := strconv.ParseUint(c.Param("bank_id"), 10, 64)

	bankType := models.BankType{}
	err = bankType.FindbyID(bankID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Tidak memiliki hak akses")
	}

	return c.JSON(http.StatusOK, bankType)
}

// BankTypePatch edit bank type by id
func BankTypePatch(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_bank_type_patch")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	bankID, _ := strconv.ParseUint(c.Param("bank_id"), 10, 64)

	bankType := models.BankType{}
	bankTypePayload := BankTypePayload{}
	err = bankType.FindbyID(bankID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Bank type %v tidak ditemukan", bankID))
	}

	payloadRules := govalidator.MapData{
		"name":        []string{},
		"description": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &bankTypePayload)
	log.Println(bankType)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
	}

	if len(bankTypePayload.Name) > 0 {
		bankType.Name = bankTypePayload.Name
	}
	if len(bankTypePayload.Description) > 0 {
		bankType.Description = bankTypePayload.Description
	}

	err = middlewares.SubmitKafkaPayload(bankType, "bank_type_update")
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank tipe %v", bankID))
	}

	return c.JSON(http.StatusOK, bankType)
}

// BankTypeDelete delete bank type by id
func BankTypeDelete(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_bank_type_delete")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	bankID, _ := strconv.ParseUint(c.Param("bank_id"), 10, 64)

	bankType := models.BankType{}
	err = bankType.FindbyID(bankID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Bank type %v tidak ditemukan", bankID))
	}

	err = middlewares.SubmitKafkaPayload(bankType, "bank_type_delete")
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank tipe %v", bankID))
	}

	return c.JSON(http.StatusOK, bankType)
}
