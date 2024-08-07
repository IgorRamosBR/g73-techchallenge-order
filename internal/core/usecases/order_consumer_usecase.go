package usecases

import (
	"encoding/json"
	"fmt"

	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/usecases/dto"
	"github.com/IgorRamosBR/g73-techchallenge-order/pkg/events"
	"github.com/IgorRamosBR/g73-techchallenge-order/pkg/events/broker"
)

type OrderConsumerUseCase interface {
	StartConsumers()
}

type orderConsumerUseCase struct {
	orderPaidConsumer  broker.Consumer
	orderReadyConsumer broker.Consumer
	orderPublisher     broker.Publisher
	orderUsecase       OrderUseCase
}

type OrderConsumerUseCaseConfig struct {
	OrderPaidConsumer  broker.Consumer
	OrderReadyConsumer broker.Consumer
	OrderPublisher     broker.Publisher
	OrderUseCase       OrderUseCase
}

func NewOrderConsumerUseCase(orderPaidConsumer, orderReadyConsumer broker.Consumer, orderPublisher broker.Publisher, orderUsecase OrderUseCase) OrderConsumerUseCase {
	return &orderConsumerUseCase{
		orderPaidConsumer:  orderPaidConsumer,
		orderReadyConsumer: orderReadyConsumer,
		orderPublisher:     orderPublisher,
		orderUsecase:       orderUsecase,
	}
}

func (u *orderConsumerUseCase) StartConsumers() {
	go u.orderPaidConsumer.StartConsumer(u.ProcessOrderMessage)
	go u.orderReadyConsumer.StartConsumer(u.ProcessOrderMessage)
}

func (u *orderConsumerUseCase) ProcessOrderMessage(message []byte) error {
	var orderEvent events.OrderStatusEventDTO
	err := json.Unmarshal(message, &orderEvent)
	if err != nil {
		return fmt.Errorf("failed to unmarshall message, error: %w", err)
	}

	err = u.orderUsecase.UpdateOrderStatus(orderEvent.OrderId, dto.OrderStatus(orderEvent.Status))
	if err != nil {
		return fmt.Errorf("failed to update order status, error: %w", err)
	}

	return nil
}
