package groups

import (
	"asira_lender/admin_handlers"
	"asira_lender/handlers"
	"asira_lender/middlewares"
	"asira_lender/reports"

	"github.com/labstack/echo"
)

// AdminGroup func
func AdminGroup(e *echo.Echo) {
	g := e.Group("/admin")
	middlewares.SetClientJWTmiddlewares(g, "users")

	// config info
	g.GET("/info", handlers.AsiraAppInfo)
	g.GET("/profile", admin_handlers.AdminProfile)

	// Client Management
	g.POST("/client", admin_handlers.CreateClient)

	// Images
	g.GET("/image/:image_id", admin_handlers.GetImageB64String)

	// Borrowers
	g.GET("/borrower", admin_handlers.BorrowerGetAll)
	g.GET("/borrower/:borrower_id", admin_handlers.BorrowerGetDetails)

	// Loans
	g.GET("/loan", admin_handlers.LoanGetAll)
	g.GET("/loan/:loan_id", admin_handlers.LoanGetDetails)

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

	// Loan Purpose
	g.GET("/loan_purposes", admin_handlers.LoanPurposeList)
	g.POST("/loan_purposes", admin_handlers.LoanPurposeNew)
	g.GET("/loan_purposes/:loan_purpose_id", admin_handlers.LoanPurposeDetail)
	g.PATCH("/loan_purposes/:loan_purpose_id", admin_handlers.LoanPurposePatch)
	g.DELETE("/loan_purposes/:loan_purpose_id", admin_handlers.LoanPurposeDelete)

	// Role
	g.GET("/roles", admin_handlers.RoleList)
	g.GET("/roles/:role_id", admin_handlers.RoleDetails)
	g.POST("/roles", admin_handlers.RoleNew)
	g.PATCH("/roles/:role_id", admin_handlers.RolePatch)
	g.GET("/roles_all", admin_handlers.RoleRange)

	// Permission
	g.GET("/permission", admin_handlers.PermissionList)

	// User
	g.GET("/users", admin_handlers.UserList)
	g.GET("/users/:user_id", admin_handlers.UserDetails)
	g.POST("/users", admin_handlers.UserNew)
	g.PATCH("/users/:user_id", admin_handlers.UserPatch)

	// Agent Provider
	g.GET("/agent_providers", admin_handlers.AgentProviderList)
	g.GET("/agent_providers/:id", admin_handlers.AgentProviderDetails)
	g.POST("/agent_providers", admin_handlers.AgentProviderNew)
	g.PATCH("/agent_providers/:id", admin_handlers.AgentProviderPatch)

	// Agent Provider
	g.GET("/agents", admin_handlers.AgentList)
	g.GET("/agents/:id", admin_handlers.AgentDetails)
	g.POST("/agents", admin_handlers.AgentNew)
	g.PATCH("/agents/:id", admin_handlers.AgentPatch)
	g.DELETE("/agents/:id", admin_handlers.AgentDelete)

	// Reports
	g.GET("/reports/convenience_fee", reports.ConvenienceFeeReport)
}
