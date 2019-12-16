package middlewares

import (
	"asira_lender/asira"
	"asira_lender/models"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

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
		var loan models.Loan

		err = json.Unmarshal([]byte(data[1]), &loan)
		if err != nil {
			return err
		}
		err = loan.Save()
		break
	case "borrower":
		var borrower models.Borrower

		err = json.Unmarshal([]byte(data[1]), &borrower)
		if err != nil {
			return err
		}
		err = borrower.Save()
		break
	}
	return err
}
