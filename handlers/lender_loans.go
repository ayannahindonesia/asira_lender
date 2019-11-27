package handlers

import (
	"asira_lender/asira"
	"asira_lender/models"
	"database/sql"
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
	// LoanRequestListCSV type
	LoanRequestListCSV struct {
		ID                uint64  `json:"id"`
		OwnerName         string  `json:"owner_name"`
		BankName          string  `json:"bank_name"`
		Status            string  `json:"status"`
		LoanAmount        float64 `json:"loan_amount"`
		Installment       int     `json:"installment"`
		Fees              string
		Interest          float64   `json:"interest"`
		TotalLoan         float64   `json:"total_loan"`
		DueDate           time.Time `json:"due_date"`
		LayawayPlan       float64   `json:"layaway_plan"`
		LoanIntention     string    `json:"loan_intention"`
		IntentionDetails  string    `json:"intention_details"`
		MonthlyIncome     int       `json:"monthly_income"`
		OtherIncome       int       `json:"other_income"`
		OtherIncomeSource string    `json:"other_incomesource"`
		BankAccountNumber string    `json:"bank_accountnumber"`
	}
	// LoanSelect select custom type
	LoanSelect struct {
		models.Loan
		BankName          string `json:"bank_name"`
		Service           string `json:"service"`
		Product           string `json:"product"`
		Category          string `json:"category"`
		AgentName         string `json:"agent_name"`
		AgentProviderName string `json:"agent_provider_name"`
	}
)

// LenderLoanRequestList load all loans
func LenderLoanRequestList(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "lender_loan_request_list")
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
		loans     []LoanSelect
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

	db = db.Table("loans l").
		Select("l.*, ba.name as bank_name, s.name as service, p.name as product, a.category, a.name as agent_name, ap.name as agent_provider_name").
		Joins("LEFT JOIN products p ON l.product = p.id").
		Joins("LEFT JOIN services s ON p.service_id = s.id").
		Joins("LEFT JOIN banks ba ON l.bank = ba.id").
		Joins("LEFT JOIN borrowers b ON b.id = l.owner").
		Joins("LEFT JOIN agents a ON b.agent_id = a.id").
		Joins("LEFT JOIN agent_providers ap ON a.agent_provider = ap.id").
		Where("l.bank = ?", bankRep.BankID)

	if status := c.QueryParam("status"); len(status) > 0 {
		db = db.Where("LOWER(l.status) LIKE ?", "%"+strings.ToLower(status)+"%")
	}

	if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
		// gorm havent support nested subquery yet.
		extraquery := fmt.Sprintf("CAST(l.owner as varchar(255)) = '%v'", searchAll) +
			fmt.Sprintf(" OR LOWER(l.owner_name) LIKE '%v'", "%"+strings.ToLower(searchAll)+"%") +
			fmt.Sprintf(" OR CAST(l.id as varchar(255)) = '%v'", searchAll) +
			fmt.Sprintf(" OR LOWER(l.disburse_status) LIKE '%v'", "%"+strings.ToLower(searchAll)+"%") +
			fmt.Sprintf(" OR LOWER(a.category) LIKE '%v'", "%"+strings.ToLower(searchAll)+"%") +
			fmt.Sprintf(" OR LOWER(a.name) LIKE '%v'", "%"+strings.ToLower(searchAll)+"%") +
			fmt.Sprintf(" OR LOWER(ap.name) LIKE '%v'", "%"+strings.ToLower(searchAll)+"%")

		db = db.Where(extraquery)
	} else {
		if owner := c.QueryParam("owner"); len(owner) > 0 {
			db = db.Where("l.owner = ?", owner)
		}
		if ownerName := c.QueryParam("owner_name"); len(ownerName) > 0 {
			db = db.Where("LOWER(l.owner_name) LIKE ?", "%"+strings.ToLower(ownerName)+"%")
		}
		if id := customSplit(c.QueryParam("id"), ","); len(id) > 0 {
			db = db.Where("l.id IN (?)", id)
		}
		if disburseStatus := c.QueryParam("disburse_status"); len(disburseStatus) > 0 {
			db = db.Where("LOWER(l.disburse_status) LIKE ?", "%"+strings.ToLower(disburseStatus)+"%")
		}
		if startDate := c.QueryParam("start_date"); len(startDate) > 0 {
			if endDate := c.QueryParam("end_date"); len(endDate) > 0 {
				db = db.Where("l.created_time BETWEEN ? AND ?", startDate, endDate)
			} else {
				db = db.Where("l.created_time BETWEEN ? AND ?", startDate, startDate)
			}
		}
		if startDisburseDate := c.QueryParam("start_disburse_date"); len(startDisburseDate) > 0 {
			if endDisburseDate := c.QueryParam("end_disburse_date"); len(endDisburseDate) > 0 {
				db = db.Where("l.disburse_date BETWEEN ? AND ?", startDisburseDate, endDisburseDate)
			} else {
				db = db.Where("l.disburse_date BETWEEN ? AND ?", startDisburseDate, startDisburseDate)
			}
		}
		if category := c.QueryParam("category"); len(category) > 0 {
			db = db.Where("LOWER(a.category) LIKE ?", "%"+strings.ToLower(category)+"%")
		}
		if agentName := c.QueryParam("agent_name"); len(agentName) > 0 {
			db = db.Where("LOWER(a.name) LIKE ?", "%"+strings.ToLower(agentName)+"%")
		}
		if agentProviderName := c.QueryParam("agent_provider_name"); len(agentProviderName) > 0 {
			db = db.Where("LOWER(ap.name) LIKE ?", "%"+strings.ToLower(agentProviderName)+"%")
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
	err = db.Find(&loans).Error
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
		Data:        loans,
	}

	return c.JSON(http.StatusOK, result)
}

// LenderLoanRequestListDetail load loan by id
func LenderLoanRequestListDetail(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "lender_loan_request_detail")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	lenderID, _ := strconv.Atoi(claims["jti"].(string))
	bankRep := models.BankRepresentatives{}
	bankRep.FindbyUserID(lenderID)

	loanID, err := strconv.Atoi(c.Param("loan_id"))
	db := asira.App.DB
	loan := models.Loan{}

	err = db.Table("loans l").
		Select("l.*, ba.name as bank_name, s.name as service, p.name as product, a.category, a.name as agent_name, ap.name as agent_provider_name").
		Joins("LEFT JOIN products p ON l.product = p.id").
		Joins("LEFT JOIN services s ON p.service_id = s.id").
		Joins("LEFT JOIN banks ba ON l.bank = ba.id").
		Joins("LEFT JOIN borrowers b ON b.id = l.owner").
		Joins("LEFT JOIN agents a ON b.agent_id = a.id").
		Joins("LEFT JOIN agent_providers ap ON a.agent_provider = ap.id").
		Where("l.bank = ?", bankRep.BankID).
		Where("l.id = ?", loanID).
		Find(&loan).Error

	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "query result error")
	}

	return c.JSON(http.StatusOK, loan)
}

// LenderLoanApproveReject approve or reject a loan
func LenderLoanApproveReject(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "lender_loan_approve_reject")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	lenderID, _ := strconv.Atoi(claims["jti"].(string))
	bankRep := models.BankRepresentatives{}
	bankRep.FindbyUserID(lenderID)

	loanID, _ := strconv.Atoi(c.Param("loan_id"))

	type Filter struct {
		Bank   sql.NullInt64 `json:"bank"`
		ID     int           `json:"id"`
		Status string        `json:"status"`
	}

	loan := models.Loan{}
	err = loan.FilterSearchSingle(&Filter{
		Bank: sql.NullInt64{
			Int64: int64(bankRep.BankID),
			Valid: true,
		},
		ID:     loanID,
		Status: "processing", // only search for processing loan.
	})

	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "query result error")
	}
	if loan.ID == 0 {
		return returnInvalidResponse(http.StatusNotFound, "", "not found")
	}

	status := c.Param("approve_reject")
	switch status {
	default:
		return returnInvalidResponse(http.StatusBadRequest, "", "not allowed status")
	case "approve":
		disburseDate, err := time.Parse("2006-01-02", c.QueryParam("disburse_date"))
		if err != nil {
			return returnInvalidResponse(http.StatusBadRequest, "", "error parsing disburse date")
		}
		loan.Approve(disburseDate)
	case "reject":
		reason := c.QueryParam("reason")
		if len(reason) < 1 {
			return returnInvalidResponse(http.StatusBadRequest, "", "please fill reject reason")
		}
		loan.Reject(reason)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("loan %v is %v", loanID, loan.Status)})
}

// LenderLoanRequestListDownload download loans csv
func LenderLoanRequestListDownload(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "lender_loan_request_list_download")
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
	var results []LoanRequestListCSV

	db = db.Table("loans l").
		Select("l.id, l.owner_name, ba.name as bank_name, l.status, l.loan_amount, l.installment, l.total_loan, l.due_date, l.layaway_plan, l.loan_intention, l.intention_details, b.monthly_income, b.other_income, b.other_incomesource, b.bank_accountnumber").
		Joins("INNER JOIN borrowers b ON b.id = l.owner").
		Joins("INNER JOIN banks ba ON ba.id = b.bank").
		Where("ba.id = ?", bankRep.BankID)

	// filters
	if status := c.QueryParam("status"); len(status) > 0 {
		db = db.Where("LOWER(l.status) = ?", strings.ToLower(status))
	}
	if ownerName := c.QueryParam("owner_name"); len(ownerName) > 0 {
		db = db.Where("LOWER(l.owner_name) = ?", strings.ToLower(ownerName))
	}
	if id := customSplit(c.QueryParam("id"), ","); len(id) > 0 {
		db = db.Where("l.id IN (?)", id)
	}
	if startDate := c.QueryParam("start_date"); len(startDate) > 0 {
		if endDate := c.QueryParam("end_date"); len(endDate) > 0 {
			db = db.Where("l.created_time BETWEEN ? AND ?", startDate, endDate)
		} else {
			db = db.Where("l.created_time BETWEEN ? AND ?", startDate, startDate)
		}
	}
	if startDisburseDate := c.QueryParam("start_disburse_date"); len(startDisburseDate) > 0 {
		if endDisburseDate := c.QueryParam("end_disburse_date"); len(endDisburseDate) > 0 {
			db = db.Where("l.disburse_date BETWEEN ? AND ?", startDisburseDate, endDisburseDate)
		} else {
			db = db.Where("l.disburse_date BETWEEN ? AND ?", startDisburseDate, startDisburseDate)
		}
	}
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")
	if len(orderby) > 0 && len(sort) > 0 {
		db = db.Order(fmt.Sprintf("%s %s", orderby, sort))
	}

	err = db.Find(&results).Error

	b, err := csvutil.Marshal(results)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, string(b))
}

// LenderLoanConfirmDisbursement confirm a loan disbursement
func LenderLoanConfirmDisbursement(c echo.Context) error {
	defer c.Request().Body.Close()

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	lenderID, _ := strconv.Atoi(claims["jti"].(string))
	bankRep := models.BankRepresentatives{}
	bankRep.FindbyUserID(lenderID)

	loanID, _ := strconv.Atoi(c.Param("loan_id"))

	type Filter struct {
		Bank           sql.NullInt64 `json:"bank"`
		ID             int           `json:"id"`
		Status         string        `json:"status"`
		DisburseStatus string        `json:"disburse_status"`
	}

	loan := models.Loan{}
	err := loan.FilterSearchSingle(&Filter{
		Bank: sql.NullInt64{
			Int64: int64(bankRep.BankID),
			Valid: true,
		},
		ID:             loanID,
		Status:         "approved",
		DisburseStatus: "processing", // only search for processing loan.
	})

	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "query result error")
	}
	if loan.ID == 0 {
		return returnInvalidResponse(http.StatusNotFound, "", "not found")
	}

	loan.DisburseConfirmed()

	return c.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("loan %v disbursement is %v", loanID, loan.DisburseStatus)})
}

// LenderLoanChangeDisburseDate change disburse date
func LenderLoanChangeDisburseDate(c echo.Context) error {
	defer c.Request().Body.Close()

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	lenderID, _ := strconv.Atoi(claims["jti"].(string))
	bankRep := models.BankRepresentatives{}
	bankRep.FindbyUserID(lenderID)

	loanID, err := strconv.Atoi(c.Param("loan_id"))

	type Filter struct {
		Bank           sql.NullInt64 `json:"bank"`
		ID             int           `json:"id"`
		DisburseStatus bool          `json:"disburse_status"`
	}

	loan := models.Loan{}
	err = loan.FilterSearchSingle(&Filter{
		Bank: sql.NullInt64{
			Int64: int64(bankRep.BankID),
			Valid: true,
		},
		ID:             loanID,
		DisburseStatus: false,
	})

	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "query result error")
	}

	disburseDate, err := time.Parse("2006-01-02", c.QueryParam("disburse_date"))
	if err != nil {
		return returnInvalidResponse(http.StatusBadRequest, err, "error parsing disburse date")
	}
	if err = loan.ChangeDisburseDate(disburseDate); err != nil {
		return returnInvalidResponse(http.StatusBadRequest, err, "error changing disburse date")
	}

	return c.JSON(http.StatusOK, loan)
}
