package handlers

import (
	"asira_lender/asira"
	"asira_lender/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ayannahindonesia/basemodel"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jszwec/csvutil"
	"github.com/labstack/echo"
)

type (
	// BorrowerCSV custom type for query
	BorrowerCSV struct {
		basemodel.BaseModel
		Status               string    `json:"status"`
		Fullname             string    `json:"fullname"`
		Gender               string    `json:"gender"`
		IDCardNumber         string    `json:"idcard_number"`
		IDCardImage          string    `json:"idcard_image"`
		TaxIDnumber          string    `json:"taxid_number"`
		TaxIDImage           string    `json:"taxid_image"`
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
		LoanCount         int    `json:"loan_count"`
		LoanStatus        string `json:"loan_status"`
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
		lastPage  int
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

	loanStatusQuery := fmt.Sprintf("CASE WHEN (SELECT COUNT(id) FROM loans l WHERE l.borrower = b.id AND status = '%s' AND (disburse_date = '0001-01-01 00:00:00+00' OR (disburse_date != '0001-01-01 00:00:00+00' AND NOW() > l.disburse_date AND NOW() < l.disburse_date + make_interval(months => l.installment) + make_interval(days => 1)))) > 0 THEN '%s' ELSE '%s' END", "approved", "active", "inactive")

	db = db.Table("borrowers b").
		Select("DISTINCT b.*, a.category, ba.name as bank_name, a.name as agent_name, ap.name as agent_provider_name, (SELECT COUNT(id) FROM loans l WHERE l.borrower = b.id AND l.status = ?) as loan_count, "+loanStatusQuery+" as loan_status", "approved").
		Joins("LEFT JOIN agents a ON b.agent_referral = a.id").
		Joins("LEFT JOIN banks ba ON ba.id = b.bank").
		Joins("LEFT JOIN agent_providers ap ON a.agent_provider = ap.id").
		Where("ba.id = ?", bankRep.BankID).
		Where("b.status != ?", "rejected")

	accountNumber := c.QueryParam("account_number")

	if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
		// gorm havent support nested subquery yet.
		searchLike := "%" + strings.ToLower(searchAll) + "%"
		extraquery := fmt.Sprintf("LOWER(b.fullname) LIKE ?") + // use searchLike #1
			fmt.Sprintf(" OR LOWER(a.category) = ?") + // use searchLike #2
			fmt.Sprintf(" OR LOWER(ba.name) LIKE ?") + // use searchLike #3
			fmt.Sprintf(" OR LOWER(a.name) LIKE ?") + // use searchLike #4
			fmt.Sprintf(" OR LOWER(ap.name) LIKE ?") + // use searchLike #5
			fmt.Sprintf(" OR CAST(b.id as varchar(255)) = ?") + // use searchAll #6
			fmt.Sprintf(" OR "+loanStatusQuery+" LIKE ?") // use searchLike #7

		if len(accountNumber) > 0 {
			if accountNumber == "null" {
				db = db.Where("b.bank_accountnumber = ?", "")
			} else if accountNumber == "not null" {
				db = db.Where("b.bank_accountnumber != ?", "")
			}
		}

		db = db.Where(extraquery, searchLike, searchLike, searchLike, searchLike, searchLike, searchAll, searchLike)
	} else {
		if fullname := c.QueryParam("fullname"); len(fullname) > 0 {
			db = db.Where("LOWER(b.fullname) LIKE ?", "%"+strings.ToLower(fullname)+"%")
		}
		if category := c.QueryParam("category"); len(category) > 0 {
			db = db.Where("LOWER(a.category) = ?", "%"+strings.ToLower(category)+"%")
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
		if len(accountNumber) > 0 {
			if accountNumber == "null" {
				db = db.Where("b.bank_accountnumber = ?", "")
			} else if accountNumber == "not null" {
				db = db.Where("b.bank_accountnumber != ?", "")
			} else {
				db = db.Where("b.bank_accountnumber LIKE ?", "%"+strings.ToLower(accountNumber)+"%")
			}
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
	tempDB.Count(&totalRows)

	if rows > 0 {
		db = db.Limit(rows).Offset(offset)
		lastPage = int(math.Ceil(float64(totalRows) / float64(rows)))
	}
	err = db.Find(&borrowers).Error
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

	loanStatusQuery := fmt.Sprintf("CASE WHEN (SELECT COUNT(id) FROM loans l WHERE l.borrower = b.id AND status = '%s' AND (disburse_date = '0001-01-01 00:00:00+00' OR (disburse_date != '0001-01-01 00:00:00+00' AND NOW() > l.disburse_date AND NOW() < l.disburse_date + make_interval(months => l.installment) + make_interval(days => 1)))) > 0 THEN '%s' ELSE '%s' END", "approved", "active", "inactive")

	err = db.Table("borrowers b").
		Select("DISTINCT b.*, a.category, ba.name as bank_name, a.name as agent_name, ap.name as agent_provider_name, (SELECT COUNT(id) FROM loans l WHERE l.borrower = b.id AND l.status = ?) as loan_count, "+loanStatusQuery+" as loan_status", "approved").
		Joins("LEFT JOIN agents a ON b.agent_referral = a.id").
		Joins("LEFT JOIN banks ba ON ba.id = b.bank").
		Joins("LEFT JOIN agent_providers ap ON a.agent_provider = ap.id").
		Where("ba.id = ?", bankRep.BankID).
		Where("b.id = ?", borrowerID).
		Where("b.status != ?", "rejected").
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

	loanStatusQuery := fmt.Sprintf("CASE WHEN (SELECT COUNT(id) FROM loans l WHERE l.borrower = b.id AND status = '%s' AND (disburse_date = '0001-01-01 00:00:00+00' OR (disburse_date != '0001-01-01 00:00:00+00' AND NOW() > l.disburse_date AND NOW() < l.disburse_date + make_interval(months => l.installment) + make_interval(days => 1)))) > 0 THEN '%s' ELSE '%s' END", "approved", "active", "inactive")

	db = db.Table("borrowers b").
		Select("DISTINCT b.*, a.category, ba.name as bank_name, a.name as agent_name, ap.name as agent_provider_name, (SELECT COUNT(id) FROM loans l WHERE l.borrower = b.id AND l.status = ?) as loan_count, "+loanStatusQuery+" as loan_status", "approved").
		Joins("LEFT JOIN agents a ON b.agent_referral = a.id").
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
	if accountNumber := c.QueryParam("account_number"); len(accountNumber) > 0 {
		if accountNumber == "null" {
			db = db.Where("b.bank_accountnumber = ?", "")
		} else if accountNumber == "not null" {
			db = db.Where("b.bank_accountnumber != ?", "")
		} else {
			db = db.Where("b.bank_accountnumber LIKE ?", "%"+strings.ToLower(accountNumber)+"%")
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

// LenderApproveRejectProspectiveBorrower approve or reject prospective borrower
func LenderApproveRejectProspectiveBorrower(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "lender_prospective_borrower_approval")
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
	type Filter struct {
		Bank              sql.NullInt64 `json:"bank"`
		ID                int           `json:"id"`
		BankAccountNumber string        `json:"bank_accountnumber"`
	}

	borrower := models.Borrower{}
	err = borrower.FilterSearchSingle(&Filter{
		Bank: sql.NullInt64{
			Int64: int64(bankRep.BankID),
			Valid: true,
		},
		ID:                borrowerID,
		BankAccountNumber: "",
	})
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "borrower not found")
	}

	approval := c.Param("approval")
	switch approval {
	default:
		if accNumber := c.QueryParam("account_number"); len(accNumber) > 0 {
			borrower.BankAccountNumber = accNumber
			borrower.Approve()
		} else {
			return returnInvalidResponse(http.StatusUnprocessableEntity, "", "invalid account number")
		}
		break
	case "reject":
		borrower.Reject()
		break
	}

	return c.JSON(http.StatusOK, borrower)
}

func mapnewBorrowerStruct(m []BorrowerSelect) []BorrowerCSV {
	var r []BorrowerCSV
	for _, v := range m {
		var unmarsh BorrowerCSV
		b, _ := json.Marshal(v)
		json.Unmarshal(b, &unmarsh)
		unmarsh.Bank = v.Bank.Int64
		unmarsh.AgentID = v.AgentReferral.Int64
		r = append(r, unmarsh)
	}
	return r
}
