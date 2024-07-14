package gateways

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/g73-techchallenge-order/internal/core/usecases/dto"
	"github.com/g73-techchallenge-order/internal/infra/drivers/broker"
)

type OrderNotify interface {
	NotifyPaymentOrder(order dto.ProductionOrderDTO) error
}

type orderNotify struct {
	publisher   broker.Publisher
	destination string
}

type OrderPaymentMessage struct {
	OrderId int `json:"orderId"`
}

func NewOrderNotify(publisher broker.Publisher, destination string) OrderNotify {
	return orderNotify{publisher: publisher, destination: destination}
}

func (o orderNotify) NotifyPaymentOrder(order dto.ProductionOrderDTO) error {
	message, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal payment order[%d], error: %v", order.ID, err)
	}

	ctx := context.Background()
	err = o.publisher.Publish(ctx, o.destination, message)
	if err != nil {
		return fmt.Errorf("failed to publish order[%d], error: %v", order.ID, err)
	}

	return nil
}
