package tests

import (
	"asira_lender/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestResetPassword(t *testing.T) {
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

	payload := map[string]interface{}{
		"email":  "testuser@ayannah.id",
		"system": "dashboard",
	}
	auth.POST("/client/forgotpassword").WithJSON(payload).
		Expect().
		Status(http.StatusNotFound).JSON().Object()

	payload = map[string]interface{}{
		"email":  "toib@ayannah.com",
		"system": "dashboard",
	}
	auth.POST("/client/forgotpassword").WithJSON(payload).
		Expect().
		Status(http.StatusOK)
}
