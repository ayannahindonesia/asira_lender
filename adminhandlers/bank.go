package adminhandlers

import (
	"asira_lender/asira"
	"asira_lender/middlewares"
	"asira_lender/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ayannahindonesia/basemodel"
	"github.com/labstack/echo"
	"github.com/lib/pq"
	"github.com/thedevsaddam/govalidator"
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
		Image               string  `json:"image"`
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
		lastPage  int
		banks     []BankSelect
	)

	// pagination parameters
	rows, _ = strconv.Atoi(c.QueryParam("rows"))
	if rows > 0 {
		page, _ = strconv.Atoi(c.QueryParam("page"))
		if page <= 0 {
			page = 1
		}
		offset = (page * rows) - rows
	}

	db = db.Table("banks").
		Select("banks.*, bt.name as bank_type_name").
		Joins("INNER JOIN bank_types bt ON banks.type = bt.id")

	if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
		db = db.Or("LOWER(bt.name) LIKE ?", "%"+strings.ToLower(searchAll)+"%").
			Or("LOWER(banks.name) LIKE ?", "%"+strings.ToLower(searchAll)+"%").
			Or("LOWER(banks.pic) LIKE ?", "%"+strings.ToLower(searchAll)+"%").
			Or("CAST(banks.id as varchar(255)) = ?", searchAll)
	} else {
		if name := c.QueryParam("name"); len(name) > 0 {
			db = db.Where("LOWER(banks.name) LIKE ?", "%"+strings.ToLower(name)+"%")
		}
		if bankType := c.QueryParam("bank_type"); len(bankType) > 0 {
			db = db.Where("LOWER(bt.name) LIKE ?", "%"+strings.ToLower(bankType)+"%")
		}
		if pic := c.QueryParam("pic"); len(pic) > 0 {
			db = db.Where("LOWER(banks.pic) LIKE ?", "%"+strings.ToLower(pic)+"%")
		}
		if id := customSplit(c.QueryParam("id"), ","); len(id) > 0 {
			db = db.Where("banks.id IN (?)", id)
		}
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

	tempDB := db
	tempDB.Where("banks.deleted_at IS NULL").Count(&totalRows)

	if rows > 0 {
		db = db.Limit(rows).Offset(offset)
		lastPage = int(math.Ceil(float64(totalRows) / float64(rows)))
	}
	err = db.Find(&banks).Error
	if err != nil {
		log.Println(err)
	}

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
		"image":          []string{},
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

	if len(bankPayload.Image) > 0 {
		unbased, _ := base64.StdEncoding.DecodeString(bankPayload.Image)
		filename := "agt" + strconv.FormatInt(time.Now().Unix(), 10)
		url, err := asira.App.S3.UploadJPEG(unbased, filename)
		if err != nil {
			return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat bank baru")
		}

		bank.Image = url
	}

	err = middlewares.SubmitKafkaPayload(bank, "bank_create")
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat bank baru")
	}

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

	db = db.Table("banks").
		Select("banks.*, bt.name as bank_type_name").
		Joins("INNER JOIN bank_types bt ON banks.type = bt.id").
		Where("banks.id = ?", bankID)

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

	bankID, _ := strconv.ParseUint(c.Param("bank_id"), 10, 64)

	bank := models.Bank{}
	bankPayload := BankPayload{}
	err = bank.FindbyID(bankID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank %v tidak ditemukan", bankID))
	}

	payloadRules := govalidator.MapData{
		"name":           []string{},
		"image":          []string{},
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
	if len(bankPayload.Image) > 0 {
		unbased, _ := base64.StdEncoding.DecodeString(bankPayload.Image)
		filename := "agt" + strconv.FormatInt(time.Now().Unix(), 10)
		url, err := asira.App.S3.UploadJPEG(unbased, filename)
		if err != nil {
			return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat bank baru")
		}

		i := strings.Split(bank.Image, "/")
		delImage := i[len(i)-1]
		err = asira.App.S3.DeleteObject(delImage)
		if err != nil {
			log.Printf("failed to delete image %v from s3 bucket", delImage)
		}

		bank.Image = url
	}

	err = middlewares.SubmitKafkaPayload(bank, "bank_update")
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

	bankID, _ := strconv.ParseUint(c.Param("bank_id"), 10, 64)

	bank := models.Bank{}
	err = bank.FindbyID(bankID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("bank type %v tidak ditemukan", bankID))
	}

	err = middlewares.SubmitKafkaPayload(bank, "bank_delete")
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank tipe %v", bankID))
	}

	return c.JSON(http.StatusOK, bank)
}
