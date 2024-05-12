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

	AuthorizerURL  string
	PaymentURL     string
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

	defaultTimeout := os.Getenv("DEFAULT_TIMEOUT")
	defaultTimeoutDuration, err := time.ParseDuration(defaultTimeout)
	if err != nil {
		panic(err)
	}
	appConfig.DefaultTimeout = defaultTimeoutDuration

	return appConfig
}
