package adminhandlers

import (
	"asira_lender/asira"
	"asira_lender/models"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/ayannahindonesia/basemodel"
	"github.com/labstack/echo"
)

// LoanSelect select custom type
type LoanSelect struct {
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

// LoanGetAll get all loans
func LoanGetAll(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_loan_get_all")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

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
		Joins("LEFT JOIN agent_providers ap ON a.agent_provider = ap.id")

	if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
		db = db.Or("CAST(loans.id as varchar(255)) = ?", searchAll).
			Or("LOWER(b.fullname) LIKE ?", "%"+strings.ToLower(searchAll)+"%").
			Or("LOWER(ba.name) LIKE ?", "%"+strings.ToLower(searchAll)+"%").
			Or("LOWER(a.category) LIKE ?", "%"+strings.ToLower(searchAll)+"%").
			Or("LOWER(s.name) LIKE ?", "%"+strings.ToLower(searchAll)+"%").
			Or("LOWER(p.name) LIKE ?", "%"+strings.ToLower(searchAll)+"%").
			Or("LOWER(loans.status) LIKE ?", "%"+strings.ToLower(searchAll)+"%")
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
		if disburseStatus := c.QueryParam("disburse_status"); len(disburseStatus) > 0 {
			db = db.Where("LOWER(loans.disburse_status) LIKE ?", "%"+strings.ToLower(disburseStatus)+"%")
		}
		if bankName := c.QueryParam("bank_name"); len(bankName) > 0 {
			db = db.Where("LOWER(ba.name) LIKE ?", "%"+strings.ToLower(bankName)+"%")
		}
		if bankAccount := c.QueryParam("bank_account"); len(bankAccount) > 0 {
			db = db.Where("b.bank_accountnumber LIKE ?", "%"+strings.ToLower(bankAccount)+"%")
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
		if status := c.QueryParam("status"); len(status) > 0 {
			db = db.Where("loans.status = ?", strings.ToLower(status))
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
	}

	tempDB := db
	tempDB.Where("loans.deleted_at IS NULL").Count(&totalRows)

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

// LoanGetDetails get loan details by id
func LoanGetDetails(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_loan_get_details")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	loan := LoanSelect{}
	db := asira.App.DB

	loanID, _ := strconv.Atoi(c.Param("loan_id"))
	err = db.Table("loans").
		Select("loans.*, b.fullname as borrower_name, ba.name as bank_name, b.bank_accountnumber as bank_account, s.name as service, p.name as product, a.category, a.name as agent_name, ap.name as agent_provider_name").
		Joins("LEFT JOIN products p ON loans.product = p.id").
		Joins("LEFT JOIN services s ON p.service_id = s.id").
		Joins("LEFT JOIN borrowers b ON b.id = loans.borrower").
		Joins("LEFT JOIN banks ba ON b.bank = ba.id").
		Joins("LEFT JOIN agents a ON b.agent_referral = a.id").
		Joins("LEFT JOIN agent_providers ap ON a.agent_provider = ap.id").
		Where("loans.id = ?", loanID).
		Find(&loan).Error

	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "query result error")
	}

	return c.JSON(http.StatusOK, loan)
}
