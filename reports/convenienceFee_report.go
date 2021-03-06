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

	"github.com/ayannahindonesia/basemodel"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// ConvenienceFeeReport for conv fee
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
		CreatedTime    time.Time `json:"created_at"`
		Plafond        float64   `json:"plafond"`
		ConvenienceFee float64   `json:"convenience_fee"`
	}
	var results []ConvenienceFeeReport
	var totalRows int
	var offset int
	var rows int
	var page int
	var lastPage int

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
		Select("ba.name as bank_name, s.name as service_name, p.name as product_name, loans.id as loan_id, loans.created_at, loan_amount as plafond, value->>'amount' as convenience_fee").
		Joins("JOIN LATERAL jsonb_array_elements(loans.fees) j ON true").
		Joins("INNER JOIN borrowers b ON b.id = loans.borrower").
		Joins("INNER JOIN banks ba ON ba.id = b.bank").
		Joins("INNER JOIN products p ON p.id = loans.product").
		Joins("INNER JOIN services s ON s.id = p.service_id").
		Where("LOWER(value->>'description') LIKE ?", "convenience%").
		Where("loans.status = ?", "approved")

	// filters
	if bankName := c.QueryParam("bank_name"); len(bankName) > 0 {
		// di frontend pakai dropdown. jadi pake =
		db = db.Where("LOWER(ba.name) = ?", strings.ToLower(bankName))
	}
	if serviceName := c.QueryParam("service_name"); len(serviceName) > 0 {
		db = db.Where("LOWER(s.name) LIKE ?", "%"+strings.ToLower(serviceName)+"%")
	}
	if productName := c.QueryParam("product_name"); len(productName) > 0 {
		db = db.Where("LOWER(p.name) LIKE ?", "%"+strings.ToLower(productName)+"%")
	}
	if loanID := c.QueryParam("loan_id"); len(loanID) > 0 {
		db = db.Where("loans.id = ?", loanID)
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
		db = db.Where("loans.created_at BETWEEN ? AND ?", startDate, endDate)
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

	tempDB := db
	tempDB.Where("loans.deleted_at IS NULL").Count(&totalRows)

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
		lastPage = int(math.Ceil(float64(totalRows) / float64(rows)))
	}
	err = db.Find(&results).Error
	if err != nil {
		log.Println(err)
	}

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
