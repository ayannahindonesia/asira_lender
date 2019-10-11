package tests

import (
	"asira_lender/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestClientLoanPurposeList(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	obj := auth.GET("/clientauth").
		Expect().
		Status(http.StatusOK).JSON().Object()

	clienttoken := obj.Value("token").String().Raw()

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+clienttoken)
	})

	// valid response
	obj = auth.GET("/client/loan_purposes").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").NotEmpty()

	auth.GET("/client/loan_purposes/1").
		Expect().
		Status(http.StatusOK).JSON().Object()

	// test without token
	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer not allowed")
	})
	auth.GET("/client/loan_purposes").
		Expect().
		Status(http.StatusUnauthorized).JSON().Object()
}
