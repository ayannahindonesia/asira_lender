package middlewares

import (
	"asira_lender/asira"
	"asira_lender/models"
	"encoding/json"
	"flag"
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
	// Filter for searching basemodel
	Filter struct {
		IDCardNumber string `json:"idcard_number"`
	}
)

var wg sync.WaitGroup

func init() {
	var err error
	topic := asira.App.Config.GetString(fmt.Sprintf("%s.kafka.topics.consumes", asira.App.ENV))

	kafka := &AsiraKafkaHandlers{}
	kafka.KafkaConsumer, err = sarama.NewConsumer([]string{asira.App.Kafka.Host}, asira.App.Kafka.Config)
	if err != nil {
		log.Printf("error while creating new kafka consumer : %v", err)
	}

	kafka.SetPartitionConsumer(topic)

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

// SubmitKafkaPayload submits payload to kafka
func SubmitKafkaPayload(i interface{}, model string) (err error) {
	// skip kafka submit when in unit testing
	if flag.Lookup("test.v") != nil {
		return nil
	}

	topic := asira.App.Config.GetString(fmt.Sprintf("%s.kafka.topics.produces", asira.App.ENV))

	var payload interface{}
	payload = kafkaPayloadBuilder(i, model)

	jMarshal, _ := json.Marshal(payload)

	kafkaProducer, err := sarama.NewAsyncProducer([]string{asira.App.Kafka.Host}, asira.App.Kafka.Config)
	if err != nil {
		return err
	}
	defer kafkaProducer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(strings.TrimSuffix(model, "_delete") + ":" + string(jMarshal)),
	}

	select {
	case kafkaProducer.Input() <- msg:
		log.Printf("Produced topic : %s", topic)
	case err := <-kafkaProducer.Errors():
		log.Printf("Fail producing topic : %s error : %v", topic, err)
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
			DisburseDate        time.Time `json:"disburse_date"`
			DisburseStatus      string    `json:"disburse_status"`
			DisburseDateChanged bool      `json:"disburse_date_changed"`
			RejectReason        string    `json:"reject_reason"`
		}
		if e, ok := i.(*models.Loan); ok {
			payload = LoanStatusUpdate{
				ID:                  e.ID,
				Status:              e.Status,
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

func processMessage(kafkaMessage []byte) (err error) {
	var arr map[string]interface{}

	data := strings.SplitN(string(kafkaMessage), ":", 2)
	switch data[0] {
	default:
		return nil
	case "agent_provider":
		var mod models.AgentProvider

		err = json.Unmarshal([]byte(data[1]), &arr)
		if err != nil {
			return err
		}

		if arr["delete"] != nil && arr["delete"].(bool) == true {
			ID := uint64(arr["id"].(float64))
			err := mod.FindbyID(ID)
			if err != nil {
				return err
			}

			err = mod.Delete()
			if err != nil {
				return err
			}
		} else {
			err = json.Unmarshal([]byte(data[1]), &mod)
			if err != nil {
				return err
			}
			err = mod.FirstOrCreate()
		}
		break
	case "agent":
		var mod models.Agent

		err = json.Unmarshal([]byte(data[1]), &arr)
		if err != nil {
			return err
		}

		if arr["delete"] != nil && arr["delete"].(bool) == true {
			ID := uint64(arr["id"].(float64))
			err := mod.FindbyID(ID)
			if err != nil {
				return err
			}

			err = mod.Delete()
			if err != nil {
				return err
			}
		} else {
			err = json.Unmarshal([]byte(data[1]), &mod)
			if err != nil {
				return err
			}
			err = mod.FirstOrCreate()
		}
		break
	case "bank_type":
		var mod models.BankType

		err = json.Unmarshal([]byte(data[1]), &arr)
		if err != nil {
			return err
		}

		if arr["delete"] != nil && arr["delete"].(bool) == true {
			ID := uint64(arr["id"].(float64))
			err := mod.FindbyID(ID)
			if err != nil {
				return err
			}

			err = mod.Delete()
			if err != nil {
				return err
			}
		} else {
			err = json.Unmarshal([]byte(data[1]), &mod)
			if err != nil {
				return err
			}
			err = mod.FirstOrCreate()
		}
		break
	case "bank":
		var mod models.Bank

		err = json.Unmarshal([]byte(data[1]), &arr)
		if err != nil {
			return err
		}

		if arr["delete"] != nil && arr["delete"].(bool) == true {
			ID := uint64(arr["id"].(float64))
			err := mod.FindbyID(ID)
			if err != nil {
				return err
			}

			err = mod.Delete()
			if err != nil {
				return err
			}
		} else {
			err = json.Unmarshal([]byte(data[1]), &mod)
			if err != nil {
				return err
			}
			err = mod.FirstOrCreate()
		}
		break
	case "loan_purpose":
		var mod models.LoanPurpose

		err = json.Unmarshal([]byte(data[1]), &arr)
		if err != nil {
			return err
		}

		if arr["delete"] != nil && arr["delete"].(bool) == true {
			ID := uint64(arr["id"].(float64))
			err := mod.FindbyID(ID)
			if err != nil {
				return err
			}

			err = mod.Delete()
			if err != nil {
				return err
			}
		} else {
			err = json.Unmarshal([]byte(data[1]), &mod)
			if err != nil {
				return err
			}
			err = mod.FirstOrCreate()
		}
		break
	case "product":
		var mod models.Product

		err = json.Unmarshal([]byte(data[1]), &arr)
		if err != nil {
			return err
		}

		if arr["delete"] != nil && arr["delete"].(bool) == true {
			ID := uint64(arr["id"].(float64))
			err := mod.FindbyID(ID)
			if err != nil {
				return err
			}

			err = mod.Delete()
			if err != nil {
				return err
			}
		} else {
			err = json.Unmarshal([]byte(data[1]), &mod)
			if err != nil {
				return err
			}
			err = mod.FirstOrCreate()
		}
		break
	case "service":
		var mod models.Service

		err = json.Unmarshal([]byte(data[1]), &arr)
		if err != nil {
			return err
		}

		if arr["delete"] != nil && arr["delete"].(bool) == true {
			ID := uint64(arr["id"].(float64))
			err := mod.FindbyID(ID)
			if err != nil {
				return err
			}

			err = mod.Delete()
			if err != nil {
				return err
			}
		} else {
			err = json.Unmarshal([]byte(data[1]), &mod)
			if err != nil {
				return err
			}
			err = mod.FirstOrCreate()
		}
		break
	case "loan":
		var mod models.Loan

		err = json.Unmarshal([]byte(data[1]), &arr)
		if err != nil {
			return err
		}

		if arr["delete"] != nil && arr["delete"].(bool) == true {
			ID := uint64(arr["id"].(float64))
			err := mod.FindbyID(ID)
			if err != nil {
				return err
			}

			err = mod.Delete()
			if err != nil {
				return err
			}
		} else {
			err = json.Unmarshal([]byte(data[1]), &mod)
			if err != nil {
				return err
			}
			err = mod.FirstOrCreate()
		}
		break
	case "borrower":
		var mod models.Borrower

		err = json.Unmarshal([]byte(data[1]), &arr)
		if err != nil {
			return err
		}

		if arr["delete"] != nil && arr["delete"].(bool) == true {
			ID := uint64(arr["id"].(float64))
			err := mod.FindbyID(ID)
			if err != nil {
				return err
			}

			err = mod.Delete()
			if err != nil {
				return err
			}
		} else {
			err = json.Unmarshal([]byte(data[1]), &mod)
			if err != nil {
				return err
			}
			err = mod.FirstOrCreate()
		}
		break
	}
	return err
}
