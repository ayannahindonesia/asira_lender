package migration

import (
	"asira_lender/asira"
	"asira_lender/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jinzhu/gorm/dialects/postgres"
)

func Seed() {
	seeder := asira.App.DB.Begin()
	defer seeder.Commit()

	if asira.App.ENV == "development" {
		// seed images
		file, _ := os.Open("migration/image_dummy.txt")
		defer file.Close()
		b64image, _ := ioutil.ReadAll(file)
		images := []models.Image{
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

		// seed bank services
		bankServices := []models.BankService{
			models.BankService{
				Name:    "Pinjaman PNS",
				ImageID: 1,
				Status:  "active",
			},
			models.BankService{
				Name:    "Pinjaman Pensiun",
				ImageID: 1,
				Status:  "active",
			},
			models.BankService{
				Name:    "Pinjaman UMKN",
				ImageID: 1,
				Status:  "active",
			},
			models.BankService{
				Name:    "Pinjaman Mikro",
				ImageID: 1,
				Status:  "inactive",
			},
			models.BankService{
				Name:    "Pinjaman Lainnya",
				ImageID: 1,
				Status:  "inactive",
			},
		}
		for _, bankService := range bankServices {
			bankService.Create()
		}

		// seed service products
		feesMarshal, _ := json.Marshal([]interface{}{map[string]interface{}{
			"description": "Admin Fee",
			"amount":      2500,
		}})
		collateralMarshal, _ := json.Marshal([]string{"Surat Tanah", "BPKB"})
		financeMarshal, _ := json.Marshal([]string{"Pendidikan"})
		serviceProducts := []models.ServiceProduct{
			models.ServiceProduct{
				Name:            "Product A",
				MinTimeSpan:     1,
				MaxTimeSpan:     6,
				Interest:        5,
				MinLoan:         1000000,
				MaxLoan:         10000000,
				Fees:            postgres.Jsonb{feesMarshal},
				ASN_Fee:         "1%",
				Service:         1,
				Collaterals:     postgres.Jsonb{collateralMarshal},
				FinancingSector: postgres.Jsonb{financeMarshal},
				Assurance:       "an Assurance",
				Status:          "active",
			},
			models.ServiceProduct{
				Name:            "Product B",
				MinTimeSpan:     3,
				MaxTimeSpan:     12,
				Interest:        5,
				MinLoan:         5000000,
				MaxLoan:         8000000,
				Fees:            postgres.Jsonb{feesMarshal},
				ASN_Fee:         "1%",
				Service:         1,
				Collaterals:     postgres.Jsonb{collateralMarshal},
				FinancingSector: postgres.Jsonb{financeMarshal},
				Assurance:       "an Assurance",
				Status:          "active",
			},
		}
		for _, serviceProduct := range serviceProducts {
			serviceProduct.Create()
		}

		// seed lenders
		rawBankServices := []string{"Pinjaman PNS", "Pinjaman Pensiun", "Pinjaman Mikro"}
		rawBankProducts := []string{"Product A", "Product B"}
		jBankServices, _ := json.Marshal(rawBankServices)
		jBankProducts, _ := json.Marshal(rawBankProducts)
		lenders := []models.Bank{
			models.Bank{
				Name:     "Bank A",
				Type:     1,
				Address:  "Bank A Address",
				Province: "Province A",
				City:     "City A",
				Services: postgres.Jsonb{jBankServices},
				Products: postgres.Jsonb{jBankProducts},
				PIC:      "Bank A PIC",
				Phone:    "081234567890",
				Username: "Banktoib",
				Password: "password",
			},
			models.Bank{
				Name:     "Bank B",
				Type:     2,
				Address:  "Bank B Address",
				Province: "Province B",
				City:     "City B",
				Services: postgres.Jsonb{jBankServices},
				Products: postgres.Jsonb{jBankProducts},
				PIC:      "Bank B PIC",
				Phone:    "081234567891",
				Username: "Banktoic",
				Password: "password",
			},
		}
		for _, lender := range lenders {
			lender.Create()
		}
	}
}

// truncate defined tables. []string{"all"} to truncate all tables.
func Truncate(tableList []string) (err error) {
	if len(tableList) > 0 {
		if tableList[0] == "all" {
			tableList = []string{
				"service_products",
				"bank_services",
				"bank_types",
				"banks",
				"borrowers",
				"loans",
				"images",
			}
		}

		tables := strings.Join(tableList, ", ")
		sqlQuery := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tables)
		err = asira.App.DB.Exec(sqlQuery).Error
		return err
	}

	return fmt.Errorf("define tables that you want to truncate")
}
