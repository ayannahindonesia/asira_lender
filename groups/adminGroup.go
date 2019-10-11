package groups

import (
	"asira_lender/admin_handlers"
	"asira_lender/handlers"
	"asira_lender/middlewares"
	"asira_lender/reports"

	"github.com/labstack/echo"
)

func AdminGroup(e *echo.Echo) {
	g := e.Group("/admin")
	middlewares.SetClientJWTmiddlewares(g, "admin")

	// config info
	g.GET("/info", handlers.AsiraAppInfo)

	// Internals Accounts Management
	g.POST("/client_config", admin_handlers.CreateInternal)

	// Images
	g.GET("/image/:image_id", admin_handlers.GetImageB64String)

	// Bank Types
	g.GET("/bank_types", admin_handlers.BankTypeList)
	g.POST("/bank_types", admin_handlers.BankTypeNew)
	g.GET("/bank_types/:bank_id", admin_handlers.BankTypeDetail)
	g.PATCH("/bank_types/:bank_id", admin_handlers.BankTypePatch)
	g.DELETE("/bank_types/:bank_id", admin_handlers.BankTypeDelete)

	// Banks
	g.GET("/banks", admin_handlers.BankList)
	g.POST("/banks", admin_handlers.BankNew)
	g.GET("/banks/:bank_id", admin_handlers.BankDetail)
	g.PATCH("/banks/:bank_id", admin_handlers.BankPatch)
	g.DELETE("/banks/:bank_id", admin_handlers.BankDelete)

	// Services
	g.GET("/services", admin_handlers.ServiceList)
	g.POST("/services", admin_handlers.ServiceNew)
	g.GET("/services/:id", admin_handlers.ServiceDetail)
	g.PATCH("/services/:id", admin_handlers.ServicePatch)
	g.DELETE("/services/:id", admin_handlers.ServiceDelete)

	// Products
	g.GET("/products", admin_handlers.ProductList)
	g.POST("/products", admin_handlers.ProductNew)
	g.GET("/products/:id", admin_handlers.ProductDetail)
	g.PATCH("/products/:id", admin_handlers.ProductPatch)
	g.DELETE("/products/:id", admin_handlers.ProductDelete)

	// Bank Services
	g.GET("/bank_services", admin_handlers.BankServiceList)
	g.GET("/bank_services/:id", admin_handlers.BankServiceDetail)
	g.POST("/bank_services", admin_handlers.BankServiceNew)
	g.PATCH("/bank_services/:id", admin_handlers.BankServicePatch)
	g.DELETE("/bank_services/:id", admin_handlers.BankServiceDelete)

	// Bank Products
	g.GET("/bank_products", admin_handlers.BankProductList)
	g.POST("/bank_products", admin_handlers.BankProductNew)
	g.GET("/bank_products/:id", admin_handlers.BankProductDetail)
	g.PATCH("/bank_products/:id", admin_handlers.BankProductPatch)
	g.DELETE("/bank_products/:id", admin_handlers.BankProductDelete)

	// Loan Purpose
	g.GET("/loan_purposes", admin_handlers.LoanPurposeList)
	g.POST("/loan_purposes", admin_handlers.LoanPurposeNew)
	g.GET("/loan_purposes/:loan_purpose_id", admin_handlers.LoanPurposeDetail)
	g.PATCH("/loan_purposes/:loan_purpose_id", admin_handlers.LoanPurposePatch)
	g.DELETE("/loan_purposes/:loan_purpose_id", admin_handlers.LoanPurposeDelete)

	// Role
	g.GET("/roles", admin_handlers.GetAllRole)
	g.GET("/roles/:role_id", admin_handlers.RoleGetDetails)
	g.POST("/roles", admin_handlers.AddRole)
	g.PATCH("/roles/:role_id", admin_handlers.UpdateRole)

	//Permission
	g.GET("/permission", admin_handlers.GetAllPermission)
	g.GET("/permission/:permission_id", admin_handlers.PermissionGetDetails)
	g.POST("/permission", admin_handlers.AddPermission)
	g.PATCH("/permission", admin_handlers.UpdatePermission)

	// Reports
	g.GET("/reports/convenience_fee", reports.ConvenienceFeeReport)
	
	// User
	g.GET("/users", admin_handlers.GetAllUser)
	g.GET("/users/:user_id", admin_handlers.UserGetDetails)
	g.POST("/users", admin_handlers.AddUser)
	g.PATCH("/users/:user_id", admin_handlers.UpdateUser)
}
