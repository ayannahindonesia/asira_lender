package models

import (
	"github.com/lib/pq"
	"gitlab.com/asira-ayannah/basemodel"
	"golang.org/x/crypto/bcrypt"
)

type (
	User struct {
		basemodel.BaseModel
		RoleID   pq.Int64Array `json:"role_id" gorm:"column:role_id"`
		Username string        `json:"username" gorm:"column:username;type:varchar(255);unique;not null"`
		Email    string        `json:"email" gorm:"column:email;type:varchar(255)"`
		Phone    string        `json:"phone" gorm:"column:phone;type:varchar(255)"`
		Password string        `json:"password" gorm:"column:password;type:text;not null"`
		Status   string        `json:"status" gorm:"column:status;type:boolean" sql:"DEFAULT:TRUE"`
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
