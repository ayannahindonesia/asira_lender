package adminhandlers

import (
	"asira_lender/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ayannahindonesia/basemodel"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

//FAQPayload payload
type FAQPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

//FAQList get FAQ list
func FAQList(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_faq_list")
	if err != nil {
		NLog("warning", "FAQList", fmt.Sprintf("unauthorized access FAQList : '%v'", err), c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")

	var (
		faq    models.FAQ
		result basemodel.PagedFindResult
	)

	if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
		type Filter struct {
			Title       string `json:"title" condition:"LIKE,optional"`
			Description string `json:"description" condition:"LIKE,optional"`
		}
		result, err = faq.PagedFindFilter(page, rows, orderby, sort, &Filter{
			Title:       searchAll,
			Description: searchAll,
		})
	} else {
		type Filter struct {
			Title       string `json:"name" condition:"LIKE"`
			Description string `json:"status"`
		}
		result, err = faq.PagedFindFilter(page, rows, orderby, sort, &Filter{
			Title:       c.QueryParam("title"),
			Description: c.QueryParam("description"),
		})
	}

	if err != nil {
		NLog("warning", "FAQList", fmt.Sprintf("error finding FAQ : %v", err), c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Pencarian tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}
