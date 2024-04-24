package configs

import (
	"os"
)

type AppConfig struct {
	Environment string

	DatabaseHost     string
	DatabasePort     string
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseSSLMode  string

	AuthorizerURL string

	SQSRegion   string
	SQSEndpoint string

	PaymentURL      string
	NotificationURL string
	SponsorId       string
}

func GetAppConfig() AppConfig {
	appConfig := AppConfig{}

	appConfig.Environment = os.Getenv("ENVIRONMENT")

	appConfig.DatabaseHost = os.Getenv("POSTGRES_HOST")
	appConfig.DatabasePort = os.Getenv("POSTGRES_PORT")
	appConfig.DatabaseName = os.Getenv("POSTGRES_DB")
	appConfig.DatabaseSSLMode = os.Getenv("POSTGRES_SSLMODE")
	appConfig.DatabaseUser = os.Getenv("POSTGRES_USER")
	appConfig.DatabasePassword = os.Getenv("POSTGRES_PASSWORD")

	appConfig.AuthorizerURL = os.Getenv("AUTHORIZER_URL")

	appConfig.PaymentURL = os.Getenv("PAYMENT_API_URL")

	return appConfig
}
