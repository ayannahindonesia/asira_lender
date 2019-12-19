package models

import (
	"log"

	"github.com/ayannahindonesia/basemodel"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type (
	User struct {
		basemodel.BaseModel
		Roles      pq.Int64Array `json:"roles" gorm:"column:roles"`
		Username   string        `json:"username" gorm:"column:username;type:varchar(255);unique;not null"`
		Email      string        `json:"email" gorm:"column:email;type:varchar(255)"`
		Phone      string        `json:"phone" gorm:"column:phone;type:varchar(255)"`
		Password   string        `json:"password" gorm:"column:password;type:text;not null"`
		Status     string        `json:"status" gorm:"column:status;type:boolean" sql:"DEFAULT:TRUE"`
		FirstLogin bool          `json:"first_login" gorm:"column:first_login;type:boolean" sql:"DEFAULT:TRUE"`
	}
)

// gorm callback hook
func (u *User) BeforeCreate() (err error) {
	passwordByte, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(passwordByte)
	return nil
}

func (u *User) Create() error {
	err := basemodel.Create(&u)
	return err
}

// gorm callback hook
func (u *User) BeforeSave() (err error) {
	return nil
}

func (u *User) Save() error {
	err := basemodel.Save(&u)
	return err
}

func (u *User) FindbyID(id int) error {
	err := basemodel.FindbyID(&u, id)
	return err
}

func (u *User) FilterSearchSingle(filter interface{}) error {
	err := basemodel.SingleFindFilter(&u, filter)
	return err
}

func (u *User) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	user := []User{}
	order := []string{orderby}
	sorts := []string{sort}
	result, err = basemodel.PagedFindFilter(&user, page, rows, order, sorts, filter)

	return result, err
}

// FirstLoginChangePassword set new password and first login to false
func (model *User) FirstLoginChangePassword(password string) {
	passwordByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}

	model.Password = string(passwordByte)
	model.FirstLogin = false

	model.Save()
}
