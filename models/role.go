package models

import (
	"github.com/ayannahindonesia/basemodel"
	"github.com/lib/pq"
)

type (
	// Roles main type
	Roles struct {
		basemodel.BaseModel
		Name        string         `json:"name" gorm:"column:name"`
		Description string         `json:"description" gorm:"column:description"`
		System      string         `json:"system" gorm:"column:system"`
		Status      string         `json:"status" gorm:"column:status" sql:"DEFAULT:active"`
		Permissions pq.StringArray `json:"permissions" gorm:"column:permissions"`
	}
)

// Create new
func (model *Roles) Create() error {
	err := basemodel.Create(&model)
	return err
}

// Save role
func (model *Roles) Save() error {
	err := basemodel.Save(&model)
	return err
}

// Delete role
func (model *Roles) Delete() error {
	err := basemodel.Delete(&model)
	return err
}

// FindbyID self explanatory
func (model *Roles) FindbyID(id uint64) error {
	err := basemodel.FindbyID(&model, id)
	return err
}

// FilterSearchSingle use filter to find one role
func (model *Roles) FilterSearchSingle(filter interface{}) error {
	err := basemodel.SingleFindFilter(&model, filter)
	return err
}

// PagedFindFilter use filter to find all matching role, return using paging format
func (model *Roles) PagedFindFilter(page int, rows int, orderby []string, sort []string, filter interface{}) (basemodel.PagedFindResult, error) {
	role := []Roles{}

	return basemodel.PagedFindFilter(&role, page, rows, orderby, sort, filter)
}

// FindFilter use filter to find all matching role
func (model *Roles) FindFilter(limit int, offset int, orderby []string, sort []string, filter interface{}) (interface{}, error) {
	role := []Roles{}

	return basemodel.FindFilter(&role, orderby, sort, limit, offset, filter)
}
