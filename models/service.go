package models

import (
	"github.com/ayannahindonesia/basemodel"
)

type (
	Service struct {
		basemodel.BaseModel
		Name   string `json:"name" gorm:"column:name;type:varchar(255)"`
		Image  string `json:"image" gorm:"column:image"`
		Status string `json:"status" gorm:"column:status;type:varchar(255)"`
	}
)

func (model *Service) Create() error {
	err := basemodel.Create(&model)
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(model, "service")

	return err
}

func (model *Service) Save() error {
	err := basemodel.Save(&model)
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(model, "service")

	return err
}

func (model *Service) Delete() error {
	err := basemodel.Delete(&model)
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(model, "service_delete")

	return err
}

func (model *Service) FindbyID(id uint64) error {
	err := basemodel.FindbyID(&model, id)
	return err
}

func (model *Service) PagedFindFilter(page int, rows int, order []string, sort []string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	services := []Service{}
	result, err = basemodel.PagedFindFilter(&services, page, rows, order, sort, filter)

	return result, err
}
