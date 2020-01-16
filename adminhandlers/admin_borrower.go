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

type (
	// BorrowerSelect select for joining
	BorrowerSelect struct {
		models.Borrower
		Category          string `json:"category"`
		LoanCount         int    `json:"loan_count"`
		LoanStatus        string `json:"loan_status"`
		BankName          string `json:"bank_name"`
		AgentName         string `json:"agent_name"`
		AgentProviderName string `json:"agent_provider_name"`
	}
)

// BorrowerGetAll get all borrowers
func BorrowerGetAll(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_borrower_get_all")
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

	loanStatusQuery := fmt.Sprintf("CASE WHEN (SELECT COUNT(id) FROM loans l WHERE l.borrower = borrowers.id AND status IN ('%s', %s) AND (due_date IS NULL OR due_date = '0001-01-01 00:00:00+00' OR (NOW() > l.disburse_date AND NOW() < l.due_date + make_interval(days => 1)))) > 0 THEN '%s' ELSE '%s' END", "approved", "processing", "active", "inactive")

	db = db.Table("borrowers").
		Select("borrowers.*, a.category, ba.name as bank_name, a.name as agent_name, ap.name as agent_provider_name, (SELECT COUNT(id) FROM loans l WHERE l.borrower = borrowers.id AND l.status = ?) as loan_count, "+loanStatusQuery+" as loan_status", "approved").
		Joins("LEFT JOIN agents a ON borrowers.agent_referral = a.id").
		Joins("LEFT JOIN banks ba ON ba.id = borrowers.bank").
		Joins("LEFT JOIN agent_providers ap ON a.agent_provider = ap.id")

	accountNumber := c.QueryParam("account_number")
	if status := c.QueryParam("status"); len(status) > 0 {
		db = db.Where("borrowers.status = ?", strings.ToLower(status))
	} else {
		db = db.Where("borrowers.status != ?", "rejected")
	}

	if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
		searchLike := "%" + strings.ToLower(searchAll) + "%"
		extraquery := fmt.Sprintf("LOWER(borrowers.fullname) LIKE ?") + // use searchLike #1
			fmt.Sprintf(" OR LOWER(a.category) = ?") + // use searchLike #2
			fmt.Sprintf(" OR LOWER(ba.name) LIKE ?") + // use searchLike #3
			fmt.Sprintf(" OR LOWER(a.name) LIKE ?") + // use searchLike #4
			fmt.Sprintf(" OR LOWER(ap.name) LIKE ?") + // use searchLike #5
			fmt.Sprintf(" OR CAST(borrowers.id as varchar(255)) = ?") + // use searchAll #6
			fmt.Sprintf(" OR "+loanStatusQuery+" LIKE ?") // use searchLike #7

		if len(accountNumber) > 0 {
			if accountNumber == "null" {
				db = db.Where("borrowers.bank_accountnumber = ?", "")
			} else if accountNumber == "not null" {
				db = db.Where("borrowers.bank_accountnumber != ?", "")
			}
		}

		db = db.Where(extraquery, searchLike, searchLike, searchLike, searchLike, searchLike, searchAll, searchLike)
	} else {
		if fullname := c.QueryParam("fullname"); len(fullname) > 0 {
			db = db.Where("LOWER(borrowers.fullname) LIKE ?", "%"+strings.ToLower(fullname)+"%")
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
			db = db.Where("borrowers.id IN (?)", id)
		}
		if len(accountNumber) > 0 {
			if accountNumber == "null" {
				db = db.Where("borrowers.bank_accountnumber = ?", "")
			} else if accountNumber == "not null" {
				db = db.Where("borrowers.bank_accountnumber != ?", "")
			} else {
				db = db.Where("borrowers.bank_accountnumber LIKE ?", "%"+strings.ToLower(accountNumber)+"%")
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
	tempDB.Where("borrowers.deleted_at IS NULL").Count(&totalRows)

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

// BorrowerGetDetails show borrower by id
func BorrowerGetDetails(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_borrower_get_details")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}
	borrowerID, _ := strconv.Atoi(c.Param("borrower_id"))

	db := asira.App.DB

	borrower := BorrowerSelect{}

	loanStatusQuery := fmt.Sprintf("CASE WHEN (SELECT COUNT(id) FROM loans l WHERE l.borrower = borrowers.id AND status IN ('%s', %s) AND (due_date IS NULL OR due_date = '0001-01-01 00:00:00+00' OR (NOW() > l.disburse_date AND NOW() < l.due_date + make_interval(days => 1)))) > 0 THEN '%s' ELSE '%s' END", "approved", "processing", "active", "inactive")

	err = db.Table("borrowers").
		Select("borrowers.*, a.category, ba.name as bank_name, a.name as agent_name, ap.name as agent_provider_name, (SELECT COUNT(id) FROM loans l WHERE l.borrower = borrowers.id AND l.status = ?) as loan_count, "+loanStatusQuery+" as loan_status", "approved").
		Joins("LEFT JOIN agents a ON borrowers.agent_referral = a.id").
		Joins("LEFT JOIN banks ba ON ba.id = borrowers.bank").
		Joins("LEFT JOIN agent_providers ap ON a.agent_provider = ap.id").
		Where("borrowers.id = ?", borrowerID).
		Find(&borrower).Error
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("id %v not found.", borrowerID))
	}

	return c.JSON(http.StatusOK, borrower)
}
