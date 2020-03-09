package migration

import (
	"asira_lender/asira"
	"asira_lender/middlewares"
	"asira_lender/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"

	"github.com/jinzhu/gorm/dialects/postgres"
)

// Seed func
func Seed() {
	if asira.App.ENV == "development" {
		// seed clients
		clients := []models.Client{
			models.Client{
				Name:   "admin",
				Key:    "adminkey",
				Secret: "adminsecret",
			},
			models.Client{
				Name:   "bank dashboard",
				Key:    "reactkey",
				Secret: "reactsecret",
			},
		}
		for _, client := range clients {
			client.Create()
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
			middlewares.SubmitKafkaPayload(bankType, "bank_type_create")
		}

		// seed services
		services := []models.Service{
			models.Service{
				Name:   "Pinjaman PNS",
				Status: "active",
				Image:  "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
			},
			models.Service{
				Name:   "Pinjaman Pensiun",
				Status: "active",
				Image:  "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
			},
			models.Service{
				Name:   "Pinjaman UMKN",
				Status: "active",
				Image:  "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
			},
			models.Service{
				Name:   "Pinjaman Mikro",
				Status: "inactive",
				Image:  "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
			},
			models.Service{
				Name:   "Pinjaman Lainnya",
				Status: "inactive",
				Image:  "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
			},
			models.Service{
				Name:   "Pinjaman Jeruk",
				Status: "active",
				Image:  "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
			},
			models.Service{
				Name:   "Pinjaman Pisang",
				Status: "active",
				Image:  "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
			},
			models.Service{
				Name:   "Pinjaman Mangga",
				Status: "active",
				Image:  "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
			},
		}
		for _, service := range services {
			service.Create()
			middlewares.SubmitKafkaPayload(service, "service_create")
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
			models.Product{
				Name:            "Product Jeruk Kecil",
				ServiceID:       6,
				MinTimeSpan:     6,
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
				Name:            "Product Jeruk Besar",
				ServiceID:       6,
				MinTimeSpan:     6,
				MaxTimeSpan:     24,
				Interest:        5,
				MinLoan:         5000000,
				MaxLoan:         20000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
			models.Product{
				Name:            "Product Pisang Kecil",
				ServiceID:       7,
				MinTimeSpan:     6,
				MaxTimeSpan:     24,
				Interest:        5,
				MinLoan:         5000000,
				MaxLoan:         10000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
			models.Product{
				Name:            "Product Pisang Raja",
				ServiceID:       7,
				MinTimeSpan:     6,
				MaxTimeSpan:     36,
				Interest:        7,
				MinLoan:         5000000,
				MaxLoan:         30000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
		}
		for _, product := range products {
			product.Create()
			middlewares.SubmitKafkaPayload(product, "product_create")
		}

		purposes := []models.LoanPurpose{
			models.LoanPurpose{
				Name:   "Lain-lain",
				Status: "active",
			},
			models.LoanPurpose{
				Name:   "Pendidikan",
				Status: "active",
			},
			models.LoanPurpose{
				Name:   "Rumah Tangga",
				Status: "active",
			},
			models.LoanPurpose{
				Name:   "Kesehatan",
				Status: "active",
			},
			models.LoanPurpose{
				Name:   "Berdagang",
				Status: "active",
			},
			models.LoanPurpose{
				Name:   "Bertani",
				Status: "active",
			},
			models.LoanPurpose{
				Name:   "Berjudi",
				Status: "inactive",
			},
		}
		for _, purpose := range purposes {
			purpose.Create()
			middlewares.SubmitKafkaPayload(purpose, "loan_purpose_create")
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
			},
			models.Bank{
				Name:                "Bank Buah",
				Type:                1,
				Address:             "jalan kaki cape sekali",
				Province:            "jambi",
				City:                "kota mati",
				AdminFeeSetup:       "potong_plafon",
				ConvenienceFeeSetup: "potong_plafon",
				PIC:                 "dindin",
				Phone:               "081234567891234",
				Services:            pq.Int64Array{6, 7, 8},
				Products:            pq.Int64Array{6, 7, 8, 9},
			},
		}
		for _, lender := range lenders {
			lender.Create()
			middlewares.SubmitKafkaPayload(lender, "bank_create")
		}

		roles := []models.Roles{
			models.Roles{
				Name:        "Administrator",
				Status:      "active",
				Description: "Super Admin",
				System:      "Core",
				Permissions: pq.StringArray{"all"},
			},
			models.Roles{
				Name:        "Ops",
				Status:      "active",
				Description: "Ops",
				System:      "Core",
				Permissions: pq.StringArray{"core_create_client", "core_view_image", "core_borrower_get_all", "core_borrower_get_details", "core_loan_get_all", "core_loan_get_details", "core_bank_type_list", "core_bank_type_new", "core_bank_type_detail", "core_bank_type_patch", "core_bank_list", "core_bank_new", "core_bank_detail", "core_bank_patch", "core_service_list", "core_service_new", "core_service_detail", "core_service_patch", "core_product_list", "core_product_new", "core_product_detail", "core_product_patch", "core_loan_purpose_list", "core_loan_purpose_new", "core_loan_purpose_detail", "core_loan_purpose_patch", "core_role_list", "core_role_details", "core_role_new", "core_role_patch", "core_role_range", "core_permission_list", "core_user_list", "core_user_details", "core_user_new", "core_user_patch", "convenience_fee_report"},
			},
			models.Roles{
				Name:        "Banker",
				Status:      "active",
				Description: "ini untuk Finance",
				System:      "Dashboard",
				Permissions: pq.StringArray{"lender_profile", "lender_profile_edit", "lender_loan_request_list", "lender_loan_request_detail", "lender_loan_approve_reject", "lender_loan_request_list_download", "lender_borrower_list", "lender_borrower_list_detail", "lender_borrower_list_download", "lender_prospective_borrower_approval", "lender_service_list", "lender_service_list_detail"},
			},
		}
		for _, role := range roles {
			role.Create()
		}

		users := []models.User{
			models.User{
				Roles:      pq.Int64Array{1},
				Username:   "adminkey",
				Password:   "adminsecret",
				Email:      "admin@ayannah.com",
				Phone:      "081234567890",
				Status:     "active",
				FirstLogin: false,
			},
			models.User{
				Roles:      pq.Int64Array{2},
				Username:   "viewer",
				Password:   "password",
				Email:      "asira@ayannah.com",
				Phone:      "081234567891",
				Status:     "active",
				FirstLogin: false,
			},
			models.User{
				Roles:      pq.Int64Array{3},
				Username:   "Banktoib",
				Email:      "toib@ayannah.com",
				Phone:      "081234567892",
				Password:   "password",
				Status:     "active",
				FirstLogin: false,
			},
			models.User{
				Roles:      pq.Int64Array{3},
				Username:   "Banktoic",
				Email:      "toic@ayannah.com",
				Phone:      "081234567893",
				Password:   "password",
				Status:     "active",
				FirstLogin: false,
			},
			models.User{
				Roles:    pq.Int64Array{3},
				Username: "Banktoid",
				Email:    "toid@ayannah.com",
				Phone:    "081234567894",
				Password: "password",
				Status:   "active",
			},
		}
		for _, user := range users {
			user.Create()
		}

		bankReps := []models.BankRepresentatives{
			models.BankRepresentatives{
				UserID: 3,
				BankID: 1,
			},
			models.BankRepresentatives{
				UserID: 4,
				BankID: 2,
			},
			models.BankRepresentatives{
				UserID: 5,
				BankID: 1,
			},
		}
		for _, bankRep := range bankReps {
			bankRep.Create()
		}

		agentProviders := []models.AgentProvider{
			models.AgentProvider{
				Name:    "Agent Provider A",
				PIC:     "PIC A",
				Phone:   "081234567890",
				Address: "address of provider a",
				Status:  "active",
			},
			models.AgentProvider{
				Name:    "Agent Provider B",
				PIC:     "PIC B",
				Phone:   "081234567891",
				Address: "address of provider b",
				Status:  "active",
			},
			models.AgentProvider{
				Name:    "Agent Provider C",
				PIC:     "PIC C",
				Phone:   "081234567892",
				Address: "address of provider c",
				Status:  "active",
			},
			models.AgentProvider{
				Name:    "Provider Buah",
				PIC:     "Dinand Buah",
				Phone:   "08123456789234",
				Address: "jalan jalan ke surabaya",
				Status:  "active",
			},
		}
		for _, agentProvider := range agentProviders {
			agentProvider.Create()
			middlewares.SubmitKafkaPayload(agentProvider, "agent_provider_create")
		}

		agents := []models.Agent{
			models.Agent{
				Name:     "Agent K",
				Username: "agentK",
				Password: "password",
				Email:    "agentk@mib.com",
				Phone:    "081234567890",
				Category: "agent",
				AgentProvider: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				Image:  "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
				Banks:  pq.Int64Array{1, 2},
				Status: "active",
			},
			models.Agent{
				Name:     "Agent J",
				Username: "agentJ",
				Password: "password",
				Email:    "agentj@mib.com",
				Phone:    "081234567891",
				Category: "account_executive",
				Image:    "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
				Banks:    pq.Int64Array{1},
				Status:   "active",
			},
			models.Agent{
				Name:     "Agent Buah Personal",
				Username: "agenbuahpsn",
				Password: "testing123",
				Email:    "agentbpsn@buah.com",
				Phone:    "08123456789112",
				Category: "agent",
				AgentProvider: sql.NullInt64{
					Int64: 4,
					Valid: true,
				},
				Banks:  pq.Int64Array{1, 3},
				Status: "active",
			},
			models.Agent{
				Name:     "Agent Buah Executive",
				Username: "agenbuahexe",
				Password: "testing123",
				Email:    "agentbexe@buah.com",
				Phone:    "08123456789122",
				Category: "account_executive",
				Banks:    pq.Int64Array{3},
				Status:   "active",
			},
		}
		for _, agent := range agents {
			agent.Create()
			middlewares.SubmitKafkaPayload(agent, "agent_create")
		}

		// Seeds for borrower
		// seed borrowers
		// borrowers := []models.Borrower{
		// 	models.Borrower{
		// 		Fullname:             "Full Name A",
		// 		Gender:               "M",
		// 		IdCardNumber:         "9876123451234567789",
		// 		TaxIDnumber:          "0987654321234567890",
		// 		Email:                "emaila@domain.com",
		// 		Birthday:             time.Now(),
		// 		Birthplace:           "a birthplace",
		// 		LastEducation:        "a last edu",
		// 		MotherName:           "a mom",
		// 		Phone:                "081234567890",
		// 		MarriedStatus:        "single",
		// 		SpouseName:           "a spouse",
		// 		SpouseBirthday:       time.Now(),
		// 		SpouseLastEducation:  "master",
		// 		Dependants:           0,
		// 		Address:              "a street address",
		// 		Province:             "a province",
		// 		City:                 "a city",
		// 		NeighbourAssociation: "a rt",
		// 		Hamlets:              "a rw",
		// 		HomePhoneNumber:      "021837163",
		// 		Subdistrict:          "a camat",
		// 		UrbanVillage:         "a lurah",
		// 		HomeOwnership:        "privately owned",
		// 		LivedFor:             5,
		// 		Occupation:           "accupation",
		// 		EmployerName:         "amployer",
		// 		EmployerAddress:      "amployer address",
		// 		Department:           "a department",
		// 		BeenWorkingFor:       2,
		// 		DirectSuperior:       "a boss",
		// 		EmployerNumber:       "02188776655",
		// 		MonthlyIncome:        5000000,
		// 		OtherIncome:          2000000,
		// 		RelatedPersonName:    "a big sis",
		// 		RelatedPhoneNumber:   "08987654321",
		// 		OTPverified:          true,
		// 		BankAccountNumber:    "520384716",
		// 		Bank: sql.NullInt64{
		// 			Int64: 1,
		// 			Valid: true,
		// 		},
		// 		AgentReferral: sql.NullInt64{
		// 			Int64: 0,
		// 			Valid: true,
		// 		},
		// 	},
		// 	models.Borrower{
		// 		Fullname:             "Full Name B",
		// 		Gender:               "F",
		// 		IdCardNumber:         "9876123451234567781",
		// 		TaxIDnumber:          "0987654321234567891",
		// 		Email:                "emailb@domain.com",
		// 		Birthday:             time.Now(),
		// 		Birthplace:           "b birthplace",
		// 		LastEducation:        "b last edu",
		// 		MotherName:           "b mom",
		// 		Phone:                "081234567891",
		// 		MarriedStatus:        "single",
		// 		SpouseName:           "b spouse",
		// 		SpouseBirthday:       time.Now(),
		// 		SpouseLastEducation:  "master",
		// 		Dependants:           0,
		// 		Address:              "b street address",
		// 		Province:             "b province",
		// 		City:                 "b city",
		// 		NeighbourAssociation: "b rt",
		// 		Hamlets:              "b rw",
		// 		HomePhoneNumber:      "021837163",
		// 		Subdistrict:          "b camat",
		// 		UrbanVillage:         "b lurah",
		// 		HomeOwnership:        "privately owned",
		// 		LivedFor:             5,
		// 		Occupation:           "bccupation",
		// 		EmployerName:         "bmployer",
		// 		EmployerAddress:      "bmployer address",
		// 		Department:           "b department",
		// 		BeenWorkingFor:       2,
		// 		DirectSuperior:       "b boss",
		// 		EmployerNumber:       "02188776655",
		// 		MonthlyIncome:        5000000,
		// 		OtherIncome:          2000000,
		// 		RelatedPersonName:    "b big sis",
		// 		RelatedPhoneNumber:   "08987654321",
		// 		RelatedAddress:       "big sis address",
		// 		OTPverified:          false,
		// 		Bank: sql.NullInt64{
		// 			Int64: 1,
		// 			Valid: true,
		// 		},
		// 		AgentReferral: sql.NullInt64{
		// 			Int64: 0,
		// 			Valid: true,
		// 		},
		// 	},
		// }
		// for _, borrower := range borrowers {
		// 	borrower.Create()
		// 	middlewares.SubmitKafkaPayload(borrower, "borrower_create")
		// }

		// loans := []models.Loan{
		// 	models.Loan{
		// 		Borrower:         1,
		// 		LoanAmount:       1000000,
		// 		Installment:      6,
		// 		LoanIntention:    "Pendidikan",
		// 		IntentionDetails: "a loan 1 intention details",
		// 		Product:          1,
		// 	},
		// 	models.Loan{
		// 		Borrower:         1,
		// 		Status:           "approved",
		// 		LoanAmount:       500000,
		// 		Installment:      2,
		// 		LoanIntention:    "Rumah Tangga",
		// 		IntentionDetails: "a loan 2 intention details",
		// 		Product:          1,
		// 		OTPverified:      true,
		// 	},
		// 	models.Loan{
		// 		Borrower:         1,
		// 		Status:           "rejected",
		// 		LoanAmount:       2000000,
		// 		Installment:      8,
		// 		LoanIntention:    "Kesehatan",
		// 		IntentionDetails: "a loan 3 intention details",
		// 		Product:          1,
		// 		OTPverified:      true,
		// 	},
		// }
		// for _, loan := range loans {
		// 	loan.Create()
		// 	middlewares.SubmitKafkaPayload(loan, "loan_create")
		// }

		faqs := []models.FAQ{
			models.FAQ{
				Title: "How to register",
				Description: `
				<html>
				<head>
				</head>
				<body>
				<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
				<strong class="accordion">Section 1</strong>
				<p>Lorem ipsum...1</p>
				</div>
				<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
				<strong class="accordion">Section 2</strong>
				<p>Lorem ipsum...2</p>
				</div>
				<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
				<strong class="accordion">Section 3</strong>
				<p>Lorem ipsum...3</p>
				</div>
				</body>
				</html>`,
			},
			models.FAQ{
				Title: "How to applying loan",
				Description: `
				<html>
				<head>
				</head>
				<body>
				<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
				<strong class="accordion">Section 1</strong>
				<p>Lorem ipsum...1</p>
				</div>
				<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
				<strong class="accordion">Section 2</strong>
				<p>Lorem ipsum...2</p>
				</div>
				<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
				<strong class="accordion">Section 3</strong>
				<p>Lorem ipsum...3</p>
				</div>
				</body>
				</html>`,
			},
		}
		for _, faq := range faqs {
			faq.Create()
		}
	}
}

// TestSeed func
func TestSeed() {
	if asira.App.ENV == "development" {
		// seed clients
		clients := []models.Client{
			models.Client{
				Name:   "admin",
				Key:    "adminkey",
				Secret: "adminsecret",
			},
			models.Client{
				Name:   "bank dashboard",
				Key:    "reactkey",
				Secret: "reactsecret",
			},
		}
		for _, client := range clients {
			client.Create()
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
				Name:   "Pinjaman PNS",
				Status: "active",
				Image:  "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
			},
			models.Service{
				Name:   "Pinjaman Pensiun",
				Status: "active",
				Image:  "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
			},
			models.Service{
				Name:   "Pinjaman UMKN",
				Status: "active",
				Image:  "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
			},
			models.Service{
				Name:   "Pinjaman Mikro",
				Status: "inactive",
				Image:  "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
			},
			models.Service{
				Name:   "Pinjaman Lainnya",
				Status: "inactive",
				Image:  "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
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
				Image:               "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
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
			},
			models.Bank{
				Name:                "Bank B",
				Image:               "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
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
			},
			models.Bank{
				Name:                "Bank ",
				Image:               "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
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

		purposes := []models.LoanPurpose{
			models.LoanPurpose{
				Name:   "Lain-lain",
				Status: "active",
			},
			models.LoanPurpose{
				Name:   "Pendidikan",
				Status: "active",
			},
			models.LoanPurpose{
				Name:   "Rumah Tangga",
				Status: "active",
			},
			models.LoanPurpose{
				Name:   "Kesehatan",
				Status: "active",
			},
			models.LoanPurpose{
				Name:   "Berdagang",
				Status: "active",
			},
			models.LoanPurpose{
				Name:   "Bertani",
				Status: "active",
			},
			models.LoanPurpose{
				Name:   "Berjudi",
				Status: "inactive",
			},
		}
		for _, purpose := range purposes {
			purpose.Create()
		}

		// seed loans
		feesMarshal, _ = json.Marshal([]interface{}{map[string]interface{}{
			"description": "Admin Fee",
			"amount":      "10000",
		}, map[string]interface{}{
			"description": "Convenience Fee",
			"amount":      "50000",
		}})
		loans := []models.Loan{
			models.Loan{
				Borrower:         1,
				LoanAmount:       5000000,
				Installment:      8,
				LoanIntention:    "a loan 1 intention",
				IntentionDetails: "a loan 1 intention details",
				Fees:             postgres.Jsonb{feesMarshal},
				Interest:         1.5,
				TotalLoan:        float64(6500000),
				LayawayPlan:      500000,
				Product:          1,
				OTPverified:      true,
			},
			models.Loan{
				Borrower:         1,
				LoanAmount:       2000000,
				Installment:      3,
				LoanIntention:    "a loan 1 intention",
				IntentionDetails: "a loan 1 intention details",
				Fees:             postgres.Jsonb{feesMarshal},
				Interest:         1.5,
				TotalLoan:        float64(3000000),
				LayawayPlan:      200000,
				Product:          1,
				OTPverified:      true,
			},
			models.Loan{
				Borrower:         1,
				LoanAmount:       29000000,
				Installment:      3,
				LoanIntention:    "a loan 1 intention",
				IntentionDetails: "a loan 1 intention details",
				Fees:             postgres.Jsonb{feesMarshal},
				Interest:         1.5,
				TotalLoan:        float64(6500000),
				LayawayPlan:      500000,
				Product:          1,
				OTPverified:      true,
			},
			models.Loan{
				Borrower:         1,
				LoanAmount:       3000000,
				Installment:      3,
				LoanIntention:    "a loan 1 intention",
				IntentionDetails: "a loan 1 intention details",
				Fees:             postgres.Jsonb{feesMarshal},
				Interest:         1.5,
				TotalLoan:        float64(3000000),
				LayawayPlan:      200000,
				Product:          1,
				OTPverified:      true,
			},
			models.Loan{
				Borrower:         1,
				LoanAmount:       9123456,
				Installment:      3,
				LoanIntention:    "a loan 3 intention",
				IntentionDetails: "a loan 5 intention details",
				Fees:             postgres.Jsonb{feesMarshal},
				Interest:         1.5,
				TotalLoan:        float64(3000000),
				LayawayPlan:      200000,
				Product:          1,
				OTPverified:      true,
			},
			models.Loan{
				Borrower:         1,
				LoanAmount:       80123456,
				Installment:      11,
				LoanIntention:    "a loan 3 intention",
				IntentionDetails: "a loan 5 intention details",
				Fees:             postgres.Jsonb{feesMarshal},
				Interest:         1.5,
				TotalLoan:        float64(3000000),
				LayawayPlan:      200000,
				Product:          1,
				OTPverified:      true,
			},
		}
		for _, loan := range loans {
			loan.Create()
		}

		roles := []models.Roles{
			models.Roles{
				Name:        "Administrator",
				Status:      "active",
				Description: "Super Admin",
				System:      "Core",
				Permissions: pq.StringArray{"all"},
			},
			models.Roles{
				Name:        "Ops",
				Status:      "active",
				Description: "Ops",
				System:      "Core",
				Permissions: pq.StringArray{"core_create_client", "core_view_image", "core_borrower_get_all", "core_borrower_get_details", "core_loan_get_all", "core_loan_get_details", "core_bank_type_list", "core_bank_type_new", "core_bank_type_detail", "core_bank_type_patch", "core_bank_list", "core_bank_new", "core_bank_detail", "core_bank_patch", "core_service_list", "core_service_new", "core_service_detail", "core_service_patch", "core_product_list", "core_product_new", "core_product_detail", "core_product_patch", "core_loan_purpose_list", "core_loan_purpose_new", "core_loan_purpose_detail", "core_loan_purpose_patch", "core_role_list", "core_role_details", "core_role_new", "core_role_patch", "core_role_range", "core_permission_list", "core_user_list", "core_user_details", "core_user_new", "core_user_patch", "convenience_fee_report"},
			},
			models.Roles{
				Name:        "Banker",
				Status:      "active",
				Description: "ini untuk Finance",
				System:      "Dashboard",
				Permissions: pq.StringArray{"lender_profile", "lender_profile_edit", "lender_loan_request_list", "lender_loan_request_detail", "lender_loan_approve_reject", "lender_loan_request_list_download", "lender_borrower_list", "lender_borrower_list_detail", "lender_borrower_list_download", "lender_prospective_borrower_approval", "lender_service_list", "lender_service_list_detail"},
			},
		}
		for _, role := range roles {
			role.Create()
		}

		users := []models.User{
			models.User{
				Roles:      pq.Int64Array{1},
				Username:   "adminkey",
				Password:   "adminsecret",
				Email:      "admin@ayannah.com",
				Phone:      "081234567890",
				Status:     "active",
				FirstLogin: false,
			},
			models.User{
				Roles:      pq.Int64Array{2},
				Username:   "viewer",
				Password:   "password",
				Email:      "asira@ayannah.com",
				Phone:      "081234567891",
				Status:     "active",
				FirstLogin: false,
			},
			models.User{
				Roles:      pq.Int64Array{3},
				Username:   "Banktoib",
				Email:      "toib@ayannah.com",
				Phone:      "081234567892",
				Password:   "password",
				Status:     "active",
				FirstLogin: false,
			},
			models.User{
				Roles:      pq.Int64Array{3},
				Username:   "Banktoic",
				Email:      "toic@ayannah.com",
				Phone:      "081234567893",
				Password:   "password",
				Status:     "active",
				FirstLogin: false,
			},
			models.User{
				Roles:    pq.Int64Array{3},
				Username: "Banktoid",
				Email:    "toid@ayannah.com",
				Phone:    "081234567894",
				Password: "password",
				Status:   "active",
			},
		}
		for _, user := range users {
			user.Create()
		}

		bankReps := []models.BankRepresentatives{
			models.BankRepresentatives{
				UserID: 3,
				BankID: 1,
			},
			models.BankRepresentatives{
				UserID: 4,
				BankID: 2,
			},
			models.BankRepresentatives{
				UserID: 5,
				BankID: 1,
			},
		}
		for _, bankRep := range bankReps {
			bankRep.Create()
		}

		agentProviders := []models.AgentProvider{
			models.AgentProvider{
				Name:    "Agent Provider A",
				PIC:     "PIC A",
				Phone:   "081234567890",
				Address: "address of provider a",
				Status:  "active",
			},
			models.AgentProvider{
				Name:    "Agent Provider B",
				PIC:     "PIC B",
				Phone:   "081234567891",
				Address: "address of provider b",
				Status:  "active",
			},
			models.AgentProvider{
				Name:    "Agent Provider C",
				PIC:     "PIC C",
				Phone:   "081234567892",
				Address: "address of provider c",
				Status:  "active",
			},
		}
		for _, agentProvider := range agentProviders {
			agentProvider.Create()
		}

		agents := []models.Agent{
			models.Agent{
				Name:     "Agent K",
				Username: "agentK",
				Password: "password",
				Email:    "agentk@mib.com",
				Phone:    "081234567890",
				Category: "agent",
				AgentProvider: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				Image:  "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
				Banks:  pq.Int64Array{1, 2},
				Status: "active",
			},
			models.Agent{
				Name:     "Agent J",
				Username: "agentJ",
				Password: "password",
				Email:    "agentj@mib.com",
				Phone:    "081234567891",
				Category: "account_executive",
				Banks:    pq.Int64Array{1},
				Image:    "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
				Status:   "active",
			},
		}
		for _, agent := range agents {
			agent.Create()
		}

		faqs := []models.FAQ{
			models.FAQ{
				Title: "How to register",
				Description: `
				<html>
				<head>
				</head>
				<body>
				<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
				<strong class="accordion">Section 1</strong>
				<p>Lorem ipsum...1</p>
				</div>
				<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
				<strong class="accordion">Section 2</strong>
				<p>Lorem ipsum...2</p>
				</div>
				<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
				<strong class="accordion">Section 3</strong>
				<p>Lorem ipsum...3</p>
				</div>
				</body>
				</html>`,
			},
			models.FAQ{
				Title: "How to applying loan",
				Description: `
				<html>
				<head>
				</head>
				<body>
				<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
				<strong class="accordion">Section 1</strong>
				<p>Lorem ipsum...1</p>
				</div>
				<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
				<strong class="accordion">Section 2</strong>
				<p>Lorem ipsum...2</p>
				</div>
				<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
				<strong class="accordion">Section 3</strong>
				<p>Lorem ipsum...3</p>
				</div>
				</body>
				</html>`,
			},
		}
		for _, faq := range faqs {
			faq.Create()
		}

	}
}

// Truncate defined tables. []string{"all"} to truncate all tables.
func Truncate(tableList []string) (err error) {
	if len(tableList) > 0 {
		if tableList[0] == "all" {
			tableList = []string{
				"clients",
				"products",
				"services",
				"banks",
				"bank_types",
				"borrowers",
				"loans",
				"roles",
				"users",
				"bank_representatives",
				"loan_purposes",
				"agent_providers",
				"faqs",
			}
		}

		tables := strings.Join(tableList, ", ")
		sqlQuery := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tables)
		err = asira.App.DB.Exec(sqlQuery).Error
		return err
	}

	return fmt.Errorf("define tables that you want to truncate")
}
