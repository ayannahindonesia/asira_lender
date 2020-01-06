package models

import (
	"github.com/ayannahindonesia/basemodel"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type (
	// User main type
	User struct {
		basemodel.BaseModel
		Roles    pq.Int64Array `json:"roles" gorm:"column:roles"`
		Username string        `json:"username" gorm:"column:username;type:varchar(255);unique;not null"`
		Email    string        `json:"email" gorm:"column:email;type:varchar(255)"`
		Phone    string        `json:"phone" gorm:"column:phone;type:varchar(255)"`
		Password string        `json:"password" gorm:"column:password;type:text;not null"`
		Status   string        `json:"status" gorm:"column:status;type:boolean" sql:"DEFAULT:TRUE"`
	}
)

// BeforCreate gorm callback hook
func (model *User) BeforeCreate() (err error) {
	err = model.ChangePassword(model.Password)
	return err
}

// Create func
func (model *User) Create() error {
	err := basemodel.Create(&model)
	return err
}

// Save func
func (model *User) Save() error {
	err := basemodel.Save(&model)
	return err
}

// FindbyID func
func (model *User) FindbyID(id int) error {
	err := basemodel.FindbyID(&model, id)
	return err
}

// FilterSearchSingle func
func (model *User) FilterSearchSingle(filter interface{}) error {
	err := basemodel.SingleFindFilter(&model, filter)
	return err
}

// PagedFilterSearch func
func (model *User) PagedFilterSearch(page int, rows int, orderby []string, sorts []string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	user := []User{}
	result, err = basemodel.PagedFindFilter(&user, page, rows, orderby, sorts, filter)

	return result, err
}

// ChangePassword update password to encrypted. does not save
func (model *User) ChangePassword(rawpassword string) error {
	passwordByte, err := bcrypt.GenerateFromPassword([]byte(rawpassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	model.Password = string(passwordByte)
	return nil
}
