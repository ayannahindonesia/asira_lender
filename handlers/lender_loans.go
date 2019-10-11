package handlers

import (
	"asira_lender/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jszwec/csvutil"
	"github.com/labstack/echo"
	"gitlab.com/asira-ayannah/basemodel"
)

type (
	LoanCSV struct {
		basemodel.BaseModel
		DeletedTime time.Time `json:"deleted_time"`
		Owner       int64
		OwnerName   string `json:"owner_name"`
		Bank        int64
		Status      string  `json:"status"`
		LoanAmount  float64 `json:"loan_amount"`
		Installment int     `json:"installment"`
		Fees        string
		Interest    float64   `json:"interest"`
		TotalLoan   float64   `json:"total_loan"`
		DueDate     time.Time `json:"due_date"`
		LayawayPlan float64   `json:"layaway_plan"`
		// Product          uint64         `json:"product"`
		// Service          uint64         `json:"service"`
		LoanIntention    string `json:"loan_intention"`
		IntentionDetails string `json:"intention_details"`
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
		loan.Reject()
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("loan %v is %v", loan_id, loan.Status)})
}

func LenderLoanRequestListDownload(c echo.Context) error {
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
	id := customSplit(c.QueryParam("id"), ",")
	start_date := c.QueryParam("start_date")
	end_date := c.QueryParam("end_date")

	type Filter struct {
		Bank        sql.NullInt64           `json:"bank"`
		Status      string                  `json:"status"`
		Owner       string                  `json:"owner"`
		OwnerName   string                  `json:"owner_name" condition:"LIKE"`
		DateBetween basemodel.CompareFilter `json:"created_time" condition:"BETWEEN"`
		ID          []string                `json:"id"`
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

	var data []LoanCSV
	data = mapnewLoanStruct(*result.Data.(*[]models.Loan))

	b, err := csvutil.Marshal(data)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, string(b))
}

func mapnewLoanStruct(m []models.Loan) []LoanCSV {
	var r []LoanCSV
	for _, v := range m {
		var unmarsh LoanCSV
		b, _ := json.Marshal(v)
		json.Unmarshal(b, &unmarsh)
		unmarsh.Owner = v.Owner.Int64
		unmarsh.Bank = v.Bank.Int64
		unmarsh.Fees = string(v.Fees.RawMessage)
		r = append(r, unmarsh)
	}
	return r
}
