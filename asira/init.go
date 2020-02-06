package asira

import (
	"asira_lender/cron"
	"asira_lender/custommodule"
	"asira_lender/validator"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/ayannahindonesia/basemodel"
	"github.com/fsnotify/fsnotify"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	// import postgres misc
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// App main var
var App *Application

type (
	// Application main variable of the app
	Application struct {
		Name       string          `json:"name"`
		Port       string          `json:"port"`
		Version    string          `json:"version"`
		ENV        string          `json:"env"`
		Config     viper.Viper     `json:"prog_config"`
		DB         *gorm.DB        `json:"db"`
		Kafka      KafkaInstance   `json:"kafka"`
		S3         custommodule.S3 `json:"s3"`
		Cron       cron.Cron       `json:"cron"`
		Permission viper.Viper     `json:"prog_permission"`
	}

	// KafkaInstance stores kafka configs
	KafkaInstance struct {
		Config *sarama.Config
		Host   string
	}
)

// Initiate asira instances
func init() {
	var err error
	App = &Application{}
	App.Name = "asira_lender"
	App.Port = os.Getenv("APPPORT")
	App.Version = os.Getenv("APPVER")
	App.loadENV()
	if err = App.LoadConfigs(); err != nil {
		log.Printf("Load config error : %v", err)
	}
	if err = App.DBinit(); err != nil {
		log.Printf("DB init error : %v", err)
	}
	if err = App.LoadPermissions(); err != nil {
		log.Printf("Load Permission error : %v", err)
	}

	App.KafkaInit()
	App.S3init()
	App.CronInit()

	// apply custom validator
	v := validator.AsiraValidator{DB: App.DB}
	v.CustomValidatorRules()
}

// Close apps
func (x *Application) Close() (err error) {
	if err = x.DB.Close(); err != nil {
		return err
	}
	x.Cron.Stop()

	return nil
}

// Loads environtment setting
func (x *Application) loadENV() {
	APPENV := os.Getenv("APPENV")

	switch APPENV {
	default:
		x.ENV = "development"
		break
	case "development":
		x.ENV = "development"
		break
	case "staging":
		x.ENV = "staging"
		break
	case "production":
		x.ENV = "production"
		break
	}
}

// LoadConfigs loads general configs
func (x *Application) LoadConfigs() error {
	var conf *viper.Viper

	conf = viper.New()
	conf.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	conf.AutomaticEnv()
	conf.SetConfigName("config")
	conf.AddConfigPath("$GOPATH/src/" + x.Name)
	conf.SetConfigType("yaml")
	if err := conf.ReadInConfig(); err != nil {
		return err
	}
	conf.WatchConfig()
	conf.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("App Config file changed %s :", e.Name)
		x.LoadConfigs()
	})
	x.Config = viper.Viper(*conf)

	return nil
}

// DBinit loads DBinit configs
func (x *Application) DBinit() error {
	dbconf := x.Config.GetStringMap(fmt.Sprintf("%s.database", x.ENV))
	Cons := basemodel.DBConfig{
		Adapter:        basemodel.PostgresAdapter,
		Host:           dbconf["host"].(string),
		Port:           dbconf["port"].(string),
		Username:       dbconf["username"].(string),
		Password:       dbconf["password"].(string),
		Table:          dbconf["table"].(string),
		Timezone:       dbconf["timezone"].(string),
		Maxlifetime:    dbconf["maxlifetime"].(int),
		IdleConnection: dbconf["idle_conns"].(int),
		OpenConnection: dbconf["open_conns"].(int),
		SSL:            dbconf["sslmode"].(string),
		Logmode:        dbconf["logmode"].(bool),
	}
	basemodel.Start(Cons)
	x.DB = basemodel.DB
	return nil
}

// KafkaInit init kafka instance
func (x *Application) KafkaInit() {
	kafkaConf := x.Config.GetStringMap(fmt.Sprintf("%s.kafka", x.ENV))

	if kafkaConf["log_verbose"].(bool) {
		sarama.Logger = log.New(os.Stdout, "[lender kafka] ", log.LstdFlags)
	}

	x.Kafka.Config = sarama.NewConfig()
	x.Kafka.Config.ClientID = kafkaConf["client_id"].(string)
	if kafkaConf["sasl"].(bool) {
		x.Kafka.Config.Net.SASL.Enable = true
	}

	x.Kafka.Config.Net.SASL.User = kafkaConf["user"].(string)
	x.Kafka.Config.Net.SASL.Password = kafkaConf["pass"].(string)

	x.Kafka.Config.Producer.Return.Successes = true
	x.Kafka.Config.Producer.Partitioner = sarama.NewRandomPartitioner
	x.Kafka.Config.Producer.RequiredAcks = sarama.WaitForAll
	x.Kafka.Config.Producer.Flush.Frequency = 500 * time.Millisecond

	x.Kafka.Config.Consumer.Return.Errors = true

	x.Kafka.Host = strings.Join([]string{kafkaConf["host"].(string), kafkaConf["port"].(string)}, ":")
}

// LoadPermissions loads general configs
func (x *Application) LoadPermissions() error {
	var conf *viper.Viper

	conf = viper.New()
	conf.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	conf.AutomaticEnv()
	conf.SetConfigName("permissions")
	conf.AddConfigPath("/go/src/")
	conf.SetConfigType("yaml")
	if err := conf.ReadInConfig(); err != nil {
		return err
	}
	conf.WatchConfig()
	conf.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("App Config file changed %s:", e.Name)
		x.LoadConfigs()
	})
	x.Permission = viper.Viper(*conf)

	return nil
}

// S3init load config for s3
func (x *Application) S3init() (err error) {
	s3conf := x.Config.GetStringMap(fmt.Sprintf("%s.s3", x.ENV))

	x.S3, err = custommodule.NewS3(s3conf["access_key"].(string), s3conf["secret_key"].(string), s3conf["host"].(string), s3conf["bucket_name"].(string), s3conf["region"].(string))

	return err
}

// CronInit load cron
func (x *Application) CronInit() (err error) {
	x.Cron.TZ = x.Config.GetString(fmt.Sprintf("%s.database.timezone", x.ENV))
	cron.DB = x.DB
	x.Cron.Time = x.Config.GetString(fmt.Sprintf("%s.cron.time", x.ENV))
	x.Cron.New()
	x.Cron.Start()

	return nil
}
