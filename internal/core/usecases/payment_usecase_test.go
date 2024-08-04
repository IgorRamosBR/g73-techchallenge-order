package usecases

import (
	"errors"
	"testing"
	"time"

	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/entities"
	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/usecases/dto"
	mock_gateways "github.com/IgorRamosBR/g73-techchallenge-order/internal/infra/gateways/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPaymentUsecase_GeneratePaymentQRCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	paymentClient := mock_gateways.NewMockPaymentClient(ctrl)

	paymentClient.EXPECT().GeneratePaymentQRCode(gomock.Any()).Times(1).Return(dto.PaymentQRCodeResponse{}, errors.New("internal server error"))

	paymentUsecase := NewPaymentUsecase(paymentClient)
	qrcode, err := paymentUsecase.GeneratePaymentQRCode(createOrder())

	assert.Empty(t, qrcode)
	assert.Error(t, err, "failed to generate payment qrcode for the order 123, error: internal server error")

	paymentClient.EXPECT().GeneratePaymentQRCode(gomock.Any()).Times(1).Return(dto.PaymentQRCodeResponse{
		QrCode: "mercadopago123456",
	}, nil)

	qrcode, err = paymentUsecase.GeneratePaymentQRCode(createOrder())

	assert.Equal(t, "mercadopago123456", qrcode)
	assert.NoError(t, err)

}

func createOrder() entities.Order {
	return entities.Order{
		ID: 123,
		Items: []entities.OrderItem{
			{
				ID:       999,
				Quantity: 1,
				Type:     "UNIT",
				Product: entities.Product{
					ID:          222,
					Name:        "Batata Frita",
					SkuId:       "333",
					Description: "Batata canoa",
					Category:    "Acompanhamento",
					Price:       9.99,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
			},
		},
		Coupon:      "APP10",
		TotalAmount: 9.99,
		Status:      "PAID",
		CreatedAt:   time.Time{},
		CustomerCPF: "111222333444",
	}
}
