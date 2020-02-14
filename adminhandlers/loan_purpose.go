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

// LoanPurposePayload handles request body
type LoanPurposePayload struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// LoanPurposeList get all loan purpose list
func LoanPurposeList(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_loan_purpose_list")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")

	var (
		purpose models.LoanPurpose
		result  basemodel.PagedFindResult
	)

	if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
		type Filter struct {
			Name   string `json:"name" condition:"LIKE,optional"`
			Status string `json:"status" condition:"optional"`
		}
		result, err = purpose.PagedFindFilter(page, rows, orderby, sort, &Filter{
			Name:   searchAll,
			Status: searchAll,
		})
	} else {
		type Filter struct {
			Name   string `json:"name" condition:"LIKE"`
			Status string `json:"status"`
		}
		result, err = purpose.PagedFindFilter(page, rows, orderby, sort, &Filter{
			Name:   c.QueryParam("name"),
			Status: c.QueryParam("status"),
		})
	}

	if err != nil {
		NLog("warning", "LoanPurposeList", fmt.Sprintf("error finding loan purposes : %v", err), c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Pencarian tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

// LoanPurposeNew create new loan purpose
func LoanPurposeNew(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_loan_purpose_new")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	purpose := models.LoanPurpose{}
	purposePayload := LoanPurposePayload{}
	payloadRules := govalidator.MapData{
		"name":   []string{"required"},
		"status": []string{"required", "loan_purpose_status"},
	}

	validate := validateRequestPayload(c, payloadRules, &purposePayload)
	if validate != nil {
		NLog("warning", "LoanPurposeNew", fmt.Sprintf("error validation : %v", validate), c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
	}

	marshal, _ := json.Marshal(purposePayload)
	json.Unmarshal(marshal, &purpose)

	err = purpose.Create()
	middlewares.SubmitKafkaPayload(purpose, "loan_purpose_create")
	if err != nil {
		NLog("error", "LoanPurposeNew", fmt.Sprintf("error create : %v", err), c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat loan purpose baru")
	}

	return c.JSON(http.StatusCreated, purpose)
}

// LoanPurposeDetail get loan purpose detail by id
func LoanPurposeDetail(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_loan_purpose_detail")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	loanPurposeID, _ := strconv.ParseUint(c.Param("loan_purpose_id"), 10, 64)

	purpose := models.LoanPurpose{}
	err = purpose.FindbyID(loanPurposeID)
	if err != nil {
		NLog("warning", "LoanPurposeDetail", fmt.Sprintf("loan purpose %v not found : %v", loanPurposeID, err), c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusNotFound, err, "Tidak memiliki hak akses")
	}

	return c.JSON(http.StatusOK, purpose)
}

// LoanPurposePatch edit loan purpose by id
func LoanPurposePatch(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_loan_purpose_patch")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	loanPurposeID, _ := strconv.ParseUint(c.Param("loan_purpose_id"), 10, 64)

	purpose := models.LoanPurpose{}
	purposePayload := LoanPurposePayload{}
	err = purpose.FindbyID(loanPurposeID)
	if err != nil {
		NLog("warning", "LoanPurposePatch", fmt.Sprintf("loan purpose %v not found : %v", loanPurposeID, err), c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusNotFound, err, "Tidak memiliki hak akses")
	}

	payloadRules := govalidator.MapData{
		"name":   []string{},
		"status": []string{"loan_purpose_status"},
	}

	validate := validateRequestPayload(c, payloadRules, &purposePayload)
	if validate != nil {
		NLog("warning", "LoanPurposePatch", fmt.Sprintf("validation error : %v", validate), c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
	}

	if len(purposePayload.Name) > 0 {
		purpose.Name = purposePayload.Name
	}
	if len(purposePayload.Status) > 0 {
		purpose.Status = purposePayload.Status
	}

	err = middlewares.SubmitKafkaPayload(purpose, "loan_purpose_update")
	if err != nil {
		NLog("error", "LoanPurposePatch", fmt.Sprintf("kafka error : %v", err), c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update loan purpose %v", loanPurposeID))
	}

	return c.JSON(http.StatusOK, purpose)
}

// LoanPurposeDelete delte loan purpose
func LoanPurposeDelete(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_loan_purpose_delete")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	loanPurposeID, _ := strconv.ParseUint(c.Param("loan_purpose_id"), 10, 64)

	purpose := models.LoanPurpose{}
	err = purpose.FindbyID(loanPurposeID)
	if err != nil {
		NLog("warning", "LoanPurposeDelete", fmt.Sprintf("delete loan purpose %v error : %v", loanPurposeID, err), c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusNotFound, err, "Tidak memiliki hak akses")
	}

	err = middlewares.SubmitKafkaPayload(purpose, "loan_purpose_delete")
	if err != nil {
		NLog("error", "LoanPurposeDelete", fmt.Sprintf("delete loan purpose %v error : %v", loanPurposeID, err), c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal delete loan purpose %v", loanPurposeID))
	}

	return c.JSON(http.StatusOK, purpose)
}
