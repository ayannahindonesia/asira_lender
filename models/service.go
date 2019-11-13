package models

import (
	"time"

	"gitlab.com/asira-ayannah/basemodel"
)

type (
	Service struct {
		basemodel.BaseModel
		DeletedTime time.Time `json:"deleted_time" gorm:"column:deleted_time"`
		Name        string    `json:"name" gorm:"column:name;type:varchar(255)"`
		ImageID     uint64    `json:"image_id" gorm:"column:image_id"`
		Status      string    `json:"status" gorm:"column:status;type:varchar(255)"`
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

func (model *Service) FindbyID(id int) error {
	err := basemodel.FindbyID(&model, id)
	return err
}

func (model *Service) PagedFindFilter(page int, rows int, order []string, sort []string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	services := []Service{}
	result, err = basemodel.PagedFindFilter(&services, page, rows, order, sort, filter)

	return result, err
}
