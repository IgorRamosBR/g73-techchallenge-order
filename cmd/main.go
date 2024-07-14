package main

import (
	"fmt"

	"github.com/g73-techchallenge-order/configs"
	"github.com/g73-techchallenge-order/internal/api"
	"github.com/g73-techchallenge-order/internal/controllers"
	"github.com/g73-techchallenge-order/internal/core/usecases"
	"github.com/g73-techchallenge-order/internal/infra/drivers/authorizer"
	"github.com/g73-techchallenge-order/internal/infra/drivers/broker"
	"github.com/g73-techchallenge-order/internal/infra/drivers/http"
	"github.com/g73-techchallenge-order/internal/infra/drivers/sql"
	"github.com/g73-techchallenge-order/internal/infra/gateways"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	appConfig := configs.GetAppConfig()

	httpClient := http.NewHttpClient(appConfig.DefaultTimeout)
	postgresSQLClient := createPostgresSQLClient(appConfig)
	err := performMigrations(postgresSQLClient, appConfig.DatabaseMigrationsPath)
	if err != nil {
		panic(err)
	}

	brokerChannel, err := NewRabbitMQBrokerChannel(appConfig.OrderEventsBrokerUrl)
	if err != nil {
		panic(err)
	}
	defer brokerChannel.Close()

	ordersPaidQueue, err := broker.NewRabbitMQConsumer(brokerChannel, appConfig.OrderEventsPaidQueue)
	if err != nil {
		panic(err)
	}

	ordersReadyQueue, err := broker.NewRabbitMQConsumer(brokerChannel, appConfig.OrderEventsReadyQueue)
	if err != nil {
		panic(err)
	}

	publisher := broker.NewRabbitMQPublisher(brokerChannel, appConfig.OrderEventsTopic)
	defer publisher.Close()

	orderNotify := gateways.NewOrderNotify(publisher, appConfig.OrderEventsReadyQueue)

	authorizer := authorizer.NewAuthorizer(httpClient, appConfig.AuthorizerURL)

	productRepositoryGateway := gateways.NewProductRepositoryGateway(postgresSQLClient)
	orderRepositoryGateway := gateways.NewOrderRepositoryGateway(postgresSQLClient)
	paymentClient := gateways.NewPaymentClient(httpClient, appConfig.PaymentURL)

	productUsecase := usecases.NewProductUsecase(productRepositoryGateway)
	paymentUsecase := usecases.NewPaymentUsecase(paymentClient)
	authorizerUsecase := usecases.NewAuthorizerUsecase(authorizer)
	orderUsecase := usecases.NewOrderUsecase(authorizerUsecase, paymentUsecase, productUsecase, orderNotify, orderRepositoryGateway)

	orderConsumerUseCase := usecases.NewOrderConsumerUseCase(ordersPaidQueue, ordersReadyQueue, publisher, orderUsecase)
	orderConsumerUseCase.StartConsumers()

	productController := controllers.NewProductController(productUsecase)
	orderController := controllers.NewOrderController(orderUsecase)

	apiParams := api.ApiParams{
		ProductController: productController,
		OrderController:   orderController,
	}
	api := api.NewApi(apiParams)
	api.Run(":" + appConfig.Port)
}

func createPostgresSQLClient(appConfig configs.AppConfig) sql.SQLClient {
	db, err := sql.NewPostgresSQLClient(appConfig.DatabaseUser, appConfig.DatabasePassword, appConfig.DatabaseHost, appConfig.DatabasePort, appConfig.DatabaseName, appConfig.DatabaseSSLMode)
	if err != nil {
		panic(fmt.Errorf("failed to connect database, error %w", err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("failed to ping database, error %w", err))
	}

	return db
}

func performMigrations(client sql.SQLClient, migrationsPath string) error {
	driver, err := postgres.WithInstance(client.GetConnection(), &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func NewRabbitMQBrokerChannel(url string) (*amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return ch, err
}
