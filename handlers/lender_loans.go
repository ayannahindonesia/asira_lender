package handlers

import (
	"asira_lender/asira"
	"asira_lender/models"
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
	// LoanRequestListCSV type
	LoanRequestListCSV struct {
		ID                uint64  `json:"id"`
		BorrowerName      string  `json:"borrower_name"`
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
		BorrowerName      string `json:"borrower_name"`
		BankName          string `json:"bank_name"`
		BankAccount       string `json:"bank_account"`
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
		Select("l.*, b.fullname as borrower_name, ba.name as bank_name, b.bank_accountnumber as bank_account, s.name as service, p.name as product, a.category, a.name as agent_name, ap.name as agent_provider_name").
		Joins("LEFT JOIN products p ON l.product = p.id").
		Joins("LEFT JOIN services s ON p.service_id = s.id").
		Joins("LEFT JOIN borrowers b ON b.id = l.borrower").
		Joins("LEFT JOIN banks ba ON b.bank = ba.id").
		Joins("LEFT JOIN agents a ON b.agent_referral = a.id").
		Joins("LEFT JOIN agent_providers ap ON a.agent_provider = ap.id").
		Where("b.bank = ?", bankRep.BankID)

	status := c.QueryParam("status")
	disburseStatus := c.QueryParam("disburse_status")
	if len(status) > 0 {
		db = db.Where("l.status = ?", strings.ToLower(status))
	}
	if len(disburseStatus) > 0 {
		db = db.Where("l.disburse_status = ?", strings.ToLower(disburseStatus))
	}

	if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
		// gorm havent support nested subquery yet.
		searchLike := "%" + strings.ToLower(searchAll) + "%"
		extraquery := fmt.Sprintf("CAST(l.id as varchar(255)) = ?") + // use searchAll #1
			fmt.Sprintf(" OR LOWER(b.fullname) LIKE ?") + // use searchLike #2
			fmt.Sprintf(" OR LOWER(s.name) LIKE ?") + // use searchLike #3
			fmt.Sprintf(" OR LOWER(p.name) LIKE ?") // use searchLike #4

		if len(status) > 0 {
			switch status {
			case "approved":
				if len(disburseStatus) < 1 {
					extraquery = extraquery + fmt.Sprintf(" OR LOWER(l.disburse_status) LIKE '%v'", searchLike)
				}
			case "rejected":
				extraquery = extraquery + fmt.Sprintf(" OR LOWER(a.category) LIKE '%v'", searchLike)
			}
		} else {
			extraquery = extraquery +
				fmt.Sprintf(" OR LOWER(l.status) LIKE '%v'", searchLike) +
				fmt.Sprintf(" OR LOWER(a.category) LIKE '%v'", searchLike)
		}

		db = db.Where(extraquery, searchAll, searchLike, searchLike, searchLike)
	} else {
		if borrower := c.QueryParam("borrower"); len(borrower) > 0 {
			db = db.Where("l.borrower = ?", borrower)
		}
		if borrowerName := c.QueryParam("borrower_name"); len(borrowerName) > 0 {
			db = db.Where("LOWER(b.fullname) LIKE ?", "%"+strings.ToLower(borrowerName)+"%")
		}
		if id := customSplit(c.QueryParam("id"), ","); len(id) > 0 {
			db = db.Where("l.id IN (?)", id)
		}
		if bankAccount := c.QueryParam("bank_account"); len(bankAccount) > 0 {
			db = db.Where("b.bank_accountnumber LIKE ?", "%"+strings.ToLower(bankAccount)+"%")
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
	loan := LoanSelect{}

	err = db.Table("loans l").
		Select("l.*, b.fullname as borrower_name, ba.name as bank_name, b.bank_accountnumber as bank_account, s.name as service, p.name as product, a.category, a.name as agent_name, ap.name as agent_provider_name").
		Joins("LEFT JOIN products p ON l.product = p.id").
		Joins("LEFT JOIN services s ON p.service_id = s.id").
		Joins("LEFT JOIN borrowers b ON b.id = l.borrower").
		Joins("LEFT JOIN banks ba ON b.bank = ba.id").
		Joins("LEFT JOIN agents a ON b.agent_referral = a.id").
		Joins("LEFT JOIN agent_providers ap ON a.agent_provider = ap.id").
		Where("b.bank = ?", bankRep.BankID).
		Where("l.id = ?", loanID).
		Find(&loan).Error

	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Pinjaman %v tidak ditemukan", loanID))
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

	db := asira.App.DB
	loan := models.Loan{}

	err = db.Table("loans l").
		Select("*").
		Joins("INNER JOIN borrowers b ON b.id = l.borrower").
		Joins("INNER JOIN banks ba ON b.bank = ba.id").
		Where("ba.id = ?", bankRep.BankID).
		Where("l.id = ?", loanID).
		Limit(1).
		Find(&loan).Error

	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Pinjaman %v tidak ditemukan", loanID))
	}
	if loan.ID == 0 {
		return returnInvalidResponse(http.StatusNotFound, "", fmt.Sprintf("Pinjaman %v tidak ditemukan", loanID))
	}

	status := c.Param("approve_reject")
	switch status {
	default:
		return returnInvalidResponse(http.StatusBadRequest, "", fmt.Sprintf("Status %v tidak dapat digunakan", status))
	case "approve":
		disburseDate, err := time.Parse("2006-01-02", c.QueryParam("disburse_date"))
		if err != nil {
			return returnInvalidResponse(http.StatusBadRequest, "", "Terjadi kesalahan")
		}
		loan.Approve(disburseDate)
	case "reject":
		reason := c.QueryParam("reason")
		if len(reason) < 1 {
			return returnInvalidResponse(http.StatusBadRequest, "", "Harap mengisi alasan menolak")
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
		Select("l.id, b.fullname as borrower_name, ba.name as bank_name, l.status, l.loan_amount, l.installment, l.total_loan, l.due_date, l.layaway_plan, l.loan_intention, l.intention_details, b.monthly_income, b.other_income, b.other_incomesource, b.bank_accountnumber").
		Joins("INNER JOIN borrowers b ON b.id = l.borrower").
		Joins("INNER JOIN banks ba ON ba.id = b.bank").
		Where("ba.id = ?", bankRep.BankID)

	// filters
	if status := c.QueryParam("status"); len(status) > 0 {
		db = db.Where("LOWER(l.status) = ?", strings.ToLower(status))
	}
	if borrowerName := c.QueryParam("borrower_name"); len(borrowerName) > 0 {
		db = db.Where("LOWER(l.borrower_name) = ?", strings.ToLower(borrowerName))
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

	db := asira.App.DB
	loan := models.Loan{}

	err = db.Table("loans l").
		Select("*").
		Joins("INNER JOIN borrowers b ON b.id = l.borrower").
		Joins("INNER JOIN banks ba ON b.bank = ba.id").
		Where("ba.id = ?", bankRep.BankID).
		Where("l.id = ?", loanID).
		Where("l.status = ?", "approved").
		Where("l.disburse_status = ?", "processing").
		Limit(1).
		Find(&loan).Error

	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Pinjaman %v tidak ditemukan", loanID))
	}
	if loan.ID == 0 {
		return returnInvalidResponse(http.StatusNotFound, "", fmt.Sprintf("Pinjaman %v tidak ditemukan", loanID))
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

	db := asira.App.DB
	loan := models.Loan{}

	err = db.Table("loans l").
		Select("*").
		Joins("INNER JOIN borrowers b ON b.id = l.borrower").
		Joins("INNER JOIN banks ba ON b.bank = ba.id").
		Where("ba.id = ?", bankRep.BankID).
		Where("l.id = ?", loanID).
		Where("l.disburse_status = ?", "processing").
		Limit(1).
		Find(&loan).Error

	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Pinjaman %v tidak ditemukan", loanID))
	}

	disburseDate, err := time.Parse("2006-01-02", c.QueryParam("disburse_date"))
	if err != nil {
		return returnInvalidResponse(http.StatusBadRequest, err, "Terjadi kesalahan")
	}
	if err = loan.ChangeDisburseDate(disburseDate); err != nil {
		return returnInvalidResponse(http.StatusBadRequest, err, "Terjadi kesalahan")
	}

	return c.JSON(http.StatusOK, loan)
}
