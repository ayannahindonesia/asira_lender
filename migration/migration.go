package migration

import (
	"asira_lender/asira"
	"asira_lender/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

func Seed() {
	seeder := asira.App.DB.Begin()
	defer seeder.Commit()

	if asira.App.ENV == "development" {
		// seed lenders
		services := []int{1, 2, 3, 5, 8}
		jMarshal, _ := json.Marshal(services)
		lenders := []models.Bank{
			models.Bank{
				Name:     "Bank A",
				Type:     "BDP",
				Address:  "Bank A Address",
				Province: "Province A",
				City:     "City A",
				Services: postgres.Jsonb{jMarshal},
				PIC:      "Bank A PIC",
				Phone:    "081234567890",
				Username: "Banktoib",
				Password: "password",
			},
			models.Bank{
				Name:     "Bank B",
				Type:     "BDP",
				Address:  "Bank B Address",
				Province: "Province B",
				City:     "City B",
				Services: postgres.Jsonb{jMarshal},
				PIC:      "Bank B PIC",
				Phone:    "081234567891",
				Username: "Banktoic",
				Password: "password",
			},
		}
		for _, lender := range lenders {
			lender.Create()
		}

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
		fees := []models.LoanFee{
			models.LoanFee{
				Description: "fee 1",
				Amount:      1000,
			},
		}
		jMarshal, _ = json.Marshal(fees)
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
				Fees:             postgres.Jsonb{jMarshal},
				Interest:         1.5,
				TotalLoan:        float64(6500000),
				LayawayPlan:      500000,
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
				OwnerName:        "Full Name A",
				LoanAmount:       2000000,
				Installment:      3,
				LoanIntention:    "a loan 1 intention",
				IntentionDetails: "a loan 1 intention details",
				Fees:             postgres.Jsonb{jMarshal},
				Interest:         1.5,
				TotalLoan:        float64(3000000),
				LayawayPlan:      200000,
			},
		}
		for _, loan := range loans {
			loan.Create()
		}
	}
}

// truncate defined tables. []string{"all"} to truncate all tables.
func Truncate(tableList []string) (err error) {
	if len(tableList) > 0 {
		if tableList[0] == "all" {
			tableList = []string{
				"banks",
				"borrowers",
				"loans",
			}
		}

		tables := strings.Join(tableList, ", ")
		sqlQuery := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tables)
		err = asira.App.DB.Exec(sqlQuery).Error
		return err
	}

	return fmt.Errorf("define tables that you want to truncate")
}
