package models

import (
	"fmt"

	"github.com/ayannahindonesia/basemodel"
)

type (
	BankRepresentatives struct {
		basemodel.BaseModel
		UserID uint64 `json:"user_id" gorm:"column:user_id"`
		BankID uint64 `json:"bank_id" gorm:"column:bank_id"`
	}
)

func (model *BankRepresentatives) Create() error {
	err := basemodel.Create(&model)
	return err
}

func (model *BankRepresentatives) Save() error {
	err := basemodel.Save(&model)
	return err
}

func (model *BankRepresentatives) FindbyID(id int) error {
	err := basemodel.FindbyID(&model, id)
	return err
}

func (model *BankRepresentatives) FindbyUserID(id int) error {
	type Filter struct {
		UserID string `json:"user_id"`
	}
	err := basemodel.SingleFindFilter(&model, &Filter{
		UserID: fmt.Sprintf("%v", id),
	})
	return err
}

func (model *BankRepresentatives) FindbyBankID(id int) error {
	type Filter struct {
		BankID string `json:"bank_id"`
	}
	err := basemodel.SingleFindFilter(&model, &Filter{
		BankID: string(id),
	})
	return err
}
