package models

import (
	"time"

	"github.com/ayannahindonesia/basemodel"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

type (
	// Loan main struct
	Loan struct {
		basemodel.BaseModel
		Borrower            uint64         `json:"borrower" gorm:"column:borrower;foreignkey"`
		Status              string         `json:"status" gorm:"column:status;type:varchar(255)" sql:"DEFAULT:'processing'"`
		LoanAmount          float64        `json:"loan_amount" gorm:"column:loan_amount;type:int;not null"`
		Installment         int            `json:"installment" gorm:"column:installment;type:int;not null"` // plan of how long loan to be paid
		InstallmentID       pq.Int64Array  `json:"installment_id" gorm:"column:installment_id"`
		Fees                postgres.Jsonb `json:"fees" gorm:"column:fees;type:jsonb"`
		Interest            float64        `json:"interest" gorm:"column:interest;type:int;not null"`
		TotalLoan           float64        `json:"total_loan" gorm:"column:total_loan;type:int;not null"`
		DisburseAmount      float64        `json:"disburse_amount" gorm:"column:disburse_amount;type:int;not null"`
		DueDate             time.Time      `json:"due_date" gorm:"column:due_date"`
		LayawayPlan         float64        `json:"layaway_plan" gorm:"column:layaway_plan;type:int;not null"` // how much borrower will pay per month
		Product             uint64         `json:"product" gorm:"column:product;foreignkey"`                  // product and service is later to be discussed
		LoanIntention       string         `json:"loan_intention" gorm:"column:loan_intention;type:varchar(255);not null"`
		IntentionDetails    string         `json:"intention_details" gorm:"column:intention_details;type:text;not null"`
		BorrowerInfo        postgres.Jsonb `json:"borrower_info" gorm:"column:borrower_info;type:jsonb"`
		OTPverified         bool           `json:"otp_verified" gorm:"column:otp_verified;type:boolean" sql:"DEFAULT:FALSE"`
		DisburseDate        time.Time      `json:"disburse_date" gorm:"column:disburse_date"`
		DisburseDateChanged bool           `json:"disburse_date_changed" gorm:"column:disburse_date_changed"`
		DisburseStatus      string         `json:"disburse_status" gorm:"column:disburse_status" sql:"DEFAULT:'processing'"`
		ApprovalDate        time.Time      `json:"approval_date" gorm:"column:approval_date"`
		RejectReason        string         `json:"reject_reason" gorm:"column:reject_reason"`
		FormInfo            postgres.Jsonb `json:"form_info" gorm:"column:form_info;type:jsonb"`
	}

	// LoanFee for loan fee
	LoanFee struct { // temporary hardcoded
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
	}
	// LoanFees slice of LoanFee
	LoanFees []LoanFee

	// LoanStatusUpdate type
	LoanStatusUpdate struct {
		ID     uint64 `json:"id"`
		Status string `json:"status"`
	}
)

// Create func
func (model *Loan) Create() error {
	return basemodel.Create(&model)
}

// Save func
func (model *Loan) Save() error {
	return basemodel.Save(&model)
}

// FirstOrCreate func
func (model *Loan) FirstOrCreate() error {
	return basemodel.FirstOrCreate(&model)
}

// Delete func
func (model *Loan) Delete() error {
	return basemodel.Delete(&model)
}

// FindbyID func
func (model *Loan) FindbyID(id uint64) error {
	return basemodel.FindbyID(&model, id)
}

// SingleFindFilter func
func (model *Loan) SingleFindFilter(filter interface{}) error {
	return basemodel.SingleFindFilter(&model, filter)
}

// PagedFilterSearch func
func (model *Loan) PagedFilterSearch(page int, rows int, orderby []string, sort []string, filter interface{}) (basemodel.PagedFindResult, error) {
	loans := []Loan{}

	return basemodel.PagedFindFilter(&loans, page, rows, orderby, sort, filter)
}
