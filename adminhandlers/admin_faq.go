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
			Description string `json:"description" condition:"LIKE"`
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

//FAQNew create new FAQ
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

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Kesalahan validasi")
	}

	marshal, _ := json.Marshal(faqPayload)
	json.Unmarshal(marshal, &faq)

	err = faq.Create()
	middlewares.SubmitKafkaPayload(faq, "faq_create")
	if err != nil {
		NLog("error", "FAQNew", fmt.Sprintf("error create : %v", err), c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat FAQ baru")
	}

	NAudittrail(models.FAQ{}, faq, c.Get("user").(*jwt.Token), "faq", fmt.Sprint(faq.ID), "create")

	return c.JSON(http.StatusCreated, faq)
}

// FAQDetail get FAQ detail by id
func FAQDetail(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_faq_detail")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	faqID, _ := strconv.ParseUint(c.Param("faq_id"), 10, 64)

	faq := models.FAQ{}
	err = faq.FindbyID(faqID)
	if err != nil {
		NLog("warning", "FAQDetail", fmt.Sprintf("FAQ %v not found : %v", faqID, err), c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusNotFound, err, "Tidak memiliki hak akses")
	}

	return c.JSON(http.StatusOK, faq)
}

// FAQPatch edit FAQ by id
func FAQPatch(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_faq_patch")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	faqID, _ := strconv.ParseUint(c.Param("faq_id"), 10, 64)

	faq := models.FAQ{}
	faqPayload := FAQPayload{}
	err = faq.FindbyID(faqID)
	if err != nil {
		NLog("warning", "FAQPatch", fmt.Sprintf("FAQ %v not found : %v", faqID, err), c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusNotFound, err, "Tidak memiliki hak akses")
	}

	origin := faq

	payloadRules := govalidator.MapData{
		"title":       []string{},
		"description": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &faqPayload)
	if validate != nil {
		NLog("warning", "FAQPatch", fmt.Sprintf("validation error : %v", validate), c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Kesalahan validasi")
	}

	if len(faqPayload.Title) > 0 {
		faq.Title = faqPayload.Title
	}
	if len(faqPayload.Description) > 0 {
		faq.Description = faqPayload.Description
	}

	err = middlewares.SubmitKafkaPayload(faq, "faq_update")
	if err != nil {
		NLog("error", "FAQPatch", fmt.Sprintf("kafka error : %v", err), c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update FAQ %v", faqID))
	}

	NAudittrail(origin, faq, c.Get("user").(*jwt.Token), "faq", fmt.Sprint(faq.ID), "update")

	return c.JSON(http.StatusOK, faq)
}

// FAQDelete delete FAQ
func FAQDelete(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_faq_delete")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	faqID, _ := strconv.ParseUint(c.Param("faq_id"), 10, 64)

	faq := models.FAQ{}
	err = faq.FindbyID(faqID)
	if err != nil {
		NLog("warning", "FAQDelete", fmt.Sprintf("delete FAQ %v error : %v", faqID, err), c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusNotFound, err, "Tidak memiliki hak akses")
	}

	err = middlewares.SubmitKafkaPayload(faq, "faq_delete")
	if err != nil {
		NLog("error", "FAQDelete", fmt.Sprintf("delete FAQ %v error : %v", faqID, err), c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal delete FAQ %v", faqID))
	}

	NAudittrail(faq, models.FAQ{}, c.Get("user").(*jwt.Token), "faq", fmt.Sprint(faq.ID), "delete")
	return c.JSON(http.StatusOK, faq)
}
