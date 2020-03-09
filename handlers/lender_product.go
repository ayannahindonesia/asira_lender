package handlers

import (
	"asira_lender/adminhandlers"
	"asira_lender/asira"
	"asira_lender/models"
	"asira_lender/modules"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

// ProductList get all product list
func LenderProductList(c echo.Context) error {
	defer c.Request().Body.Close()

	var services []models.Service

	jti := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["jti"].(string)
	lenderID, _ := strconv.ParseUint(jti, 10, 64)
	bankRep := models.BankRepresentatives{}

	//get bank representatives
	err = bankRep.FindbyUserID(int(lenderID))
	if err != nil {
		adminhandlers.NLog("warning", "LenderServiceList", map[string]interface{}{"message": "error listing services", "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	//Extended Query
	QPaged := modules.QueryPaged{}
	QPaged.Init(c)
	fmt.Println("QPaged = ", QPaged)

	db := asira.App.DB
	db = db.Table("products").
		Select("products.*").
		Joins("INNER JOIN banks b ON products.id IN (SELECT UNNEST(b.products)) ").
		Where("b.id = ?", bankRep.BankID)

	//execute anonymous function
	err = QPaged.Exec(db, &services, func(db *gorm.DB, srv interface{}) error {
		//manual type casting :)
		err := db.Find(srv.(*[]models.Service)).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		adminhandlers.NLog("warning", "LenderServiceList", map[string]interface{}{"message": "error listing services", "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Pencarian tidak ditemukan")
	}

	//get result format
	result := QPaged.GetPage(services)

	return c.JSON(http.StatusOK, result)

	// if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
	// 	type Filter struct {
	// 		ID              int64  `json:"id" condition:"optional"`
	// 		Name            string `json:"name" condition:"LIKE,optional"`
	// 		Interest        string `json:"interest" condition:"LIKE,optional"`
	// 		Fees            string `json:"fees" condition:"LIKE,optional"`
	// 		Collaterals     string `json:"collaterals" condition:"LIKE,optional"`
	// 		FinancingSector string `json:"financing_sector" condition:"LIKE,optional"`
	// 		Assurance       string `json:"assurance" condition:"LIKE,optional"`
	// 		Status          string `json:"status" condition:"LIKE,optional"`
	// 	}
	// 	id, _ := strconv.ParseInt(searchAll, 10, 64)
	// 	result, err = product.PagedFindFilter(page, rows, order, sort, &Filter{
	// 		ID:              id,
	// 		Name:            searchAll,
	// 		Interest:        searchAll,
	// 		Fees:            searchAll,
	// 		Collaterals:     searchAll,
	// 		FinancingSector: searchAll,
	// 		Assurance:       searchAll,
	// 		Status:          searchAll,
	// 	})
	// } else {
	// 	type Filter struct {
	// 		ID              []string `json:"id"`
	// 		Name            string   `json:"name" condition:"LIKE"`
	// 		ServiceID       []string `json:"service_id"`
	// 		MinTimeSpan     string   `json:"min_timespan"`
	// 		MaxTimeSpan     string   `json:"max_timespan"`
	// 		Interest        string   `json:"interest" condition:"LIKE"`
	// 		MinLoan         string   `json:"min_loan"`
	// 		MaxLoan         string   `json:"max_loan"`
	// 		Fees            string   `json:"fees" condition:"LIKE"`
	// 		Collaterals     string   `json:"collaterals" condition:"LIKE"`
	// 		FinancingSector string   `json:"financing_sector" condition:"LIKE"`
	// 		Assurance       string   `json:"assurance" condition:"LIKE"`
	// 		Status          string   `json:"status" condition:"LIKE"`
	// 	}
	// 	result, err = product.PagedFindFilter(page, rows, order, sort, &Filter{
	// 		ID:              customSplit(c.QueryParam("id"), ","),
	// 		Name:            c.QueryParam("name"),
	// 		ServiceID:       customSplit(c.QueryParam("service_id"), ","),
	// 		MinTimeSpan:     c.QueryParam("min_timespan"),
	// 		MaxTimeSpan:     c.QueryParam("max_timespan"),
	// 		Interest:        c.QueryParam("interest"),
	// 		MinLoan:         c.QueryParam("min_loan"),
	// 		MaxLoan:         c.QueryParam("max_loan"),
	// 		Fees:            c.QueryParam("fee"),
	// 		Collaterals:     c.QueryParam("collaterals"),
	// 		FinancingSector: c.QueryParam("financing_sector"),
	// 		Assurance:       c.QueryParam("assurance"),
	// 		Status:          c.QueryParam("status"),
	// 	})
	// }

	// if err != nil {
	// 	NLog("warning", "ProductList", map[string]interface{}{"message": "error listing products", "error": err}, c.Get("user").(*jwt.Token), "", false)

	// 	return returnInvalidResponse(http.StatusInternalServerError, err, "Pencarian tidak ditemukan")
	// }

	// return c.JSON(http.StatusOK, "ok")
}

// ProductDetail get product detail by id
func LenderProductDetail(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "lender_product_list_detail")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	productID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	product := models.Product{}
	err = product.FindbyID(productID)
	if err != nil {
		adminhandlers.NLog("warning", "ProductDetail", map[string]interface{}{"message": fmt.Sprintf("find product %v error", productID), "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Product %v tidak ditemukan", productID))
	}

	return c.JSON(http.StatusOK, product)
}
