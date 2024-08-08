package usecases

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/usecases/dto"
	mock_usecases "github.com/IgorRamosBR/g73-techchallenge-order/internal/core/usecases/mocks"
	"github.com/IgorRamosBR/g73-techchallenge-order/pkg/events"
	mock_broker "github.com/IgorRamosBR/g73-techchallenge-order/pkg/events/broker/mocks"
	"go.uber.org/mock/gomock"
)

func TestStartConsumers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderPaidConsumer := mock_broker.NewMockConsumer(ctrl)
	mockOrderReadyConsumer := mock_broker.NewMockConsumer(ctrl)
	mockOrderPublisher := mock_broker.NewMockPublisher(ctrl)
	mockOrderUsecase := mock_usecases.NewMockOrderUseCase(ctrl)

	uc := NewOrderConsumerUseCase(mockOrderPaidConsumer, mockOrderReadyConsumer, mockOrderPublisher, mockOrderUsecase)

	mockOrderPaidConsumer.EXPECT().StartConsumer(gomock.Any()).Times(1)
	mockOrderReadyConsumer.EXPECT().StartConsumer(gomock.Any()).Times(1)

	uc.StartConsumers()
	time.Sleep(1 * time.Second)
}

func TestProcessOrderMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderPaidConsumer := mock_broker.NewMockConsumer(ctrl)
	mockOrderReadyConsumer := mock_broker.NewMockConsumer(ctrl)
	mockOrderPublisher := mock_broker.NewMockPublisher(ctrl)
	mockOrderUsecase := mock_usecases.NewMockOrderUseCase(ctrl)

	uc := &orderConsumerUseCase{
		orderPaidConsumer:  mockOrderPaidConsumer,
		orderReadyConsumer: mockOrderReadyConsumer,
		orderPublisher:     mockOrderPublisher,
		orderUsecase:       mockOrderUsecase,
	}

	t.Run("successful processing", func(t *testing.T) {
		orderEvent := events.OrderStatusEventDTO{
			OrderId: 123,
			Status:  "Paid",
		}
		message, _ := json.Marshal(orderEvent)

		mockOrderUsecase.EXPECT().UpdateOrderStatus(orderEvent.OrderId, dto.OrderStatus(orderEvent.Status)).Return(nil).Times(1)

		err := uc.ProcessOrderMessage(message)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("failed to unmarshal message", func(t *testing.T) {
		invalidMessage := []byte("invalid")

		err := uc.ProcessOrderMessage(invalidMessage)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("failed to update order status", func(t *testing.T) {
		orderEvent := events.OrderStatusEventDTO{
			OrderId: 123,
			Status:  "Paid",
		}
		message, _ := json.Marshal(orderEvent)

		mockOrderUsecase.EXPECT().UpdateOrderStatus(orderEvent.OrderId, dto.OrderStatus(orderEvent.Status)).Return(errors.New("update failed")).Times(1)

		err := uc.ProcessOrderMessage(message)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
