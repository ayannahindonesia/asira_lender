package models

import (
	"github.com/ayannahindonesia/basemodel"
	"github.com/google/uuid"
)

// Client struct
type Client struct {
	basemodel.BaseModel
	Name   string `json:"name" gorm:"column:name"`
	Key    string `json:"key" gorm:"column:key"`
	Secret string `json:"secret" gorm:"column:secret"`
}

// BeforeCreate callback
func (model *Client) BeforeCreate() (err error) {
	if len(model.Secret) < 1 {
		model.Secret = uuid.New().String()
	}
	return nil
}

// Create func
func (model *Client) Create() (err error) {
	return basemodel.Create(&model)
}

// Save func
func (model *Client) Save() (err error) {
	return basemodel.Save(&model)
}

// Delete func
func (model *Client) Delete() (err error) {
	return basemodel.Delete(&model)
}

// FindbyID func
func (model *Client) FindbyID(id uint64) error {
	return basemodel.FindbyID(&model, id)
}

// FilterSearchSingle func
func (model *Client) SingleFindFilter(filter interface{}) (err error) {
	return basemodel.SingleFindFilter(&model, filter)
}
