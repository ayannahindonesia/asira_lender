package models

import (
	"database/sql"
	"time"

	"gitlab.com/asira-ayannah/basemodel"
)

type (
	Borrower struct {
		basemodel.BaseModel
		DeletedTime          time.Time     `json:"deleted_time" gorm:"column:deleted_time"`
		Status               string        `json:"status" gorm:"column:status"`
		Fullname             string        `json:"fullname" gorm:"column:fullname;type:varchar(255);not_null" csv:"fullname"`
		Gender               string        `json:"gender" gorm:"column:gender;type:varchar(1);not null csv:"gender"`
		IdCardNumber         string        `json:"idcard_number" gorm:"column:idcard_number;type:varchar(255);unique;not null" csv:"idcard_number"`
		IdCardImageID        string        `json:"idcard_imageid" gorm:"column:idcard_imageid;type:varchar(255)" csv:"idcard_imageid"`
		TaxIDnumber          string        `json:"taxid_number" gorm:"column:taxid_number;type:varchar(255)" csv:"taxid_number"`
		TaxIDImageID         string        `json:"taxid_imageid" gorm:"column:taxid_imageid;type:varchar(255)" csv:"taxid_imageid"`
		Email                string        `json:"email" gorm:"column:email;type:varchar(255);unique" csv:"email"`
		Birthday             time.Time     `json:"birthday" gorm:"column:birthday;not null" csv:"birthday"`
		Birthplace           string        `json:"birthplace" gorm:"column:birthplace;type:varchar(255);not null" csv:"birthplace"`
		LastEducation        string        `json:"last_education" gorm:"column:last_education;type:varchar(255);not null" csv:"last_education"`
		MotherName           string        `json:"mother_name" gorm:"column:mother_name;type:varchar(255);not null" csv:"mother_name"`
		Phone                string        `json:"phone" gorm:"column:phone;type:varchar(255);unique;not null" csv:"phone"`
		MarriedStatus        string        `json:"marriage_status" gorm:"column:marriage_status;type:varchar(255);not null" csv:"marriage_status"`
		SpouseName           string        `json:"spouse_name" gorm:"column:spouse_name;type:varchar(255)" csv:"spouse_name"`
		SpouseBirthday       time.Time     `json:"spouse_birthday" gorm:"column:spouse_birthday" csv:"spouse_birthday"`
		SpouseLastEducation  string        `json:"spouse_lasteducation" gorm:"column:spouse_lasteducation;type:varchar(255)" csv:"spouse_lasteducation"`
		Dependants           int           `json:"dependants,omitempty" gorm:"column:dependants;type:int" sql:"DEFAULT:0" csv:"dependants,omitempty"`
		Address              string        `json:"address" gorm:"column:address;type:varchar(255);not null" csv:"address"`
		Province             string        `json:"province" gorm:"column:province;type:varchar(255);not null" csv:"province"`
		City                 string        `json:"city" gorm:"column:city;type:varchar(255);not null" csv:"city"`
		NeighbourAssociation string        `json:"neighbour_association" gorm:"column:neighbour_association;type:varchar(255);not null" csv:"neighbour_association"`
		Hamlets              string        `json:"hamlets" gorm:"column:hamlets;type:varchar(255);not null" csv:"hamlets"`
		HomePhoneNumber      string        `json:"home_phonenumber" gorm:"column:home_phonenumber;type:varchar(255)" csv:"home_phonenumber"`
		Subdistrict          string        `json:"subdistrict" gorm:"column:subdistrict;type:varchar(255)";not null csv:"subdistrict"`
		UrbanVillage         string        `json:"urban_village" gorm:"column:urban_village;type:varchar(255)";not null csv:"urban_village"`
		HomeOwnership        string        `json:"home_ownership" gorm:"column:home_ownership;type:varchar(255);not null csv:"home_ownership"`
		LivedFor             int           `json:"lived_for" gorm:"column:lived_for;type:int;not null" csv:"lived_for"`
		Occupation           string        `json:"occupation" gorm:"column:occupation;type:varchar(255);not null" csv:"occupation"`
		EmployeeID           string        `json:"employee_id" gorm:"column:employee_id;type:varchar(255)" csv:"employee_id"`
		EmployerName         string        `json:"employer_name" gorm:"column:employer_name;type:varchar(255);not null" csv:"employer_name"`
		EmployerAddress      string        `json:"employer_address" gorm:"column:employer_address;type:varchar(255);not null" csv:"employer_address"`
		Department           string        `json:"department" gorm:"column:department;type:varchar(255);not null" csv:"department"`
		BeenWorkingFor       int           `json:"been_workingfor" gorm:"column:been_workingfor;type:int;not null" csv:"been_workingfor"`
		DirectSuperior       string        `json:"direct_superiorname" gorm:"column:direct_superiorname;type:varchar(255)" csv:"direct_superiorname"`
		EmployerNumber       string        `json:"employer_number" gorm:"column:employer_number;type:varchar(255);not null" csv:"employer_number"`
		MonthlyIncome        int           `json:"monthly_income" gorm:"column:monthly_income;type:int;not null" csv:"monthly_income"`
		OtherIncome          int           `json:"other_income" gorm:"column:other_income;type:int" csv:"other_income"`
		OtherIncomeSource    string        `json:"other_incomesource" gorm:"column:other_incomesource;type:varchar(255)" csv:"other_incomesource"`
		FieldOfWork          string        `json:"field_of_work" gorm:"column:field_of_work;type:varchar(255);not null" csv:"field_of_work"`
		RelatedPersonName    string        `json:"related_personname" gorm:"column:related_personname;type:varchar(255);not null" csv:"related_personname"`
		RelatedRelation      string        `json:"related_relation" gorm:"column:related_relation;type:varchar(255);not null" csv:"related_relation"`
		RelatedPhoneNumber   string        `json:"related_phonenumber" gorm:"column:related_phonenumber;type:varchar(255);not null" csv:"related_phonenumber"`
		RelatedHomePhone     string        `json:"related_homenumber" gorm:"column:related_homenumber;type:varchar(255)" csv:"related_phonenumber"`
		RelatedAddress       string        `json:"related_address" gorm:"column:related_address;type:text" csv:"related_address"`
		Bank                 sql.NullInt64 `json:"bank" gorm:"column:bank" sql:"DEFAULT:NULL" csv:"bank"`
		BankAccountNumber    string        `json:"bank_accountnumber" gorm:"column:bank_accountnumber" csv:"bank_accountnumber"`
	}
)

func (b *Borrower) Create() error {
	err := basemodel.Create(&b)
	return err
}

func (b *Borrower) Save() error {
	err := basemodel.Save(&b)
	return err
}

func (b *Borrower) Delete() error {
	b.DeletedTime = time.Now()
	err := basemodel.Save(&b)

	return err
}

func (b *Borrower) FindbyID(id int) error {
	err := basemodel.FindbyID(&b, id)
	return err
}

func (b *Borrower) FilterSearchSingle(filter interface{}) error {
	err := basemodel.SingleFindFilter(&b, filter)
	return err
}

func (b *Borrower) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	borrowers := []Borrower{}

	order := []string{orderby}
	sorts := []string{sort}
	result, err = basemodel.PagedFindFilter(&borrowers, page, rows, order, sorts, filter)

	return result, err
}
