package reports

import (
	"asira_lender/asira"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitlab.com/asira-ayannah/basemodel"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func ConvenienceFeeReport(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "convenience_fee_report")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	db := asira.App.DB

	type ConvenienceFeeReport struct {
		BankName       string    `json:"bank_name"`
		ServiceName    string    `json:"service_name"`
		ProductName    string    `json:"product_name"`
		LoanID         string    `json:"loan_id"`
		CreatedTime    time.Time `json:"created_time"`
		Plafond        float64   `json:"plafond"`
		ConvenienceFee float64   `json:"convenience_fee"`
	}
	var results []ConvenienceFeeReport
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

	db = db.Table("loans l").
		Select("b.name as bank_name, s.name as service_name, p.name as product_name, l.id as loan_id, l.created_time, loan_amount as plafond, value->>'amount' as convenience_fee").
		Joins("JOIN LATERAL jsonb_array_elements(l.fees) j ON true").
		Joins("INNER JOIN banks b ON b.id = l.bank").
		Joins("INNER JOIN products p ON p.id = l.product").
		Joins("INNER JOIN services s ON s.id = p.service_id").
		Where("LOWER(value->>'description') LIKE 'convenience%'")

	// filters
	if bankName := c.QueryParam("bank_name"); len(bankName) > 0 {
		db = db.Where("LOWER(b.name) = ?", strings.ToLower(bankName))
	}
	if serviceName := c.QueryParam("service_name"); len(serviceName) > 0 {
		db = db.Where("LOWER(s.name) LIKE ?", "%"+strings.ToLower(serviceName)+"%")
	}
	if productName := c.QueryParam("product_name"); len(productName) > 0 {
		db = db.Where("LOWER(p.name) LIKE ?", "%"+strings.ToLower(productName)+"%")
	}
	if loanID := c.QueryParam("loan_id"); len(loanID) > 0 {
		db = db.Where("l.id = ?", loanID)
	}
	if plafond := c.QueryParam("plafond"); len(plafond) > 0 {
		db = db.Where("loan_amount = ?", plafond)
	}
	if convenienceFee := c.QueryParam("convenience_fee"); len(convenienceFee) > 0 {
		db = db.Where("value->>'amount' = ?", convenienceFee)
	}
	if startDate := c.QueryParam("start_date"); len(startDate) > 0 {
		endDate := c.QueryParam("end_date")
		if len(endDate) < 1 {
			endDate = startDate
		}
		db = db.Where("l.created_time BETWEEN ? AND ?", startDate, endDate)
	}
	if start_disburse_date := c.QueryParam("start_disburse_date"); len(start_disburse_date) > 0 {
		if end_disburse_date := c.QueryParam("end_disburse_date"); len(end_disburse_date) > 0 {
			db = db.Where("l.disburse_date BETWEEN ? AND ?", start_disburse_date, end_disburse_date)
		} else {
			db = db.Where("l.disburse_date BETWEEN ? AND ?", start_disburse_date, start_disburse_date)
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

	response := basemodel.PagedFindResult{
		TotalData:   totalRows,
		Rows:        rows,
		CurrentPage: page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + rows,
		Data:        results,
	}

	return c.JSON(http.StatusOK, response)
}

func validatePermission(c echo.Context, permission string) error {
	user := c.Get("user")
	token := user.(*jwt.Token)

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if claimPermissions, ok := claims["permissions"]; ok {
			s := strings.Split(strings.Trim(fmt.Sprintf("%v", claimPermissions), "[]"), " ")
			for _, v := range s {
				if strings.ToLower(v) == strings.ToLower(permission) || strings.ToLower(v) == "all" {
					return nil
				}
			}
		}
		return fmt.Errorf("Permission Denied")
	}

	return fmt.Errorf("Permission Denied")
}

func returnInvalidResponse(httpcode int, details interface{}, message string) error {
	responseBody := map[string]interface{}{
		"message": message,
		"details": details,
	}

	return echo.NewHTTPError(httpcode, responseBody)
}
