package models

import (
	"database/sql"
	"time"

	"github.com/ayannahindonesia/basemodel"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Agent struct {
	basemodel.BaseModel
	DeletedTime   time.Time     `json:"deleted_time" gorm:"column:deleted_time"`
	Name          string        `json:"name" gorm:"column:name"`
	Username      string        `json:"username" gorm:"column:username"`
	Password      string        `json:"password" gorm:"column:password"`
	ImageID       sql.NullInt64 `json:"image_id" gorm:"column:image_id"`
	Email         string        `json:"email" gorm:"column:email"`
	Phone         string        `json:"phone" gorm:"column:phone"`
	Category      string        `json:"category" gorm:"column:category"`
	AgentProvider sql.NullInt64 `json:"agent_provider" gorm:"column:agent_provider"`
	Banks         pq.Int64Array `json:"banks" gorm:"column:banks"`
	Status        string        `json:"status" gorm:"column:status"`
}

// BeforeCreate gorm callback
func (model *Agent) BeforeCreate() (err error) {
	passwordByte, err := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	model.Password = string(passwordByte)
	return nil
}

// Create new agent
func (model *Agent) Create() error {
	err := basemodel.Create(&model)
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(model, "agent")

	return err
}

// Save update agent
func (model *Agent) Save() error {
	err := basemodel.Save(&model)
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(model, "agent")

	return err
}

// Delete agent
func (model *Agent) Delete() error {
	err := basemodel.Delete(&model)
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(model, "agent_delete")

	return err
}

// FindbyID find agent with id
func (model *Agent) FindbyID(id int) error {
	err := basemodel.FindbyID(&model, id)
	return err
}

// FilterSearchSingle search using filter and return last
func (model *Agent) FilterSearchSingle(filter interface{}) error {
	err := basemodel.SingleFindFilter(&model, filter)
	return err
}

// PagedFilterSearch search using filter and return with pagination format
func (model *Agent) PagedFilterSearch(page int, rows int, order []string, sort []string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	agents := []Agent{}
	result, err = basemodel.PagedFindFilter(&agents, page, rows, order, sort, filter)

	return result, err
}
