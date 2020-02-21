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
		return createUnitTestModels(i, model)
	}

	topic := asira.App.Config.GetString(fmt.Sprintf("%s.kafka.topics.produces", asira.App.ENV))

	var payload interface{}

	payload = kafkaPayloadBuilder(i, &model)

	jMarshal, _ := json.Marshal(payload)

	kafkaProducer, err := sarama.NewAsyncProducer([]string{asira.App.Kafka.Host}, asira.App.Kafka.Config)
	if err != nil {
		return err
	}
	defer kafkaProducer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(model + ":" + string(jMarshal)),
	}

	select {
	case kafkaProducer.Input() <- msg:
		log.Printf("Produced topic : %s", topic)
	case err := <-kafkaProducer.Errors():
		log.Printf("Fail producing topic : %s error : %v", topic, err)
	}

	return nil
}

func kafkaPayloadBuilder(i interface{}, model *string) (payload interface{}) {
	type KafkaModelPayload struct {
		ID      float64     `json:"id"`
		Payload interface{} `json:"payload"`
		Mode    string      `json:"mode"`
	}
	var mode string

	checkSuffix := *model
	if strings.HasSuffix(checkSuffix, "_delete") {
		mode = "delete"
		*model = strings.TrimSuffix(checkSuffix, "_delete")
	} else if strings.HasSuffix(checkSuffix, "_create") {
		mode = "create"
		*model = strings.TrimSuffix(checkSuffix, "_create")
	} else if strings.HasSuffix(checkSuffix, "_update") {
		mode = "update"
		*model = strings.TrimSuffix(checkSuffix, "_update")
	}

	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(i)
	json.Unmarshal(inrec, &inInterface)
	if modelID, ok := inInterface["id"].(float64); ok {
		payload = KafkaModelPayload{
			ID:      modelID,
			Payload: i,
			Mode:    mode,
		}
	}

	return payload
}

func processMessage(kafkaMessage []byte) (err error) {
	var arr map[string]interface{}

	data := strings.SplitN(string(kafkaMessage), ":", 2)

	err = json.Unmarshal([]byte(data[1]), &arr)
	if err != nil {
		return err
	}

	switch data[0] {
	default:
		return nil
	case "agent_provider":
		mod := models.AgentProvider{}

		marshal, _ := json.Marshal(arr["payload"])
		json.Unmarshal(marshal, &mod)

		switch arr["mode"] {
		default:
			err = fmt.Errorf("invalid payload")
			break
		case "create":
			err = mod.FirstOrCreate()
			break
		case "update":
			err = mod.Save()
			break
		case "delete":
			err = mod.Delete()
			break
		}
		break
	case "agent":
		mod := models.Agent{}

		marshal, _ := json.Marshal(arr["payload"])
		json.Unmarshal(marshal, &mod)

		switch arr["mode"] {
		default:
			err = fmt.Errorf("invalid payload")
			break
		case "create":
			err = mod.FirstOrCreate()
			break
		case "update":
			err = mod.Save()
			break
		case "delete":
			err = mod.Delete()
			break
		}
		break
	case "bank_type":
		mod := models.BankType{}

		marshal, _ := json.Marshal(arr["payload"])
		json.Unmarshal(marshal, &mod)

		switch arr["mode"] {
		default:
			err = fmt.Errorf("invalid payload")
			break
		case "create":
			err = mod.FirstOrCreate()
			break
		case "update":
			err = mod.Save()
			break
		case "delete":
			err = mod.Delete()
			break
		}
		break
	case "bank":
		mod := models.Bank{}

		marshal, _ := json.Marshal(arr["payload"])
		json.Unmarshal(marshal, &mod)

		switch arr["mode"] {
		default:
			err = fmt.Errorf("invalid payload")
			break
		case "create":
			err = mod.FirstOrCreate()
			break
		case "update":
			err = mod.Save()
			break
		case "delete":
			err = mod.Delete()
			break
		}
		break
	case "loan_purpose":
		mod := models.LoanPurpose{}

		marshal, _ := json.Marshal(arr["payload"])
		json.Unmarshal(marshal, &mod)

		switch arr["mode"] {
		default:
			err = fmt.Errorf("invalid payload")
			break
		case "create":
			err = mod.FirstOrCreate()
			break
		case "update":
			err = mod.Save()
			break
		case "delete":
			err = mod.Delete()
			break
		}
		break
	case "product":
		mod := models.Product{}

		marshal, _ := json.Marshal(arr["payload"])
		json.Unmarshal(marshal, &mod)

		switch arr["mode"] {
		default:
			err = fmt.Errorf("invalid payload")
			break
		case "create":
			err = mod.FirstOrCreate()
			break
		case "update":
			err = mod.Save()
			break
		case "delete":
			err = mod.Delete()
			break
		}
		break
	case "service":
		mod := models.Service{}

		marshal, _ := json.Marshal(arr["payload"])
		json.Unmarshal(marshal, &mod)

		switch arr["mode"] {
		default:
			err = fmt.Errorf("invalid payload")
			break
		case "create":
			err = mod.FirstOrCreate()
			break
		case "update":
			err = mod.Save()
			break
		case "delete":
			err = mod.Delete()
			break
		}
		break
	case "loan":
		mod := models.Loan{}

		marshal, _ := json.Marshal(arr["payload"])
		json.Unmarshal(marshal, &mod)

		switch arr["mode"] {
		default:
			err = fmt.Errorf("invalid payload")
			break
		case "create":
			err = mod.FirstOrCreate()
			break
		case "update":
			err = mod.Save()
			break
		case "delete":
			err = mod.Delete()
			break
		}
		break
	case "borrower":
		mod := models.Borrower{}

		marshal, _ := json.Marshal(arr["payload"])
		json.Unmarshal(marshal, &mod)

		switch arr["mode"] {
		default:
			err = fmt.Errorf("invalid payload")
			break
		case "create":
			err = mod.FirstOrCreate()
			break
		case "update":
			err = mod.Save()
			break
		case "delete":
			err = mod.Delete()
			break
		}
		break
	case "faq":
		mod := models.FAQ{}

		marshal, _ := json.Marshal(arr["payload"])
		json.Unmarshal(marshal, &mod)

		switch arr["mode"] {
		default:
			err = fmt.Errorf("invalid payload")
			break
		case "create":
			err = mod.FirstOrCreate()
			break
		case "update":
			err = mod.Save()
			break
		case "delete":
			err = mod.Delete()
			break
		}
	}
	return err
}

func createUnitTestModels(i interface{}, model string) error {
	var (
		mode string
		err  error
	)
	if strings.HasSuffix(model, "_delete") {
		mode = "delete"
		model = strings.TrimSuffix(model, "_delete")
	} else if strings.HasSuffix(model, "_create") {
		mode = "create"
		model = strings.TrimSuffix(model, "_create")
	} else if strings.HasSuffix(model, "_update") {
		mode = "update"
		model = strings.TrimSuffix(model, "_update")
	}

	switch model {
	default:
		return fmt.Errorf("invalid model")
	case "agent_provider":
		if x, ok := i.(models.AgentProvider); ok {
			switch mode {
			default:
				return fmt.Errorf("invalid model")
			case "create":
				err = x.Create()
				break
			case "update":
				err = x.Save()
				break
			case "delete":
				err = x.Delete()
				break
			}
		}
		break
	case "agent":
		if x, ok := i.(models.Agent); ok {
			switch mode {
			default:
				return fmt.Errorf("invalid model")
			case "create":
				err = x.Create()
				break
			case "update":
				err = x.Save()
				break
			case "delete":
				err = x.Delete()
				break
			}
		}
		break
	case "bank_type":
		if x, ok := i.(models.BankType); ok {
			switch mode {
			default:
				return fmt.Errorf("invalid model")
			case "create":
				err = x.Create()
				break
			case "update":
				err = x.Save()
				break
			case "delete":
				err = x.Delete()
				break
			}
		}
		break
	case "bank":
		if x, ok := i.(models.Bank); ok {
			switch mode {
			default:
				return fmt.Errorf("invalid model")
			case "create":
				err = x.Create()
				break
			case "update":
				err = x.Save()
				break
			case "delete":
				err = x.Delete()
				break
			}
		}
		break
	case "loan_purpose":
		if x, ok := i.(models.LoanPurpose); ok {
			switch mode {
			default:
				return fmt.Errorf("invalid model")
			case "create":
				err = x.Create()
				break
			case "update":
				err = x.Save()
				break
			case "delete":
				err = x.Delete()
				break
			}
		}
		break
	case "product":
		if x, ok := i.(models.Product); ok {
			switch mode {
			default:
				return fmt.Errorf("invalid model")
			case "create":
				err = x.Create()
				break
			case "update":
				err = x.Save()
				break
			case "delete":
				err = x.Delete()
				break
			}
		}
		break
	case "service":
		if x, ok := i.(models.Service); ok {
			switch mode {
			default:
				return fmt.Errorf("invalid model")
			case "create":
				err = x.Create()
				break
			case "update":
				err = x.Save()
				break
			case "delete":
				err = x.Delete()
				break
			}
		}
		break
	case "loan":
		if x, ok := i.(models.Loan); ok {
			switch mode {
			default:
				return fmt.Errorf("invalid model")
			case "create":
				err = x.Create()
				break
			case "update":
				err = x.Save()
				break
			case "delete":
				err = x.Delete()
				break
			}
		}
		break
	case "borrower":
		if x, ok := i.(models.Borrower); ok {
			switch mode {
			default:
				return fmt.Errorf("invalid model")
			case "create":
				err = x.Create()
				break
			case "update":
				err = x.Save()
				break
			case "delete":
				err = x.Delete()
				break
			}
		}
		break
	case "faq":
		if x, ok := i.(models.FAQ); ok {
			switch mode {
			default:
				return fmt.Errorf("invalid model")
			case "create":
				err = x.Create()
				break
			case "update":
				err = x.Save()
				break
			case "delete":
				err = x.Delete()
				break
			}
		}
	}

	return err
}
