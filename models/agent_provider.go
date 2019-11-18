package models

import (
	"time"

	"gitlab.com/asira-ayannah/basemodel"
)

// AgentProvider model
type AgentProvider struct {
	basemodel.BaseModel
	DeletedTime time.Time `json:"deleted_time" gorm:"column:deleted_time" sql:"DEFAULT:current_timestamp"`
	Name        string    `json:"name" gorm:"column:name"`
	PIC         string    `json:"pic" gorm:"column:pic"`
	Phone       string    `json:"phone" gorm:"column:phone"`
	Address     string    `json:"address" gorm:"column:address"`
	Status      string    `json:"status" gorm:"column:status"`
}

// Create new
func (model *AgentProvider) Create() error {
	err := basemodel.Create(&model)
	if err != nil {
		return err
	}

	return err
}

// Save update
func (model *AgentProvider) Save() error {
	err := basemodel.Save(&model)
	if err != nil {
		return err
	}

	return err
}

// Delete model
func (model *AgentProvider) Delete() error {
	err := basemodel.Delete(&model)
	if err != nil {
		return err
	}

	return err
}

// FindbyID func
func (model *AgentProvider) FindbyID(id int) error {
	err := basemodel.FindbyID(&model, id)
	return err
}

// PagedFilterSearch paged list
func (model *AgentProvider) PagedFilterSearch(page int, rows int, order []string, sorts []string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	agentProviders := []AgentProvider{}
	result, err = basemodel.PagedFindFilter(&agentProviders, page, rows, order, sorts, filter)

	return result, err
}
