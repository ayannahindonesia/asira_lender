package models

import (
	"database/sql"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type (
	Loan struct {
		BaseModel
		DeletedTime time.Time      `json:"deleted_time" gorm:"column:deleted_time"`
		Owner       sql.NullInt64  `json:"owner" gorm:"column:owner;foreignkey"`
		OwnerName   string         `json:"owner_name" gorm:"column:owner_name"`
		Bank        sql.NullInt64  `json:"bank" gorm:"column:bank;foreignkey"`
		Status      string         `json:"status" gorm:"column:status;type:varchar(255)" sql:"DEFAULT:'processing'"`
		LoanAmount  float64        `json:"loan_amount" gorm:"column:loan_amount;type:int;not null"`
		Installment int            `json:"installment" gorm:"column:installment;type:int;not null"` // plan of how long loan to be paid
		Fees        postgres.Jsonb `json:"fees" gorm:"column:fees;type:jsonb"`
		Interest    float64        `json:"interest" gorm:"column:interest;type:int;not null"`
		TotalLoan   float64        `json:"total_loan" gorm:"column:total_loan;type:int;not null"`
		DueDate     time.Time      `json:"due_date" gorm:"column:due_date"`
		LayawayPlan float64        `json:"layaway_plan" gorm:"column:layaway_plan;type:int;not null"` // how much borrower will pay per month
		// Product          uint64         `json:"product" gorm:"column:product;foreignkey"` // product and service is later to be discussed
		// Service          uint64         `json:"service" gorm:"column:service;foreignkey"`
		LoanIntention    string `json:"loan_intention" gorm:"column:loan_intention;type:varchar(255);not null"`
		IntentionDetails string `json:"intention_details" gorm:"column:intention_details;type:text;not null"`
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

func (l *Loan) Create() (*Loan, error) {
	err := Create(&l)
	return l, err
}

func (l *Loan) Save() (*Loan, error) {
	err := Save(&l)
	return l, err
}

func (l *Loan) Delete() (*Loan, error) {
	l.DeletedTime = time.Now()
	err := Save(&l)

	return l, err
}

func (l *Loan) FindbyID(id int) (*Loan, error) {
	err := FindbyID(&l, id)
	return l, err
}

func (l *Loan) FilterSearchSingle(filter interface{}) (*Loan, error) {
	err := FilterSearchSingle(&l, filter)
	return l, err
}

func (l *Loan) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result PagedSearchResult, err error) {
	loans := []Loan{}
	result, err = PagedFilterSearch(&loans, page, rows, orderby, sort, filter)

	return result, err
}

func (l *Loan) Approve() error {
	l.Status = "approved"

	_, err := l.Save()
	if err != nil {
		return err
	}

	KafkaSubmitModel(l, "loan")

	return nil
}

func (l *Loan) Reject() error {
	l.Status = "rejected"

	_, err := l.Save()
	if err != nil {
		return err
	}

	KafkaSubmitModel(l, "loan")

	return nil
}
