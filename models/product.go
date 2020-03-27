package models

import (
	"github.com/ayannahindonesia/basemodel"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

// Product main type
type Product struct {
	basemodel.BaseModel
	Name            string         `json:"name" gorm:"column:name;type:varchar(255)"`
	ServiceID       uint64         `json:"service_id" gorm:"column:service_id"`
	MinTimeSpan     int            `json:"min_timespan" gorm:"column:min_timespan"`
	MaxTimeSpan     int            `json:"max_timespan" gorm:"column:max_timespan"`
	Interest        float64        `json:"interest" gorm:"column:interest"`
	InterestType    string         `json:"interest_type" gorm:"column:interest_type"`
	MinLoan         int            `json:"min_loan" gorm:"column:min_loan"`
	MaxLoan         int            `json:"max_loan" gorm:"column:max_loan"`
	Fees            postgres.Jsonb `json:"fees" gorm:"column:fees"`
	Collaterals     pq.StringArray `json:"collaterals" gorm:"column:collaterals"`
	FinancingSector pq.StringArray `json:"financing_sector" gorm:"column:financing_sector"`
	Assurance       string         `json:"assurance" gorm:"column:assurance"`
	Status          string         `json:"status" gorm:"column:status;type:varchar(255)"`
	Form            postgres.Jsonb `json:"form" gorm:"column:form;type:text"`
	Description     string         `json:"description" gorm:"column:description;type:text"`
}

// Create func
func (model *Product) Create() error {
	return basemodel.Create(&model)
}

// Save func
func (model *Product) Save() error {
	return basemodel.Save(&model)
}

// FirstOrCreate func
func (model *Product) FirstOrCreate() error {
	return basemodel.FirstOrCreate(&model)
}

// Delete func
func (model *Product) Delete() error {
	return basemodel.Delete(&model)
}

// FindbyID func
func (model *Product) FindbyID(id uint64) error {
	return basemodel.FindbyID(&model, id)
}

// PagedFindFilter func
func (model *Product) PagedFindFilter(page int, rows int, order []string, sort []string, filter interface{}) (basemodel.PagedFindResult, error) {
	products := []Product{}

	return basemodel.PagedFindFilter(&products, page, rows, order, sort, filter)
}
