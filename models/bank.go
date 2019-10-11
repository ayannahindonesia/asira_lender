package models

import (
	"time"

	"github.com/google/uuid"
	"gitlab.com/asira-ayannah/basemodel"
	"golang.org/x/crypto/bcrypt"
)

type (
	Bank struct {
		basemodel.BaseModel
		DeletedTime         time.Time `json:"deleted_time" gorm:"column:deleted_time" sql:"DEFAULT:current_timestamp"`
		Name                string    `json:"name" gorm:"column:name;type:varchar(255)"`
		Type                uint64    `json:"type" gorm:"column:type;type:bigserial"`
		Address             string    `json:"address" gorm:"column:address;type:text"`
		Province            string    `json:"province" gorm:"column:province;type:varchar(255)"`
		City                string    `json:"city" gorm:"column:city;type:varchar(255)"`
		PIC                 string    `json:"pic" gorm:"column:pic;type:varchar(255)"`
		Phone               string    `json:"phone" gorm:"column:phone;type:varchar(255)"`
		AdminFeeSetup       string    `json:"adminfee_setup" gorm:"column:adminfee_setup;type:varchar(255)"`
		ConvenienceFeeSetup string    `json:"convfee_setup" gorm:"column:convfee_setup;type:varchar(255)"`
		Username            string    `json:"username" gorm:"column:username;type:varchar(255);unique;not null"`
		Password            string    `json:"password" gorm:"column:password;type:text;not null"`
	}
)

// gorm callback hook
func (model *Bank) BeforeCreate() (err error) {
	if len(model.Username) < 1 {
		model.Username = uuid.New().String()
	}
	if len(model.Username) < 1 {
		model.Password = uuid.New().String()
	}
	passwordByte, err := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	model.Password = string(passwordByte)
	return nil
}

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

func (model *Bank) FindbyID(id int) error {
	err := basemodel.FindbyID(&model, id)
	return err
}

func (model *Bank) PagedFindFilter(page int, rows int, order []string, sort []string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	bank_type := []Bank{}
	result, err = basemodel.PagedFindFilter(&bank_type, page, rows, order, sort, filter)

	return result, err
}
