package router

import (
	"asira_lender/groups"
	"asira_lender/handlers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func NewRouter() *echo.Echo {
	e := echo.New()

	// ignore /api-lender
	e.Pre(middleware.Rewrite(map[string]string{
		"/api-lender/*": "/$1",
	}))

	// e.GET("/test", handlers.Test)
	e.GET("/clientauth", handlers.ClientLogin)
	e.GET("/ping", handlers.Ping)

	groups.AdminGroup(e)
	groups.ClientGroup(e)
	groups.LenderGroup(e)

	return e
}
