package gateways

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/usecases/dto"
	"github.com/IgorRamosBR/g73-techchallenge-order/internal/infra/drivers/http"
)

type PaymentClient interface {
	GeneratePaymentQRCode(dto.PaymentRequest) (dto.PaymentQRCodeResponse, error)
}

type paymentAPIClient struct {
	httpClient http.HttpClient
	apiUrl     string
}

func NewPaymentClient(httpClient http.HttpClient, apiUrl string) PaymentClient {
	return paymentAPIClient{
		httpClient: httpClient,
		apiUrl:     apiUrl,
	}
}

func (p paymentAPIClient) GeneratePaymentQRCode(request dto.PaymentRequest) (dto.PaymentQRCodeResponse, error) {
	reqBody, err := json.Marshal(&request)
	if err != nil {
		return dto.PaymentQRCodeResponse{}, fmt.Errorf("failed to marshal payment request, error: %v", err)
	}

	response, err := p.httpClient.DoPost(p.apiUrl, reqBody)
	if err != nil {
		return dto.PaymentQRCodeResponse{}, fmt.Errorf("failed to call mercado pago broker, error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return dto.PaymentQRCodeResponse{}, errors.New("failed to pay order")
	}

	var paymentQRCodeResponse dto.PaymentQRCodeResponse
	err = json.NewDecoder(response.Body).Decode(&paymentQRCodeResponse)
	if err != nil {
		return dto.PaymentQRCodeResponse{}, fmt.Errorf("failed to decode mercado pago response, error: %v", err)
	}

	return paymentQRCodeResponse, nil
}
