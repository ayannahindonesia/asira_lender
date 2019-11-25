package adminhandlers

import (
	"asira_lender/models"
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/asira-ayannah/basemodel"

	"github.com/labstack/echo"
)

// BorrowerGetAll get all borrowers
func BorrowerGetAll(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_borrower_get_all")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	var (
		borrower models.Borrower
		result   basemodel.PagedFindResult
	)

	if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
		type Filter struct {
			Fullname string `json:"fullname" condition:"LIKE,optional"`
			ID       string `json:"id" condition:"optional"`
		}
		result, err = borrower.PagedFilterSearch(page, rows, orderby, sort, &Filter{
			Fullname: searchAll,
			ID:       searchAll,
		})
	} else {
		type Filter struct {
			Fullname string `json:"fullname" condition:"LIKE"`
			ID       string `json:"id"`
		}
		result, err = borrower.PagedFilterSearch(page, rows, orderby, sort, &Filter{
			Fullname: c.QueryParam("fullname"),
			ID:       c.QueryParam("id"),
		})
	}

	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Borrower tidak Ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

// BorrowerGetDetails get borrower details by id
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
