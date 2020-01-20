package models

import (
	"asira_lender/asira"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

func KafkaSubmitModel(i interface{}, model string) (err error) {
	// skip kafka submit when in unit testing
	if flag.Lookup("test.v") != nil {
		return nil
	}

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
			ID                  uint64    `json:"id"`
			Status              string    `json:"status"`
			DueDate             time.Time `json:"due_date"`
			DisburseDate        time.Time `json:"disburse_date"`
			DisburseStatus      string    `json:"disburse_status"`
			DisburseDateChanged bool      `json:"disburse_date_changed"`
			RejectReason        string    `json:"reject_reason"`
		}
		if e, ok := i.(*Loan); ok {
			payload = LoanStatusUpdate{
				ID:                  e.ID,
				Status:              e.Status,
				DueDate:             e.DueDate,
				DisburseDate:        e.DisburseDate,
				DisburseStatus:      e.DisburseStatus,
				DisburseDateChanged: e.DisburseDateChanged,
				RejectReason:        e.RejectReason,
			}
		}
		break
	}

	return payload
}
