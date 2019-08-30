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

func BankList(c echo.Context) error {
	defer c.Request().Body.Close()

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	// filters
	name := c.QueryParam("name")
	id := c.QueryParam("id")

	type Filter struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	bank := models.Bank{}
	result, err := bank.PagedFilterSearch(page, rows, orderby, sort, &Filter{
		ID:   id,
		Name: name,
	})
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "pencarian tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

func BankNew(c echo.Context) error {
	defer c.Request().Body.Close()

	bank := models.Bank{}

	payloadRules := govalidator.MapData{
		"name":     []string{"required"},
		"type":     []string{"bank_type_id"},
		"address":  []string{"required"},
		"province": []string{"required"},
		"city":     []string{"required"},
		"services": []string{},
		"products": []string{},
		"pic":      []string{"required"},
		"phone":    []string{"required"},
		// "username": []string{"required", "unique:banks,username"},
		// "password": []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &bank)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	newBankType, err := bank.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat tipe bank baru")
	}

	return c.JSON(http.StatusCreated, newBankType)
}

func BankDetail(c echo.Context) error {
	defer c.Request().Body.Close()

	bank_id, _ := strconv.Atoi(c.Param("bank_id"))

	bank := models.Bank{}
	result, err := bank.FindbyID(bank_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank type %v tidak ditemukan", bank_id))
	}

	return c.JSON(http.StatusOK, result)
}

func BankPatch(c echo.Context) error {
	defer c.Request().Body.Close()

	bank_id, _ := strconv.Atoi(c.Param("bank_id"))

	bank := models.Bank{}
	result, err := bank.FindbyID(bank_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank type %v tidak ditemukan", bank_id))
	}

	// dont allow admin to change bank credentials
	tempUsername := result.Username
	tempPassword := result.Password

	payloadRules := govalidator.MapData{
		"name":     []string{},
		"type":     []string{"bank_type_id"},
		"address":  []string{},
		"province": []string{},
		"city":     []string{},
		"services": []string{},
		"products": []string{},
		"pic":      []string{},
		"phone":    []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &result)
	log.Println(result)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	result.Username = tempUsername
	result.Password = tempPassword

	_, err = result.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank tipe %v", bank_id))
	}

	return c.JSON(http.StatusOK, result)
}

func BankDelete(c echo.Context) error {
	defer c.Request().Body.Close()

	bank_id, _ := strconv.Atoi(c.Param("bank_id"))

	bank := models.Bank{}
	result, err := bank.FindbyID(bank_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank type %v tidak ditemukan", bank_id))
	}

	_, err = result.Delete()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank tipe %v", bank_id))
	}

	return c.JSON(http.StatusOK, result)
}
