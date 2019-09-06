package models

import "time"

type (
	BankType struct {
		BaseModel
		DeletedTime time.Time `json:"deleted_time" gorm:"column:deleted_time" sql:"DEFAULT:current_timestamp"`
		Name        string    `json:"name" gorm:"name"`
		Description string    `json:"description" gorm:"description"`
	}
)

func (b *BankType) Create() (*BankType, error) {
	err := Create(&b)
	if err != nil {
		return nil, err
	}

	err = KafkaSubmitModel(b, "bank_type")

	return b, err
}

func (b *BankType) Save() (*BankType, error) {
	err := Save(&b)
	if err != nil {
		return nil, err
	}

	err = KafkaSubmitModel(b, "bank_type")

	return b, err
}

func (b *BankType) Delete() (*BankType, error) {
	err := Delete(&b)
	if err != nil {
		return nil, err
	}

	err = KafkaSubmitModel(b, "bank_type_delete")

	return b, err
}

func (b *BankType) FindbyID(id int) (*BankType, error) {
	err := FindbyID(&b, id)
	return b, err
}

func (b *BankType) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result PagedSearchResult, err error) {
	bank_type := []BankType{}
	result, err = PagedFilterSearch(&bank_type, page, rows, orderby, sort, filter)

	return result, err
}
