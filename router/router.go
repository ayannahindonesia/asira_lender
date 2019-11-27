package router

import (
	"asira_lender/groups"
	"asira_lender/handlers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// NewRouter func
func NewRouter() *echo.Echo {
	e := echo.New()

	// ignore /api-lender
	e.Pre(middleware.Rewrite(map[string]string{
		"/api-lender/*":       "/$1",
		"/api-lender-devel/*": "/$1",
	}))

	// e.GET("/test", handlers.Test)
	e.GET("/clientauth", handlers.ClientLogin)
	e.GET("/ping", handlers.Ping)

	e.POST("/test", handlers.S3test)

	groups.AdminGroup(e)
	groups.ClientGroup(e)
	groups.LenderGroup(e)

	return e
}
