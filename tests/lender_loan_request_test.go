package tests

import (
	"asira_lender/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestLenderGetLoanRequestList(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	lendertoken := getLenderLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+lendertoken)
	})

	// valid response
	obj := auth.GET("/lender/loanrequest_list").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 6)

	// wrong token
	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer thisisinvalidtoken")
	})

	auth.GET("/lender/loanrequest_list").
		Expect().
		Status(http.StatusUnauthorized).JSON().Object()
}

func TestLenderGetLoanRequestListDetail(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	lendertoken := getLenderLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+lendertoken)
	})

	// valid response
	auth.GET("/lender/loanrequest_list/1/detail").
		Expect().
		Status(http.StatusOK).JSON().Object()

	// not owned by lender
	auth.GET("/lender/loanrequest_list/7/detail").
		Expect().
		Status(http.StatusNotFound).JSON().Object()
}

func TestLenderApproveRejectLoan(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	lendertoken := getLenderLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+lendertoken)
	})

	// valid approve
	auth.GET("/lender/loanrequest_list/1/detail/approve").WithQuery("disburse_date", "2019-10-11").
		Expect().
		Status(http.StatusOK).JSON().Object()

	// valid reject
	auth.GET("/lender/loanrequest_list/3/detail/reject").WithQuery("reason", "reject reason").
		Expect().
		Status(http.StatusOK).JSON().Object()
}

func TestLenderChangeDisburseDate(t *testing.T) {
	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	lendertoken := getLenderLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+lendertoken)
	})

	auth.GET("/lender/loanrequest_list/1/detail/change_disburse_date").WithQuery("disburse_date", "2019-10-11").
		Expect().
		Status(http.StatusOK).JSON().Object()
}

func TestLenderConfirmDisbursement(t *testing.T) {
	// RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	lendertoken := getLenderLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+lendertoken)
	})

	// confirm disbursement
	auth.GET("/lender/loanrequest_list/1/detail/confirm_disbursement").
		Expect().
		Status(http.StatusOK).JSON().Object()
}
