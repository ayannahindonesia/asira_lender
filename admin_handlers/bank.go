package admin_handlers

import (
	"asira_lender/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/lib/pq"
	"github.com/thedevsaddam/govalidator"
)

// type BankPayload struct {
// 	Name                string        `json:"name"`
// 	Type                uint64        `json:"type"`
// 	Address             string        `json:"address"`
// 	Province            string        `json:"province"`
// 	City                string        `json:"city"`
// 	Services            pq.Int64Array `json:"services"`
// 	Products            pq.Int64Array `json:"products"`
// 	PIC                 string        `json:"pic"`
// 	Phone               string        `json:"phone"`
// 	AdminFeeSetup       string        `json:"adminfee_setup"`
// 	ConvenienceFeeSetup string        `json:"convfee_setup"`
// 	Username            string        `json:"username"`
// 	Password            string        `json:"password"`
// }

func BankList(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_bank_list")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	order := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")

	// filters
	name := c.QueryParam("name")
	id := customSplit(c.QueryParam("id"), ",")

	type Filter struct {
		ID   []string `json:"id"`
		Name string   `json:"name" condition:"LIKE"`
	}

	bank := models.Bank{}
	result, err := bank.PagedFindFilter(page, rows, order, sort, &Filter{
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
	err := validatePermission(c, "core_bank_new")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	bank := models.Bank{}

	payloadRules := govalidator.MapData{
		"name":           []string{"required"},
		"type":           []string{"required", "valid_id:bank_types"},
		"address":        []string{"required"},
		"province":       []string{"required"},
		"city":           []string{"required"},
		"services":       []string{"required"},
		"products":       []string{"required"},
		"pic":            []string{"required"},
		"phone":          []string{"required"},
		"adminfee_setup": []string{"required"},
		"convfee_setup":  []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &bank)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err = bank.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat bank baru")
	}

	// @ToDo remodel this flow
	user := models.User{
		Username: bank.Name,
		Roles:    pq.Int64Array{3},
		Phone:    bank.Phone,
		Status:   "active",
	}
	err = user.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat user")
	}

	bankRep := models.BankRepresentatives{
		UserID: user.ID,
		BankID: bank.ID,
	}
	err = bankRep.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat bank representative")
	}
	// ------

	return c.JSON(http.StatusCreated, bank)
}

func BankDetail(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_bank_detail")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	bank_id, _ := strconv.Atoi(c.Param("bank_id"))

	bank := models.Bank{}
	err = bank.FindbyID(bank_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank type %v tidak ditemukan", bank_id))
	}

	return c.JSON(http.StatusOK, bank)
}

func BankPatch(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_bank_patch")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	bank_id, _ := strconv.Atoi(c.Param("bank_id"))

	bank := models.Bank{}
	err = bank.FindbyID(bank_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank %v tidak ditemukan", bank_id))
	}

	payloadRules := govalidator.MapData{
		"name":           []string{},
		"type":           []string{"valid_id:bank_types"},
		"address":        []string{},
		"province":       []string{},
		"city":           []string{},
		"services":       []string{},
		"products":       []string{},
		"pic":            []string{},
		"phone":          []string{},
		"adminfee_setup": []string{},
		"convfee_setup":  []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &bank)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err = bank.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank %v", bank_id))
	}

	return c.JSON(http.StatusOK, bank)
}

func BankDelete(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_bank_delete")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	bank_id, _ := strconv.Atoi(c.Param("bank_id"))

	bank := models.Bank{}
	err = bank.FindbyID(bank_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank type %v tidak ditemukan", bank_id))
	}

	err = bank.Delete()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank tipe %v", bank_id))
	}

	return c.JSON(http.StatusOK, bank)
}
