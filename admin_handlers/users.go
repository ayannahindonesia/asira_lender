package admin_handlers

import (
	"asira_lender/asira"
	"asira_lender/email"
	"asira_lender/models"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/lib/pq"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
	"gitlab.com/asira-ayannah/basemodel"
)

type UserSelect struct {
	models.User
	RolesName pq.StringArray `json:"roles_name"`
}

func UserList(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_user_list")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	db := asira.App.DB

	var results []UserSelect
	var totalRows int
	var offset int
	var rows int
	var page int

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
	db = db.Table("users u").
		Select("DISTINCT u.*, (SELECT ARRAY_AGG(r.name) FROM roles r WHERE id IN (SELECT UNNEST(u.roles))) as roles_name").
		Joins("INNER JOIN roles r ON r.id IN (SELECT UNNEST(u.roles))")

	if name := c.QueryParam("username"); len(name) > 0 {
		db = db.Where("u.username LIKE ?", name)
	}
	if id := customSplit(c.QueryParam("id"), ","); len(id) > 0 {
		db = db.Where("u.id IN (?)", id)
	}
	if email := c.QueryParam("email"); len(email) > 0 {
		db = db.Where("u.email LIKE ?", email)
	}
	if phone := c.QueryParam("phone"); len(phone) > 0 {
		db = db.Where("u.phone LIKE ?", phone)
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
	err = db.Find(&results).Count(&totalRows).Error
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
		Data:        results,
	}

	return c.JSON(http.StatusOK, result)
}

func UserDetails(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_user_details")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	db := asira.App.DB

	user := UserSelect{}

	userID, _ := strconv.Atoi(c.Param("user_id"))

	err = db.Table("users u").
		Select("DISTINCT u.*, (SELECT ARRAY_AGG(r.name) FROM roles r WHERE id IN (SELECT UNNEST(u.roles))) as roles_name").
		Joins("INNER JOIN roles r ON r.id IN (SELECT UNNEST(u.roles))").
		Where("u.id = ?", userID).Find(&user).Error
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "User ID tidak ditemukan")
	}

	return c.JSON(http.StatusOK, user)
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
		"username": []string{},
		"email":    []string{},
		"phone":    []string{},
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
