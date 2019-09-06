package tests

import (
	"asira_lender/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestLenderGetBankServiceList(t *testing.T) {
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
	auth.GET("/admin/bank_services").
		Expect().
		Status(http.StatusOK).JSON().Object()

	// test query found
	obj := auth.GET("/admin/bank_services").WithQuery("name", "Pinjaman PNS").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 1)
	// test query found with part name
	obj := auth.GET("/admin/bank_services").WithQuery("name", "pinjaman").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 5)

	// test query invalid
	obj = auth.GET("/admin/bank_services").WithQuery("name", "should not found this").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 0)
}

func TestNewBankService(t *testing.T) {
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
		"name":   "Test New Bank Service",
		"image":  "this is a long long base64 encoded image string",
		"status": "active",
	}

	// normal scenario
	obj := auth.POST("/admin/bank_services").WithJSON(payload).
		Expect().
		Status(http.StatusCreated).JSON().Object()
	obj.ContainsKey("name").ValueEqual("name", "Test New Bank Service")

	// invalid status
	payload = map[string]interface{}{
		"name":   "Test New Bank Service",
		"image":  "this is a long long base64 encoded image string",
		"status": "not valid",
	}
	auth.PATCH("/admin/bank_services/1").WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()

	// test invalid
	payload = map[string]interface{}{
		"name": "",
	}
	auth.POST("/admin/bank_services").WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()
}

func TestGetBankServicebyID(t *testing.T) {
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
	obj := auth.GET("/admin/bank_services/1").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("id").ValueEqual("id", 1)

	// not found
	auth.GET("/admin/bank_services/9999").
		Expect().
		Status(http.StatusNotFound).JSON().Object()
}

func TestPatchBankService(t *testing.T) {
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
		"name": "Test Service Patch",
	}

	// valid response
	obj := auth.PATCH("/admin/bank_services/1").WithJSON(payload).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("name").ValueEqual("name", "Test Service Patch")

	// invalid status
	payload = map[string]interface{}{
		"status": "not valid",
	}
	auth.PATCH("/admin/bank_services/1").WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()

	// test invalid token
	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer wrong token")
	})
	auth.PATCH("/admin/bank_services/1").WithJSON(payload).
		Expect().
		Status(http.StatusUnauthorized).JSON().Object()
}
