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

type (
	// UserSelect for query
	UserSelect struct {
		models.User
		RolesName pq.StringArray `json:"roles_name"`
		BankID    uint64         `json:"bank_id"`
		BankName  string         `json:"bank_name"`
	}
	// UserFields for post and patch
	UserFields struct {
		Username string        `json:"username"`
		Email    string        `json:"email"`
		Phone    string        `json:"phone"`
		Roles    pq.Int64Array `json:"roles"`
		Status   string        `json:"status"`
		Bank     uint64        `json:"bank"`
	}
)

// UserList gets all users
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
		Select("DISTINCT u.*, (SELECT ARRAY_AGG(r.name) FROM roles r WHERE id IN (SELECT UNNEST(u.roles))) as roles_name, b.id as bank_id, b.name as bank_name").
		Joins("INNER JOIN roles r ON r.id IN (SELECT UNNEST(u.roles))").
		Joins("LEFT JOIN bank_representatives br ON br.user_id = u.id").
		Joins("LEFT JOIN banks b ON br.bank_id = b.id")

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
	if bankName := c.QueryParam("bank_name"); len(bankName) > 0 {
		db = db.Where("bank_name LIKE ?", bankName)
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

// UserDetails get user detail by id
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
		Select("DISTINCT u.*, (SELECT ARRAY_AGG(r.name) FROM roles r WHERE id IN (SELECT UNNEST(u.roles))) as roles_name, b.id as bank_id, b.name as bank_name").
		Joins("INNER JOIN roles r ON r.id IN (SELECT UNNEST(u.roles))").
		Joins("LEFT JOIN bank_representatives br ON br.user_id = u.id").
		Joins("LEFT JOIN banks b ON br.bank_id = b.id").
		Where("u.id = ?", userID).Find(&user).Error
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "User ID tidak ditemukan")
	}

	return c.JSON(http.StatusOK, user)
}

// UserNew create new user
func UserNew(c echo.Context) error {
	bankRepsFlag := false
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_user_new")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	userF := UserFields{}

	payloadRules := govalidator.MapData{
		"username": []string{"required", "unique:users,username"},
		"email":    []string{"required", "unique:users,email"},
		"phone":    []string{"required", "unique:users,phone"},
		"bank":     []string{"valid_id:banks"},
		"roles":    []string{},
		"status":   []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &userF)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	if userF.Bank > 0 {
		db := asira.App.DB
		var count int
		db.Table("roles r").Select("*").
			Where("r.id IN (?)", []int64(userF.Roles)).
			Where("r.system = ?", "Dashboard").Count(&count)

		if len(userF.Roles) != count {
			return returnInvalidResponse(http.StatusInternalServerError, nil, "Roles tidak valid.")
		}

		bankRepsFlag = true
	}

	tempPW := RandString(8)
	newUser := models.User{
		Username: userF.Username,
		Email:    userF.Email,
		Phone:    userF.Phone,
		Roles:    userF.Roles,
		Status:   userF.Status,
		Password: tempPW,
	}

	err = newUser.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat User")
	}

	if bankRepsFlag {
		bankRep := models.BankRepresentatives{
			UserID: newUser.ID,
			BankID: userF.Bank,
		}
		err = bankRep.Create()
		if err != nil {
			return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat Bank User")
		}
	}

	to := newUser.Email
	subject := "[NO REPLY] - Password Aplikasi ASIRA"
	message := "Selamat Pagi,\n\nIni adalah password anda untuk login " + tempPW + " \n\n\n Ayannah Solusi Nusantara Team"

	err = email.SendMail(to, subject, message)
	if err != nil {
		log.Println(err.Error())
	}

	return c.JSON(http.StatusCreated, newUser)
}

// UserPatch edit user
func UserPatch(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_user_patch")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	userID, _ := strconv.Atoi(c.Param("user_id"))

	userM := models.User{}
	userF := UserFields{}
	err = userM.FindbyID(userID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("User %v tidak ditemukan", userID))
	}

	payloadRules := govalidator.MapData{
		"username": []string{"unique:users,username,1"},
		"email":    []string{"unique:users,email,1"},
		"phone":    []string{"unique:users,phone,1"},
		"bank":     []string{"valid_id:banks"},
		"roles":    []string{},
		"status":   []string{},
	}
	validate := validateRequestPayload(c, payloadRules, &userF)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	bankRep := models.BankRepresentatives{}
	bankRep.FindbyUserID(int(userM.ID))
	if len(userF.Roles) > 0 && bankRep.ID != 0 {
		db := asira.App.DB
		var count int
		db.Table("roles r").Select("*").
			Where("r.id IN (?)", []int64(userF.Roles)).
			Where("r.system = ?", "Dashboard").Count(&count)

		if len(userF.Roles) != count {
			return returnInvalidResponse(http.StatusUnprocessableEntity, nil, "Roles tidak valid.")
		}
	}

	if len(userF.Username) > 0 {
		userM.Username = userF.Username
	}
	if len(userF.Email) > 0 {
		userM.Email = userF.Email
	}
	if len(userF.Phone) > 0 {
		userM.Phone = userF.Phone
	}
	if len(userF.Status) > 0 {
		userM.Status = userF.Status
	}
	if userF.Bank != 0 {
		bankRep.BankID = userF.Bank
		bankRep.Save()
	}
	if len(userF.Roles) > 0 {
		userM.Roles = userF.Roles
	}

	err = userM.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update User %v", userID))
	}

	return c.JSON(http.StatusOK, userM)
}
