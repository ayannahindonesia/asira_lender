package tests

import (
	"asira_lender/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestGetProductList(t *testing.T) {
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
	auth.GET("/admin/products").
		Expect().
		Status(http.StatusOK).JSON().Object()

	// test query found
	obj := auth.GET("/admin/products").WithQuery("name", "Product A").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 1)
	//with part of name
	obj = auth.GET("/admin/products").WithQuery("name", "prod").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 5)

	// test query invalid
	obj = auth.GET("/admin/products").WithQuery("name", "should not found this").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 0)
}

func TestNewProduct(t *testing.T) {
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
		"name":         "Test Product",
		"service_id":   1,
		"min_timespan": 3,
		"max_timespan": 12,
		"interest":     5,
		"min_loan":     1000000,
		"max_loan":     10000000,
		"fees": []interface{}{
			map[string]interface{}{
				"description": "Admin Fee",
				"amount":      "1%",
			},
			map[string]interface{}{
				"description": "Convenience Fee",
				"amount":      "2%",
			},
		},
		"collaterals":      []string{"Surat Tanah", "BPKB"},
		"financing_sector": []string{"Pendidikan"},
		"assurance":        "assuransi",
		"status":           "active",
	}

	// normal scenario
	obj := auth.POST("/admin/products").WithJSON(payload).
		Expect().
		Status(http.StatusCreated).JSON().Object()
	obj.ContainsKey("name").ValueEqual("name", "Test Product")

	// test invalid
	payload = map[string]interface{}{
		"name": "",
	}
	auth.POST("/admin/products").WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()
}

func TestGetProductbyID(t *testing.T) {
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
	obj := auth.GET("/admin/products/1").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("id").ValueEqual("id", 1)

	// not found
	auth.GET("/admin/products/9999").
		Expect().
		Status(http.StatusNotFound).JSON().Object()
}

func TestPatchProduct(t *testing.T) {
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
		"name": "Test Service Product Patch",
	}

	// valid response
	obj := auth.PATCH("/admin/products/1").WithJSON(payload).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("name").ValueEqual("name", "Test Service Product Patch")

	// valid response
	payload = map[string]interface{}{
		"status": "invalid",
	}
	auth.PATCH("/admin/products/1").WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()

	// test invalid token
	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer wrong token")
	})
	auth.PATCH("/admin/products/1").WithJSON(payload).
		Expect().
		Status(http.StatusUnauthorized).JSON().Object()
}
