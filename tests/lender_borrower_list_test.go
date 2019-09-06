package tests

import (
	"asira_lender/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestBorrowerGetAll(t *testing.T) {
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

	// valid response of borrowers
	obj := auth.GET("/lender/borrower_list").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("to").ValueEqual("to", 25)

	// valid response of borrowers
	obj = auth.GET("/lender/borrower_list").WithQuery("fullname", "ame").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("to").ValueEqual("to", 25)
}

func TestBorrowerGetDetail(t *testing.T) {
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

	// valid response of borrowers
	obj := auth.GET("/lender/borrower_list/1/detail").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("id").ValueEqual("id", 1)

	// invalid response of borrowers
	obj = auth.GET("/lender/borrower_list/99/detail").
		Expect().
		Status(http.StatusInternalServerError).JSON().Object()

}
