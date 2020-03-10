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

//ProductFilter for generate parameter filter
type ProductFilter struct {
	ID              int64  `json:"id" condition:"optional"`
	Name            string `json:"name" condition:"LIKE,optional"`
	Interest        string `json:"interest" condition:"LIKE,optional"`
	Fees            string `json:"fees" condition:"LIKE,optional"`
	Collaterals     string `json:"collaterals" condition:"LIKE,optional"`
	FinancingSector string `json:"financing_sector" condition:"LIKE,optional"`
	Assurance       string `json:"assurance" condition:"LIKE,optional"`
	Status          string `json:"status" condition:"LIKE,optional"`
}

// ProductList get all product list
func LenderProductList(c echo.Context) error {
	defer c.Request().Body.Close()

	const LogTag = "LenderProductList"

	var services []models.Service

	//get token & jti
	token := c.Get("user").(*jwt.Token)
	jti := token.Claims.(jwt.MapClaims)["jti"].(string)
	lenderID, _ := strconv.ParseUint(jti, 10, 64)
	bankRep := models.BankRepresentatives{}

	//get bank representatives
	err = bankRep.FindbyUserID(int(lenderID))
	if err != nil {
		adminhandlers.NLog("error", LogTag, map[string]interface{}{
			"message": "invalid lender id",
			"error":   err}, token, "", false)

		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	//Extended Query
	QPaged := modules.QueryPaged{}
	QPaged.Init(c)
	fmt.Println("QPaged = ", QPaged)

	//custom query
	db := asira.App.DB
	db = db.Table("products").
		Select("products.*").
		Joins("INNER JOIN banks b ON products.id IN (SELECT UNNEST(b.products)) ").
		Where("b.id = ?", bankRep.BankID)

	//generate filter, return db and error
	db, err = QPaged.GenerateFilters(db, ProductFilter{}, "products")
	if err != nil {
		adminhandlers.NLog("warning", LogTag, map[string]interface{}{
			"message": "error listing services",
			"error":   err}, token, "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, "kesalahan dalam filters")
	}

	//execute anonymous function pass db and data pass by reference (services)
	err = QPaged.Exec(db, &services, func(DB *gorm.DB, rows interface{}) error {
		//manual type casting :)
		err := DB.Find(rows.(*[]models.Service)).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		adminhandlers.NLog("error", LogTag, map[string]interface{}{
			"message": "error listing services",
			"error":   err}, token, "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Pencarian tidak ditemukan")
	}

	//get result format
	result := QPaged.GetPage(services)

	return c.JSON(http.StatusOK, result)

}

//LenderProductDetail get product detail by id
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
