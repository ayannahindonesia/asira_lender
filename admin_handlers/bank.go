package admin_handlers

import (
	"asira_lender/asira"
	"asira_lender/models"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/lib/pq"
	"github.com/thedevsaddam/govalidator"
	"gitlab.com/asira-ayannah/basemodel"
)

type (
	// BankSelect for custom query
	BankSelect struct {
		models.Bank
		BankTypeName string `json:"bank_type_name"`
	}
	// BankPayload request body container
	BankPayload struct {
		Name                string  `json:"name"`
		Type                uint64  `json:"type"`
		Address             string  `json:"address"`
		Province            string  `json:"province"`
		City                string  `json:"city"`
		PIC                 string  `json:"pic"`
		Phone               string  `json:"phone"`
		Services            []int64 `json:"services"`
		Products            []int64 `json:"products"`
		AdminFeeSetup       string  `json:"adminfee_setup"`
		ConvenienceFeeSetup string  `json:"convfee_setup"`
	}
)

// BankList get all bank list
func BankList(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_bank_list")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	db := asira.App.DB
	var (
		totalRows int
		offset    int
		rows      int
		page      int
		banks     []BankSelect
	)

	// pagination parameters
	if c.QueryParam("rows") != "all" {
		rows, _ = strconv.Atoi(c.QueryParam("rows"))
		page, _ = strconv.Atoi(c.QueryParam("page"))
		if page <= 0 {
			page = 1
		}
		if rows <= 0 {
			rows = 25
		}
		offset = (page * rows) - rows
	}

	db = db.Table("banks b").
		Select("b.*, bt.name as bank_type_name").
		Joins("INNER JOIN bank_types bt ON b.type = bt.id")

	if name := c.QueryParam("name"); len(name) > 0 {
		db = db.Where("b.name LIKE ?", name)
	}
	if id := customSplit(c.QueryParam("id"), ","); len(id) > 0 {
		db = db.Where("b.id IN (?)", id)
	}

	if order := strings.Split(c.QueryParam("orderby"), ","); len(order) > 0 {
		if sort := strings.Split(c.QueryParam("sort"), ","); len(sort) > 0 {
			for k, v := range order {
				q := v
				if len(sort) > k {
					value := sort[k]
					if strings.ToUpper(value) == "ASC" || strings.ToUpper(value) == "DESC" {
						q = v + " " + strings.ToUpper(value)
					}
				}
				db = db.Order(q)
			}
		}
	}

	if rows > 0 && offset > 0 {
		db = db.Limit(rows).Offset(offset)
	}
	err = db.Find(&banks).Count(&totalRows).Error
	if err != nil {
		log.Println(err)
	}

	lastPage := int(math.Ceil(float64(totalRows) / float64(rows)))

	result := basemodel.PagedFindResult{
		TotalData:   totalRows,
		Rows:        rows,
		CurrentPage: page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + rows,
		Data:        banks,
	}

	return c.JSON(http.StatusOK, result)
}

// BankNew create new bank
func BankNew(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_bank_new")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	bank := models.Bank{}
	bankPayload := BankPayload{}

	payloadRules := govalidator.MapData{
		"name":           []string{"required"},
		"type":           []string{"required", "valid_id:bank_types"},
		"address":        []string{"required"},
		"province":       []string{"required"},
		"city":           []string{"required"},
		"services":       []string{"required", "valid_id:services"},
		"products":       []string{"required", "valid_id:products"},
		"pic":            []string{"required"},
		"phone":          []string{"required"},
		"adminfee_setup": []string{"required"},
		"convfee_setup":  []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &bankPayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	marshal, _ := json.Marshal(bankPayload)
	json.Unmarshal(marshal, &bank)

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

// BankDetail get bank detail by id
func BankDetail(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_bank_detail")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	db := asira.App.DB

	bankID, _ := strconv.Atoi(c.Param("bank_id"))

	db = db.Table("banks b").
		Select("b.*, bt.name as bank_type_name").
		Joins("INNER JOIN bank_types bt ON b.type = bt.id").
		Where("b.id = ?", bankID)

	bank := BankSelect{}
	err = db.Find(&bank).Error
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank type %v tidak ditemukan", bankID))
	}

	return c.JSON(http.StatusOK, bank)
}

// BankPatch edit bank by id
func BankPatch(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_bank_patch")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	bankID, _ := strconv.Atoi(c.Param("bank_id"))

	bank := models.Bank{}
	bankPayload := BankPayload{}
	err = bank.FindbyID(bankID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank %v tidak ditemukan", bankID))
	}

	payloadRules := govalidator.MapData{
		"name":           []string{},
		"type":           []string{"valid_id:bank_types"},
		"address":        []string{},
		"province":       []string{},
		"city":           []string{},
		"services":       []string{"valid_id:services"},
		"products":       []string{"valid_id:products"},
		"pic":            []string{},
		"phone":          []string{},
		"adminfee_setup": []string{},
		"convfee_setup":  []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &bankPayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	if len(bankPayload.Name) > 0 {
		bank.Name = bankPayload.Name
	}
	if bankPayload.Type > 0 {
		bank.Type = bankPayload.Type
	}
	if len(bankPayload.Address) > 0 {
		bank.Address = bankPayload.Address
	}
	if len(bankPayload.Province) > 0 {
		bank.Province = bankPayload.Province
	}
	if len(bankPayload.City) > 0 {
		bank.City = bankPayload.City
	}
	if len(bankPayload.Services) > 0 {
		bank.Services = pq.Int64Array(bankPayload.Services)
	}
	if len(bankPayload.Products) > 0 {
		bank.Products = pq.Int64Array(bankPayload.Products)
	}
	if len(bankPayload.PIC) > 0 {
		bank.PIC = bankPayload.PIC
	}
	if len(bankPayload.Phone) > 0 {
		bank.Phone = bankPayload.Phone
	}
	if len(bankPayload.AdminFeeSetup) > 0 {
		bank.AdminFeeSetup = bankPayload.AdminFeeSetup
	}
	if len(bankPayload.ConvenienceFeeSetup) > 0 {
		bank.ConvenienceFeeSetup = bankPayload.ConvenienceFeeSetup
	}

	err = bank.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank %v", bankID))
	}

	return c.JSON(http.StatusOK, bank)
}

// BankDelete delete bank
func BankDelete(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_bank_delete")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	bankID, _ := strconv.Atoi(c.Param("bank_id"))

	bank := models.Bank{}
	err = bank.FindbyID(bankID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank type %v tidak ditemukan", bankID))
	}

	err = bank.Delete()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank tipe %v", bankID))
	}

	return c.JSON(http.StatusOK, bank)
}
