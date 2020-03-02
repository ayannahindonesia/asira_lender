package groups

import (
	"asira_lender/adminhandlers"
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
	g.GET("/profile", adminhandlers.AdminProfile)

	// Client Management
	g.POST("/client", adminhandlers.CreateClient)

	// Borrowers
	g.GET("/borrower", adminhandlers.BorrowerGetAll)
	g.GET("/borrower/:borrower_id", adminhandlers.BorrowerGetDetails)

	// Loans
	g.GET("/loan", adminhandlers.LoanGetAll)
	g.GET("/loan/:loan_id", adminhandlers.LoanGetDetails)

	// Bank Types
	g.GET("/bank_types", adminhandlers.BankTypeList)
	g.POST("/bank_types", adminhandlers.BankTypeNew)
	g.GET("/bank_types/:bank_id", adminhandlers.BankTypeDetail)
	g.PATCH("/bank_types/:bank_id", adminhandlers.BankTypePatch)
	g.DELETE("/bank_types/:bank_id", adminhandlers.BankTypeDelete)

	// Banks
	g.GET("/banks", adminhandlers.BankList)
	g.POST("/banks", adminhandlers.BankNew)
	g.GET("/banks/:bank_id", adminhandlers.BankDetail)
	g.PATCH("/banks/:bank_id", adminhandlers.BankPatch)
	g.DELETE("/banks/:bank_id", adminhandlers.BankDelete)

	// Services
	g.GET("/services", adminhandlers.ServiceList)
	g.POST("/services", adminhandlers.ServiceNew)
	g.GET("/services/:id", adminhandlers.ServiceDetail)
	g.PATCH("/services/:id", adminhandlers.ServicePatch)
	g.DELETE("/services/:id", adminhandlers.ServiceDelete)

	// Products
	g.GET("/products", adminhandlers.ProductList)
	g.POST("/products", adminhandlers.ProductNew)
	g.GET("/products/:id", adminhandlers.ProductDetail)
	g.PATCH("/products/:id", adminhandlers.ProductPatch)
	g.DELETE("/products/:id", adminhandlers.ProductDelete)

	// Loan Purpose
	g.GET("/loan_purposes", adminhandlers.LoanPurposeList)
	g.POST("/loan_purposes", adminhandlers.LoanPurposeNew)
	g.GET("/loan_purposes/:loan_purpose_id", adminhandlers.LoanPurposeDetail)
	g.PATCH("/loan_purposes/:loan_purpose_id", adminhandlers.LoanPurposePatch)
	g.DELETE("/loan_purposes/:loan_purpose_id", adminhandlers.LoanPurposeDelete)

	// Role
	g.GET("/roles", adminhandlers.RoleList)
	g.GET("/roles/:id", adminhandlers.RoleDetails)
	g.POST("/roles", adminhandlers.RoleNew)
	g.PATCH("/roles/:id", adminhandlers.RolePatch)
	g.GET("/roles_all", adminhandlers.RoleRange)

	// Permission
	g.GET("/permission", adminhandlers.PermissionList)

	// User
	g.GET("/users", adminhandlers.UserList)
	g.GET("/users/:id", adminhandlers.UserDetails)
	g.POST("/users", adminhandlers.UserNew)
	g.PATCH("/users/:id", adminhandlers.UserPatch)

	// Agent Provider
	g.GET("/agent_providers", adminhandlers.AgentProviderList)
	g.GET("/agent_providers/:id", adminhandlers.AgentProviderDetails)
	g.POST("/agent_providers", adminhandlers.AgentProviderNew)
	g.PATCH("/agent_providers/:id", adminhandlers.AgentProviderPatch)

	// Agent Provider
	g.GET("/agents", adminhandlers.AgentList)
	g.GET("/agents/:id", adminhandlers.AgentDetails)
	g.POST("/agents", adminhandlers.AgentNew)
	g.PATCH("/agents/:id", adminhandlers.AgentPatch)
	g.DELETE("/agents/:id", adminhandlers.AgentDelete)

	// Reports
	g.GET("/reports/convenience_fee", reports.ConvenienceFeeReport)

	// FAQ
	g.GET("/faq", adminhandlers.FAQList)
	g.POST("/faq", adminhandlers.FAQNew)
	g.GET("/faq/:faq_id", adminhandlers.FAQDetail)
	g.PATCH("/faq/:faq_id", adminhandlers.FAQPatch)
	g.DELETE("/faq/:faq_id", adminhandlers.FAQDelete)
}
