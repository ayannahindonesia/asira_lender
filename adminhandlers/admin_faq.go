package adminhandlers

import (
	"asira_lender/middlewares"
	"asira_lender/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ayannahindonesia/basemodel"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
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
			Title       string `json:"title" condition:"LIKE"`
			Description string `json:"description"`
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

// LoanPurposeNew create new loan purpose
func FAQNew(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_faq_new")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	faq := models.FAQ{}
	faqPayload := FAQPayload{}
	payloadRules := govalidator.MapData{
		"title":       []string{"required"},
		"description": []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &faqPayload)
	if validate != nil {
		NLog("warning", "FAQNew", fmt.Sprintf("error validation : %v", validate), c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
	}

	marshal, _ := json.Marshal(faqPayload)
	json.Unmarshal(marshal, &faq)

	err = faq.Create()
	middlewares.SubmitKafkaPayload(faq, "faq_create")
	if err != nil {
		NLog("error", "FAQNew", fmt.Sprintf("error create : %v", err), c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat loan purpose baru")
	}

	return c.JSON(http.StatusCreated, faq)
}
