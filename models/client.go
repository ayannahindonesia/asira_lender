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
func (model *Client) BeforeCreate() error {
	if len(model.Secret) < 1 {
		model.Secret = uuid.New().String()
	}
	return nil
}

// Create func
func (model *Client) Create() error {
	return basemodel.Create(&model)
}

// Save func
func (model *Client) Save() error {
	return basemodel.Save(&model)
}

// Delete func
func (model *Client) Delete() error {
	return basemodel.Delete(&model)
}

// FindbyID func
func (model *Client) FindbyID(id uint64) error {
	return basemodel.FindbyID(&model, id)
}

// SingleFindFilter func
func (model *Client) SingleFindFilter(filter interface{}) error {
	return basemodel.SingleFindFilter(&model, filter)
}
