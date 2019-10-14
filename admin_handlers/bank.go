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

type BankPayload struct {
	Name          string   `json:"name"`
	Type          uint64   `json:"type"`
	Address       string   `json:"address"`
	Province      string   `json:"province"`
	City          string   `json:"city"`
	Services      []uint64 `json:"services"`
	Products      []uint64 `json:"products"`
	PIC           string   `json:"pic"`
	Phone         string   `json:"phone"`
	AdminFeeSetup string   `json:"adminfee_setup"`
	ConvFeeSetup  string   `json:"convfee_setup"`
}

func BankList(c echo.Context) error {
	defer c.Request().Body.Close()

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

	var v uint64

	bankPayload := BankPayload{}

	payloadRules := govalidator.MapData{
		"name":           []string{"required"},
		"type":           []string{"required", "valid_id:bank_types"},
		"address":        []string{"required"},
		"province":       []string{"required"},
		"city":           []string{"required"},
		"services":       []string{"required", "valid_id:services", "unique_distinct:bank_services,bank_id,service_id,0"},
		"products":       []string{"required", "valid_id:products", "unique_distinct:bank_products,bank_id,product_id,0"},
		"pic":            []string{"required"},
		"phone":          []string{"required"},
		"adminfee_setup": []string{"required"},
		"convfee_setup":  []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &bankPayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	bank := models.Bank{
		Name:                bankPayload.Name,
		Type:                bankPayload.Type,
		Address:             bankPayload.Address,
		Province:            bankPayload.Province,
		City:                bankPayload.City,
		PIC:                 bankPayload.PIC,
		Phone:               bankPayload.Phone,
		AdminFeeSetup:       bankPayload.AdminFeeSetup,
		ConvenienceFeeSetup: bankPayload.ConvFeeSetup,
	}
	err := bank.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat tipe bank baru")
	}

	for _, v = range bankPayload.Services {
		bankService := models.BankService{
			ServiceID: v,
			BankID:    bank.ID,
		}
		bankService.Create()
	}

	for _, v = range bankPayload.Products {
		bankProduct := models.BankProduct{
			ProductID: v,
			BankID:    bank.ID,
		}
		bankProduct.Create()
	}

	return c.JSON(http.StatusCreated, bank)
}

func BankDetail(c echo.Context) error {
	defer c.Request().Body.Close()

	bank_id, _ := strconv.Atoi(c.Param("bank_id"))

	bank := models.Bank{}
	err := bank.FindbyID(bank_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank type %v tidak ditemukan", bank_id))
	}

	return c.JSON(http.StatusOK, bank)
}

func BankPatch(c echo.Context) error {
	defer c.Request().Body.Close()

	var v uint64
	var err error

	bankPayload := BankPayload{}

	bank_id, _ := strconv.Atoi(c.Param("bank_id"))

	bank := models.Bank{}
	err = bank.FindbyID(bank_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank type %v tidak ditemukan", bank_id))
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
		"services":       []string{"valid_id:services", "unique_distinct:bank_services,bank_id,service_id,1"},
		"products":       []string{"valid_id:products", "unique_distinct:bank_products,bank_id,product_id,1"},
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
	if len(bankPayload.PIC) > 0 {
		bank.PIC = bankPayload.PIC
	}
	if len(bankPayload.Phone) > 0 {
		bank.Phone = bankPayload.Phone
	}
	if len(bankPayload.AdminFeeSetup) > 0 {
		bank.AdminFeeSetup = bankPayload.AdminFeeSetup
	}
	if len(bankPayload.ConvFeeSetup) > 0 {
		bank.ConvenienceFeeSetup = bankPayload.ConvFeeSetup
	}

	type Filter struct {
		BankID uint64 `json:"bank_id"`
	}
	if len(bankPayload.Services) > 0 {
		bankService := models.BankService{}
		bankServices, _ := bankService.FindFilter([]string{}, []string{}, 0, 0, &Filter{
			BankID: bank.ID,
		})

		for _, bs := range bankServices {
			bs.Delete()
		}
		for _, v = range bankPayload.Services {
			bankService = models.BankService{
				ServiceID: v,
				BankID:    bank.ID,
			}
			bankService.Create()
		}
	}
	if len(bankPayload.Products) > 0 {
		bankProduct := models.BankProduct{}
		bankProducts, _ := bankProduct.FindFilter([]string{}, []string{}, 0, 0, &Filter{
			BankID: bank.ID,
		})

		for _, bp := range bankProducts {
			bp.Delete()
		}
		for _, v = range bankPayload.Products {
			bankProduct = models.BankProduct{
				ProductID: v,
				BankID:    bank.ID,
			}
			bankProduct.Create()
		}
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
