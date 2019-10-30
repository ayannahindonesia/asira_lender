package groups

import (
	"asira_lender/admin_handlers"
	"asira_lender/handlers"
	"asira_lender/middlewares"

	"github.com/labstack/echo"
)

func ClientGroup(e *echo.Echo) {
	g := e.Group("/client")
	middlewares.SetClientJWTmiddlewares(g, "client")
	g.POST("/lender_login", handlers.LenderLogin)
	g.POST("/admin_login", admin_handlers.AdminLogin)

	// loan purposes
	g.GET("/loan_purposes", handlers.LoanPurposeList)
	g.GET("/loan_purposes/:loan_purpose_id", handlers.LoanPurposeDetail)
}
