package tests

import (
	"asira_lender/asira"
	"asira_lender/migration"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/gavv/httpexpect"
)

var (
	clientBasicToken string = base64.StdEncoding.EncodeToString([]byte("reactkey:reactsecret"))
	adminBasicToken  string = base64.StdEncoding.EncodeToString([]byte("adminkey:adminsecret"))
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
	migration.TestSeed()
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

func getAdminLoginToken(e *httpexpect.Expect, auth *httpexpect.Expect, admin_id string) string {
	obj := auth.GET("/clientauth").
		Expect().
		Status(http.StatusOK).JSON().Object()

	admintoken := obj.Value("token").String().Raw()

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+admintoken)
	})

	var payload map[string]interface{}
	switch admin_id {
	case "1":
		payload = map[string]interface{}{
			"key":      "adminkey",
			"password": "adminsecret",
		}
	}

	obj = auth.POST("/client/admin_login").WithJSON(payload).
		Expect().
		Status(http.StatusOK).JSON().Object()

	return obj.Value("token").String().Raw()
}

func getLenderAdminToken(e *httpexpect.Expect, auth *httpexpect.Expect) string {
	obj := auth.GET("/clientauth").
		Expect().
		Status(http.StatusOK).JSON().Object()

	admintoken := obj.Value("token").String().Raw()

	return admintoken
}
