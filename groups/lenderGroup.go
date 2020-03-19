package groups

import (
	"asira_lender/handlers"
	"asira_lender/middlewares"

	"github.com/labstack/echo"
)

// LenderGroup group
func LenderGroup(e *echo.Echo) {
	g := e.Group("/lender")
	middlewares.SetClientJWTmiddlewares(g, "users")

	// Profile endpoints
	g.GET("/profile", handlers.LenderProfile)
	g.PATCH("/profile", handlers.LenderProfileEdit)
	g.POST("/first_login", handlers.UserFirstLoginChangePassword)

	// Loans endpoints
	g.GET("/loanrequest_list", handlers.LenderLoanRequestList)
	g.GET("/loanrequest_list/:loan_id/detail", handlers.LenderLoanRequestListDetail)
	g.GET("/loanrequest_list/:loan_id/detail/:approve_reject", handlers.LenderLoanApproveReject)
	g.GET("/loanrequest_list/:loan_id/detail/confirm_disbursement", handlers.LenderLoanConfirmDisbursement)
	g.GET("/loanrequest_list/:loan_id/detail/change_disburse_date", handlers.LenderLoanChangeDisburseDate)
	g.GET("/loanrequest_list/download", handlers.LenderLoanRequestListDownload)
	g.GET("/loanrequest_list/:loan_id/detail/installments", handlers.LenderLoanInstallmentList)
	g.PATCH("/loanrequest_list/:loan_id/detail/installment_approve/:installment_id", handlers.LenderLoanInstallmentsApprove)
	g.PATCH("/loanrequest_list/:loan_id/detail/installment_approve/bulk", handlers.LenderLoanInstallmentsApproveBulk)

	// Borrowers endpoints
	g.GET("/borrower_list", handlers.LenderBorrowerList)
	g.GET("/borrower_list/:borrower_id/detail", handlers.LenderBorrowerListDetail)
	g.GET("/borrower_list/download", handlers.LenderBorrowerListDownload)
	g.GET("/borrower_list/:borrower_id/:approval", handlers.LenderApproveRejectProspectiveBorrower)

	// services owned by bank (lender)
	g.GET("/services", handlers.LenderServiceList)
	g.GET("/services/:service_id", handlers.LenderServiceLListDetail)

	g.GET("/products", handlers.LenderProductList)
	g.GET("/products/:product_id", handlers.LenderProductDetail)
}
