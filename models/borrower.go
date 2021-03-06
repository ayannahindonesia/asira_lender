package models

import (
	"database/sql"
	"time"

	"github.com/ayannahindonesia/basemodel"
)

// Borrower main type
type Borrower struct {
	basemodel.BaseModel
	Status               string        `json:"status" gorm:"column:status"`
	Fullname             string        `json:"fullname" gorm:"column:fullname;type:varchar(255);not_null" csv:"fullname"`
	Nickname             string        `json:"nickname" gorm:"column:nickname;type:varchar(255)"`
	Gender               string        `json:"gender" gorm:"column:gender;type:varchar(1);not null csv:"gender"`
	ImageProfile         string        `json:"image_profile" gorm:"column:image_profile"`
	IdCardNumber         string        `json:"idcard_number" gorm:"column:idcard_number;type:varchar(255);unique;not null" csv:"idcard_number"`
	IdCardImage          string        `json:"idcard_image" gorm:"column:idcard_image;type:varchar(255)" csv:"idcard_image"`
	TaxIDnumber          string        `json:"taxid_number" gorm:"column:taxid_number;type:varchar(255)" csv:"taxid_number"`
	TaxIDImage           string        `json:"taxid_image" gorm:"column:taxid_image;type:varchar(255)" csv:"taxid_image"`
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
	Address              string        `json:"address" gorm:"column:address;type:varchar(255)" csv:"address"`
	Province             string        `json:"province" gorm:"column:province;type:varchar(255)" csv:"province"`
	City                 string        `json:"city" gorm:"column:city;type:varchar(255)" csv:"city"`
	NeighbourAssociation string        `json:"neighbour_association" gorm:"column:neighbour_association;type:varchar(255)" csv:"neighbour_association"`
	Hamlets              string        `json:"hamlets" gorm:"column:hamlets;type:varchar(255)" csv:"hamlets"`
	HomePhoneNumber      string        `json:"home_phonenumber" gorm:"column:home_phonenumber;type:varchar(255)" csv:"home_phonenumber"`
	Subdistrict          string        `json:"subdistrict" gorm:"column:subdistrict;type:varchar(255)" csv:"subdistrict"`
	UrbanVillage         string        `json:"urban_village" gorm:"column:urban_village;type:varchar(255)" csv:"urban_village"`
	HomeOwnership        string        `json:"home_ownership" gorm:"column:home_ownership;type:varchar(255)" csv:"home_ownership"`
	LivedFor             int           `json:"lived_for" gorm:"column:lived_for;type:int" csv:"lived_for"`
	Occupation           string        `json:"occupation" gorm:"column:occupation;type:varchar(255)" csv:"occupation"`
	EmployeeID           string        `json:"employee_id" gorm:"column:employee_id;type:varchar(255)" csv:"employee_id"`
	EmployerName         string        `json:"employer_name" gorm:"column:employer_name;type:varchar(255)" csv:"employer_name"`
	EmployerAddress      string        `json:"employer_address" gorm:"column:employer_address;type:varchar(255)" csv:"employer_address"`
	Department           string        `json:"department" gorm:"column:department;type:varchar(255)" csv:"department"`
	BeenWorkingFor       int           `json:"been_workingfor" gorm:"column:been_workingfor;type:int" csv:"been_workingfor"`
	DirectSuperior       string        `json:"direct_superiorname" gorm:"column:direct_superiorname;type:varchar(255)" csv:"direct_superiorname"`
	EmployerNumber       string        `json:"employer_number" gorm:"column:employer_number;type:varchar(255)" csv:"employer_number"`
	MonthlyIncome        int           `json:"monthly_income" gorm:"column:monthly_income;type:int" csv:"monthly_income"`
	OtherIncome          int           `json:"other_income" gorm:"column:other_income;type:int" csv:"other_income"`
	OtherIncomeSource    string        `json:"other_incomesource" gorm:"column:other_incomesource;type:varchar(255)" csv:"other_incomesource"`
	FieldOfWork          string        `json:"field_of_work" gorm:"column:field_of_work;type:varchar(255)" csv:"field_of_work"`
	RelatedPersonName    string        `json:"related_personname" gorm:"column:related_personname;type:varchar(255)" csv:"related_personname"`
	RelatedRelation      string        `json:"related_relation" gorm:"column:related_relation;type:varchar(255)" csv:"related_relation"`
	RelatedPhoneNumber   string        `json:"related_phonenumber" gorm:"column:related_phonenumber;type:varchar(255)" csv:"related_phonenumber"`
	RelatedHomePhone     string        `json:"related_homenumber" gorm:"column:related_homenumber;type:varchar(255)" csv:"related_phonenumber"`
	RelatedAddress       string        `json:"related_address" gorm:"column:related_address;type:text" csv:"related_address"`
	Bank                 sql.NullInt64 `json:"bank" gorm:"column:bank" sql:"DEFAULT:NULL" csv:"bank"`
	BankAccountNumber    string        `json:"bank_accountnumber" gorm:"column:bank_accountnumber" csv:"bank_accountnumber"`
	AgentReferral        sql.NullInt64 `json:"agent_referral" gorm:"column:agent_referral"`
	OTPverified          bool          `json:"otp_verified" gorm:"column:otp_verified;type:boolean" sql:"DEFAULT:FALSE"`
}

// Create func
func (model *Borrower) Create() error {
	return basemodel.Create(&model)
}

// Save func
func (model *Borrower) Save() error {
	return basemodel.Save(&model)
}

// FirstOrCreate func
func (model *Borrower) FirstOrCreate() error {
	return basemodel.FirstOrCreate(&model)
}

// Delete func
func (model *Borrower) Delete() error {
	return basemodel.Delete(&model)
}

// FindbyID func
func (model *Borrower) FindbyID(id uint64) error {
	return basemodel.FindbyID(&model, id)
}

// SingleFindFilter func
func (model *Borrower) SingleFindFilter(filter interface{}) error {
	return basemodel.SingleFindFilter(&model, filter)
}

// PagedFindFilter func
func (model *Borrower) PagedFindFilter(page int, rows int, orderby []string, sort []string, filter interface{}) (basemodel.PagedFindResult, error) {
	borrowers := []Borrower{}

	return basemodel.PagedFindFilter(&borrowers, page, rows, orderby, sort, filter)
}
