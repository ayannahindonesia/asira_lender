package tests

import (
	"asira_lender/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestLenderGetBankList(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+adminBasicToken)
	})

	adminToken := getLenderAdminToken(e, auth)

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+adminToken)
	})

	// valid response
	auth.GET("/admin/banks").
		Expect().
		Status(http.StatusOK).JSON().Object()

	// test query found
	obj := auth.GET("/admin/banks").WithQuery("name", "Bank A").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 1)
	// test query found with part name
	obj := auth.GET("/admin/banks").WithQuery("name", "bank").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 2)

	// test query invalid
	obj = auth.GET("/admin/banks").WithQuery("name", "should not found this").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 0)
}

func TestNewBank(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+adminBasicToken)
	})

	adminToken := getLenderAdminToken(e, auth)

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+adminToken)
	})

	payload := map[string]interface{}{
		"name":     "Test New Bank",
		"type":     1,
		"address":  "testing st.",
		"province": "test province",
		"city":     "test city",
		"services": []string{"Service 1", "Service 2", "Service 3"},
		"products": []string{"Products 1", "Products 2", "Products 3"},
		"pic":      "test pic",
		"phone":    "08123454321",
	}

	// normal scenario
	obj := auth.POST("/admin/banks").WithJSON(payload).
		Expect().
		Status(http.StatusCreated).JSON().Object()
	obj.ContainsKey("name").ValueEqual("name", "Test New Bank")

	// test invalid
	payload = map[string]interface{}{
		"name": "",
	}
	auth.POST("/admin/banks").WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()
}

func TestGetBankbyID(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+adminBasicToken)
	})

	adminToken := getLenderAdminToken(e, auth)

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+adminToken)
	})

	// valid response
	obj := auth.GET("/admin/banks/1").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("id").ValueEqual("id", 1)

	// not found
	auth.GET("/admin/banks/9999").
		Expect().
		Status(http.StatusNotFound).JSON().Object()
}

func TestPatchBank(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+adminBasicToken)
	})

	adminToken := getLenderAdminToken(e, auth)

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+adminToken)
	})

	payload := map[string]interface{}{
		"name": "Test Patch",
	}

	// valid response
	obj := auth.PATCH("/admin/banks/1").WithJSON(payload).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("name").ValueEqual("name", "Test Patch")

	// test invalid token
	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer wrong token")
	})
	auth.PATCH("/admin/banks/1").WithJSON(payload).
		Expect().
		Status(http.StatusUnauthorized).JSON().Object()
}
