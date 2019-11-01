package admin_handlers

import (
	"asira_lender/email"
	"asira_lender/models"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

func UserList(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_user_list")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	userM := models.User{}
	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	name := c.QueryParam("name")
	id := c.QueryParam("id")
	email := c.QueryParam("email")
	phone := c.QueryParam("phone")

	type Filter struct {
		ID       string `json:"id"`
		Username string `json:"username" condition:"LIKE"`
		Email    string `json:"email" condition:"LIKE"`
		Phone    string `json:"phone" condition:"LIKE"`
	}

	result, err := userM.PagedFilterSearch(page, rows, orderby, sort, &Filter{
		ID:       id,
		Username: name,
		Email:    email,
		Phone:    phone,
	})

	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "User tidak Ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

func UserDetails(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_user_details")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	userM := models.User{}

	userID, _ := strconv.Atoi(c.Param("user_id"))
	err = userM.FindbyID(userID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "User ID tidak ditemukan")
	}

	return c.JSON(http.StatusOK, userM)
}

func UserNew(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_user_new")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	userM := models.User{}

	payloadRules := govalidator.MapData{
		"username": []string{"required", "unique:users,username"},
		"email":    []string{"required", "unique:users,email"},
		"phone":    []string{"required", "unique:users,phone"},
		"roles":    []string{},
		"status":   []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &userM)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}
	tempPW := RandString(8)
	userM.Password = tempPW

	err = userM.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat User")
	}

	to := userM.Email
	subject := "[NO REPLY] - Password Aplikasi ASIRA"
	message := "Selamat Pagi,\n\nIni adalah password anda untuk login " + tempPW + " \n\n\n Ayannah Solusi Nusantara Team"

	err = email.SendMail(to, subject, message)
	if err != nil {
		log.Println(err.Error())
	}

	return c.JSON(http.StatusCreated, userM)
}

func UserPatch(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_user_patch")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	userID, _ := strconv.Atoi(c.Param("user_id"))

	userM := models.User{}
	err = userM.FindbyID(userID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("User %v tidak ditemukan", userID))
	}
	tempPassword := userM.Password
	payloadRules := govalidator.MapData{
		"username": []string{"required", "unique:users,username,1"},
		"email":    []string{"required", "unique:users,email,1"},
		"phone":    []string{"required", "unique:users,phone,1"},
		"roles":    []string{},
		"status":   []string{},
	}
	validate := validateRequestPayload(c, payloadRules, &userM)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	userM.Password = tempPassword
	err = userM.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update User %v", userID))
	}

	return c.JSON(http.StatusOK, userM)
}
