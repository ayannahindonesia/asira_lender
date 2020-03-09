package tests

import (
	"asira_lender/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestLenderServiceList(t *testing.T) {
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
	auth.GET("/lender/services").
		Expect().
		Status(http.StatusOK).JSON().Object()

	// test query found
	obj := auth.GET("/lender/services").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 2)

	// test query found
	obj = auth.GET("/lender/services").WithQuery("name", "Pinjaman PNS").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 1)

	// test query found with part name
	obj = auth.GET("/lender/services").WithQuery("name", "pinjaman").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 2)

	// test query invalid
	obj = auth.GET("/lender/services").WithQuery("name", "should not found this").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 0)
}

func TestLenderGetServicebyID(t *testing.T) {
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
	obj := auth.GET("/lender/services/1").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("id").ValueEqual("id", 1)

	// not found
	auth.GET("/lender/services/9999").
		Expect().
		Status(http.StatusNotFound).JSON().Object()
}
