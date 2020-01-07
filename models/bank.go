package models

import (
	"github.com/lib/pq"

	"github.com/ayannahindonesia/basemodel"
)

type (
	Bank struct {
		basemodel.BaseModel
		Name                string        `json:"name" gorm:"column:name;type:varchar(255)"`
		Image               string        `json:"image" gorm:"column:image;type:text"`
		Type                uint64        `json:"type" gorm:"column:type;type:bigserial"`
		Address             string        `json:"address" gorm:"column:address;type:text"`
		Province            string        `json:"province" gorm:"column:province;type:varchar(255)"`
		City                string        `json:"city" gorm:"column:city;type:varchar(255)"`
		PIC                 string        `json:"pic" gorm:"column:pic;type:varchar(255)"`
		Phone               string        `json:"phone" gorm:"column:phone;type:varchar(255)"`
		Services            pq.Int64Array `json:"services" gorm "column:services"`
		Products            pq.Int64Array `json:"products" gorm "column:products"`
		AdminFeeSetup       string        `json:"adminfee_setup" gorm:"column:adminfee_setup;type:varchar(255)"`
		ConvenienceFeeSetup string        `json:"convfee_setup" gorm:"column:convfee_setup;type:varchar(255)"`
	}
)

func (model *Bank) Create() error {
	err := basemodel.Create(&model)
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(model, "bank")

	return err
}

func (model *Bank) Save() error {
	err := basemodel.Save(&model)
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(model, "bank")

	return err
}

func (model *Bank) Delete() error {
	err := basemodel.Delete(&model)
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(model, "bank_delete")

	return err
}

func (model *Bank) FindbyID(id uint64) error {
	err := basemodel.FindbyID(&model, id)
	return err
}

func (model *Bank) FindFilter(order []string, sort []string, limit int, offset int, filter interface{}) ([]Bank, error) {
	banks := []Bank{}
	_, err := basemodel.FindFilter(&banks, order, sort, limit, offset, filter)
	return banks, err
}

func (model *Bank) PagedFindFilter(page int, rows int, order []string, sort []string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	bank_type := []Bank{}
	result, err = basemodel.PagedFindFilter(&bank_type, page, rows, order, sort, filter)

	return result, err
}
