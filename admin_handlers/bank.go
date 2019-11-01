package admin_handlers

import (
	"asira_lender/asira"
	"asira_lender/models"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
	"gitlab.com/asira-ayannah/basemodel"
)

type BankSelect struct {
	models.Bank
	BankTypeName string `json:"bank_type_name"`
}

func BankList(c echo.Context) error {
	defer c.Request().Body.Close()

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
	err := db.Find(&banks).Count(&totalRows).Error
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

func BankNew(c echo.Context) error {
	defer c.Request().Body.Close()

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

	err := bank.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat bank baru")
	}

	return c.JSON(http.StatusCreated, bank)
}

func BankDetail(c echo.Context) error {
	defer c.Request().Body.Close()

	db := asira.App.DB

	bank_id, _ := strconv.Atoi(c.Param("bank_id"))

	db = db.Table("banks b").
		Select("b.*, bt.name as bank_type_name").
		Joins("INNER JOIN bank_types bt ON b.type = bt.id").
		Where("b.id = ?", bank_id)

	bank := BankSelect{}
	err := db.Find(&bank).Error
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank type %v tidak ditemukan", bank_id))
	}

	return c.JSON(http.StatusOK, bank)
}

func BankPatch(c echo.Context) error {
	defer c.Request().Body.Close()

	var err error

	bank_id, _ := strconv.Atoi(c.Param("bank_id"))

	bank := models.Bank{}
	err = bank.FindbyID(bank_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank %v tidak ditemukan", bank_id))
	}

	// dont allow admin to change bank credentials
	tempUsername := bank.Username
	tempPassword := bank.Password

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

	bank.Username = tempUsername
	bank.Password = tempPassword

	err = bank.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank %v", bank_id))
	}

	return c.JSON(http.StatusOK, bank)
}

func BankDelete(c echo.Context) error {
	defer c.Request().Body.Close()

	bank_id, _ := strconv.Atoi(c.Param("bank_id"))

	bank := models.Bank{}
	err := bank.FindbyID(bank_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank type %v tidak ditemukan", bank_id))
	}

	err = bank.Delete()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank tipe %v", bank_id))
	}

	return c.JSON(http.StatusOK, bank)
}
