package handlers

import (
	"asira_lender/models"
	"database/sql"
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func LenderBorrowerList(c echo.Context) error {
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
	fullname := c.QueryParam("fullname")
	status := c.QueryParam("status")
	id := c.QueryParam("id")

	type Filter struct {
		Bank     sql.NullInt64 `json:"bank"`
		Fullname string        `json:"fullname" condition:"LIKE"`
		Status   string        `json:"status"`
		ID       string        `json:"id"`
	}

	borrower := models.Borrower{}
	result, err := borrower.PagedFilterSearch(page, rows, orderby, sort, &Filter{
		Bank: sql.NullInt64{
			Int64: int64(lenderID),
			Valid: true,
		},
		Fullname: fullname,
		Status:   status,
		ID:       id,
	})

	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "query result error")
	}

	return c.JSON(http.StatusOK, result)
}

func LenderBorrowerListDetail(c echo.Context) error {
	defer c.Request().Body.Close()

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	lenderID, _ := strconv.Atoi(claims["jti"].(string))

	borrower_id, err := strconv.Atoi(c.Param("borrower_id"))
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "error parsing borrower id")
	}
	type Filter struct {
		Bank sql.NullInt64 `json:"bank"`
		ID   int           `json:"id"`
	}

	borrower := models.Borrower{}
	result, err := borrower.FilterSearchSingle(&Filter{
		Bank: sql.NullInt64{
			Int64: int64(lenderID),
			Valid: true,
		},
		ID: borrower_id,
	})

	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "query result error")
	}

	return c.JSON(http.StatusOK, result)
}
