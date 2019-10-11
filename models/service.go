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

func (s *Service) Create() error {
	err := basemodel.Create(&s)
	return err
}

func (s *Service) Save() error {
	err := basemodel.Save(&s)
	return err
}

func (s *Service) Delete() error {
	err := basemodel.Delete(&s)
	return err
}

func (s *Service) FindbyID(id int) error {
	err := basemodel.FindbyID(&s, id)
	return err
}

func (s *Service) PagedFindFilter(page int, rows int, order []string, sort []string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	services := []Service{}
	result, err = basemodel.PagedFindFilter(&services, page, rows, order, sort, filter)

	return result, err
}
