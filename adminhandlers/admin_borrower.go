package adminhandlers

import (
	"asira_lender/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func BorrowerGetAll(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_borrower_get_all")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	borrower := models.Borrower{}
	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")
	// filters
	fullname := c.QueryParam("fullname")
	id := c.QueryParam("id")

	type Filter struct {
		Fullname string `json:"fullname" condition:"LIKE"`
		ID       string `json:"id"`
	}
	result, err := borrower.PagedFilterSearch(page, rows, orderby, sort, &Filter{
		Fullname: fullname,
		ID:       id,
	})

	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Borrower tidak Ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

func BorrowerGetDetails(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_borrower_get_details")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	borrowerModel := models.Borrower{}

	borrowerID, _ := strconv.Atoi(c.Param("borrower_id"))
	err = borrowerModel.FindbyID(borrowerID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Borrower ID tidak ditemukan")
	}

	return c.JSON(http.StatusOK, borrowerModel)
}
