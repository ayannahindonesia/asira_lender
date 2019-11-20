package middlewares

import (
	"asira_lender/asira"
	"asira_lender/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

type (
	// AsiraKafkaHandlers type
	AsiraKafkaHandlers struct {
		KafkaConsumer     sarama.Consumer
		PartitionConsumer sarama.PartitionConsumer
	}
	// BorrowerInfo borrower info passed to lender
	BorrowerInfo struct {
		Info interface{} `json:"borrower_info"`
	}
)

var wg sync.WaitGroup

func init() {
	var err error
	topics := asira.App.Config.GetStringMap(fmt.Sprintf("%s.kafka.topics.consumes", asira.App.ENV))

	kafka := &AsiraKafkaHandlers{}
	kafka.KafkaConsumer, err = sarama.NewConsumer([]string{asira.App.Kafka.Host}, asira.App.Kafka.Config)
	if err != nil {
		log.Printf("error while creating new kafka consumer : %v", err)
	}

	kafka.SetPartitionConsumer(topics["for_lender"].(string))

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer kafka.KafkaConsumer.Close()
		for {
			message, err := kafka.Listen()
			if err != nil {
				log.Printf("error occured when listening kafka : %v", err)
			}
			if message != nil {
				err := processMessage(message)
				if err != nil {
					log.Printf("%v . message : %v", err, string(message))
				}
			}
		}
	}()
}

// SetPartitionConsumer func
func (k *AsiraKafkaHandlers) SetPartitionConsumer(topic string) (err error) {
	k.PartitionConsumer, err = k.KafkaConsumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		return err
	}

	return nil
}

// Listen to kafka
func (k *AsiraKafkaHandlers) Listen() ([]byte, error) {
	select {
	case err := <-k.PartitionConsumer.Errors():
		return nil, err
	case msg := <-k.PartitionConsumer.Messages():
		return msg.Value, nil
	}
}

func processMessage(kafkaMessage []byte) (err error) {
	data := strings.SplitN(string(kafkaMessage), ":", 2)
	switch data[0] {
	default:
		return nil
	case "loan":
		// create borrower first
		var borrowerInfo BorrowerInfo
		err = json.Unmarshal([]byte(data[1]), &borrowerInfo)
		if err != nil {
			return err
		}

		marshal, err := json.Marshal(borrowerInfo.Info)
		if err != nil {
			return err
		}

		var borrower models.Borrower
		err = json.Unmarshal(marshal, &borrower)
		if err != nil {
			return err
		}

		err = borrower.FirstOrCreate() // finish borrower create
		if err != nil {
			return err
		}

		// create loan
		var loan models.Loan
		err = json.Unmarshal([]byte(data[1]), &loan)
		if err != nil {
			return err
		}

		loan.Bank = borrower.Bank
		loan.OwnerName = borrower.Fullname

		err = loan.Save() // finish create loan
		break
	case "agent_borrower":
		type AgentBorrowerContainer struct {
			Fullname             string        `json:"fullname"`
			Nickname             string        `json:"nickname"`
			Gender               string        `json:"gender"`
			IDCardNumber         string        `json:"idcard_number"`
			IDCardImage          sql.NullInt64 `json:"idcard_image"`
			TaxIDnumber          string        `json:"taxid_number"`
			TaxIDImage           sql.NullInt64 `json:"taxid_image"`
			Nationality          string        `json:"nationality"`
			Email                string        `json:"email"`
			Birthday             time.Time     `json:"birthday"`
			Birthplace           string        `json:"birthplace"`
			LastEducation        string        `json:"last_education"`
			MotherName           string        `json:"mother_name"`
			Phone                string        `json:"phone"`
			MarriedStatus        string        `json:"marriage_status"`
			SpouseName           string        `json:"spouse_name"`
			SpouseBirthday       time.Time     `json:"spouse_birthday"`
			SpouseLastEducation  string        `json:"spouse_lasteducation"`
			Dependants           int           `json:"dependants"`
			Address              string        `json:"address"`
			Province             string        `json:"province"`
			City                 string        `json:"city"`
			NeighbourAssociation string        `json:"neighbour_association"`
			Hamlets              string        `json:"hamlets"`
			HomePhoneNumber      string        `json:"home_phonenumber"`
			Subdistrict          string        `json:"subdistrict"`
			UrbanVillage         string        `json:"urban_village"`
			HomeOwnership        string        `json:"home_ownership" `
			LivedFor             int           `json:"lived_for"`
			Occupation           string        `json:"occupation"`
			EmployeeID           string        `json:"employee_id"`
			EmployerName         string        `json:"employer_name"`
			EmployerAddress      string        `json:"employer_address"`
			Department           string        `json:"department"`
			BeenWorkingFor       int           `json:"been_workingfor"`
			DirectSuperior       string        `json:"direct_superiorname"`
			EmployerNumber       string        `json:"employer_number"`
			MonthlyIncome        int           `json:"monthly_income"`
			OtherIncome          int           `json:"other_income"`
			OtherIncomeSource    string        `json:"other_incomesource"`
			FieldOfWork          string        `json:"field_of_work"`
			RelatedPersonName    string        `json:"related_personname"`
			RelatedRelation      string        `json:"related_relation"`
			RelatedPhoneNumber   string        `json:"related_phonenumber"`
			RelatedHomePhone     string        `json:"related_homenumber"`
			RelatedAddress       string        `json:"related_address"`
			Bank                 sql.NullInt64 `json:"bank" gorm:"column:bank"`
			BankAccountNumber    string        `json:"bank_accountnumber"`
			AgentID              int64         `json:"agent_id"`
		}

		type SearchID struct {
			IDCardNumber string `json:"idcard_number"`
		}

		var (
			agentContainer AgentBorrowerContainer
			borrower       models.Borrower
		)

		err = json.Unmarshal([]byte(data[1]), &agentContainer)
		if err != nil {
			return err
		}
		marshal, _ := json.Marshal(agentContainer)

		err := borrower.FilterSearchSingle(&SearchID{
			IDCardNumber: agentContainer.IDCardNumber,
		})
		json.Unmarshal(marshal, &borrower)
		if err != nil {
			borrower = models.Borrower{}
			err = borrower.Create()
		} else {
			err = borrower.Save()
		}
		break
	}
	return err
}
