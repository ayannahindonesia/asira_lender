package models

import (
	"fmt"
	"time"

	"github.com/ayannahindonesia/basemodel"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type (
	Loan struct {
		basemodel.BaseModel
		Borrower            uint64         `json:"borrower" gorm:"column:borrower;foreignkey"`
		Status              string         `json:"status" gorm:"column:status;type:varchar(255)" sql:"DEFAULT:'processing'"`
		LoanAmount          float64        `json:"loan_amount" gorm:"column:loan_amount;type:int;not null"`
		Installment         int            `json:"installment" gorm:"column:installment;type:int;not null"` // plan of how long loan to be paid
		Fees                postgres.Jsonb `json:"fees" gorm:"column:fees;type:jsonb"`
		Interest            float64        `json:"interest" gorm:"column:interest;type:int;not null"`
		TotalLoan           float64        `json:"total_loan" gorm:"column:total_loan;type:int;not null"`
		DueDate             time.Time      `json:"due_date" gorm:"column:due_date"`
		LayawayPlan         float64        `json:"layaway_plan" gorm:"column:layaway_plan;type:int;not null"` // how much borrower will pay per month
		Product             uint64         `json:"product" gorm:"column:product;foreignkey"`                  // product and service is later to be discussed
		LoanIntention       string         `json:"loan_intention" gorm:"column:loan_intention;type:varchar(255);not null"`
		IntentionDetails    string         `json:"intention_details" gorm:"column:intention_details;type:text;not null"`
		BorrowerInfo        postgres.Jsonb `json:"borrower_info" gorm:"column:borrower_info;type:jsonb"`
		DisburseDate        time.Time      `json:"disburse_date" gorm:"column:disburse_date"`
		DisburseDateChanged bool           `json:"disburse_date_changed" gorm:"column:disburse_date_changed"`
		DisburseStatus      string         `json:"disburse_status" gorm:"column:disburse_status" sql:"DEFAULT:'processing'"`
		ApprovalDate        time.Time      `json:"approval_date" gorm:"column:approval_date"`
		RejectReason        string         `json:"reject_reason" gorm:"column:reject_reason"`
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
	err := basemodel.Delete(&l)

	return err
}

func (l *Loan) FindbyID(id uint64) error {
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
	l.ApprovalDate = time.Now()

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
	l.ApprovalDate = time.Now()

	err := l.Save()
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(l, "loan")

	return err
}

// DisburseConfirmed confirm disburse
func (model *Loan) DisburseConfirmed() error {
	model.DisburseStatus = "confirmed"

	err := model.Save()
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(model, "loan")

	return err
}

// ChangeDisburseDate func
func (l *Loan) ChangeDisburseDate(disburseDate time.Time) (err error) {
	if l.DisburseDateChanged != true {
		l.DisburseDate = disburseDate
		l.DisburseDateChanged = true

		err = l.Save()
		if err != nil {
			return err
		}

		err = KafkaSubmitModel(l, "loan")
	} else {
		err = fmt.Errorf("disburse date already changed before.")
	}

	return err
}
