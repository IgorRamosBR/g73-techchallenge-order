package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/g73-techchallenge-order/internal/core/usecases/dto"
	mock_usecases "github.com/g73-techchallenge-order/internal/core/usecases/mocks"
	"github.com/g73-techchallenge-order/internal/infra/drivers/authorizer"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
)

var orderRequestMissingStatus, _ = os.ReadFile("./testdata/order_request_missing_status.json")
var orderRequestWrongCpf, _ = os.ReadFile("./testdata/order_request_wrong_cpf.json")
var orderRequestValid, _ = os.ReadFile("./testdata/order_request_valid.json")

func TestOrderController_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	orderUseCase := mock_usecases.NewMockOrderUsecase(ctrl)
	orderController := NewOrderController(orderUseCase)

	type args struct {
		reqBody string
	}
	type want struct {
		statusCode int
		respBody   string
	}
	type orderUseCaseCall struct {
		times         int
		orderResponse dto.OrderCreationResponse
		err           error
	}
	tests := []struct {
		name string
		args
		want
		orderUseCaseCall
	}{
		{
			name: "should return bad request when req body is not a json",
			args: args{
				reqBody: "<invalidJson>",
			},
			want: want{
				statusCode: 400,
				respBody:   `{"message":"failed to bind order payload","error":"invalid character '\u003c' looking for beginning of value"}`,
			},
		},
		{
			name: "should return bad request when status is missing in the request",
			args: args{
				reqBody: string(orderRequestMissingStatus),
			},
			want: want{
				statusCode: 400,
				respBody:   `{"message":"invalid order payload","error":"Status is invalid"}`,
			},
		},
		{
			name: "should return bad request when cpf is wrong in the request",
			args: args{
				reqBody: string(orderRequestWrongCpf),
			},
			want: want{
				statusCode: 400,
				respBody:   `{"message":"invalid order payload","error":"invalid CPF [11122233344]"}`,
			},
		},
		{
			name: "should not authorize request when the user is not authorized",
			args: args{
				reqBody: string(orderRequestValid),
			},
			want: want{
				statusCode: 403,
				respBody:   `{"message":"customer cpf invalid","error":"customer unauthorized"}`,
			},
			orderUseCaseCall: orderUseCaseCall{
				times:         1,
				orderResponse: dto.OrderCreationResponse{},
				err:           authorizer.ErrUnauthorized,
			},
		},
		{
			name: "should not create order when the user case returns error",
			args: args{
				reqBody: string(orderRequestValid),
			},
			want: want{
				statusCode: 500,
				respBody:   `{"message":"failed to create order","error":"internal server error"}`,
			},
			orderUseCaseCall: orderUseCaseCall{
				times:         1,
				orderResponse: dto.OrderCreationResponse{},
				err:           errors.New("internal server error"),
			},
		},
		{
			name: "should create order succesfully",
			args: args{
				reqBody: string(orderRequestValid),
			},
			want: want{
				statusCode: 200,
				respBody:   `{"qrCode":"mercadopago123456","orderId":98765}`,
			},
			orderUseCaseCall: orderUseCaseCall{
				times: 1,
				orderResponse: dto.OrderCreationResponse{
					QRCode:  "mercadopago123456",
					OrderID: 98765,
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		orderUseCase.
			EXPECT().
			CreateOrder(gomock.Any()).
			Times(tt.orderUseCaseCall.times).
			Return(tt.orderUseCaseCall.orderResponse, tt.orderUseCaseCall.err)

		router := createRouter(orderController)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/order", strings.NewReader(tt.args.reqBody))
		router.ServeHTTP(w, req)

		assert.Equal(t, tt.want.statusCode, w.Code)
		assert.Equal(t, tt.want.respBody, w.Body.String())
	}
}

func createRouter(orderController OrderController) *gin.Engine {

	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.POST("/order", orderController.CreateOrder)
	}
	return router
}
