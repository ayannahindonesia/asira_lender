package models

import (
	"github.com/ayannahindonesia/basemodel"
)

// Service main type
type Service struct {
	basemodel.BaseModel
	Name        string `json:"name" gorm:"column:name;type:varchar(255)"`
	Image       string `json:"image" gorm:"column:image"`
	Status      string `json:"status" gorm:"column:status;type:varchar(255)"`
	Description string `json:"description" gorm:"column:description;type:varchar(255)"`
}

// Create func
func (model *Service) Create() error {
	return basemodel.Create(&model)
}

// Save func
func (model *Service) Save() error {
	return basemodel.Save(&model)
}

// FirstOrCreate func
func (model *Service) FirstOrCreate() error {
	return basemodel.FirstOrCreate(&model)
}

// Delete func
func (model *Service) Delete() error {
	return basemodel.Delete(&model)
}

// FindbyID func
func (model *Service) FindbyID(id uint64) error {
	return basemodel.FindbyID(&model, id)
}

// PagedFindFilter func
func (model *Service) PagedFindFilter(page int, rows int, order []string, sort []string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	services := []Service{}

	return basemodel.PagedFindFilter(&services, page, rows, order, sort, filter)
}
