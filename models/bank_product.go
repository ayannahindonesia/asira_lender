package models

import (
	"time"

	"gitlab.com/asira-ayannah/basemodel"
)

type (
	BankProduct struct {
		basemodel.BaseModel
		DeletedTime time.Time `json:"deleted_time" gorm:"column:deleted_time"`
		ProductID   uint64    `json:"product_id" gorm:"column:product_id"`
		BankID      uint64    `json:"bank_id" gorm:"column:bank_id"`
	}
)

func (model *BankProduct) Create() error {
	err := basemodel.Create(&model)
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(model, "bank_product")

	return err
}

func (model *BankProduct) Save() error {
	err := basemodel.Save(&model)
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(model, "bank_product")

	return err
}

func (model *BankProduct) Delete() error {
	err := basemodel.Delete(&model)
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(model, "bank_product_delete")

	return err
}

func (model *BankProduct) FindbyID(id int) error {
	err := basemodel.FindbyID(&model, id)
	return err
}

func (model *BankProduct) PagedFindFilter(page int, rows int, order []string, sort []string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	products := []BankProduct{}
	result, err = basemodel.PagedFindFilter(&products, page, rows, order, sort, filter)

	return result, err
}
