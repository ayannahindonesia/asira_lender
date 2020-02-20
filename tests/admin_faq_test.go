package tests

import (
	"asira_lender/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestFAQList(t *testing.T) {
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
	obj := auth.GET("/admin/faq").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").NotEmpty()

	obj = auth.GET("/admin/faq").WithQuery("title", "How to").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 2)

	// test query not found
	obj = auth.GET("/admin/faq").WithQuery("title", "should not be found").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 0)
}

func TestNewFAQ(t *testing.T) {
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
		"title": "Test new FAQ",
		"description": `<html>
		<head>
		</head>
		<body>
		<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
		<strong class="accordion">Section 1</strong>
		<p>Lorem ipsum...1</p>
		</div>
		<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
		<strong class="accordion">Section 2</strong>
		<p>Lorem ipsum...2</p>
		</div>
		<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
		<strong class="accordion">Section 3</strong>
		<p>Lorem ipsum...3</p>
		</div>
		</body>
		</html>`,
	}

	// normal scenario
	obj := auth.POST("/admin/faq").WithJSON(payload).
		Expect().
		Status(http.StatusCreated).JSON().Object()
	obj.ContainsKey("title").ValueEqual("title", "Test new FAQ")

	// test invalid
	payload = map[string]interface{}{
		"title": "",
	}
	auth.POST("/admin/faq").WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()
}

func TestGetFAQByID(t *testing.T) {
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
	obj := auth.GET("/admin/loan_purposes/1").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("id").ValueEqual("id", 1)

	// not found
	auth.GET("/admin/loan_purposes/9999").
		Expect().
		Status(http.StatusNotFound).JSON().Object()
}

// func TestFAQPatch(t *testing.T) {
// 	RebuildData()

// 	api := router.NewRouter()

// 	server := httptest.NewServer(api)

// 	defer server.Close()

// 	e := httpexpect.New(t, server.URL)

// 	auth := e.Builder(func(req *httpexpect.Request) {
// 		req.WithHeader("Authorization", "Basic "+clientBasicToken)
// 	})

// 	adminToken := getAdminLoginToken(e, auth, "1")

// 	auth = e.Builder(func(req *httpexpect.Request) {
// 		req.WithHeader("Authorization", "Bearer "+adminToken)
// 	})

// 	payload := map[string]interface{}{
// 		"name": "Test Patch",
// 	}

// 	// valid response
// 	obj := auth.PATCH("/admin/loan_purposes/1").WithJSON(payload).
// 		Expect().
// 		Status(http.StatusOK).JSON().Object()
// 	obj.ContainsKey("name").ValueEqual("name", "Test Patch")

// 	// test invalid token
// 	auth = e.Builder(func(req *httpexpect.Request) {
// 		req.WithHeader("Authorization", "Bearer wrong token")
// 	})
// 	auth.PATCH("/admin/loan_purposes/1").WithJSON(payload).
// 		Expect().
// 		Status(http.StatusUnauthorized).JSON().Object()
// }

// func TestDeleteFAQ(t *testing.T) {
// 	RebuildData()

// 	api := router.NewRouter()

// 	server := httptest.NewServer(api)

// 	defer server.Close()

// 	e := httpexpect.New(t, server.URL)

// 	auth := e.Builder(func(req *httpexpect.Request) {
// 		req.WithHeader("Authorization", "Basic "+clientBasicToken)
// 	})

// 	adminToken := getAdminLoginToken(e, auth, "1")

// 	auth = e.Builder(func(req *httpexpect.Request) {
// 		req.WithHeader("Authorization", "Bearer "+adminToken)
// 	})

// 	// valid response
// 	auth.DELETE("/admin/loan_purposes/1").
// 		Expect().
// 		Status(http.StatusOK).JSON().Object()

// 	auth.GET("/admin/loan_purposes/1").
// 		Expect().
// 		Status(http.StatusNotFound).JSON().Object()
// }
