package admin_handlers

import (
	"asira_lender/models"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

func BankTypeList(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_bank_type_list")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

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

	bank_type := models.BankType{}
	result, err := bank_type.PagedFilterSearch(page, rows, orderby, sort, &Filter{
		Name: name,
	})
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "pencarian tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

func BankTypeNew(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_bank_type_new")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	bank_type := models.BankType{}

	payloadRules := govalidator.MapData{
		"name": []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &bank_type)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err = bank_type.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat tipe bank baru")
	}

	return c.JSON(http.StatusCreated, bank_type)
}

func BankTypeDetail(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_bank_type_detail")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	bank_id, _ := strconv.Atoi(c.Param("bank_id"))

	bankType := models.BankType{}
	err = bankType.FindbyID(bank_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank type %v tidak ditemukan", bank_id))
	}

	return c.JSON(http.StatusOK, bankType)
}

func BankTypePatch(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_bank_type_patch")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	bank_id, _ := strconv.Atoi(c.Param("bank_id"))

	bankType := models.BankType{}
	err = bankType.FindbyID(bank_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank type %v tidak ditemukan", bank_id))
	}

	payloadRules := govalidator.MapData{
		"name": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &bankType)
	log.Println(bankType)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err = bankType.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank tipe %v", bank_id))
	}

	return c.JSON(http.StatusOK, bankType)
}

func BankTypeDelete(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_bank_type_delete")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	bank_id, _ := strconv.Atoi(c.Param("bank_id"))

	bankType := models.BankType{}
	err = bankType.FindbyID(bank_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank type %v tidak ditemukan", bank_id))
	}

	err = bankType.Delete()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank tipe %v", bank_id))
	}

	return c.JSON(http.StatusOK, bankType)
}
