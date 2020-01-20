package models

import (
	"asira_lender/asira"

	"github.com/ayannahindonesia/basemodel"
)

// AgentProvider model
type AgentProvider struct {
	basemodel.BaseModel
	Name    string `json:"name" gorm:"column:name"`
	PIC     string `json:"pic" gorm:"column:pic"`
	Phone   string `json:"phone" gorm:"column:phone"`
	Address string `json:"address" gorm:"column:address"`
	Status  string `json:"status" gorm:"column:status"`
}

// Create new
func (model *AgentProvider) Create() error {
	err := basemodel.Create(&model)
	if err != nil {
		return err
	}

	return err
}

// BeforeSave gorm callback
func (model *AgentProvider) BeforeSave() error {
	if model.Status == "inactive" {
		deactivateAgents(model.ID)
	}
	return nil
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
func (model *AgentProvider) FindbyID(id uint64) error {
	err := basemodel.FindbyID(&model, id)
	return err
}

// PagedFilterSearch paged list
func (model *AgentProvider) PagedFilterSearch(page int, rows int, order []string, sorts []string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	agentProviders := []AgentProvider{}
	result, err = basemodel.PagedFindFilter(&agentProviders, page, rows, order, sorts, filter)

	return result, err
}

func deactivateAgents(providerID uint64) {
	db := asira.App.DB

	db.Model(&Agent{}).Where("agent_provider = ?", providerID).Update("status", "inactive")
}
