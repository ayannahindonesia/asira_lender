package handlers

import (
	"asira_lender/asira"
	"asira_lender/models"
	"database/sql"
	"fmt"
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
)

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

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	// filters
	status := c.QueryParam("status")
	owner := c.QueryParam("owner")
	ownerName := c.QueryParam("owner_name")
	id := c.QueryParam("id")
	start_date := c.QueryParam("start_date")
	end_date := c.QueryParam("end_date")
	start_disburse_date := c.QueryParam("start_disburse_date")
	end_disburse_date := c.QueryParam("end_disburse_date")

	type Filter struct {
		Bank                sql.NullInt64           `json:"bank"`
		Status              string                  `json:"status"`
		Owner               string                  `json:"owner"`
		OwnerName           string                  `json:"owner_name" condition:"LIKE"`
		DateBetween         basemodel.CompareFilter `json:"created_time" condition:"BETWEEN"`
		DisburseDateBetween basemodel.CompareFilter `json:"disburse_date" condition:"BETWEEN"`
		ID                  string                  `json:"id"`
	}

	loan := models.Loan{}
	result, err := loan.PagedFilterSearch(page, rows, orderby, sort, &Filter{
		Bank: sql.NullInt64{
			Int64: int64(bankRep.BankID),
			Valid: true,
		},
		Owner:     owner,
		Status:    status,
		OwnerName: ownerName,
		ID:        id,
		DateBetween: basemodel.CompareFilter{
			Value1: start_date,
			Value2: end_date,
		},
		DisburseDateBetween: basemodel.CompareFilter{
			Value1: start_disburse_date,
			Value2: end_disburse_date,
		},
	})

	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "query result error")
	}

	return c.JSON(http.StatusOK, result)
}

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

	loan_id, err := strconv.Atoi(c.Param("loan_id"))

	type Filter struct {
		Bank sql.NullInt64 `json:"bank"`
		ID   int           `json:"id"`
	}

	loan := models.Loan{}
	err = loan.FilterSearchSingle(&Filter{
		Bank: sql.NullInt64{
			Int64: int64(bankRep.BankID),
			Valid: true,
		},
		ID: loan_id,
	})

	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "query result error")
	}

	return c.JSON(http.StatusOK, loan)
}

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

	loan_id, _ := strconv.Atoi(c.Param("loan_id"))

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
		ID:     loan_id,
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

	return c.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("loan %v is %v", loan_id, loan.Status)})
}

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
		Where("b.id = ?", bankRep.BankID)

	// filters
	if status := c.QueryParam("status"); len(status) > 0 {
		db = db.Where("LOWER(l.status) = ?", strings.ToLower(status))
	}
	if ownerName := c.QueryParam("owner_name"); len(ownerName) > 0 {
		db = db.Where("LOWER(l.owner_name) = ?", strings.ToLower(ownerName))
	}
	if id := customSplit(c.QueryParam("id"), ","); len(id) > 0 {
		db = db.Where("LOWER(l.id) IN ?", id)
	}
	if start_date := c.QueryParam("start_date"); len(start_date) > 0 {
		if end_date := c.QueryParam("end_date"); len(end_date) > 0 {
			db = db.Where("l.created_time BETWEEN ? AND ?", start_date, end_date)
		} else {
			db = db.Where("l.created_time BETWEEN ? AND ?", start_date, start_date)
		}
	}
	if start_disburse_date := c.QueryParam("start_disburse_date"); len(start_disburse_date) > 0 {
		if end_disburse_date := c.QueryParam("end_disburse_date"); len(end_disburse_date) > 0 {
			db = db.Where("l.disburse_date BETWEEN ? AND ?", start_disburse_date, end_disburse_date)
		} else {
			db = db.Where("l.disburse_date BETWEEN ? AND ?", start_disburse_date, start_disburse_date)
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

func LenderLoanConfirmDisbursement(c echo.Context) error {
	defer c.Request().Body.Close()

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	lenderID, _ := strconv.Atoi(claims["jti"].(string))
	bankRep := models.BankRepresentatives{}
	bankRep.FindbyUserID(lenderID)

	loan_id, _ := strconv.Atoi(c.Param("loan_id"))

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
		ID:             loan_id,
		Status:         "approved",
		DisburseStatus: "processing", // only search for processing loan.
	})

	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "query result error")
	}
	if loan.ID == 0 {
		return returnInvalidResponse(http.StatusNotFound, "", "not found")
	}

	loan.DisburseStatus = "confirmed"
	loan.Save()

	return c.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("loan %v disbursement is %v", loan_id, loan.DisburseStatus)})
}

// LenderLoanChangeDisburseDate func
func LenderLoanChangeDisburseDate(c echo.Context) error {
	defer c.Request().Body.Close()

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	lenderID, _ := strconv.Atoi(claims["jti"].(string))
	bankRep := models.BankRepresentatives{}
	bankRep.FindbyUserID(lenderID)

	loan_id, err := strconv.Atoi(c.Param("loan_id"))

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
		ID:             loan_id,
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
