package groups

import (
	"asira_lender/adminhandlers"
	"asira_lender/handlers"
	"asira_lender/middlewares"

	"github.com/labstack/echo"
)

// ClientGroup group
func ClientGroup(e *echo.Echo) {
	g := e.Group("/client")
	middlewares.SetClientJWTmiddlewares(g, "client")
	g.POST("/lender_login", handlers.LenderLogin)
	g.POST("/admin_login", adminhandlers.AdminLogin)
	g.POST("/forgotpassword", handlers.UserResetPasswordRequest)
	g.POST("/resetpassword", handlers.UserResetPasswordVerify)

	g.GET("/serviceinfo", handlers.ServiceInfo)
}
