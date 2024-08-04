package gateways

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IgorRamosBR/g73-techchallenge-order/pkg/events"
	"github.com/IgorRamosBR/g73-techchallenge-order/pkg/events/broker"
)

type OrderNotify interface {
	NotifyPaymentOrder(order events.OrderProductionDTO) error
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

func (o orderNotify) NotifyPaymentOrder(order events.OrderProductionDTO) error {
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
