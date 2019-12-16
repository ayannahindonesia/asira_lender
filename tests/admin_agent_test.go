package tests

import (
	"asira_lender/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestAgentListandDetails(t *testing.T) {
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
	auth.GET("/admin/agents").
		Expect().
		Status(http.StatusOK).JSON().Object()

	// test query found
	obj := auth.GET("/admin/agents").WithQuery("name", "Agent K").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 1)

	// test query invalid
	obj = auth.GET("/admin/agents").WithQuery("name", "should not found this").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 0)

	// get by id
	obj = auth.GET("/admin/agents/1").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("name").ValueEqual("name", "Agent K")

	// test query found
	auth.GET("/admin/agents/9999").
		Expect().
		Status(http.StatusNotFound).JSON().Object()
}

func TestAgentNew(t *testing.T) {
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
		"name":           "Test Agent",
		"username":       "testagent",
		"email":          "agent@test.com",
		"phone":          "0812345567890",
		"category":       "agent",
		"agent_provider": 1,
		"image":          "base64 very long string",
		"banks":          []int{1},
		"status":         "active",
	}

	// normal scenario
	obj := auth.POST("/admin/agents").WithJSON(payload).
		Expect().
		Status(http.StatusCreated).JSON().Object()
	obj.ContainsKey("name").ValueEqual("name", "Test Agent")

	// test invalid
	auth.POST("/admin/agents").WithJSON(map[string]interface{}{}).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()
}

func TestAgentPatch(t *testing.T) {
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
		"agent_provider": 2,
		"status":         "inactive",
	}

	// normal scenario
	obj := auth.PATCH("/admin/agents/3").WithJSON(payload).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("status").ValueEqual("status", "inactive")
	obj.Value("agent_provider").Object().ValueEqual("Int64", 2)

	// uniques
	auth.PATCH("/admin/agents/3").WithJSON(map[string]interface{}{
		"phone": "081234567890",
	}).
		Expect().
		Status(http.StatusInternalServerError).JSON().Object()
	auth.PATCH("/admin/agents/3").WithJSON(map[string]interface{}{
		"email": "agentk@mib.com",
	}).
		Expect().
		Status(http.StatusInternalServerError).JSON().Object()

	// test invalid
	auth.POST("/admin/agents").WithJSON(map[string]interface{}{
		"banks": []int{99},
	}).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()
}
