package handlers

import (
	"asira_lender/asira"
	"asira_lender/models"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jszwec/csvutil"
	"github.com/labstack/echo"
	"gitlab.com/asira-ayannah/basemodel"
)

type (
	// BorrowerCSV custom type for query
	BorrowerCSV struct {
		basemodel.BaseModel
		DeletedTime          time.Time `json:"deleted_time"`
		Status               string    `json:"status"`
		Fullname             string    `json:"fullname"`
		Gender               string    `json:"gender"`
		IDCardNumber         string    `json:"idcard_number"`
		IDCardImageID        string    `json:"idcard_imageid"`
		TaxIDnumber          string    `json:"taxid_number"`
		TaxIDImageID         string    `json:"taxid_imageid"`
		Email                string    `json:"email"`
		Birthday             time.Time `json:"birthday"`
		Birthplace           string    `json:"birthplace"`
		LastEducation        string    `json:"last_education"`
		MotherName           string    `json:"mother_name"`
		Phone                string    `json:"phone"`
		MarriedStatus        string    `json:"marriage_status"`
		SpouseName           string    `json:"spouse_name"`
		SpouseBirthday       time.Time `json:"spouse_birthday"`
		SpouseLastEducation  string    `json:"spouse_lasteducation"`
		Dependants           int       `json:"dependants"`
		Address              string    `json:"address"`
		Province             string    `json:"province"`
		City                 string    `json:"city"`
		NeighbourAssociation string    `json:"neighbour_association"`
		Hamlets              string    `json:"hamlets"`
		HomePhoneNumber      string    `json:"home_phonenumber"`
		Subdistrict          string    `json:"subdistrict"`
		UrbanVillage         string    `json:"urban_village"`
		HomeOwnership        string    `json:"home_ownership"`
		LivedFor             int       `json:"lived_for"`
		Occupation           string    `json:"occupation"`
		EmployeeID           string    `json:"employee_id"`
		EmployerName         string    `json:"employer_name"`
		EmployerAddress      string    `json:"employer_address"`
		Department           string    `json:"department"`
		BeenWorkingFor       int       `json:"been_workingfor"`
		DirectSuperior       string    `json:"direct_superiorname"`
		EmployerNumber       string    `json:"employer_number"`
		MonthlyIncome        int       `json:"monthly_income"`
		OtherIncome          int       `json:"other_income"`
		OtherIncomeSource    string    `json:"other_incomesource"`
		FieldOfWork          string    `json:"field_of_work"`
		RelatedPersonName    string    `json:"related_personname"`
		RelatedRelation      string    `json:"related_relation"`
		RelatedPhoneNumber   string    `json:"related_phonenumber"`
		RelatedHomePhone     string    `json:"related_homenumber"`
		RelatedAddress       string    `json:"related_address"`
		Bank                 int64     `json:"bank"`
		BankAccountNumber    string    `json:"bank_accountnumber"`
		AgentID              int64     `json:"agent_id"`
		Category             string    `json:"category"`
		BankName             string    `json:"bank_name"`
		AgentName            string    `json:"agent_name"`
		AgentProviderName    string    `json:"agent_provider_name"`
	}
	// BorrowerSelect select for joining
	BorrowerSelect struct {
		models.Borrower
		Category          string `json:"category"`
		BankName          string `json:"bank_name"`
		AgentName         string `json:"agent_name"`
		AgentProviderName string `json:"agent_provider_name"`
	}
)

// LenderBorrowerList load all borrowers
func LenderBorrowerList(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "lender_borrower_list")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	lenderID, _ := strconv.Atoi(claims["jti"].(string))
	bankRep := models.BankRepresentatives{}
	bankRep.FindbyUserID(lenderID)

	db := asira.App.DB
	var (
		totalRows int
		offset    int
		rows      int
		page      int
		borrowers []BorrowerSelect
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

	db = db.Table("borrowers b").
		Select("b.*, a.category, ba.name as bank_name, a.name as agent_name, ap.name as agent_provider_name").
		Joins("LEFT JOIN agents a ON b.agent_id = a.id").
		Joins("LEFT JOIN banks ba ON ba.id = b.bank").
		Joins("LEFT JOIN agent_providers ap ON a.agent_provider = ap.id").
		Where("ba.id = ?", bankRep.BankID)

	if fullname := c.QueryParam("fullname"); len(fullname) > 0 {
		db = db.Where("LOWER(b.fullname) LIKE ?", "%"+strings.ToLower(fullname)+"%")
	}
	if category := c.QueryParam("category"); len(category) > 0 {
		db = db.Where("LOWER(a.category) = ?", strings.ToLower(category))
	}
	if bankName := c.QueryParam("bank_name"); len(bankName) > 0 {
		db = db.Where("LOWER(ba.name) LIKE ?", "%"+strings.ToLower(bankName)+"%")
	}
	if agentName := c.QueryParam("agent_name"); len(agentName) > 0 {
		db = db.Where("LOWER(a.name) LIKE ?", "%"+strings.ToLower(agentName)+"%")
	}
	if agentProviderName := c.QueryParam("agent_provider_name"); len(agentProviderName) > 0 {
		db = db.Where("LOWER(ap.name) LIKE ?", "%"+strings.ToLower(agentProviderName)+"%")
	}
	if id := customSplit(c.QueryParam("id"), ","); len(id) > 0 {
		db = db.Where("b.id IN (?)", id)
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
	tempDB.Count(&totalRows)

	if rows > 0 {
		db = db.Limit(rows).Offset(offset)
	}
	err = db.Find(&borrowers).Error
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
		Data:        borrowers,
	}

	return c.JSON(http.StatusOK, result)
}

// LenderBorrowerListDetail load borrower detail by id
func LenderBorrowerListDetail(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "lender_borrower_list_detail")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	lenderID, _ := strconv.Atoi(claims["jti"].(string))
	bankRep := models.BankRepresentatives{}
	bankRep.FindbyUserID(lenderID)

	borrowerID, err := strconv.Atoi(c.Param("borrower_id"))
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "error parsing borrower id")
	}

	db := asira.App.DB

	borrower := BorrowerSelect{}

	err = db.Table("borrowers b").
		Select("b.*, a.category, ba.name as bank_name, a.name as agent_name, ap.name as agent_provider_name").
		Joins("LEFT JOIN agents a ON b.agent_id = a.id").
		Joins("LEFT JOIN banks ba ON ba.id = b.bank").
		Joins("LEFT JOIN agent_providers ap ON a.agent_provider = ap.id").
		Where("ba.id = ?", bankRep.BankID).
		Where("b.id = ?", borrowerID).
		Find(&borrower).Error
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("id %v not found.", borrowerID))
	}

	return c.JSON(http.StatusOK, borrower)
}

// LenderBorrowerListDownload download borrower list csv
func LenderBorrowerListDownload(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "lender_borrower_list_download")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	lenderID, _ := strconv.Atoi(claims["jti"].(string))
	bankRep := models.BankRepresentatives{}
	bankRep.FindbyUserID(lenderID)

	db := asira.App.DB
	var (
		offset    int
		rows      int
		page      int
		borrowers []BorrowerSelect
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

	db = db.Table("borrowers b").
		Select("b.*, a.category, ba.name as bank_name, a.name as agent_name, ap.name as agent_provider_name").
		Joins("LEFT JOIN agents a ON b.agent_id = a.id").
		Joins("LEFT JOIN banks ba ON ba.id = b.bank").
		Joins("LEFT JOIN agent_providers ap ON a.agent_provider = ap.id").
		Where("ba.id = ?", bankRep.BankID)

	if fullname := c.QueryParam("fullname"); len(fullname) > 0 {
		db = db.Where("LOWER(b.fullname) LIKE ?", "%"+strings.ToLower(fullname)+"%")
	}
	if category := c.QueryParam("category"); len(category) > 0 {
		db = db.Where("LOWER(a.category) = ?", strings.ToLower(category))
	}
	if bankName := c.QueryParam("bank_name"); len(bankName) > 0 {
		db = db.Where("LOWER(ba.name) LIKE ?", "%"+strings.ToLower(bankName)+"%")
	}
	if agentName := c.QueryParam("agent_name"); len(agentName) > 0 {
		db = db.Where("LOWER(a.name) LIKE ?", "%"+strings.ToLower(agentName)+"%")
	}
	if agentProviderName := c.QueryParam("agent_provider_name"); len(agentProviderName) > 0 {
		db = db.Where("LOWER(ap.name) LIKE ?", "%"+strings.ToLower(agentProviderName)+"%")
	}
	if id := customSplit(c.QueryParam("id"), ","); len(id) > 0 {
		db = db.Where("b.id IN (?)", id)
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

	if rows > 0 {
		db = db.Limit(rows).Offset(offset)
	}
	err = db.Find(&borrowers).Error
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "internal error")
	}

	data := mapnewBorrowerStruct(borrowers)

	b, err := csvutil.Marshal(data)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "internal error")
	}

	return c.JSON(http.StatusOK, string(b))
}

func mapnewBorrowerStruct(m []BorrowerSelect) []BorrowerCSV {
	var r []BorrowerCSV
	for _, v := range m {
		var unmarsh BorrowerCSV
		b, _ := json.Marshal(v)
		json.Unmarshal(b, &unmarsh)
		unmarsh.Bank = v.Bank.Int64
		unmarsh.AgentID = v.AgentID.Int64
		r = append(r, unmarsh)
	}
	return r
}
