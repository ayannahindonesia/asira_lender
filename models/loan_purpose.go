package models

import (
	"github.com/ayannahindonesia/basemodel"
)

// LoanPurpose main type
type LoanPurpose struct {
	basemodel.BaseModel
	Name   string `json:"name" gorm:"column:name"`
	Status string `json:"status" gorm:"column:status"`
}

// Create func
func (model *LoanPurpose) Create() error {
	return basemodel.Create(&model)
}

// Save func
func (model *LoanPurpose) Save() error {
	return basemodel.Save(&model)
}

// FirstOrCreate func
func (model *LoanPurpose) FirstOrCreate() error {
	return basemodel.FirstOrCreate(&model)
}

// Delete func
func (model *LoanPurpose) Delete() error {
	return basemodel.Delete(&model)
}

// FindbyID func
func (model *LoanPurpose) FindbyID(id uint64) error {
	return basemodel.FindbyID(&model, id)
}

// SingleFindFilter func
func (model *LoanPurpose) SingleFindFilter(filter interface{}) error {
	return basemodel.SingleFindFilter(&model, filter)
}

// PagedFindFilter func
func (model *LoanPurpose) PagedFindFilter(page int, rows int, orderby []string, sort []string, filter interface{}) (basemodel.PagedFindResult, error) {
	loanpurposes := []LoanPurpose{}

	return basemodel.PagedFindFilter(&loanpurposes, page, rows, orderby, sort, filter)
}
