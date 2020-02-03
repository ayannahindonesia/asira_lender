package models

import (
	"fmt"

	"github.com/ayannahindonesia/basemodel"
)

// BankRepresentatives main type
type BankRepresentatives struct {
	basemodel.BaseModel
	UserID uint64 `json:"user_id" gorm:"column:user_id"`
	BankID uint64 `json:"bank_id" gorm:"column:bank_id"`
}

// Create func
func (model *BankRepresentatives) Create() error {
	return basemodel.Create(&model)
}

// Save func
func (model *BankRepresentatives) Save() error {
	return basemodel.Save(&model)
}

// FindbyID func
func (model *BankRepresentatives) FindbyID(id uint64) error {
	return basemodel.FindbyID(&model, id)
}

// FindbyUserID func
func (model *BankRepresentatives) FindbyUserID(id int) error {
	type Filter struct {
		UserID string `json:"user_id"`
	}
	return basemodel.SingleFindFilter(&model, &Filter{
		UserID: fmt.Sprintf("%v", id),
	})
}

// FindbyBankID func
func (model *BankRepresentatives) FindbyBankID(id int) error {
	type Filter struct {
		BankID string `json:"bank_id"`
	}
	return basemodel.SingleFindFilter(&model, &Filter{
		BankID: string(id),
	})
}
