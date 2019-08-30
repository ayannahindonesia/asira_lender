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
	AsiraKafkaHandlers struct {
		KafkaProducer     sarama.AsyncProducer
		KafkaConsumer     sarama.Consumer
		PartitionConsumer sarama.PartitionConsumer
	}
	BorrowerInfo struct {
		Info interface{} `json:"borrower_info"`
	}
)

var wg sync.WaitGroup

func init() {
	topics := asira.App.Config.GetStringMap(fmt.Sprintf("%s.kafka.topics.consumes", asira.App.ENV))

	kafka := &AsiraKafkaHandlers{}
	kafka.KafkaProducer = asira.App.Kafka.Producer
	kafka.KafkaConsumer = asira.App.Kafka.Consumer

	kafka.SetPartitionConsumer(topics["for_lender"].(string))

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			message, err := kafka.Listen()
			if err != nil {
				log.Printf("error occured when listening kafka : %v", err)
			}
			if message != nil {
				err := getEntity(message)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
}

func (k *AsiraKafkaHandlers) SetPartitionConsumer(topic string) (err error) {
	k.PartitionConsumer, err = k.KafkaConsumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		return err
	}

	return nil
}

func (k *AsiraKafkaHandlers) Listen() ([]byte, error) {
	select {
	case err := <-k.PartitionConsumer.Errors():
		return nil, err
	case msg := <-k.PartitionConsumer.Messages():
		return msg.Value, nil
	}

	return nil, fmt.Errorf("unidentified error while listening")
}

func getEntity(kafkaMessage []byte) (err error) {
	data := strings.SplitN(string(kafkaMessage), ":", 2)
	switch data[0] {
	default:
		return nil
	case "loan":
		// create borrower first
		var borrowerInfo BorrowerInfo
		err = json.Unmarshal(kafkaMessage, &borrowerInfo)
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

		_, err = borrower.Save() // finish borrower create
		if err != nil {
			return err
		}

		// create loan
		var loan models.Loan
		err = json.Unmarshal(kafkaMessage, &loan)
		if err != nil {
			return err
		}

		loan.Bank = borrower.Bank
		loan.OwnerName = borrower.Fullname

		_, err = loan.Save() // finish create loan
		break
	}
	return err
}
