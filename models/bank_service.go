package models

import (
	"time"

	"gitlab.com/asira-ayannah/basemodel"
)

type (
	BankService struct {
		basemodel.BaseModel
		DeletedTime time.Time `json:"deleted_time" gorm:"column:deleted_time"`
		ServiceID   uint64    `json:"service_id" gorm:"column:service_id"`
		BankID      uint64    `json:"bank_id" gorm:"column:bank_id"`
	}
)

func (model *BankService) Create() error {
	err := basemodel.Create(&model)
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(model, "bank_service")

	return err
}

func (model *BankService) Save() error {
	err := basemodel.Save(&model)
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(model, "bank_service")

	return err
}

func (model *BankService) Delete() error {
	err := basemodel.Delete(&model)
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(model, "bank_service_delete")

	return err
}

func (model *BankService) FindbyID(id int) error {
	err := basemodel.FindbyID(&model, id)
	return err
}

func (model *BankService) PagedFindFilter(page int, rows int, order []string, sort []string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	bank_type := []BankService{}
	result, err = basemodel.PagedFindFilter(&bank_type, page, rows, order, sort, filter)

	return result, err
}
