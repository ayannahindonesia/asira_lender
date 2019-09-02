package models

import "github.com/google/uuid"

type Internals struct {
	BaseModel
	Name   string `json:"name" gorm:"column:name"`
	Secret string `json:"secret" gorm:"column:secret"`
	Key    string `json:"key" gorm:"column:key"`
	Role   string `json:"role" gorm:"column:role"`
}

func (i *Internals) BeforeCreate() (err error) {
	if len(i.Secret) < 1 {
		i.Secret = uuid.New().String()
	}
	return nil
}

func (i *Internals) Create() (err error) {
	err = Create(&i)
	return err
}

func (i *Internals) Save() (err error) {
	err = Save(&i)
	return err
}

func (i *Internals) Delete() (err error) {
	err = Delete(&i)
	return err
}

func (l *Internals) FilterSearchSingle(filter interface{}) (err error) {
	err = FilterSearchSingle(&l, filter)
	return err
}
