package handlers

import (
	"asira_lender/adminhandlers"
	"asira_lender/asira"
	"asira_lender/middlewares"
	"asira_lender/models"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ayannahindonesia/basemodel"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

type (
	// LoanRequestListCSV type
	LoanRequestListCSV struct {
		ID                uint64         `json:"id"`
		BorrowerName      string         `json:"borrower_name"`
		BankName          string         `json:"bank_name"`
		Status            string         `json:"status"`
		LoanAmount        float64        `json:"loan_amount"`
		Installment       string         `json:"installment"`
		Fees              postgres.Jsonb `json:"fees"`
		Interest          float64        `json:"interest"`
		TotalLoan         float64        `json:"total_loan"`
		DueDate           time.Time      `json:"due_date"`
		LayawayPlan       float64        `json:"layaway_plan"`
		LoanIntention     string         `json:"loan_intention"`
		IntentionDetails  string         `json:"intention_details"`
		MonthlyIncome     string         `json:"monthly_income"`
		OtherIncome       string         `json:"other_income"`
		OtherIncomesource string         `json:"other_incomesource"`
		BankAccount       string         `json:"bank_account"`
	}
	// LoanSelect select custom type
	LoanSelect struct {
		models.Loan
		BorrowerName      string               `json:"borrower_name"`
		BankName          string               `json:"bank_name"`
		BankAccount       string               `json:"bank_account"`
		Service           string               `json:"service"`
		Product           string               `json:"product"`
		Category          string               `json:"category"`
		AgentName         string               `json:"agent_name"`
		AgentProviderName string               `json:"agent_provider_name"`
		Installments      []models.Installment `json:"installment_details"`
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

	db = db.Table("loans").
		Select("loans.*, b.fullname as borrower_name, ba.name as bank_name, b.bank_accountnumber as bank_account, s.name as service, p.name as product, a.category, a.name as agent_name, ap.name as agent_provider_name").
		Joins("LEFT JOIN products p ON loans.product = p.id").
		Joins("LEFT JOIN services s ON p.service_id = s.id").
		Joins("LEFT JOIN borrowers b ON b.id = loans.borrower").
		Joins("LEFT JOIN banks ba ON b.bank = ba.id").
		Joins("LEFT JOIN agents a ON b.agent_referral = a.id").
		Joins("LEFT JOIN agent_providers ap ON a.agent_provider = ap.id").
		Where("loans.otp_verified = ?", true).
		Where("b.bank = ?", bankRep.BankID)

	status := c.QueryParam("status")
	disburseStatus := c.QueryParam("disburse_status")
	if len(status) > 0 {
		db = db.Where("loans.status = ?", strings.ToLower(status))
	}
	if len(disburseStatus) > 0 {
		db = db.Where("loans.disburse_status = ?", strings.ToLower(disburseStatus))
	}

	if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
		// gorm havent support nested subquery yet.
		searchLike := "%" + strings.ToLower(searchAll) + "%"
		extraquery := fmt.Sprintf("CAST(loans.id as varchar(255)) = ?") + // use searchAll #1
			fmt.Sprintf(" OR LOWER(b.fullname) LIKE ?") + // use searchLike #2
			fmt.Sprintf(" OR LOWER(s.name) LIKE ?") + // use searchLike #3
			fmt.Sprintf(" OR LOWER(p.name) LIKE ?") // use searchLike #4

		if len(status) > 0 {
			switch status {
			case "approved":
				if len(disburseStatus) < 1 {
					extraquery = extraquery + fmt.Sprintf(" OR LOWER(loans.disburse_status) LIKE '%v'", searchLike)
				}
			case "rejected":
				extraquery = extraquery + fmt.Sprintf(" OR LOWER(a.category) LIKE '%v'", searchLike)
			}
		} else {
			extraquery = extraquery +
				fmt.Sprintf(" OR LOWER(loans.status) LIKE '%v'", searchLike) +
				fmt.Sprintf(" OR LOWER(a.category) LIKE '%v'", searchLike)
		}

		db = db.Where(extraquery, searchAll, searchLike, searchLike, searchLike)
	} else {
		if borrower := c.QueryParam("borrower"); len(borrower) > 0 {
			db = db.Where("loans.borrower = ?", borrower)
		}
		if borrowerName := c.QueryParam("borrower_name"); len(borrowerName) > 0 {
			db = db.Where("LOWER(b.fullname) LIKE ?", "%"+strings.ToLower(borrowerName)+"%")
		}
		if id := customSplit(c.QueryParam("id"), ","); len(id) > 0 {
			db = db.Where("loans.id IN (?)", id)
		}
		if bankAccount := c.QueryParam("bank_account"); len(bankAccount) > 0 {
			db = db.Where("b.bank_accountnumber LIKE ?", "%"+strings.ToLower(bankAccount)+"%")
		}
		if disburseStatus := c.QueryParam("disburse_status"); len(disburseStatus) > 0 {
			db = db.Where("LOWER(loans.disburse_status) LIKE ?", "%"+strings.ToLower(disburseStatus)+"%")
		}
		if startDate := c.QueryParam("start_date"); len(startDate) > 0 {
			if endDate := c.QueryParam("end_date"); len(endDate) > 0 {
				db = db.Where("loans.created_at BETWEEN ? AND ?", startDate, endDate)
			} else {
				db = db.Where("loans.created_at BETWEEN ? AND ?", startDate, startDate)
			}
		}
		if startDisburseDate := c.QueryParam("start_disburse_date"); len(startDisburseDate) > 0 {
			if endDisburseDate := c.QueryParam("end_disburse_date"); len(endDisburseDate) > 0 {
				db = db.Where("loans.disburse_date BETWEEN ? AND ?", startDisburseDate, endDisburseDate)
			} else {
				db = db.Where("loans.disburse_date BETWEEN ? AND ?", startDisburseDate, startDisburseDate)
			}
		}
		if startApprovalDate := c.QueryParam("start_approval_date"); len(startApprovalDate) > 0 {
			if endApprovalDate := c.QueryParam("end_approval_date"); len(endApprovalDate) > 0 {
				db = db.Where("loans.approval_date BETWEEN ? AND ?", startApprovalDate, endApprovalDate)
			} else {
				db = db.Where("loans.approval_date BETWEEN ? AND ?", startApprovalDate, startApprovalDate)
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
	tempDB.Where("loans.deleted_at IS NULL").Count(&totalRows)

	if rows > 0 {
		db = db.Limit(rows).Offset(offset)
		lastPage = int(math.Ceil(float64(totalRows) / float64(rows)))
	}
	err = db.Find(&loans).Error
	if err != nil {
		adminhandlers.NLog("warning", "LenderLoanRequestList", map[string]interface{}{"message": "error listing loans", "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusNotFound, err, "Tidak ada loan yang ditemukan.")
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
	installments := []models.Installment{}

	err = db.Table("loans").
		Select("loans.*, b.fullname as borrower_name, ba.name as bank_name, b.bank_accountnumber as bank_account, s.name as service, p.name as product, a.category, a.name as agent_name, ap.name as agent_provider_name").
		Joins("LEFT JOIN products p ON loans.product = p.id").
		Joins("LEFT JOIN services s ON p.service_id = s.id").
		Joins("LEFT JOIN borrowers b ON b.id = loans.borrower").
		Joins("LEFT JOIN banks ba ON b.bank = ba.id").
		Joins("LEFT JOIN agents a ON b.agent_referral = a.id").
		Joins("LEFT JOIN agent_providers ap ON a.agent_provider = ap.id").
		Where("loans.otp_verified = ?", true).
		Where("b.bank = ?", bankRep.BankID).
		Where("loans.id = ?", loanID).
		Find(&loan).Error

	if err != nil {
		adminhandlers.NLog("warning", "LenderLoanRequestListDetail", map[string]interface{}{"message": fmt.Sprintf("error while finding loan %v", loanID), "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Pinjaman %v tidak ditemukan", loanID))
	}

	err = db.Table("installments").
		Select("*").
		Where("id IN (?)", strings.Fields(strings.Trim(fmt.Sprint(loan.InstallmentDetails), "[]"))).
		Scan(&installments).Error
	if err != nil {
		adminhandlers.NLog("warning", "LenderLoanRequestListDetail", map[string]interface{}{"message": "query not found : '%v' error : %v", "query": db.QueryExpr(), "error": err}, c.Get("user").(*jwt.Token), "", false)
	}
	loan.Installments = installments

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

	err = db.Table("loans").
		Select("*").
		Joins("INNER JOIN borrowers b ON b.id = loans.borrower").
		Joins("INNER JOIN banks ba ON b.bank = ba.id").
		Where("loans.otp_verified = ?", true).
		Where("ba.id = ?", bankRep.BankID).
		Where("loans.id = ?", loanID).
		Limit(1).
		Find(&loan).Error

	if err != nil {
		adminhandlers.NLog("warning", "LenderLoanApproveReject", map[string]interface{}{"message": fmt.Sprintf("error while finding loan %v", loanID), "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Pinjaman %v tidak ditemukan", loanID))
	}
	origin := loan
	if loan.ID == 0 {
		adminhandlers.NLog("warning", "LenderLoanApproveReject", map[string]interface{}{"message": "loan id is 0"}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusNotFound, "", fmt.Sprintf("Pinjaman %v tidak ditemukan", loanID))
	}

	status := c.Param("approve_reject")
	switch status {
	default:
		adminhandlers.NLog("warning", "LenderLoanApproveReject", map[string]interface{}{"message": fmt.Sprintf("status %v is not allowed", status)}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusBadRequest, "", fmt.Sprintf("Status %v tidak dapat digunakan", status))
	case "approve":
		disburseDate, err := time.Parse("2006-01-02", c.QueryParam("disburse_date"))
		if err != nil {
			adminhandlers.NLog("warning", "LenderLoanApproveReject", map[string]interface{}{"message": fmt.Sprintf("error parsing disburse date %v", c.QueryParam("disburse_date")), "error": err}, c.Get("user").(*jwt.Token), "", false)

			return returnInvalidResponse(http.StatusBadRequest, "", "Terjadi kesalahan")
		}
		loan.Status = "approved"
		loan.DisburseDate = disburseDate
		loan.ApprovalDate = time.Now()

		err = middlewares.SubmitKafkaPayload(loan, "loan_update")
		if err != nil {
			adminhandlers.NLog("error", "LenderLoanApproveReject", map[string]interface{}{"message": "error submitting kafka", "error": err}, c.Get("user").(*jwt.Token), "", false)

			return returnInvalidResponse(http.StatusBadRequest, err, "Gagal approve pinjaman")
		}
	case "reject":
		reason := c.QueryParam("reason")
		if len(reason) < 1 {
			return returnInvalidResponse(http.StatusBadRequest, "", "Harap mengisi alasan menolak")
		}
		loan.Status = "rejected"
		loan.RejectReason = reason
		loan.ApprovalDate = time.Now()

		err = middlewares.SubmitKafkaPayload(loan, "loan_update")
		if err != nil {
			adminhandlers.NLog("error", "LenderLoanApproveReject", map[string]interface{}{"message": "error submitting kafka", "error": err, "loan": loan}, c.Get("user").(*jwt.Token), "", false)

			return returnInvalidResponse(http.StatusBadRequest, err, "Gagal reject pinjaman")
		}
	}

	adminhandlers.NAudittrail(origin, loan, c.Get("user").(*jwt.Token), "loan", fmt.Sprint(loan.ID), "approve borrower")

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

	db = db.Table("loans").
		Select("loans.*, b.fullname as borrower_name, ba.name as bank_name, b.monthly_income, b.other_income, b.other_incomesource, b.bank_accountnumber as bank_account").
		Joins("LEFT JOIN products p ON loans.product = p.id").
		Joins("LEFT JOIN services s ON p.service_id = s.id").
		Joins("LEFT JOIN borrowers b ON b.id = loans.borrower").
		Joins("LEFT JOIN banks ba ON b.bank = ba.id").
		Joins("LEFT JOIN agents a ON b.agent_referral = a.id").
		Joins("LEFT JOIN agent_providers ap ON a.agent_provider = ap.id").
		Where("loans.otp_verified = ?", true).
		Where("ba.id = ?", bankRep.BankID)

	// filters
	if status := c.QueryParam("status"); len(status) > 0 {
		db = db.Where("loans.status = ?", strings.ToLower(status))
	}
	if disburseStatus := c.QueryParam("disburse_status"); len(disburseStatus) > 0 {
		db = db.Where("loans.disburse_status = ?", strings.ToLower(disburseStatus))
	}
	if borrower := c.QueryParam("borrower"); len(borrower) > 0 {
		db = db.Where("loans.borrower = ?", borrower)
	}
	if borrowerName := c.QueryParam("borrower_name"); len(borrowerName) > 0 {
		db = db.Where("LOWER(b.fullname) LIKE ?", "%"+strings.ToLower(borrowerName)+"%")
	}
	if id := customSplit(c.QueryParam("id"), ","); len(id) > 0 {
		db = db.Where("loans.id IN (?)", id)
	}
	if bankAccount := c.QueryParam("bank_account"); len(bankAccount) > 0 {
		db = db.Where("b.bank_accountnumber LIKE ?", "%"+strings.ToLower(bankAccount)+"%")
	}
	if disburseStatus := c.QueryParam("disburse_status"); len(disburseStatus) > 0 {
		db = db.Where("LOWER(loans.disburse_status) LIKE ?", "%"+strings.ToLower(disburseStatus)+"%")
	}
	if startDate := c.QueryParam("start_date"); len(startDate) > 0 {
		if endDate := c.QueryParam("end_date"); len(endDate) > 0 {
			db = db.Where("loans.created_at BETWEEN ? AND ?", startDate, endDate)
		} else {
			db = db.Where("loans.created_at BETWEEN ? AND ?", startDate, startDate)
		}
	}
	if startDisburseDate := c.QueryParam("start_disburse_date"); len(startDisburseDate) > 0 {
		if endDisburseDate := c.QueryParam("end_disburse_date"); len(endDisburseDate) > 0 {
			db = db.Where("loans.disburse_date BETWEEN ? AND ?", startDisburseDate, endDisburseDate)
		} else {
			db = db.Where("loans.disburse_date BETWEEN ? AND ?", startDisburseDate, startDisburseDate)
		}
	}
	if startApprovalDate := c.QueryParam("start_approval_date"); len(startApprovalDate) > 0 {
		if endApprovalDate := c.QueryParam("end_approval_date"); len(endApprovalDate) > 0 {
			db = db.Where("loans.approval_date BETWEEN ? AND ?", startApprovalDate, endApprovalDate)
		} else {
			db = db.Where("loans.approval_date BETWEEN ? AND ?", startApprovalDate, startApprovalDate)
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
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")
	if len(orderby) > 0 && len(sort) > 0 {
		db = db.Order(fmt.Sprintf("%s %s", orderby, sort))
	}

	err = db.Scan(&results).Error

	// b, err := csvutil.Marshal(results)
	if err != nil {
		adminhandlers.NLog("warning", "LenderLoanRequestListDownload", map[string]interface{}{"message": "query error", "error": err, "query": db.QueryExpr()}, c.Get("user").(*jwt.Token), "", false)

		return err
	}
	adminhandlers.NLog("info", "LenderLoanRequestListDownload", map[string]interface{}{"message": "loan download"}, c.Get("user").(*jwt.Token), "", false)

	return c.JSON(http.StatusOK, results)
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

	err = db.Table("loans").
		Select("*").
		Joins("INNER JOIN borrowers b ON b.id = loans.borrower").
		Joins("INNER JOIN banks ba ON b.bank = ba.id").
		Where("loans.otp_verified = ?", true).
		Where("ba.id = ?", bankRep.BankID).
		Where("loans.id = ?", loanID).
		Where("loans.status = ?", "approved").
		Where("loans.disburse_status = ?", "processing").
		Limit(1).
		Find(&loan).Error

	if err != nil {
		adminhandlers.NLog("warning", "LenderLoanConfirmDisbursement", map[string]interface{}{"message": "error confirm disburse", "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Pinjaman %v tidak ditemukan", loanID))
	}
	origin := loan
	if loan.ID == 0 {
		adminhandlers.NLog("warning", "LenderLoanConfirmDisbursement", map[string]interface{}{"message": "loan id is 0"}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusNotFound, "", fmt.Sprintf("Pinjaman %v tidak ditemukan", loanID))
	}

	loan.DisburseStatus = "confirmed"

	err = middlewares.SubmitKafkaPayload(loan, "loan_update")
	if err != nil {
		adminhandlers.NLog("error", "LenderLoanConfirmDisbursement", map[string]interface{}{"message": "error submitting kafka", "error": err, "loan": loan}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusBadRequest, err, "Gagal confirm disbursement pinjaman")
	}
	adminhandlers.NLog("info", "LenderLoanConfirmDisbursement", map[string]interface{}{"message": fmt.Sprintf("confirmed loan %v", loan.ID)}, c.Get("user").(*jwt.Token), "", false)

	adminhandlers.NAudittrail(origin, loan, c.Get("user").(*jwt.Token), "loan", fmt.Sprint(loan.ID), "confirm loan")

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

	err = db.Table("loans").
		Select("*").
		Joins("INNER JOIN borrowers b ON b.id = loans.borrower").
		Joins("INNER JOIN banks ba ON b.bank = ba.id").
		Where("loans.otp_verified = ?", true).
		Where("ba.id = ?", bankRep.BankID).
		Where("loans.id = ?", loanID).
		Where("loans.disburse_status = ?", "processing").
		Limit(1).
		Find(&loan).Error

	if err != nil {
		adminhandlers.NLog("warning", "LenderLoanChangeDisburseDate", map[string]interface{}{"message": "error finding loan to change disburse date", "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Pinjaman %v tidak ditemukan", loanID))
	}
	origin := loan

	disburseDate, err := time.Parse("2006-01-02", c.QueryParam("disburse_date"))
	if err != nil {
		adminhandlers.NLog("warning", "LenderLoanChangeDisburseDate", map[string]interface{}{"message": fmt.Sprintf("error parsing disburse date %v", c.QueryParam("disburse_date")), "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusBadRequest, err, "Terjadi kesalahan")
	}

	loan.DisburseDate = disburseDate
	loan.DisburseDateChanged = true

	err = middlewares.SubmitKafkaPayload(loan, "loan_update")
	if err != nil {
		adminhandlers.NLog("error", "LenderLoanChangeDisburseDate", map[string]interface{}{"message": "error submitting kafka", "error": err, "loan": loan}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusBadRequest, err, "Gagal confirm disbursement pinjaman")
	}
	adminhandlers.NLog("info", "LenderLoanChangeDisburseDate", map[string]interface{}{"message": fmt.Sprintf("loan %v disburse date changed", loan.ID)}, c.Get("user").(*jwt.Token), "", false)

	adminhandlers.NAudittrail(origin, loan, c.Get("user").(*jwt.Token), "loan", fmt.Sprint(loan.ID), "change loan disburse date")

	return c.JSON(http.StatusOK, loan)
}

// LenderLoanInstallmentsApprove func
func LenderLoanInstallmentsApprove(c echo.Context) error {
	type InstallmentPayload struct {
		PaidStatus   bool    `json:"paid_status"`
		Underpayment float64 `json:"underpayment"`
		Note         string  `json:"note"`
	}
	var installmentPayload InstallmentPayload

	defer c.Request().Body.Close()
	err := validatePermission(c, "lender_loan_installment_approve")
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
	installmentID, err := strconv.Atoi(c.Param("installment_id"))
	db := asira.App.DB
	installment := models.Installment{}

	loansQ := db.Table("loans").
		Select("UNNEST(loans.installment_details)").
		Joins("LEFT JOIN borrowers b ON b.id = loans.borrower").
		Joins("LEFT JOIN banks ba ON b.bank = ba.id").
		Where("loans.otp_verified = ?", true).
		Where("b.bank = ?", bankRep.BankID).
		Where("loans.id = ?", loanID).
		QueryExpr()

	err = db.Table("installments").
		Select("*").
		Where("id IN (?)", loansQ).
		Where("id = ?", installmentID).
		Scan(&installment).Error
	if err != nil {
		adminhandlers.NLog("warning", "LenderLoanInstallmentsApprove", map[string]interface{}{"message": "query not found : '%v' error : %v", "query": db.QueryExpr(), "error": err}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Cicilan %v pada pinjaman %v tidak ditemukan", installmentID, loanID))
	}

	payloadRules := govalidator.MapData{
		"paid_status":  []string{"bool"},
		"underpayment": []string{"numeric"},
		"note":         []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &installmentPayload)
	if validate != nil {
		adminhandlers.NLog("warning", "LenderLoanInstallmentsApprove", map[string]interface{}{"message": "error validation", "error": validate}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
	}

	now := time.Now()
	if installmentPayload.PaidStatus {
		installment.PaidStatus = true
		installment.PaidDate = &now
	}
	if installmentPayload.Underpayment > 0 {
		installment.Underpayment = installmentPayload.Underpayment
	}
	if len(installmentPayload.Note) > 0 {
		installment.Note = installmentPayload.Note
	}

	err = middlewares.SubmitKafkaPayload(installment, "installment_update")
	if err != nil {
		adminhandlers.NLog("error", "LenderLoanInstallmentsApprove", map[string]interface{}{"message": "error submitting kafka borrower", "error": err, "installment": installment}, user.(*jwt.Token), "", false)

		returnInvalidResponse(http.StatusUnprocessableEntity, err, "Gagal reject borrower")
	}

	return c.JSON(http.StatusOK, installment)
}
