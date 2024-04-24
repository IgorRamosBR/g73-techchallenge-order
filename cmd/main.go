package main

import (
	"fmt"

	"github.com/g73-techchallenge-order/configs"
	"github.com/g73-techchallenge-order/internal/api"
	"github.com/g73-techchallenge-order/internal/controllers"
	"github.com/g73-techchallenge-order/internal/core/usecases"
	"github.com/g73-techchallenge-order/internal/infra/drivers/authorizer"
	"github.com/g73-techchallenge-order/internal/infra/drivers/http"
	"github.com/g73-techchallenge-order/internal/infra/drivers/sql"
	"github.com/g73-techchallenge-order/internal/infra/gateways"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	appConfig := configs.GetAppConfig()

	httpClient := http.NewMockHttpClient()
	postgresSQLClient := createPostgresSQLClient(appConfig)
	err := performMigrations(postgresSQLClient)
	if err != nil {
		panic(err)
	}

	authorizerClient := http.NewHttpClient()
	authorizer := authorizer.NewAuthorizer(authorizerClient, appConfig.AuthorizerURL)

	customerRepositoryGateway := gateways.NewCustomerRepositoryGateway(postgresSQLClient)
	productRepositoryGateway := gateways.NewProductRepositoryGateway(postgresSQLClient)
	orderRepositoryGateway := gateways.NewOrderRepositoryGateway(postgresSQLClient)
	paymentClient := gateways.NewPaymentClient(httpClient, appConfig.PaymentURL)

	customerUsecase := usecases.NewCustomerUsecase(customerRepositoryGateway)
	productUsecase := usecases.NewProductUsecase(productRepositoryGateway)
	paymentUsecase := usecases.NewPaymentUsecase(paymentClient)
	authorizerUsecase := usecases.NewAuthorizerUsecase(authorizer)
	orderUsecase := usecases.NewOrderUsecase(authorizerUsecase, paymentUsecase, productUsecase, orderRepositoryGateway)

	customerController := controllers.NewCustomerController(customerUsecase)
	productController := controllers.NewProductController(productUsecase)
	orderController := controllers.NewOrderController(orderUsecase)

	apiParams := api.ApiParams{
		CustomerController: customerController,
		ProductController:  productController,
		OrderController:    orderController,
	}
	api := api.NewApi(apiParams)
	api.Run(":8080")
}

func createPostgresSQLClient(appConfig configs.AppConfig) sql.SQLClient {
	db, err := sql.NewPostgresSQLClient(appConfig.DatabaseUser, appConfig.DatabasePassword, appConfig.DatabaseHost, appConfig.DatabasePort, appConfig.DatabaseName)
	if err != nil {
		panic(fmt.Errorf("failed to connect database, error %w", err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("failed to ping database, error %w", err))
	}

	return db
}

func performMigrations(client sql.SQLClient) error {
	driver, err := postgres.WithInstance(client.GetConnection(), &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
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
