package adminhandlers

import (
	"asira_lender/asira"
	"asira_lender/email"
	"asira_lender/models"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/lib/pq"

	"github.com/ayannahindonesia/basemodel"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

type (
	// UserSelect custom select container
	UserSelect struct {
		// UserSelect custom query
		models.User
		RolesName pq.StringArray `json:"roles_name"`
		BankID    uint64         `json:"bank_id"`
		BankName  string         `json:"bank_name"`
	}
	// UserPayload handle user request body
	UserPayload struct {
		Roles    []int64 `json:"roles"`
		Username string  `json:"username"`
		Email    string  `json:"email"`
		Phone    string  `json:"phone"`
		Status   string  `json:"status"`
		Bank     uint64  `json:"bank"`
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
	rows, _ = strconv.Atoi(c.QueryParam("rows"))
	if rows > 0 {
		page, _ = strconv.Atoi(c.QueryParam("page"))
		if page <= 0 {
			page = 1
		}
		offset = (page * rows) - rows
	}
	db = db.Table("users").
		Select("DISTINCT users.*, (SELECT ARRAY_AGG(r.name) FROM roles r WHERE id IN (SELECT UNNEST(users.roles))) as roles_name, b.id as bank_id, b.name as bank_name").
		Joins("INNER JOIN roles r ON r.id IN (SELECT UNNEST(users.roles))").
		Joins("LEFT JOIN bank_representatives br ON br.user_id = users.id").
		Joins("LEFT JOIN banks b ON br.bank_id = b.id")

	if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
		db = db.Or("users.username LIKE ?", "%"+searchAll+"%").Or("users.id = ?", searchAll).Or("users.email LIKE ?", "%"+searchAll+"%").Or("users.phone LIKE ?", "%"+searchAll+"%").Or("bank_name LIKE ?", "%"+searchAll+"%")
	} else {
		if name := c.QueryParam("username"); len(name) > 0 {
			db = db.Where("users.username LIKE ?", "%"+name+"%")
		}
		if id := customSplit(c.QueryParam("id"), ","); len(id) > 0 {
			db = db.Where("users.id IN (?)", id)
		}
		if email := c.QueryParam("email"); len(email) > 0 {
			db = db.Where("users.email LIKE ?", "%"+email+"%")
		}
		if phone := c.QueryParam("phone"); len(phone) > 0 {
			db = db.Where("users.phone LIKE ?", "%"+phone+"%")
		}
		if bankName := c.QueryParam("bank_name"); len(bankName) > 0 {
			db = db.Where("bank_name LIKE ?", "%"+bankName+"%")
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
	tempDB.Where("users.deleted_at IS NULL").Count(&totalRows)

	if rows > 0 {
		db = db.Limit(rows).Offset(offset)
	}
	err = db.Find(&results).Error
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

	userID, _ := strconv.Atoi(c.Param("id"))

	err = db.Table("users").
		Select("DISTINCT users.*, (SELECT ARRAY_AGG(r.name) FROM roles r WHERE id IN (SELECT UNNEST(users.roles))) as roles_name, b.id as bank_id, b.name as bank_name").
		Joins("INNER JOIN roles r ON r.id IN (SELECT UNNEST(users.roles))").
		Joins("LEFT JOIN bank_representatives br ON br.user_id = users.id").
		Joins("LEFT JOIN banks b ON br.bank_id = b.id").
		Where("users.id = ?", userID).Find(&user).Error
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

	userM := models.User{}
	userPayload := UserPayload{}

	payloadRules := govalidator.MapData{
		"username": []string{"required", "unique:users,username"},
		"email":    []string{"required", "unique:users,email"},
		"phone":    []string{"required", "unique:users,phone"},
		"bank":     []string{"valid_id:banks"},
		"roles":    []string{"valid_id:roles"},
		"status":   []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &userPayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
	}

	if userPayload.Bank > 0 {
		db := asira.App.DB
		var count int
		db.Table("roles").Select("*").
			Where("roles.id IN (?)", []int64(userPayload.Roles)).
			Where("roles.system = ?", "Dashboard").Count(&count)

		if len(userPayload.Roles) != count {
			return returnInvalidResponse(http.StatusInternalServerError, nil, "Roles tidak valid.")
		}

		bankRepsFlag = true
	}

	marshal, _ := json.Marshal(userPayload)
	json.Unmarshal(marshal, &userM)

	tempPW := RandString(8)
	newUser := models.User{
		Username: userPayload.Username,
		Email:    userPayload.Email,
		Phone:    userPayload.Phone,
		Roles:    userPayload.Roles,
		Status:   userPayload.Status,
		Password: tempPW,
	}

	err = newUser.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat User")
	}

	if bankRepsFlag {
		bankRep := models.BankRepresentatives{
			UserID: newUser.ID,
			BankID: userPayload.Bank,
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

// UserPatch edit user by id
func UserPatch(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_user_patch")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	userM := models.User{}
	userPayload := UserPayload{}
	err = userM.FindbyID(userID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("User %v tidak ditemukan", userID))
	}

	payloadRules := govalidator.MapData{
		"username": []string{},
		"email":    []string{},
		"phone":    []string{},
		"bank":     []string{"valid_id:banks"},
		"roles":    []string{"valid_id:roles"},
		"status":   []string{},
	}
	validate := validateRequestPayload(c, payloadRules, &userPayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
	}

	bankRep := models.BankRepresentatives{}
	bankRep.FindbyUserID(int(userM.ID))
	if len(userPayload.Roles) > 0 && bankRep.ID != 0 {
		db := asira.App.DB
		var count int
		db.Table("roles").Select("*").
			Where("roles.id IN (?)", []int64(userPayload.Roles)).
			Where("roles.system = ?", "Dashboard").Count(&count)

		if len(userPayload.Roles) != count {
			return returnInvalidResponse(http.StatusUnprocessableEntity, nil, "Roles tidak valid.")
		}
	}

	if len(userPayload.Username) > 0 {
		userM.Username = userPayload.Username
	}
	if len(userPayload.Email) > 0 {
		userM.Email = userPayload.Email
	}
	if len(userPayload.Phone) > 0 {
		userM.Phone = userPayload.Phone
	}
	if len(userPayload.Status) > 0 {
		userM.Status = userPayload.Status
	}
	if len(userPayload.Roles) > 0 {
		userM.Roles = pq.Int64Array(userPayload.Roles)
	}
	if userPayload.Bank != 0 {
		bankRep.BankID = userPayload.Bank
		bankRep.Save()
	}

	err = userM.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update User %v", userID))
	}

	return c.JSON(http.StatusOK, userM)
}
