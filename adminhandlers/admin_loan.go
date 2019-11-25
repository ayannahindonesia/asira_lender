package adminhandlers

import (
	"asira_lender/models"
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/asira-ayannah/basemodel"

	"github.com/labstack/echo"
)

// LoanGetAll get all loans
func LoanGetAll(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_loan_get_all")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	var (
		loan   models.Loan
		result basemodel.PagedFindResult
	)

	if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
		type Filter struct {
			Owner     string `json:"owner" condition:"optional"`
			ID        string `json:"id" condition:"optional"`
			Status    string `json:"status" condition:"optional"`
			OwnerName string `json:"owner_name" condition:"LIKE,optional"`
		}
		result, err = loan.PagedFilterSearch(page, rows, orderby, sort, &Filter{
			Owner:     searchAll,
			ID:        searchAll,
			Status:    searchAll,
			OwnerName: searchAll,
		})
	} else {
		type Filter struct {
			Owner     []string `json:"owner"`
			ID        []string `json:"id"`
			Status    string   `json:"status"`
			OwnerName string   `json:"owner_name" condition:"LIKE"`
		}
		result, err = loan.PagedFilterSearch(page, rows, orderby, sort, &Filter{
			Owner:     customSplit(c.QueryParam("owner"), ","),
			ID:        customSplit(c.QueryParam("id"), ","),
			Status:    c.QueryParam("status"),
			OwnerName: c.QueryParam("owner_name"),
		})
	}

	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Loan tidak Ditemukan")
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

	loanModel := models.Loan{}

	loanID, _ := strconv.Atoi(c.Param("loan_id"))
	err = loanModel.FindbyID(loanID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Loan ID tidak ditemukan")
	}

	return c.JSON(http.StatusOK, loanModel)
}
