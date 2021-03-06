package handlers

import (
	"asira_lender/adminhandlers"
	"asira_lender/asira"
	"asira_lender/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lib/pq"
	"github.com/thedevsaddam/govalidator"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// LenderProfilePayload type
type LenderProfilePayload struct {
	Name     string  `json:"name"`
	Type     uint64  `json:"type"`
	Address  string  `json:"address"`
	Province string  `json:"province"`
	City     string  `json:"city"`
	PIC      string  `json:"pic"`
	Phone    string  `json:"phone"`
	Services []int64 `json:"services"`
	Products []int64 `json:"products"`
}

// TemporalSelect select sementara karena harusnya disini yang d select user bukan bank
type TemporalSelect struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	FirstLogin bool   `json:"first_login"`
}

// LenderProfile show current lender info
func LenderProfile(c echo.Context) error {
	defer c.Request().Body.Close()
	// err := validatePermission(c, "lender_profile")
	// if err != nil {
	// 	return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	// }

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	lenderID, _ := strconv.Atoi(claims["jti"].(string))
	bankRep := models.BankRepresentatives{}
	bankRep.FindbyUserID(lenderID)

	temporal := TemporalSelect{}

	db := asira.App.DB
	err := db.Table("bank_representatives").
		Select("u.id, b.name, b.image, u.first_login").
		Joins("INNER JOIN users u ON u.id = bank_representatives.user_id").
		Joins("INNER JOIN banks b ON b.id = bank_representatives.bank_id").
		Where("bank_representatives.user_id = ?", lenderID).Find(&temporal).Error

	if err != nil {
		adminhandlers.NLog("warning", "LenderProfile", map[string]interface{}{"message": "error finding profile", "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusForbidden, err, "Tidak memiliki hak akses")
	}

	return c.JSON(http.StatusOK, temporal)
}

// LenderProfileEdit edit current lender profile
func LenderProfileEdit(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "lender_profile_edit")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	lenderModel := models.Bank{}
	lenderPayload := LenderProfilePayload{}

	lenderID, _ := strconv.Atoi(claims["jti"].(string))
	bankRep := models.BankRepresentatives{}
	bankRep.FindbyUserID(lenderID)

	err = lenderModel.FindbyID(bankRep.BankID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Tidak memiliki hak akses")
	}
	origin := lenderModel

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

	validate := validateRequestPayload(c, payloadRules, &lenderPayload)
	if validate != nil {
		adminhandlers.NLog("warning", "LenderProfileEdit", map[string]interface{}{"message": "error validation", "error": validate}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
	}

	if len(lenderPayload.Name) > 0 {
		lenderModel.Name = lenderPayload.Name
	}
	if lenderPayload.Type > 0 {
		lenderModel.Type = lenderPayload.Type
	}
	if len(lenderPayload.Address) > 0 {
		lenderModel.Address = lenderPayload.Address
	}
	if len(lenderPayload.Province) > 0 {
		lenderModel.Province = lenderPayload.Province
	}
	if len(lenderPayload.City) > 0 {
		lenderModel.City = lenderPayload.City
	}
	if len(lenderPayload.Services) > 0 {
		lenderModel.Services = pq.Int64Array(lenderPayload.Services)
	}
	if len(lenderPayload.Products) > 0 {
		lenderModel.Products = pq.Int64Array(lenderPayload.Products)
	}
	if len(lenderPayload.PIC) > 0 {
		lenderModel.PIC = lenderPayload.PIC
	}
	if len(lenderPayload.Phone) > 0 {
		lenderModel.Phone = lenderPayload.Phone
	}

	err = lenderModel.Save()
	if err != nil {
		adminhandlers.NLog("error", "LenderProfileEdit", map[string]interface{}{"message": "error saving profile", "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "Terjadi kesalahan")
	}

	adminhandlers.NAudittrail(origin, lenderModel, c.Get("user").(*jwt.Token), "user", fmt.Sprint(lenderModel.ID), "update")

	return c.JSON(http.StatusOK, lenderModel)
}

// UserFirstLoginChangePassword check if user is first login and change the password
func UserFirstLoginChangePassword(c echo.Context) error {
	defer c.Request().Body.Close()

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	userModel := models.User{}

	lenderID, _ := strconv.Atoi(claims["jti"].(string))
	bankRep := models.BankRepresentatives{}
	bankRep.FindbyUserID(lenderID)

	err = userModel.FindbyID(bankRep.UserID)
	if err != nil {
		adminhandlers.NLog("error", "UserFirstLoginChangePassword", map[string]interface{}{"message": fmt.Sprintf("error finding profile %v", lenderID), "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusForbidden, err, "Tidak memiliki hak akses")
	}
	origin := userModel

	if userModel.FirstLogin {
		type Password struct {
			Pass string `json:"password"`
		}
		var pass Password
		payloadRules := govalidator.MapData{
			"password": []string{"required"},
		}

		validate := validateRequestPayload(c, payloadRules, &pass)
		if validate != nil {
			return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
		}
		userModel.FirstLoginChangePassword(pass.Pass)

		adminhandlers.NLog("info", "UserFirstLoginChangePassword", map[string]interface{}{"message": fmt.Sprintf("lender %v changed password", lenderID)}, c.Get("user").(*jwt.Token), "", false)

		adminhandlers.NAudittrail(origin, userModel, c.Get("user").(*jwt.Token), "user", fmt.Sprint(userModel.ID), "first login change password")

		return c.JSON(http.StatusOK, "Password anda telah diganti.")
	}

	adminhandlers.NLog("error", "UserFirstLoginChangePassword", map[string]interface{}{"message": fmt.Sprintf("lender %v is not new account, therefore cant change password", lenderID)}, c.Get("user").(*jwt.Token), "", false)

	return c.JSON(http.StatusUnauthorized, "Akun anda bukan akun baru.")
}
