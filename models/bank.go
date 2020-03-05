package models

import (
	"github.com/lib/pq"

	"github.com/ayannahindonesia/basemodel"
)

// Bank main type
type Bank struct {
	basemodel.BaseModel
	Name     string        `json:"name" gorm:"column:name;type:varchar(255)"`
	Image    string        `json:"image" gorm:"column:image;type:text"`
	Type     uint64        `json:"type" gorm:"column:type;type:bigserial"`
	Address  string        `json:"address" gorm:"column:address;type:text"`
	Province string        `json:"province" gorm:"column:province;type:varchar(255)"`
	City     string        `json:"city" gorm:"column:city;type:varchar(255)"`
	PIC      string        `json:"pic" gorm:"column:pic;type:varchar(255)"`
	Phone    string        `json:"phone" gorm:"column:phone;type:varchar(255)"`
	Services pq.Int64Array `json:"services" gorm "column:services"`
	Products pq.Int64Array `json:"products" gorm "column:products"`
}

// Create func
func (model *Bank) Create() error {
	return basemodel.Create(&model)
}

// Save func
func (model *Bank) Save() error {
	return basemodel.Save(&model)
}

// FirstOrCreate func
func (model *Bank) FirstOrCreate() error {
	return basemodel.FirstOrCreate(&model)
}

// Delete func
func (model *Bank) Delete() error {
	return basemodel.Delete(&model)
}

// FindbyID func
func (model *Bank) FindbyID(id uint64) error {
	return basemodel.FindbyID(&model, id)
}

// FindFilter func
func (model *Bank) FindFilter(order []string, sort []string, limit int, offset int, filter interface{}) ([]Bank, error) {
	banks := []Bank{}
	_, err := basemodel.FindFilter(&banks, order, sort, limit, offset, filter)
	return banks, err
}

// PagedFindFilter func
func (model *Bank) PagedFindFilter(page int, rows int, order []string, sort []string, filter interface{}) (basemodel.PagedFindResult, error) {
	banks := []Bank{}

	return basemodel.PagedFindFilter(&banks, page, rows, order, sort, filter)
}
