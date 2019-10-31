package admin_handlers

import (
	"asira_lender/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func LoanGetAll(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_loan_get_all")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	loan := models.Loan{}
	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")
	//owner ID / Borrower ID
	owner := c.QueryParam("owner")
	id := c.QueryParam("id")
	fullname := c.QueryParam("fullname")

	type Filter struct {
		Owner        string `json:"owner"`
		ID           string `json:"id"`
		BorrowerInfo string `json:"borrower_info::text" condition:"LIKE"`
	}
	result, err := loan.PagedFilterSearch(page, rows, orderby, sort, &Filter{
		Owner:        owner,
		ID:           id,
		BorrowerInfo: fullname,
	})
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Loan tidak Ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

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
