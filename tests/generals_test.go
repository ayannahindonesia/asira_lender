package tests

import (
	"asira_lender/asira"
	"asira_lender/migration"
	"fmt"
	"net/http"
	"os"

	"github.com/gavv/httpexpect"
)

var (
	clientBasicToken string = "cmVhY3RrZXk6cmVhY3RwYXNz"
)

func init() {
	// restrict test to development environment only.
	if asira.App.ENV != "development" {
		fmt.Printf("test aren't allowed in %s environment.", asira.App.ENV)
		os.Exit(1)
	}
}

func RebuildData() {
	migration.Truncate([]string{"all"})
	migration.Seed()
}

func getLenderLoginToken(e *httpexpect.Expect, auth *httpexpect.Expect, lender_id string) string {
	obj := auth.GET("/clientauth").
		Expect().
		Status(http.StatusOK).JSON().Object()

	admintoken := obj.Value("token").String().Raw()

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+admintoken)
	})

	var payload map[string]interface{}
	switch lender_id {
	case "1":
		payload = map[string]interface{}{
			"key":      "Banktoib",
			"password": "password",
		}
	}

	obj = auth.POST("/client/lender_login").WithJSON(payload).
		Expect().
		Status(http.StatusOK).JSON().Object()

	return obj.Value("token").String().Raw()
}
