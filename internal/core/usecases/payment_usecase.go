package usecases

import (
	"github.com/g73-techchallenge-order/internal/core/entities"
	"github.com/g73-techchallenge-order/internal/core/usecases/dto"
	"github.com/g73-techchallenge-order/internal/infra/gateways"

	log "github.com/sirupsen/logrus"
)

type PaymentUsecase interface {
	GeneratePaymentQRCode(order entities.Order) (string, error)
}

type paymentUsecase struct {
	paymentClient gateways.PaymentClient
}

func NewPaymentUsecase(paymentClient gateways.PaymentClient) PaymentUsecase {
	return paymentUsecase{
		paymentClient: paymentClient,
	}
}

func (u paymentUsecase) GeneratePaymentQRCode(order entities.Order) (string, error) {
	paymentRequest := u.createPaymentRequest(order)
	paymentResponse, err := u.paymentClient.GeneratePaymentQRCode(paymentRequest)
	if err != nil {
		log.Errorf("failed to generate payment qrcode for the order [%d], error: %v", order.ID, err)
		return "", err
	}

	return paymentResponse.QrCode, nil
}

func (u paymentUsecase) createPaymentRequest(order entities.Order) dto.PaymentRequest {
	var items []dto.PaymentItemRequest
	for _, item := range order.Items {
		items = append(items, createPaymentItem(item))
	}

	return dto.PaymentRequest{
		OrderId:     order.ID,
		CustomerCpf: order.CustomerCPF,
		TotalAmount: order.TotalAmount,
		Items:       items,
	}
}

func createPaymentItem(item entities.OrderItem) dto.PaymentItemRequest {
	paymentItem := dto.PaymentItemRequest{
		Quantity: item.Quantity,
		Product: dto.PaymentProductRequest{
			Name:        item.Product.Name,
			SkuId:       item.Product.SkuId,
			Description: item.Product.Description,
			Category:    item.Product.Category,
			Type:        item.Type,
			Price:       item.Product.Price,
		},
	}

	return paymentItem
}
