package configs

import (
	"os"
	"time"
)

type AppConfig struct {
	Port                   string
	DatabaseHost           string
	DatabasePort           string
	DatabaseName           string
	DatabaseUser           string
	DatabasePassword       string
	DatabaseSSLMode        string
	DatabaseMigrationsPath string

	AuthorizerURL string
	PaymentURL    string

	OrderEventsBrokerUrl             string
	OrderEventsTopic                 string
	OrderEventsPaidQueue             string
	OrderEventsReadyQueue            string
	OrderEventsInProgressDestination string

	DefaultTimeout time.Duration
}

func GetAppConfig() AppConfig {
	appConfig := AppConfig{}

	appConfig.Port = os.Getenv("PORT")

	appConfig.DatabaseHost = os.Getenv("POSTGRES_HOST")
	appConfig.DatabasePort = os.Getenv("POSTGRES_PORT")
	appConfig.DatabaseName = os.Getenv("POSTGRES_DB")
	appConfig.DatabaseSSLMode = os.Getenv("POSTGRES_SSLMODE")
	appConfig.DatabaseUser = os.Getenv("POSTGRES_USER")
	appConfig.DatabasePassword = os.Getenv("POSTGRES_PASSWORD")
	appConfig.DatabaseMigrationsPath = os.Getenv("MIGRATIONS_PATH")

	appConfig.AuthorizerURL = os.Getenv("AUTHORIZER_URL")
	appConfig.PaymentURL = os.Getenv("PAYMENT_URL")

	appConfig.OrderEventsBrokerUrl = os.Getenv("ORDER_EVENTS_BROKER_URL")
	appConfig.OrderEventsTopic = os.Getenv("ORDER_EVENTS_TOPIC")
	appConfig.OrderEventsPaidQueue = os.Getenv("ORDER_EVENTS_PAID_QUEUE")
	appConfig.OrderEventsReadyQueue = os.Getenv("ORDER_EVENTS_READY_QUEUE")
	appConfig.OrderEventsInProgressDestination = os.Getenv("ORDER_EVENTS_IN_PROGRESS_DESTINATION")

	defaultTimeout := os.Getenv("DEFAULT_TIMEOUT")
	defaultTimeoutDuration, err := time.ParseDuration(defaultTimeout)
	if err != nil {
		panic(err)
	}
	appConfig.DefaultTimeout = defaultTimeoutDuration

	return appConfig
}
