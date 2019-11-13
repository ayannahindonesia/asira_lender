package tests

import (
	"asira_lender/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestLenderGetBankTypeList(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	adminToken := getAdminLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+adminToken)
	})

	// valid response
	auth.GET("/admin/bank_types").
		Expect().
		Status(http.StatusOK).JSON().Object()

	// test query found
	obj := auth.GET("/admin/bank_types").WithQuery("name", "Koperasi").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 1)
	// test query found with part
	obj = auth.GET("/admin/bank_types").WithQuery("name", "operasi").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 1)

	// test query invalid
	obj = auth.GET("/admin/bank_types").WithQuery("name", "should not found this").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 0)
}

func TestNewBankType(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	adminToken := getAdminLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+adminToken)
	})

	payload := map[string]interface{}{
		"name": "Test New Type",
	}

	// normal scenario
	obj := auth.POST("/admin/bank_types").WithJSON(payload).
		Expect().
		Status(http.StatusCreated).JSON().Object()
	obj.ContainsKey("name").ValueEqual("name", "Test New Type")

	// test invalid
	payload = map[string]interface{}{
		"name": "",
	}
	auth.POST("/admin/bank_types").WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()
}

func TestGetBankTypebyID(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	adminToken := getAdminLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+adminToken)
	})

	// valid response
	obj := auth.GET("/admin/bank_types/1").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("id").ValueEqual("id", 1)

	// not found
	auth.GET("/admin/bank_types/9999").
		Expect().
		Status(http.StatusNotFound).JSON().Object()
}

func TestPatchBankType(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	adminToken := getAdminLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+adminToken)
	})

	payload := map[string]interface{}{
		"name": "Test Patch",
	}

	// valid response
	obj := auth.PATCH("/admin/bank_types/1").WithJSON(payload).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("name").ValueEqual("name", "Test Patch")

	// test invalid token
	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer wrong token")
	})
	auth.PATCH("/admin/bank_types/1").WithJSON(payload).
		Expect().
		Status(http.StatusUnauthorized).JSON().Object()
}

// @ToDo must delete foreign reference first
// func TestDeleteBankType(t *testing.T) {
// 	RebuildData()

// 	api := router.NewRouter()

// 	server := httptest.NewServer(api)

// 	defer server.Close()

// 	e := httpexpect.New(t, server.URL)

// 	auth := e.Builder(func(req *httpexpect.Request) {
// 		req.WithHeader("Authorization", "Basic "+adminBasicToken)
// 	})

// 	adminToken := getLenderAdminToken(e, auth)

// 	auth = e.Builder(func(req *httpexpect.Request) {
// 		req.WithHeader("Authorization", "Bearer "+adminToken)
// 	})

// 	// valid response
// 	obj := auth.DELETE("/admin/bank_types/1").
// 		Expect().
// 		Status(http.StatusOK).JSON().Object()
// 	obj.ContainsKey("name").ValueEqual("name", "Test Patch")
// 	auth.GET("/admin/bank_types/1").
// 		Expect().
// 		Status(http.StatusNotFound).JSON().Object()
// }
