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
	obj.ContainsKey("total_data").ValueEqual("total_data", 0)

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
	// auth.GET("/lender/loanrequest_list/1/detail").
	// 	Expect().
	// 	Status(http.StatusOK).JSON().Object()

	// not owned by lender
	auth.GET("/lender/loanrequest_list/2/detail").
		Expect().
		Status(http.StatusInternalServerError).JSON().Object()
}
