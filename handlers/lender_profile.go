package handlers

import (
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

// TemporalSelect select sementara karena harusnya disini yang d select user bukan bank
type TemporalSelect struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
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

	lenderModel := models.Bank{}
	userModel := models.User{}

	lenderID, _ := strconv.Atoi(claims["jti"].(string))
	bankRep := models.BankRepresentatives{}
	bankRep.FindbyUserID(lenderID)

	err = lenderModel.FindbyID(int(bankRep.BankID))
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "unauthorized")
	}
	err = userModel.FindbyID(int(bankRep.UserID))
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "unauthorized")
	}

	temporal := TemporalSelect{
		ID:         lenderModel.ID,
		Name:       lenderModel.Name,
		FirstLogin: userModel.FirstLogin,
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

	err = lenderModel.FindbyID(int(bankRep.BankID))
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "unauthorized")
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

	validate := validateRequestPayload(c, payloadRules, &lenderPayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
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
	if len(lenderPayload.AdminFeeSetup) > 0 {
		lenderModel.AdminFeeSetup = lenderPayload.AdminFeeSetup
	}
	if len(lenderPayload.ConvenienceFeeSetup) > 0 {
		lenderModel.ConvenienceFeeSetup = lenderPayload.ConvenienceFeeSetup
	}

	err = lenderModel.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "error saving profile")
	}

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

	err = userModel.FindbyID(int(bankRep.UserID))
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "unauthorized")
	}

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
			return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
		}
		userModel.FirstLoginChangePassword(pass.Pass)

		return c.JSON(http.StatusOK, "Password anda telah diganti.")
	}

	return c.JSON(http.StatusUnauthorized, "Akun anda bukan akun baru.")
}
