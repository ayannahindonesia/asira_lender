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

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	lenderID, _ := strconv.Atoi(claims["jti"].(string))

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

	type Filter struct {
		Bank        sql.NullInt64           `json:"bank"`
		Status      string                  `json:"status"`
		Owner       string                  `json:"owner"`
		OwnerName   string                  `json:"owner_name" condition:"LIKE"`
		DateBetween basemodel.CompareFilter `json:"created_time" condition:"BETWEEN"`
		ID          string                  `json:"id"`
	}

	loan := models.Loan{}
	result, err := loan.PagedFilterSearch(page, rows, orderby, sort, &Filter{
		Bank: sql.NullInt64{
			Int64: int64(lenderID),
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
	})

	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "query result error")
	}

	return c.JSON(http.StatusOK, result)
}

func LenderLoanRequestListDetail(c echo.Context) error {
	defer c.Request().Body.Close()

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	lenderID, _ := strconv.Atoi(claims["jti"].(string))

	loan_id, err := strconv.Atoi(c.Param("loan_id"))

	type Filter struct {
		Bank sql.NullInt64 `json:"bank"`
		ID   int           `json:"id"`
	}

	loan := models.Loan{}
	err = loan.FilterSearchSingle(&Filter{
		Bank: sql.NullInt64{
			Int64: int64(lenderID),
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

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	lenderID, _ := strconv.Atoi(claims["jti"].(string))

	loan_id, _ := strconv.Atoi(c.Param("loan_id"))

	type Filter struct {
		Bank   sql.NullInt64 `json:"bank"`
		ID     int           `json:"id"`
		Status string        `json:"status"`
	}

	loan := models.Loan{}
	err := loan.FilterSearchSingle(&Filter{
		Bank: sql.NullInt64{
			Int64: int64(lenderID),
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

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	lenderID, _ := strconv.Atoi(claims["jti"].(string))

	db := asira.App.DB
	var results []LoanRequestListCSV

	db = db.Table("loans l").
		Select("l.id, l.owner_name, ba.name as bank_name, l.status, l.loan_amount, l.installment, l.total_loan, l.due_date, l.layaway_plan, l.loan_intention, l.intention_details, b.monthly_income, b.other_income, b.other_incomesource, b.bank_accountnumber").
		Joins("INNER JOIN borrowers b ON b.id = l.owner").
		Joins("INNER JOIN banks ba ON ba.id = b.bank").
		Where("ba.id = ?", lenderID)

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
	if start_date := c.QueryParam("start_date"); len(start_date) > 0 {
		if end_date := c.QueryParam("end_date"); len(end_date) > 0 {
			db = db.Where("l.created_time BETWEEN ? AND ?", start_date, end_date)
		} else {
			db = db.Where("l.created_time BETWEEN ? AND ?", start_date, start_date)
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
