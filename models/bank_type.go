package models

import (
	"github.com/ayannahindonesia/basemodel"
)

// BankType main type
type BankType struct {
	basemodel.BaseModel
	Name        string `json:"name" gorm:"name"`
	Description string `json:"description" gorm:"description"`
}

// Create func
func (model *BankType) Create() error {
	return basemodel.Create(&model)
}

// Save func
func (model *BankType) Save() error {
	return basemodel.Save(&model)
}

// FirstOrCreate func
func (model *BankType) FirstOrCreate() error {
	return basemodel.FirstOrCreate(&model)
}

// Delete func
func (model *BankType) Delete() error {
	return basemodel.Delete(&model)
}

// FindbyID func
func (model *BankType) FindbyID(id uint64) error {
	return basemodel.FindbyID(&model, id)
}

// PagedFilterSearch func
func (model *BankType) PagedFilterSearch(page int, rows int, orderby []string, sort []string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	banktype := []BankType{}

	return basemodel.PagedFindFilter(&banktype, page, rows, orderby, sort, filter)
}
