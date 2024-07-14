package gateways

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/usecases/dto"
	mock_http "github.com/IgorRamosBR/g73-techchallenge-order/internal/infra/drivers/http/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPaymentClient_GeneratePaymentQRCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := mock_http.NewMockHttpClient(ctrl)

	type args struct {
		paymentRequest dto.PaymentRequest
	}
	type want struct {
		paymentResponse dto.PaymentQRCodeResponse
		err             error
	}
	type httpCall struct {
		apiUrl   string
		times    int
		response *http.Response
		err      error
	}
	tests := []struct {
		name string
		args
		want
		httpCall
	}{
		{
			name: "should fail to generate payment qr code when http client returns error",
			args: args{
				paymentRequest: createPaymentRequest(),
			},
			want: want{
				paymentResponse: dto.PaymentQRCodeResponse{},
				err:             errors.New("failed to call mercado pago broker, error: internal server error"),
			},
			httpCall: httpCall{
				apiUrl: "/payments",
				times:  1,
				err:    errors.New("internal server error"),
			},
		},
		{
			name: "should fail to generate payment qr code when http client returns non-2xx response",
			args: args{
				paymentRequest: createPaymentRequest(),
			},
			want: want{
				paymentResponse: dto.PaymentQRCodeResponse{},
				err:             errors.New("failed to pay order"),
			},
			httpCall: httpCall{
				apiUrl: "/payments",
				times:  1,
				response: &http.Response{
					StatusCode: 500,
					Body:       io.NopCloser(strings.NewReader("")),
				},
				err: nil,
			},
		},
		{
			name: "should fail to generate payment qr code when json decoder fails",
			args: args{
				paymentRequest: createPaymentRequest(),
			},
			want: want{
				paymentResponse: dto.PaymentQRCodeResponse{},
				err:             errors.New("failed to decode mercado pago response, error: invalid character '<' looking for beginning of value"),
			},
			httpCall: httpCall{
				apiUrl: "/payments",
				times:  1,
				response: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader("<invalid json>")),
				},
				err: nil,
			},
		},
		{
			name: "should generate payment qr code succesfully",
			args: args{
				paymentRequest: createPaymentRequest(),
			},
			want: want{
				paymentResponse: dto.PaymentQRCodeResponse{
					QrCode: "mercadopago123456",
				},
				err: nil,
			},
			httpCall: httpCall{
				apiUrl: "/payments",
				times:  1,
				response: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(`{"qrcode":"mercadopago123456"}`)),
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		httpClient.
			EXPECT().
			DoPost(gomock.Eq(tt.httpCall.apiUrl), gomock.Any()).
			Times(tt.httpCall.times).
			Return(tt.httpCall.response, tt.httpCall.err)

		paymentClient := NewPaymentClient(httpClient, "/payments")
		paymentResponse, err := paymentClient.GeneratePaymentQRCode(tt.args.paymentRequest)

		assert.Equal(t, tt.want.paymentResponse, paymentResponse)
		if tt.want.err != nil {
			assert.Equal(t, tt.want.err, err)
		}
	}
}

func createPaymentRequest() dto.PaymentRequest {
	return dto.PaymentRequest{
		OrderId: 123,
		Items: []dto.PaymentItemRequest{
			{
				Quantity: 1,
				Product: dto.PaymentProductRequest{
					Name:        "Batata Frita",
					SkuId:       "333",
					Description: "Batata canoa",
					Category:    "Acompanhamento",
					Price:       9.99,
				},
			},
		},
		TotalAmount: 9.99,
		CustomerCpf: "111222333444",
	}
}
