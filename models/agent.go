package models

import (
	"asira_lender/email"
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/ayannahindonesia/basemodel"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Agent main type
type Agent struct {
	basemodel.BaseModel
	Name          string        `json:"name" gorm:"column:name"`
	Username      string        `json:"username" gorm:"column:username"`
	Password      string        `json:"password" gorm:"column:password"`
	Image         string        `json:"image" gorm:"column:image"`
	Email         string        `json:"email" gorm:"column:email"`
	Phone         string        `json:"phone" gorm:"column:phone"`
	Category      string        `json:"category" gorm:"column:category"`
	AgentProvider sql.NullInt64 `json:"agent_provider" gorm:"column:agent_provider"`
	Banks         pq.Int64Array `json:"banks" gorm:"column:banks"`
	Status        string        `json:"status" gorm:"column:status"`
}

// BeforeCreate gorm callback
func (model *Agent) BeforeCreate() (err error) {
	if len(model.Password) < 1 {
		model.Password = randString(8)
	}
	passwordByte, err := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	model.SendPasswordEmail(model.Password)

	model.Password = string(passwordByte)
	return nil
}

// Create new agent
func (model *Agent) Create() error {
	return basemodel.Create(&model)
}

// Save update agent
func (model *Agent) Save() error {
	return basemodel.Save(&model)
}

// FirstOrCreate saves, or create if not exist
func (model *Agent) FirstOrCreate() error {
	return basemodel.FirstOrCreate(&model)
}

// Delete agent
func (model *Agent) Delete() error {
	return basemodel.Delete(&model)
}

// FindbyID find agent with id
func (model *Agent) FindbyID(id uint64) error {
	return basemodel.FindbyID(&model, id)
}

// SingleFindFilter search using filter and return last
func (model *Agent) SingleFindFilter(filter interface{}) error {
	return basemodel.SingleFindFilter(&model, filter)
}

// PagedFilterSearch search using filter and return with pagination format
func (model *Agent) PagedFilterSearch(page int, rows int, order []string, sort []string, filter interface{}) (basemodel.PagedFindResult, error) {
	agents := []Agent{}

	return basemodel.PagedFindFilter(&agents, page, rows, order, sort, filter)
}

// SendPasswordEmail sends password to agent
func (model *Agent) SendPasswordEmail(password string) {
	message := fmt.Sprintf("Selamat bergabung dengan asira sebagai Agent. anda dapat login menggunakan username dengan password : %v", password)

	email.SendMail(model.Email, "Selamat Bergabung dengan Asira", message)
}

func randString(n int) string {
	var letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
