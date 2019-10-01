package models

import (
	"asira_lender/asira"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

func KafkaSubmitModel(i interface{}, model string) (err error) {
	topics := asira.App.Config.GetStringMap(fmt.Sprintf("%s.kafka.topics.produces", asira.App.ENV))

	var payload interface{}
	payload = kafkaPayloadBuilder(i, model)

	jMarshal, _ := json.Marshal(payload)

	kafkaProducer, err := sarama.NewAsyncProducer([]string{asira.App.Kafka.Host}, asira.App.Kafka.Config)
	if err != nil {
		return err
	}
	defer kafkaProducer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topics["for_borrower"].(string),
		Value: sarama.StringEncoder(strings.TrimSuffix(model, "_delete") + ":" + string(jMarshal)),
	}

	select {
	case kafkaProducer.Input() <- msg:
		log.Printf("Produced topic : %s", topics["for_borrower"].(string))
	case err := <-kafkaProducer.Errors():
		log.Printf("Fail producing topic : %s error : %v", topics["for_borrower"].(string), err)
	}

	return nil
}

func kafkaPayloadBuilder(i interface{}, model string) (payload interface{}) {
	switch model {
	default:
		if strings.HasSuffix(model, "_delete") {
			type ModelDelete struct {
				ID     float64 `json:"id"`
				Model  string  `json:"model"`
				Delete bool    `json:"delete"`
			}
			var inInterface map[string]interface{}
			inrec, _ := json.Marshal(i)
			json.Unmarshal(inrec, &inInterface)
			if modelID, ok := inInterface["id"].(float64); ok {
				payload = ModelDelete{
					ID:     modelID,
					Model:  strings.TrimSuffix(model, "_delete"),
					Delete: true,
				}
			}
		} else {
			payload = i
		}
		break
	case "loan":
		type LoanStatusUpdate struct {
			ID           uint64    `json:"id"`
			Status       string    `json:"status"`
			DisburseDate time.Time `json:"disburse_date"`
		}
		if e, ok := i.(*Loan); ok {
			payload = LoanStatusUpdate{
				ID:           e.ID,
				Status:       e.Status,
				DisburseDate: e.DisburseDate,
			}
		}
		break
	case "bank_service":
		type BankServiceUpdate struct {
			ID      uint64 `json:"id"`
			Name    string `json:"name"`
			BankID  uint64 `json:"bank_id"`
			ImageID int    `json:"image_id"`
			Status  string `json:"status"`
		}
		if e, ok := i.(*BankService); ok {
			service := Service{}
			service.FindbyID(int(e.BankID))
			payload = BankServiceUpdate{
				ID:      e.ID,
				Name:    service.Name,
				BankID:  e.BankID,
				ImageID: e.ImageID,
				Status:  e.Status,
			}
		}
		break
	case "bank_product":
		type BankProductUpdate struct {
			ID              uint64         `json:"id"`
			Name            string         `json:"name"`
			BankServiceID   uint64         `json:"bank_service_id"`
			MinTimeSpan     int            `json:"min_timespan"`
			MaxTimeSpan     int            `json:"max_timespan"`
			Interest        float64        `json:"interest"`
			MinLoan         int            `json:"min_loan"`
			MaxLoan         int            `json:"max_loan"`
			Fees            postgres.Jsonb `json:"fees"`
			Collaterals     pq.StringArray `json:"collaterals"`
			FinancingSector pq.StringArray `json:"financing_sector"`
			Assurance       string         `json:"assurance"`
			Status          string         `json:"status"`
		}
		if e, ok := i.(*BankProduct); ok {
			product := Product{}
			product.FindbyID(int(e.ProductID))
			payload = BankProductUpdate{
				ID:              e.ID,
				Name:            product.Name,
				BankServiceID:   e.BankServiceID,
				MinTimeSpan:     e.MinTimeSpan,
				MaxTimeSpan:     e.MaxTimeSpan,
				Interest:        e.Interest,
				MinLoan:         e.MinLoan,
				MaxLoan:         e.MaxLoan,
				Fees:            e.Fees,
				Collaterals:     e.Collaterals,
				FinancingSector: e.FinancingSector,
				Assurance:       e.Assurance,
				Status:          e.Status,
			}
		}
	}

	return payload
}
