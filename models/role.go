package models

import (
	"github.com/lib/pq"
	"gitlab.com/asira-ayannah/basemodel"
)

type (
	Roles struct {
		basemodel.BaseModel
		Name        string         `json:"name" gorm:"column:name"`
		Description string         `json:"description" gorm:"column:description"`
		System      string         `json:"system" gorm:"column:system"`
		Status      string         `json:"status" gorm:"column:status;type:boolean" sql:"DEFAULT:TRUE"`
		Permissions pq.StringArray `json:"permissions" gorm:"column:permissions"`
	}
)

func (b *Roles) Create() error {
	err := basemodel.Create(&b)
	return err
}

func (b *Roles) Save() error {
	err := basemodel.Save(&b)
	return err
}

func (b *Roles) Delete() error {
	err := basemodel.Delete(&b)
	return err
}

func (b *Roles) FindbyID(id int) error {
	err := basemodel.FindbyID(&b, id)
	return err
}

func (b *Roles) FilterSearchSingle(filter interface{}) error {
	err := basemodel.SingleFindFilter(&b, filter)
	return err
}

func (b *Roles) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	role := []Roles{}
	order := []string{orderby}
	sorts := []string{sort}
	result, err = basemodel.PagedFindFilter(&role, page, rows, order, sorts, filter)

	return result, err
}

func (b *Roles) FilterSearch(limit int, offset int, orderby string, sort string, filter interface{}) (result interface{}, err error) {
	role := []Roles{}
	order := []string{orderby}
	sorts := []string{sort}
	result, err = basemodel.FindFilter(&role, order, sorts, limit, offset, filter)

	return result, err
}
