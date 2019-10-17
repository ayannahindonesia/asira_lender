package migration

import (
	"asira_lender/asira"
	"asira_lender/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/lib/pq"

	"github.com/jinzhu/gorm/dialects/postgres"
)

func Seed() {
	if asira.App.ENV == "development" {
		// seed internals
		internals := []models.Internals{
			models.Internals{
				Name:   "admin",
				Key:    "adminkey",
				Role:   "admin",
				Secret: "adminsecret",
			},
			models.Internals{
				Name:   "bank dashboard",
				Key:    "reactkey",
				Role:   "client",
				Secret: "reactsecret",
			},
		}
		for _, internal := range internals {
			internal.Create()
		}

		// seed images
		file, _ := os.Open("migration/image_dummy.txt")
		defer file.Close()
		b64image, _ := ioutil.ReadAll(file)
		images := []models.Image{
			models.Image{
				Image_string: string(b64image),
			},
			models.Image{
				Image_string: string(b64image),
			},
			models.Image{
				Image_string: string(b64image),
			},
			models.Image{
				Image_string: string(b64image),
			},
			models.Image{
				Image_string: string(b64image),
			},
		}
		for _, image := range images {
			image.Create()
		}

		// seed bank types
		bankTypes := []models.BankType{
			models.BankType{
				Name:        "BPD",
				Description: "Description of BPD bank type",
			},
			models.BankType{
				Name:        "BPR",
				Description: "Description of BPR bank type",
			},
			models.BankType{
				Name:        "Koperasi",
				Description: "Description of Koperasi bank type",
			},
		}
		for _, bankType := range bankTypes {
			bankType.Create()
		}

		// seed services
		services := []models.Service{
			models.Service{
				Name:    "Pinjaman PNS",
				Status:  "active",
				ImageID: 1,
			},
			models.Service{
				Name:    "Pinjaman Pensiun",
				Status:  "active",
				ImageID: 2,
			},
			models.Service{
				Name:    "Pinjaman UMKN",
				Status:  "active",
				ImageID: 3,
			},
			models.Service{
				Name:    "Pinjaman Mikro",
				Status:  "inactive",
				ImageID: 4,
			},
			models.Service{
				Name:    "Pinjaman Lainnya",
				Status:  "inactive",
				ImageID: 5,
			},
		}
		for _, service := range services {
			service.Create()
		}

		// seed products
		feesMarshal, _ := json.Marshal([]interface{}{map[string]interface{}{
			"description": "Admin Fee",
			"amount":      "1%",
		}, map[string]interface{}{
			"description": "Convenience Fee",
			"amount":      "2%",
		}})
		products := []models.Product{
			models.Product{
				Name:            "Product A",
				ServiceID:       1,
				MinTimeSpan:     3,
				MaxTimeSpan:     12,
				Interest:        5,
				MinLoan:         5000000,
				MaxLoan:         8000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
			models.Product{
				Name:            "Product B",
				ServiceID:       2,
				MinTimeSpan:     3,
				MaxTimeSpan:     12,
				Interest:        5,
				MinLoan:         5000000,
				MaxLoan:         8000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
			models.Product{
				Name:            "Product C",
				ServiceID:       3,
				MinTimeSpan:     3,
				MaxTimeSpan:     12,
				Interest:        5,
				MinLoan:         5000000,
				MaxLoan:         8000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
			models.Product{
				Name:            "Product D",
				ServiceID:       4,
				MinTimeSpan:     3,
				MaxTimeSpan:     12,
				Interest:        5,
				MinLoan:         5000000,
				MaxLoan:         8000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
			models.Product{
				Name:            "Product E",
				ServiceID:       5,
				MinTimeSpan:     3,
				MaxTimeSpan:     12,
				Interest:        5,
				MinLoan:         5000000,
				MaxLoan:         8000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
		}
		for _, product := range products {
			product.Create()
		}

		// seed lenders
		lenders := []models.Bank{
			models.Bank{
				Name:                "Bank A",
				Type:                1,
				Address:             "Bank A Address",
				Province:            "Province A",
				City:                "City A",
				AdminFeeSetup:       "potong_plafon",
				ConvenienceFeeSetup: "potong_plafon",
				PIC:                 "Bank A PIC",
				Phone:               "081234567890",
				Services:            pq.Int64Array{1, 2},
				Products:            pq.Int64Array{1, 2},
				Username:            "Banktoib",
				Password:            "password",
			},
			models.Bank{
				Name:                "Bank B",
				Type:                2,
				Address:             "Bank B Address",
				Province:            "Province B",
				City:                "City B",
				AdminFeeSetup:       "potong_plafon",
				ConvenienceFeeSetup: "potong_plafon",
				PIC:                 "Bank B PIC",
				Phone:               "081234567891",
				Services:            pq.Int64Array{1, 2},
				Products:            pq.Int64Array{1, 2},
				Username:            "Banktoic",
				Password:            "password",
			},
		}
		for _, lender := range lenders {
			lender.Create()
		}

		roles := []models.Roles{
			models.Roles{
				Name:        "Core",
				Status:      true,
				Description: "ini Super Admin",
				System:      "Core",
			},
			models.Roles{
				Name:        "Manager",
				Status:      true,
				Description: "ini untuk Finance",
				System:      "Core",
			},
		}
		for _, role := range roles {
			role.Create()
		}

		permis := []models.Permissions{
			models.Permissions{
				RoleID:      2,
				Permissions: "Bank_List",
			},
			models.Permissions{
				RoleID:      2,
				Permissions: "Bank_Add",
			},
			models.Permissions{
				RoleID:      2,
				Permissions: "Bank_Edit",
			},
			models.Permissions{
				RoleID:      2,
				Permissions: "Role_List",
			},
			models.Permissions{
				RoleID:      2,
				Permissions: "Role_Add",
			},
			models.Permissions{
				RoleID:      2,
				Permissions: "Role_Edit",
			},
			models.Permissions{
				RoleID:      2,
				Permissions: "Permission_List",
			},
			models.Permissions{
				RoleID:      2,
				Permissions: "Permission_Add",
			},
			models.Permissions{
				RoleID:      2,
				Permissions: "Permission_Edit",
			},
			models.Permissions{
				RoleID:      1,
				Permissions: "All",
			},
		}
		for _, per := range permis {
			per.Create()
		}

		users := []models.User{
			models.User{
				RoleID:   1,
				Username: "adminkey",
				Password: "adminsecret",
			},
			models.User{
				RoleID:   2,
				Username: "manager",
				Password: "password",
			},
		}
		for _, user := range users {
			user.Create()
		}
	}
}

func TestSeed() {
	if asira.App.ENV == "development" {
		// seed internals
		internals := []models.Internals{
			models.Internals{
				Name:   "admin",
				Key:    "adminkey",
				Role:   "admin",
				Secret: "adminsecret",
			},
			models.Internals{
				Name:   "bank dashboard",
				Key:    "reactkey",
				Role:   "client",
				Secret: "reactsecret",
			},
		}
		for _, internal := range internals {
			internal.Create()
		}

		// seed images
		file, _ := os.Open("migration/image_dummy.txt")
		defer file.Close()
		b64image, _ := ioutil.ReadAll(file)
		images := []models.Image{
			models.Image{
				Image_string: string(b64image),
			},
			models.Image{
				Image_string: string(b64image),
			},
			models.Image{
				Image_string: string(b64image),
			},
			models.Image{
				Image_string: string(b64image),
			},
			models.Image{
				Image_string: string(b64image),
			},
		}
		for _, image := range images {
			image.Create()
		}

		// seed bank types
		bankTypes := []models.BankType{
			models.BankType{
				Name:        "BPD",
				Description: "Description of BPD bank type",
			},
			models.BankType{
				Name:        "BPR",
				Description: "Description of BPR bank type",
			},
			models.BankType{
				Name:        "Koperasi",
				Description: "Description of Koperasi bank type",
			},
		}
		for _, bankType := range bankTypes {
			bankType.Create()
		}

		// seed services
		services := []models.Service{
			models.Service{
				Name:    "Pinjaman PNS",
				Status:  "active",
				ImageID: 1,
			},
			models.Service{
				Name:    "Pinjaman Pensiun",
				Status:  "active",
				ImageID: 2,
			},
			models.Service{
				Name:    "Pinjaman UMKN",
				Status:  "active",
				ImageID: 3,
			},
			models.Service{
				Name:    "Pinjaman Mikro",
				Status:  "inactive",
				ImageID: 4,
			},
			models.Service{
				Name:    "Pinjaman Lainnya",
				Status:  "inactive",
				ImageID: 5,
			},
		}
		for _, service := range services {
			service.Create()
		}

		// seed products
		feesMarshal, _ := json.Marshal([]interface{}{map[string]interface{}{
			"description": "Admin Fee",
			"amount":      "1%",
		}, map[string]interface{}{
			"description": "Convenience Fee",
			"amount":      "2%",
		}})
		products := []models.Product{
			models.Product{
				Name:            "Product A",
				ServiceID:       1,
				MinTimeSpan:     3,
				MaxTimeSpan:     12,
				Interest:        5,
				MinLoan:         5000000,
				MaxLoan:         8000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
			models.Product{
				Name:            "Product B",
				ServiceID:       2,
				MinTimeSpan:     3,
				MaxTimeSpan:     12,
				Interest:        5,
				MinLoan:         5000000,
				MaxLoan:         8000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
			models.Product{
				Name:            "Product C",
				ServiceID:       3,
				MinTimeSpan:     3,
				MaxTimeSpan:     12,
				Interest:        5,
				MinLoan:         5000000,
				MaxLoan:         8000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
			models.Product{
				Name:            "Product D",
				ServiceID:       4,
				MinTimeSpan:     3,
				MaxTimeSpan:     12,
				Interest:        5,
				MinLoan:         5000000,
				MaxLoan:         8000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
			models.Product{
				Name:            "Product E",
				ServiceID:       5,
				MinTimeSpan:     3,
				MaxTimeSpan:     12,
				Interest:        5,
				MinLoan:         5000000,
				MaxLoan:         8000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
		}
		for _, product := range products {
			product.Create()
		}

		// seed lenders
		lenders := []models.Bank{
			models.Bank{
				Name:                "Bank A",
				Type:                1,
				Address:             "Bank A Address",
				Province:            "Province A",
				City:                "City A",
				AdminFeeSetup:       "potong_plafon",
				ConvenienceFeeSetup: "potong_plafon",
				PIC:                 "Bank A PIC",
				Phone:               "081234567890",
				Services:            pq.Int64Array{1, 2},
				Products:            pq.Int64Array{1, 2},
				Username:            "Banktoib",
				Password:            "password",
			},
			models.Bank{
				Name:                "Bank B",
				Type:                2,
				Address:             "Bank B Address",
				Province:            "Province B",
				City:                "City B",
				AdminFeeSetup:       "potong_plafon",
				ConvenienceFeeSetup: "potong_plafon",
				PIC:                 "Bank B PIC",
				Phone:               "081234567891",
				Services:            pq.Int64Array{1, 2},
				Products:            pq.Int64Array{1, 2},
				Username:            "Banktoic",
				Password:            "password",
			},
		}
		for _, lender := range lenders {
			lender.Create()
		}

		// @ToDo borrower and loans should be get from borrower platform
		// seed borrowers
		borrowers := []models.Borrower{
			models.Borrower{
				Fullname:             "Full Name A",
				Gender:               "M",
				IdCardNumber:         "9876123451234567789",
				TaxIDnumber:          "0987654321234567890",
				Email:                "emaila@domain.com",
				Birthday:             time.Now(),
				Birthplace:           "a birthplace",
				LastEducation:        "a last edu",
				MotherName:           "a mom",
				Phone:                "081234567890",
				MarriedStatus:        "single",
				SpouseName:           "a spouse",
				SpouseBirthday:       time.Now(),
				SpouseLastEducation:  "master",
				Dependants:           0,
				Address:              "a street address",
				Province:             "a province",
				City:                 "a city",
				NeighbourAssociation: "a rt",
				Hamlets:              "a rw",
				HomePhoneNumber:      "021837163",
				Subdistrict:          "a camat",
				UrbanVillage:         "a lurah",
				HomeOwnership:        "privately owned",
				LivedFor:             5,
				Occupation:           "accupation",
				EmployerName:         "amployer",
				EmployerAddress:      "amployer address",
				Department:           "a department",
				BeenWorkingFor:       2,
				DirectSuperior:       "a boss",
				EmployerNumber:       "02188776655",
				MonthlyIncome:        5000000,
				OtherIncome:          2000000,
				RelatedPersonName:    "a big sis",
				RelatedPhoneNumber:   "08987654321",
				BankAccountNumber:    "520384716",
				Bank: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
			},
			models.Borrower{
				Fullname:             "Full Name B",
				Gender:               "F",
				IdCardNumber:         "9876123451234567781",
				TaxIDnumber:          "0987654321234567891",
				Email:                "emailb@domain.com",
				Birthday:             time.Now(),
				Birthplace:           "b birthplace",
				LastEducation:        "b last edu",
				MotherName:           "b mom",
				Phone:                "081234567891",
				MarriedStatus:        "single",
				SpouseName:           "b spouse",
				SpouseBirthday:       time.Now(),
				SpouseLastEducation:  "master",
				Dependants:           0,
				Address:              "b street address",
				Province:             "b province",
				City:                 "b city",
				NeighbourAssociation: "b rt",
				Hamlets:              "b rw",
				HomePhoneNumber:      "021837163",
				Subdistrict:          "b camat",
				UrbanVillage:         "b lurah",
				HomeOwnership:        "privately owned",
				LivedFor:             5,
				Occupation:           "bccupation",
				EmployerName:         "bmployer",
				EmployerAddress:      "bmployer address",
				Department:           "b department",
				BeenWorkingFor:       2,
				DirectSuperior:       "b boss",
				EmployerNumber:       "02188776655",
				MonthlyIncome:        5000000,
				OtherIncome:          2000000,
				RelatedPersonName:    "b big sis",
				RelatedPhoneNumber:   "08987654321",
				RelatedAddress:       "big sis address",
				Bank: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
			},
		}
		for _, borrower := range borrowers {
			borrower.Create()
		}

		// seed loans
		loans := []models.Loan{
			models.Loan{
				Bank: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				Owner: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				OwnerName:        "Full Name A",
				LoanAmount:       5000000,
				Installment:      8,
				LoanIntention:    "a loan 1 intention",
				IntentionDetails: "a loan 1 intention details",
				Fees:             postgres.Jsonb{feesMarshal},
				Interest:         1.5,
				TotalLoan:        float64(6500000),
				LayawayPlan:      500000,
				Product:          1,
			},
			models.Loan{
				Bank: sql.NullInt64{
					Int64: 2,
					Valid: true,
				},
				Owner: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				OwnerName:        "Full Name B",
				LoanAmount:       2000000,
				Installment:      3,
				LoanIntention:    "a loan 1 intention",
				IntentionDetails: "a loan 1 intention details",
				Fees:             postgres.Jsonb{feesMarshal},
				Interest:         1.5,
				TotalLoan:        float64(3000000),
				LayawayPlan:      200000,
				Product:          1,
			},
			models.Loan{
				Bank: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				Owner: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				OwnerName:        "Full Name C",
				LoanAmount:       29000000,
				Installment:      3,
				LoanIntention:    "a loan 1 intention",
				IntentionDetails: "a loan 1 intention details",
				Fees:             postgres.Jsonb{feesMarshal},
				Interest:         1.5,
				TotalLoan:        float64(6500000),
				LayawayPlan:      500000,
				Product:          1,
			},
			models.Loan{
				Bank: sql.NullInt64{
					Int64: 2,
					Valid: true,
				},
				Owner: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				OwnerName:        "Full Name D",
				LoanAmount:       3000000,
				Installment:      3,
				LoanIntention:    "a loan 1 intention",
				IntentionDetails: "a loan 1 intention details",
				Fees:             postgres.Jsonb{feesMarshal},
				Interest:         1.5,
				TotalLoan:        float64(3000000),
				LayawayPlan:      200000,
				Product:          1,
			},
			models.Loan{
				Bank: sql.NullInt64{
					Int64: 2,
					Valid: true,
				},
				Owner: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				OwnerName:        "Full Name E",
				LoanAmount:       9123456,
				Installment:      3,
				LoanIntention:    "a loan 3 intention",
				IntentionDetails: "a loan 5 intention details",
				Fees:             postgres.Jsonb{feesMarshal},
				Interest:         1.5,
				TotalLoan:        float64(3000000),
				LayawayPlan:      200000,
				Product:          1,
			},
			models.Loan{
				Bank: sql.NullInt64{
					Int64: 2,
					Valid: true,
				},
				Owner: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				OwnerName:        "Full Name F",
				LoanAmount:       80123456,
				Installment:      11,
				LoanIntention:    "a loan 3 intention",
				IntentionDetails: "a loan 5 intention details",
				Fees:             postgres.Jsonb{feesMarshal},
				Interest:         1.5,
				TotalLoan:        float64(3000000),
				LayawayPlan:      200000,
				Product:          1,
			},
		}
		for _, loan := range loans {
			loan.Create()
		}

		roles := []models.Roles{
			models.Roles{
				Name:        "Core",
				Status:      true,
				Description: "ini Super Admin",
				System:      "Core",
			},
			models.Roles{
				Name:        "Manager",
				Status:      true,
				Description: "ini untuk Finance",
				System:      "Core",
			},
		}
		for _, role := range roles {
			role.Create()
		}

		permis := []models.Permissions{
			models.Permissions{
				RoleID:      2,
				Permissions: "Bank_List",
			},
			models.Permissions{
				RoleID:      2,
				Permissions: "Bank_Add",
			},
			models.Permissions{
				RoleID:      2,
				Permissions: "Bank_Edit",
			},
			models.Permissions{
				RoleID:      2,
				Permissions: "Role_List",
			},
			models.Permissions{
				RoleID:      2,
				Permissions: "Role_Add",
			},
			models.Permissions{
				RoleID:      2,
				Permissions: "Role_Edit",
			},
			models.Permissions{
				RoleID:      2,
				Permissions: "Permission_List",
			},
			models.Permissions{
				RoleID:      2,
				Permissions: "Permission_Add",
			},
			models.Permissions{
				RoleID:      2,
				Permissions: "Permission_Edit",
			},
			models.Permissions{
				RoleID:      1,
				Permissions: "All",
			},
		}
		for _, per := range permis {
			per.Create()
		}

		users := []models.User{
			models.User{
				RoleID:   1,
				Username: "adminkey",
				Password: "adminsecret",
			},
			models.User{
				RoleID:   2,
				Username: "manager",
				Password: "password",
			},
		}
		for _, user := range users {
			user.Create()
		}
	}
}

// truncate defined tables. []string{"all"} to truncate all tables.
func Truncate(tableList []string) (err error) {
	if len(tableList) > 0 {
		if tableList[0] == "all" {
			tableList = []string{
				"internals",
				"products",
				"services",
				"banks",
				"bank_types",
				"borrowers",
				"loans",
				"images",
				"roles",
				"permissions",
				"users",
				"user_relations",
			}
		}

		tables := strings.Join(tableList, ", ")
		sqlQuery := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tables)
		err = asira.App.DB.Exec(sqlQuery).Error
		return err
	}

	return fmt.Errorf("define tables that you want to truncate")
}
