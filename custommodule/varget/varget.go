package varget

import "github.com/ayannahindonesia/basemodel"

// Varget main type
type Varget struct {
	basemodel.BaseModel
	Identifier string `json:"identifier" gorm:"column:identifier;type:varchar(255)"`
	Value      string `json:"value" gorm:"column:value;type:text"`
}

// Create func
func (model *Varget) Create() error {
	return basemodel.Create(&model)
}

// Save func
func (model *Varget) Save() error {
	return basemodel.Save(&model)
}

// Delete func
func (model *Varget) Delete() error {
	return basemodel.Delete(&model)
}

// Get func
func Get(key string, returnDefault string) (result string) {
	model := Varget{}
	if basemodel.SingleFindFilter(&model, &Varget{Identifier: key}) != nil {
		return returnDefault
	}
	return model.Value
}

// Set func
func Set(key string, value string) (err error) {
	model := Varget{}
	if basemodel.SingleFindFilter(&model, &Varget{Identifier: key}) != nil {
		return basemodel.Create(&Varget{Identifier: key, Value: value})
	}
	model.Identifier = key
	model.Value = value
	return basemodel.Save(&model)
}

// PagedFindFilter func
func (model *Varget) PagedFindFilter(page int, rows int, orderby []string, sort []string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	models := []Varget{}

	return basemodel.PagedFindFilter(&models, page, rows, orderby, sort, filter)
}
