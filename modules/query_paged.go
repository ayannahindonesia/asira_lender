package modules

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/ayannahindonesia/basemodel"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type QueryPaged struct {
	Result    basemodel.PagedFindResult
	TotalRows int
	Offset    int
	Rows      int
	Page      int
	LastPage  int
	Order     []string
	Sort      []string
	c         echo.Context
}

//QueryFunc user defined func must implement before call Exec
type QueryFunc func(*gorm.DB, interface{}) error

//Init all atribute
func (mod *QueryPaged) Init(c echo.Context) error {

	//store context
	mod.c = c

	// pagination parameters
	mod.Rows, _ = strconv.Atoi(c.QueryParam("rows"))
	mod.Page, _ = strconv.Atoi(c.QueryParam("page"))
	mod.Order = strings.Split(c.QueryParam("orderby"), ",")
	mod.Sort = strings.Split(c.QueryParam("sort"), ",")

	// pagination parameters
	if mod.Rows > 0 {
		if mod.Page <= 0 {
			mod.Page = 1
		}
		mod.Offset = (mod.Page * mod.Rows) - mod.Rows
	}

	return nil
}

//Exec custom query
func (mod *QueryPaged) Exec(db *gorm.DB, data interface{}, qFunc QueryFunc) error {

	//generate query sorting
	if len(mod.Order) > 0 {
		if len(mod.Sort) > 0 {
			for k, v := range mod.Order {
				q := v
				if len(mod.Sort) > k {
					value := mod.Sort[k]
					if strings.ToUpper(value) == "ASC" || strings.ToUpper(value) == "DESC" {
						q = v + " " + strings.ToUpper(value)
					}
				}
				db = db.Order(q)
			}
		}
	}

	//new instance
	tempDB := db
	tempDB.Count(&mod.TotalRows)

	if mod.Rows > 0 {
		db = db.Limit(mod.Rows).Offset(mod.Offset)
		mod.LastPage = int(math.Ceil(float64(mod.TotalRows) / float64(mod.Rows)))
	}

	//call user defined function
	return qFunc(db, data)
}

//GetPage result
func (mod *QueryPaged) GetPage(data interface{}) basemodel.PagedFindResult {

	result := basemodel.PagedFindResult{
		TotalData:   mod.TotalRows,
		Rows:        mod.Rows,
		CurrentPage: mod.Page,
		LastPage:    mod.LastPage,
		From:        mod.Offset + 1,
		To:          mod.Offset + mod.Rows,
		Data:        data,
	}

	return result
}

//GenerateFilters generate parameters filter
func (mod *QueryPaged) GenerateFilters(db *gorm.DB, filter interface{}, tableName string) (*gorm.DB, error) {

	//get reflect data
	val := reflect.ValueOf(filter)

	//loop over fields
	for i := 0; i < val.Type().NumField(); i++ {
		t := val.Type().Field(i)
		fieldName := t.Name

		//if not empty json
		if jsonTag := t.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			//explode by comma
			if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {
				fieldName = jsonTag[:commaIdx]
			} else {
				fieldName = jsonTag
			}

			//cek datatype per field
			switch val.Type().Field(i).Type.Kind() {
			case reflect.String:
				if field := mod.c.QueryParam(fieldName); len(field) > 0 {
					db = db.Where("LOWER("+tableName+"."+fieldName+") LIKE ?", "%"+strings.ToLower(field)+"%")
				}
				break
			case reflect.Int64:
				if field := mod.c.QueryParam(fieldName); len(field) > 0 {
					db = db.Where(tableName+"."+fieldName+" IN (?)", field)
				}
				break
			}
		}

		// fmt.Println(fieldName + " : " + val.Type().Field(i).Type.String())
		// fmt.Printf("%+v\n\n", val.Type().Field(i))
		fmt.Printf("+ %+v\n\n", fieldName)
	}

	fmt.Printf(">>>>>>>>> %+v\n\n", db.QueryExpr())

	return db, nil
}
