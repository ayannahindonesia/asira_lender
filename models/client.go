package models

import (
	"github.com/google/uuid"
	"gitlab.com/asira-ayannah/basemodel"
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
	err = basemodel.Create(&model)
	return err
}

// Save func
func (model *Client) Save() (err error) {
	err = basemodel.Save(&model)
	return err
}

// Delete func
func (model *Client) Delete() (err error) {
	err = basemodel.Delete(&model)
	return err
}

// FindbyID func
func (model *Client) FindbyID(id int) error {
	err := basemodel.FindbyID(&model, id)
	return err
}

// FilterSearchSingle func
func (model *Client) FilterSearchSingle(filter interface{}) (err error) {
	err = basemodel.SingleFindFilter(&model, filter)
	return err
}
