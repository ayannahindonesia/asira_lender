package models

import (
	"database/sql"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	"gitlab.com/asira-ayannah/basemodel"
)

type (
	Loan struct {
		basemodel.BaseModel
		DeletedTime      time.Time      `json:"deleted_time" gorm:"column:deleted_time"`
		Owner            sql.NullInt64  `json:"owner" gorm:"column:owner;foreignkey"`
		OwnerName        string         `json:"owner_name" gorm:"column:owner_name"`
		Bank             sql.NullInt64  `json:"bank" gorm:"column:bank;foreignkey"`
		Status           string         `json:"status" gorm:"column:status;type:varchar(255)" sql:"DEFAULT:'processing'"`
		LoanAmount       float64        `json:"loan_amount" gorm:"column:loan_amount;type:int;not null"`
		Installment      int            `json:"installment" gorm:"column:installment;type:int;not null"` // plan of how long loan to be paid
		Fees             postgres.Jsonb `json:"fees" gorm:"column:fees;type:jsonb"`
		Interest         float64        `json:"interest" gorm:"column:interest;type:int;not null"`
		TotalLoan        float64        `json:"total_loan" gorm:"column:total_loan;type:int;not null"`
		DueDate          time.Time      `json:"due_date" gorm:"column:due_date"`
		LayawayPlan      float64        `json:"layaway_plan" gorm:"column:layaway_plan;type:int;not null"` // how much borrower will pay per month
		Product          uint64         `json:"product" gorm:"column:product;foreignkey"`                  // product and service is later to be discussed
		LoanIntention    string         `json:"loan_intention" gorm:"column:loan_intention;type:varchar(255);not null"`
		IntentionDetails string         `json:"intention_details" gorm:"column:intention_details;type:text;not null"`
		DisburseDate     time.Time      `json:"disburse_date" gorm:"column:disburse_date"`
		RejectReason     string         `json:"reject_reason" gorm:"column:reject_reason"`
	}

	LoanFee struct { // temporary hardcoded
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
	}
	LoanFees []LoanFee

	LoanStatusUpdate struct {
		ID     uint64 `json:"id"`
		Status string `json:"status"`
	}
)

func (l *Loan) Create() error {
	err := basemodel.Create(&l)
	return err
}

func (l *Loan) Save() error {
	err := basemodel.Save(&l)
	return err
}

func (l *Loan) Delete() error {
	l.DeletedTime = time.Now()
	err := basemodel.Save(&l)

	return err
}

func (l *Loan) FindbyID(id int) error {
	err := basemodel.FindbyID(&l, id)
	return err
}

func (l *Loan) FilterSearchSingle(filter interface{}) error {
	err := basemodel.SingleFindFilter(&l, filter)
	return err
}

func (l *Loan) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	loans := []Loan{}

	order := []string{orderby}
	sorts := []string{sort}
	result, err = basemodel.PagedFindFilter(&loans, page, rows, order, sorts, filter)

	return result, err
}

func (l *Loan) Approve(disburseDate time.Time) error {
	l.Status = "approved"
	l.DisburseDate = disburseDate

	err := l.Save()
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(l, "loan")

	return err
}

func (l *Loan) Reject(reason string) error {
	l.Status = "rejected"
	l.RejectReason = reason

	err := l.Save()
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(l, "loan")

	return err
}
