package usecases

import (
	"errors"
	"testing"
	"time"

	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/entities"
	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/usecases/dto"
	mock_usecases "github.com/IgorRamosBR/g73-techchallenge-order/internal/core/usecases/mocks"
	"github.com/IgorRamosBR/g73-techchallenge-order/internal/infra/drivers/authorizer"
	mock_gateways "github.com/IgorRamosBR/g73-techchallenge-order/internal/infra/gateways/mocks"
	"github.com/IgorRamosBR/g73-techchallenge-order/pkg/events"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestOrderUsecase_GetAllOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	orderRepository := mock_gateways.NewMockOrderRepositoryGateway(ctrl)

	orderUsecase := NewOrderUsecase(nil, nil, nil, nil, orderRepository)

	pageParams := dto.NewPageParams(20, 10)

	orderRepository.EXPECT().
		FindAllOrders(gomock.Eq(pageParams)).
		Times(1).
		Return(nil, errors.New("internal server error"))

	orders, err := orderUsecase.GetAllOrders(pageParams)

	assert.Empty(t, orders)
	assert.EqualError(t, err, "internal server error")

	returnedOrders := []entities.Order{createOrder()}
	expectedOrders := dto.Page[entities.Order]{
		Result: returnedOrders,
		Next:   nil,
	}
	orderRepository.EXPECT().
		FindAllOrders(gomock.Eq(pageParams)).
		Times(1).
		Return(returnedOrders, nil)

	orders, err = orderUsecase.GetAllOrders(pageParams)

	assert.Equal(t, expectedOrders, orders)
	assert.NoError(t, err)
}

func TestOrderUsecase_GetOrderStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	orderRepository := mock_gateways.NewMockOrderRepositoryGateway(ctrl)

	orderUsecase := NewOrderUsecase(nil, nil, nil, nil, orderRepository)

	orderId := 123

	orderRepository.EXPECT().
		GetOrderStatus(gomock.Eq(orderId)).
		Times(1).
		Return("", errors.New("internal server error"))

	orderStatus, err := orderUsecase.GetOrderStatus(orderId)

	assert.Empty(t, orderStatus)
	assert.EqualError(t, err, "internal server error")

	expectedOrderStatus := dto.OrderStatusDTO{
		Status: dto.OrderStatusCreated,
	}

	orderRepository.EXPECT().
		GetOrderStatus(gomock.Eq(orderId)).
		Times(1).
		Return("CREATED", nil)

	orderStatus, err = orderUsecase.GetOrderStatus(orderId)

	assert.Equal(t, expectedOrderStatus, orderStatus)
	assert.NoError(t, err)
}

func TestOrderUsecase_UpdateOrderStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	orderRepository := mock_gateways.NewMockOrderRepositoryGateway(ctrl)
	orderNotify := mock_gateways.NewMockOrderNotify(ctrl)
	orderUsecase := NewOrderUsecase(nil, nil, nil, orderNotify, orderRepository)

	type args struct {
		id          int
		orderStatus dto.OrderStatus
	}
	type want struct {
		err error
	}
	type updateOrderStatusCall struct {
		id          int
		orderStatus string
		times       int
		err         error
	}
	type getOrderCall struct {
		id    int
		order entities.Order
		times int
		err   error
	}
	type orderNotifyCall struct {
		productionOrder events.OrderProductionDTO
		times           int
		err             error
	}
	tests := []struct {
		name string
		args
		want
		updateOrderStatusCall
		getOrderCall
		orderNotifyCall
	}{
		{
			name: "should fail to update order status when repository returns error",
			args: args{
				id:          123,
				orderStatus: "CREATED",
			},
			want: want{
				err: errors.New("internal server error"),
			},
			updateOrderStatusCall: updateOrderStatusCall{
				id:          123,
				orderStatus: "CREATED",
				times:       1,
				err:         errors.New("internal server error"),
			},
		},
		{
			name: "should update order status and notify order succesfully",
			args: args{
				id:          123,
				orderStatus: "PAID",
			},
			want: want{
				err: nil,
			},
			updateOrderStatusCall: updateOrderStatusCall{
				id:          123,
				orderStatus: "PAID",
				times:       1,
				err:         nil,
			},
			getOrderCall: getOrderCall{
				id: 123,
				order: entities.Order{
					ID: 123,
					Items: []entities.OrderItem{
						{
							ID:       222,
							Quantity: 1,
							Type:     "UNIT",
							Product: entities.Product{
								ID:          11,
								Name:        "Batata",
								SkuId:       "",
								Description: "Frita",
								Category:    "Acompanhamento",
								Price:       99.99,
								CreatedAt:   time.Time{},
								UpdatedAt:   time.Time{},
							},
						},
					},
					TotalAmount: 99.99,
					CustomerCPF: "123456789",
				},
				times: 1,
				err:   nil,
			},
			orderNotifyCall: orderNotifyCall{
				productionOrder: events.OrderProductionDTO{
					ID:     123,
					Status: "IN_PROGRESS",
					Items: []events.OrderItemProductionDTO{
						{
							Quantity: 1,
							Products: events.OrderProductionProductDTO{
								Name:        "Batata",
								Description: "Frita",
								Category:    "Acompanhamento",
							},
							Type: "UNIT",
						},
					},
				},
				times: 1,
				err:   nil,
			},
		},
		{
			name: "should update order status and failed get order to notify",
			args: args{
				id:          123,
				orderStatus: "PAID",
			},
			want: want{
				err: errors.New("internal server error"),
			},
			updateOrderStatusCall: updateOrderStatusCall{
				id:          123,
				orderStatus: "PAID",
				times:       1,
				err:         nil,
			},
			getOrderCall: getOrderCall{
				id:    123,
				order: entities.Order{},
				times: 1,
				err:   errors.New("internal server error"),
			},
			orderNotifyCall: orderNotifyCall{
				productionOrder: events.OrderProductionDTO{},
				times:           0,
				err:             nil,
			},
		},
		{
			name: "should update order status and failed to notify",
			args: args{
				id:          123,
				orderStatus: "PAID",
			},
			want: want{
				err: errors.New("internal server error"),
			},
			updateOrderStatusCall: updateOrderStatusCall{
				id:          123,
				orderStatus: "PAID",
				times:       1,
				err:         nil,
			},
			getOrderCall: getOrderCall{
				id: 123,
				order: entities.Order{
					ID:     123,
					Status: "PAID",
				},
				times: 1,
				err:   nil,
			},
			orderNotifyCall: orderNotifyCall{
				productionOrder: events.OrderProductionDTO{
					ID:     123,
					Status: "IN_PROGRESS",
					Items:  []events.OrderItemProductionDTO{},
				},
				times: 1,
				err:   errors.New("internal server error"),
			},
		},
	}

	for _, tt := range tests {
		orderRepository.EXPECT().
			UpdateOrderStatus(gomock.Eq(tt.updateOrderStatusCall.id), gomock.Eq(tt.updateOrderStatusCall.orderStatus)).
			Times(tt.updateOrderStatusCall.times).
			Return(tt.updateOrderStatusCall.err)

		orderRepository.EXPECT().
			FindOrderById(gomock.Eq(tt.getOrderCall.id)).
			Times(tt.getOrderCall.times).
			Return(tt.getOrderCall.order, tt.getOrderCall.err)

		orderNotify.EXPECT().
			NotifyPaymentOrder(tt.orderNotifyCall.productionOrder).
			Times(tt.orderNotifyCall.times).
			Return(tt.orderNotifyCall.err)

		err := orderUsecase.UpdateOrderStatus(tt.args.id, tt.args.orderStatus)

		if err != nil {
			assert.EqualError(t, tt.want.err, err.Error())
		} else {
			assert.NoError(t, tt.want.err)
		}
	}
}

func TestOrderUsecase_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	authorizerUsecase := mock_usecases.NewMockAuthorizerUsecase(ctrl)
	paymentUsecase := mock_usecases.NewMockPaymentUsecase(ctrl)
	productUsecase := mock_usecases.NewMockProductUsecase(ctrl)
	orderRepository := mock_gateways.NewMockOrderRepositoryGateway(ctrl)

	orderUsecase := NewOrderUsecase(authorizerUsecase, paymentUsecase, productUsecase, nil, orderRepository)

	type args struct {
		orderDTO dto.OrderDTO
	}
	type want struct {
		orderCreation dto.OrderCreationResponse
		err           error
	}
	type authorizerCall struct {
		cpf   string
		times int
		err   error
	}
	type productUseCaseCall struct {
		id      int
		times   int
		product entities.Product
		err     error
	}
	type repositoryCall struct {
		times   int
		orderId int
		err     error
	}
	type paymentCall struct {
		times  int
		qrcode string
		err    error
	}
	tests := []struct {
		name string
		args
		want
		authorizerCall
		productUseCaseCall
		repositoryCall
		paymentCall
	}{
		{
			name: "should not create order when user is not authorized",
			args: args{
				orderDTO: createOrderDTO(),
			},
			want: want{
				orderCreation: dto.OrderCreationResponse{},
				err:           authorizer.ErrUnauthorized,
			},
			authorizerCall: authorizerCall{
				cpf:   "111222333444",
				times: 1,
				err:   authorizer.ErrUnauthorized,
			},
		},
		{
			name: "should not create order when is not able to calculate products",
			args: args{
				orderDTO: createOrderDTO(),
			},
			want: want{
				orderCreation: dto.OrderCreationResponse{},
				err:           errors.New("internal server error"),
			},
			authorizerCall: authorizerCall{
				cpf:   "111222333444",
				times: 1,
				err:   nil,
			},
			productUseCaseCall: productUseCaseCall{
				id:      222,
				times:   1,
				product: entities.Product{},
				err:     errors.New("internal server error"),
			},
		},
		{
			name: "should not create order when order repository returns error",
			args: args{
				orderDTO: createOrderDTO(),
			},
			want: want{
				orderCreation: dto.OrderCreationResponse{},
				err:           errors.New("internal server error"),
			},
			authorizerCall: authorizerCall{
				cpf:   "111222333444",
				times: 1,
				err:   nil,
			},
			productUseCaseCall: productUseCaseCall{
				id:      222,
				times:   1,
				product: createOrder().Items[0].Product,
				err:     nil,
			},
			repositoryCall: repositoryCall{
				times:   1,
				orderId: -1,
				err:     errors.New("internal server error"),
			},
		},
		{
			name: "should not create order when payment qrcode generation returns error",
			args: args{
				orderDTO: createOrderDTO(),
			},
			want: want{
				orderCreation: dto.OrderCreationResponse{},
				err:           errors.New("internal server error"),
			},
			authorizerCall: authorizerCall{
				cpf:   "111222333444",
				times: 1,
				err:   nil,
			},
			productUseCaseCall: productUseCaseCall{
				id:      222,
				times:   1,
				product: createOrder().Items[0].Product,
				err:     nil,
			},
			repositoryCall: repositoryCall{
				times:   1,
				orderId: 123,
				err:     nil,
			},
			paymentCall: paymentCall{
				times:  1,
				qrcode: "",
				err:    errors.New("internal server error"),
			},
		},
		{
			name: "should create order successfully",
			args: args{
				orderDTO: createOrderDTO(),
			},
			want: want{
				orderCreation: dto.OrderCreationResponse{
					QRCode:  "mercadopago123456",
					OrderID: 123,
				},
				err: nil,
			},
			authorizerCall: authorizerCall{
				cpf:   "111222333444",
				times: 1,
				err:   nil,
			},
			productUseCaseCall: productUseCaseCall{
				id:      222,
				times:   1,
				product: createOrder().Items[0].Product,
				err:     nil,
			},
			repositoryCall: repositoryCall{
				times:   1,
				orderId: 123,
				err:     nil,
			},
			paymentCall: paymentCall{
				times:  1,
				qrcode: "mercadopago123456",
				err:    nil,
			},
		},
	}

	for _, tt := range tests {
		authorizerUsecase.
			EXPECT().
			AuthorizeUser(gomock.Eq(tt.authorizerCall.cpf)).
			Times(tt.authorizerCall.times).
			Return(dto.AuthorizedUser{}, tt.authorizerCall.err)

		productUsecase.
			EXPECT().
			GetProductById(gomock.Eq(tt.productUseCaseCall.id)).
			Times(tt.productUseCaseCall.times).
			Return(tt.productUseCaseCall.product, tt.productUseCaseCall.err)

		orderRepository.
			EXPECT().
			SaveOrder(gomock.Any()).
			Times(tt.repositoryCall.times).
			Return(tt.repositoryCall.orderId, tt.repositoryCall.err)

		paymentUsecase.
			EXPECT().
			GeneratePaymentQRCode(gomock.Any()).
			Times(tt.paymentCall.times).
			Return(tt.paymentCall.qrcode, tt.paymentCall.err)

		orderResp, err := orderUsecase.CreateOrder(tt.args.orderDTO)

		assert.Equal(t, tt.want.orderCreation, orderResp)
		assert.Equal(t, tt.want.err, err)
	}
}

func createOrderDTO() dto.OrderDTO {
	return dto.OrderDTO{
		Items: []dto.OrderItemDTO{
			{
				ProductId: 222,
				Quantity:  1,
				Type:      "UNIT",
			},
		},
		Coupon:      "APP10",
		CustomerCPF: "111222333444",
		Status:      "PAID",
	}
}
