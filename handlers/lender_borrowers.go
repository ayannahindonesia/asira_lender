package handlers

import (
	"asira_lender/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jszwec/csvutil"
	"github.com/labstack/echo"
	"gitlab.com/asira-ayannah/basemodel"
)

type (
	// BorrowerCSV custom type for query
	BorrowerCSV struct {
		basemodel.BaseModel
		DeletedTime          time.Time `json:"deleted_time"`
		Status               string    `json:"status"`
		Fullname             string    `json:"fullname"`
		Gender               string    `json:"gender"`
		IDCardNumber         string    `json:"idcard_number"`
		IDCardImageID        string    `json:"idcard_imageid"`
		TaxIDnumber          string    `json:"taxid_number"`
		TaxIDImageID         string    `json:"taxid_imageid"`
		Email                string    `json:"email"`
		Birthday             time.Time `json:"birthday"`
		Birthplace           string    `json:"birthplace"`
		LastEducation        string    `json:"last_education"`
		MotherName           string    `json:"mother_name"`
		Phone                string    `json:"phone"`
		MarriedStatus        string    `json:"marriage_status"`
		SpouseName           string    `json:"spouse_name"`
		SpouseBirthday       time.Time `json:"spouse_birthday"`
		SpouseLastEducation  string    `json:"spouse_lasteducation"`
		Dependants           int       `json:"dependants"`
		Address              string    `json:"address"`
		Province             string    `json:"province"`
		City                 string    `json:"city"`
		NeighbourAssociation string    `json:"neighbour_association"`
		Hamlets              string    `json:"hamlets"`
		HomePhoneNumber      string    `json:"home_phonenumber"`
		Subdistrict          string    `json:"subdistrict"`
		UrbanVillage         string    `json:"urban_village"`
		HomeOwnership        string    `json:"home_ownership"`
		LivedFor             int       `json:"lived_for"`
		Occupation           string    `json:"occupation"`
		EmployeeID           string    `json:"employee_id"`
		EmployerName         string    `json:"employer_name"`
		EmployerAddress      string    `json:"employer_address"`
		Department           string    `json:"department"`
		BeenWorkingFor       int       `json:"been_workingfor"`
		DirectSuperior       string    `json:"direct_superiorname"`
		EmployerNumber       string    `json:"employer_number"`
		MonthlyIncome        int       `json:"monthly_income"`
		OtherIncome          int       `json:"other_income"`
		OtherIncomeSource    string    `json:"other_incomesource"`
		FieldOfWork          string    `json:"field_of_work"`
		RelatedPersonName    string    `json:"related_personname"`
		RelatedRelation      string    `json:"related_relation"`
		RelatedPhoneNumber   string    `json:"related_phonenumber"`
		RelatedHomePhone     string    `json:"related_homenumber"`
		RelatedAddress       string    `json:"related_address"`
		Bank                 int64
		BankAccountNumber    string `json:"bank_accountnumber"`
	}
)

// LenderBorrowerList load all borrowers
func LenderBorrowerList(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "lender_borrower_list")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	lenderID, _ := strconv.Atoi(claims["jti"].(string))
	bankRep := models.BankRepresentatives{}
	bankRep.FindbyUserID(lenderID)

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
			Int64: int64(bankRep.BankID),
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

// LenderBorrowerListDetail load borrower detail by id
func LenderBorrowerListDetail(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "lender_borrower_list_detail")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	lenderID, _ := strconv.Atoi(claims["jti"].(string))
	bankRep := models.BankRepresentatives{}
	bankRep.FindbyUserID(lenderID)

	borrowerID, err := strconv.Atoi(c.Param("borrower_id"))
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "error parsing borrower id")
	}
	type Filter struct {
		Bank sql.NullInt64 `json:"bank"`
		ID   int           `json:"id"`
	}

	borrower := models.Borrower{}
	err = borrower.FilterSearchSingle(&Filter{
		Bank: sql.NullInt64{
			Int64: int64(bankRep.BankID),
			Valid: true,
		},
		ID: borrowerID,
	})

	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "query result error")
	}

	return c.JSON(http.StatusOK, borrower)
}

// LenderBorrowerListDownload download borrower list csv
func LenderBorrowerListDownload(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "lender_borrower_list_download")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	lenderID, _ := strconv.Atoi(claims["jti"].(string))

	// pagination parameters
	rows, _ := strconv.Atoi(c.QueryParam("rows"))
	page, _ := strconv.Atoi(c.QueryParam("page"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	// filters
	fullname := c.QueryParam("fullname")
	status := c.QueryParam("status")
	id := c.QueryParam("id")

	type Filter struct {
		Bank     sql.NullInt64 `json:"bank"`
		Fullname string        `json:"fullname"`
		Status   string        `json:"status"`
		ID       string        `json:"id"`
	}

	borrower := models.Borrower{}
	result, _ := borrower.PagedFilterSearch(page, rows, orderby, sort, &Filter{
		Bank: sql.NullInt64{
			Int64: int64(lenderID),
			Valid: true,
		},
		Fullname: fullname,
		Status:   status,
		ID:       id,
	})

	var data []BorrowerCSV
	data = mapnewBorrowerStruct(*result.Data.(*[]models.Borrower))

	b, err := csvutil.Marshal(data)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, string(b))
}

// LenderApproveRejectProspectiveBorrower approve or reject prospective borrower
func LenderApproveRejectProspectiveBorrower(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "lender_prospective_borrower_approval")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	lenderID, _ := strconv.Atoi(claims["jti"].(string))
	bankRep := models.BankRepresentatives{}
	bankRep.FindbyUserID(lenderID)

	borrowerID, err := strconv.Atoi(c.Param("borrower_id"))
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "error parsing borrower id")
	}
	type Filter struct {
		Bank              sql.NullInt64 `json:"bank"`
		ID                int           `json:"id"`
		BankAccountNumber string        `json:"bank_accountnumber"`
	}

	borrower := models.Borrower{}
	err = borrower.FilterSearchSingle(&Filter{
		Bank: sql.NullInt64{
			Int64: int64(bankRep.BankID),
			Valid: true,
		},
		ID:                borrowerID,
		BankAccountNumber: "",
	})
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "borrower not found")
	}

	if accNumber := c.QueryParam("account_number"); len(accNumber) > 0 {
		borrower.BankAccountNumber = accNumber
		approval := c.Param("approval")
		switch approval {
		default:
			borrower.Approve()
		case "reject":
			borrower.Reject()
		}
	} else {
		return returnInvalidResponse(http.StatusUnprocessableEntity, "", "invalid account number")
	}

	return c.JSON(http.StatusOK, borrower)
}

func mapnewBorrowerStruct(m []models.Borrower) []BorrowerCSV {
	var r []BorrowerCSV
	for _, v := range m {
		var unmarsh BorrowerCSV
		b, _ := json.Marshal(v)
		json.Unmarshal(b, &unmarsh)
		unmarsh.Bank = v.Bank.Int64
		r = append(r, unmarsh)
	}
	return r
}
