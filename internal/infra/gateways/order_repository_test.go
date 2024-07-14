package gateways

import (
	"errors"
	"testing"
	"time"

	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/entities"
	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/usecases/dto"
	"github.com/IgorRamosBR/g73-techchallenge-order/internal/infra/drivers/sql"
	mock_sql "github.com/IgorRamosBR/g73-techchallenge-order/internal/infra/drivers/sql/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestOrderRepositoryGateway_FindAllOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	sqlClient := mock_sql.NewMockSQLClient(ctrl)

	type args struct {
		pageParams dto.PageParams
	}
	type want struct {
		orders []entities.Order
		err    error
	}
	type findOrderCall struct {
		limit  int
		offset int
		times  int
		orders []entities.Order
		err    error
	}
	type findOrderItemCall struct {
		orderId    int
		times      int
		orderItems []entities.OrderItem
		err        error
	}
	tests := []struct {
		name string
		args
		want
		findOrderCall
		findOrderItemCall
	}{
		{
			name: "should fail to find all orders when client returns error",
			args: args{
				pageParams: dto.NewPageParams(20, 10),
			},
			want: want{
				orders: nil,
				err:    errors.New("failed to find all orders, error internal error"),
			},
			findOrderCall: findOrderCall{
				times:  1,
				limit:  10,
				offset: 20,
				err:    errors.New("internal error"),
			},
		},
		{
			name: "should fail to find order items when client returns error",
			args: args{
				pageParams: dto.NewPageParams(20, 10),
			},
			want: want{
				orders: nil,
				err:    errors.New("failed to scan order items, error internal error"),
			},
			findOrderCall: findOrderCall{
				times:  1,
				limit:  10,
				offset: 20,
				orders: []entities.Order{{
					ID:          123,
					CustomerCPF: "123456",
				}},
				err: nil,
			},
			findOrderItemCall: findOrderItemCall{
				orderId:    123,
				times:      1,
				orderItems: []entities.OrderItem{},
				err:        errors.New("internal error"),
			},
		},
		{
			name: "should find all orders successfully",
			args: args{
				pageParams: dto.NewPageParams(20, 10),
			},
			want: want{
				orders: []entities.Order{{
					ID:          123,
					CustomerCPF: "123456",
					Items: []entities.OrderItem{
						{
							ID:       999,
							Quantity: 1,
						},
					},
				}},
				err: nil,
			},
			findOrderCall: findOrderCall{
				times:  1,
				limit:  10,
				offset: 20,
				orders: []entities.Order{{
					ID:          123,
					CustomerCPF: "123456",
				}},
				err: nil,
			},
			findOrderItemCall: findOrderItemCall{
				orderId: 123,
				times:   1,
				orderItems: []entities.OrderItem{
					{
						ID:       999,
						Quantity: 1,
					},
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		sqlClient.EXPECT().
			Find(gomock.Any(), gomock.Any(), gomock.Eq(tt.findOrderCall.limit), gomock.Eq(tt.findOrderCall.offset)).
			SetArg(0, tt.findOrderCall.orders).
			Times(tt.findOrderCall.times).
			Return(tt.findOrderCall.err)

		sqlClient.EXPECT().
			Find(gomock.Any(), gomock.Any(), gomock.Eq(tt.findOrderItemCall.orderId)).
			SetArg(0, tt.findOrderItemCall.orderItems).
			Times(tt.findOrderItemCall.times).
			Return(tt.findOrderItemCall.err)

		orderRepository := NewOrderRepositoryGateway(sqlClient)
		orders, err := orderRepository.FindAllOrders(tt.args.pageParams)

		assert.Equal(t, tt.want.orders, orders)
		if tt.want.err != nil {
			assert.EqualError(t, err, tt.want.err.Error())
		}
	}
}

func TestOrderRepositoryGateway_GetOrderStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	sqlClient := mock_sql.NewMockSQLClient(ctrl)

	type args struct {
		orderId int
	}
	type want struct {
		orderStatus string
		err         error
	}
	type findOrderStatusCall struct {
		orderId     int
		times       int
		orderStatus string
		err         error
	}

	tests := []struct {
		name string
		args
		want
		findOrderStatusCall
	}{
		{
			name: "should fail to get order status when client returns error",
			args: args{
				orderId: 123,
			},
			want: want{
				orderStatus: "",
				err:         errors.New("failed to find order status, error internal error"),
			},
			findOrderStatusCall: findOrderStatusCall{
				orderId: 123,
				times:   1,
				err:     errors.New("internal error"),
			},
		},
		{
			name: "should find all orders successfully",
			args: args{
				orderId: 123,
			},
			want: want{
				orderStatus: "PAID",
				err:         nil,
			},
			findOrderStatusCall: findOrderStatusCall{
				orderId:     123,
				orderStatus: "PAID",
				times:       1,
				err:         nil,
			},
		},
	}

	for _, tt := range tests {
		sqlClient.EXPECT().
			FindOne(gomock.Any(), gomock.Any(), gomock.Eq(tt.findOrderStatusCall.orderId)).
			SetArg(0, tt.findOrderStatusCall.orderStatus).
			Times(tt.findOrderStatusCall.times).
			Return(tt.findOrderStatusCall.err)

		orderRepository := NewOrderRepositoryGateway(sqlClient)
		orderStatus, err := orderRepository.GetOrderStatus(tt.args.orderId)

		assert.Equal(t, tt.want.orderStatus, orderStatus)
		if tt.want.err != nil {
			assert.EqualError(t, err, tt.want.err.Error())
		}
	}
}

func TestOrderRepositoryGateway_SaveOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	sqlClient := mock_sql.NewMockSQLClient(ctrl)
	tx := mock_sql.NewMockTransactionWrapper(ctrl)
	row := mock_sql.NewMockRowWrapper(ctrl)
	result := mock_sql.NewMockResultWrapper(ctrl)

	type args struct {
		order entities.Order
	}
	type want struct {
		orderId int
		err     error
	}
	type beginTxCall struct {
		tx    sql.TransactionWrapper
		times int
		err   error
	}
	type rollbackTxCall struct {
		times int
		err   error
	}
	type insertOrderExecCall struct {
		order entities.Order
		times int
		row   sql.RowWrapper
	}
	type insertOrderScanCall struct {
		orderId int
		times   int
		err     error
	}
	type insertOrderItemsExecCall struct {
		orderItem entities.OrderItem
		times     int
		err       error
	}
	type commitTxCall struct {
		times int
		err   error
	}
	tests := []struct {
		name string
		args
		want
		beginTxCall
		rollbackTxCall
		insertOrderExecCall
		insertOrderScanCall
		insertOrderItemsExecCall
		commitTxCall
	}{
		{
			name: "should fail to save orders when client fails to start a transaction",
			args: args{
				order: createOrder(),
			},
			want: want{
				orderId: -1,
				err:     errors.New("failed to create a transaction, error internal server error"),
			},
			beginTxCall: beginTxCall{
				tx:    tx,
				times: 1,
				err:   errors.New("internal server error"),
			},
		},
		{
			name: "should fail to save orders when client fails to scan orderId",
			args: args{
				order: createOrder(),
			},
			want: want{
				orderId: -1,
				err:     errors.New("failed to save order, error internal server error"),
			},
			beginTxCall: beginTxCall{
				tx:    tx,
				times: 1,
				err:   nil,
			},
			rollbackTxCall: rollbackTxCall{
				times: 1,
				err:   nil,
			},
			insertOrderExecCall: insertOrderExecCall{
				order: createOrder(),
				times: 1,
				row:   row,
			},
			insertOrderScanCall: insertOrderScanCall{
				orderId: 0,
				times:   1,
				err:     errors.New("internal server error"),
			},
		},
		{
			name: "should fail to save orders when client fails to save order items",
			args: args{
				order: createOrder(),
			},
			want: want{
				orderId: -1,
				err:     errors.New("failed to save order items associations, error internal server error"),
			},
			beginTxCall: beginTxCall{
				tx:    tx,
				times: 1,
				err:   nil,
			},
			rollbackTxCall: rollbackTxCall{
				times: 1,
				err:   nil,
			},
			insertOrderExecCall: insertOrderExecCall{
				order: createOrder(),
				times: 1,
				row:   row,
			},
			insertOrderScanCall: insertOrderScanCall{
				orderId: 123,
				times:   1,
				err:     nil,
			},
			insertOrderItemsExecCall: insertOrderItemsExecCall{
				orderItem: createOrder().Items[0],
				times:     1,
				err:       errors.New("internal server error"),
			},
		},
		{
			name: "should fail to save orders when client fails to commit transaction",
			args: args{
				order: createOrder(),
			},
			want: want{
				orderId: -1,
				err:     errors.New("failed to commit the transaction, error internal server error"),
			},
			beginTxCall: beginTxCall{
				tx:    tx,
				times: 1,
				err:   nil,
			},
			rollbackTxCall: rollbackTxCall{
				times: 1,
				err:   nil,
			},
			insertOrderExecCall: insertOrderExecCall{
				order: createOrder(),
				times: 1,
				row:   row,
			},
			insertOrderScanCall: insertOrderScanCall{
				orderId: 123,
				times:   1,
				err:     nil,
			},
			insertOrderItemsExecCall: insertOrderItemsExecCall{
				orderItem: createOrder().Items[0],
				times:     1,
				err:       nil,
			},
			commitTxCall: commitTxCall{
				times: 1,
				err:   errors.New("internal server error"),
			},
		},
		{
			name: "should save orders successfully",
			args: args{
				order: createOrder(),
			},
			want: want{
				orderId: 123,
				err:     nil,
			},
			beginTxCall: beginTxCall{
				tx:    tx,
				times: 1,
				err:   nil,
			},
			rollbackTxCall: rollbackTxCall{
				times: 1,
				err:   nil,
			},
			insertOrderExecCall: insertOrderExecCall{
				order: createOrder(),
				times: 1,
				row:   row,
			},
			insertOrderScanCall: insertOrderScanCall{
				orderId: 123,
				times:   1,
				err:     nil,
			},
			insertOrderItemsExecCall: insertOrderItemsExecCall{
				orderItem: createOrder().Items[0],
				times:     1,
				err:       nil,
			},
			commitTxCall: commitTxCall{
				times: 1,
				err:   nil,
			},
		},
	}

	for _, tt := range tests {
		sqlClient.EXPECT().
			Begin().
			Times(tt.beginTxCall.times).
			Return(tt.beginTxCall.tx, tt.beginTxCall.err)

		tx.EXPECT().
			Rollback().
			Times(tt.rollbackTxCall.times).
			Return(tt.rollbackTxCall.err)

		tx.EXPECT().
			ExecWithReturn(gomock.Any(), gomock.Eq(tt.insertOrderExecCall.order.Coupon), gomock.Eq(tt.insertOrderExecCall.order.TotalAmount), gomock.Eq(tt.insertOrderExecCall.order.CustomerCPF), gomock.Eq(tt.insertOrderExecCall.order.Status), gomock.Eq(tt.insertOrderExecCall.order.CreatedAt)).
			Times(tt.insertOrderExecCall.times).
			Return(tt.insertOrderExecCall.row)

		row.EXPECT().
			Scan(gomock.Any()).
			SetArg(0, tt.insertOrderScanCall.orderId).
			Times(tt.insertOrderScanCall.times).
			Return(tt.insertOrderScanCall.err)

		tx.EXPECT().
			Exec(gomock.Any(), gomock.Eq(tt.insertOrderScanCall.orderId), gomock.Eq(tt.insertOrderItemsExecCall.orderItem.Product.ID), gomock.Eq(tt.insertOrderItemsExecCall.orderItem.Quantity), gomock.Eq(tt.insertOrderItemsExecCall.orderItem.Type)).
			Times(tt.insertOrderItemsExecCall.times).
			Return(result, tt.insertOrderItemsExecCall.err)

		tx.EXPECT().
			Commit().
			Times(tt.commitTxCall.times).
			Return(tt.commitTxCall.err)

		orderRepository := NewOrderRepositoryGateway(sqlClient)
		orderId, err := orderRepository.SaveOrder(tt.args.order)

		assert.Equal(t, tt.want.orderId, orderId)
		if tt.want.err != nil {
			assert.EqualError(t, err, tt.want.err.Error())
		}
	}
}

func TestOrderRepositoryGateway_UpdateOrderStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	sqlClient := mock_sql.NewMockSQLClient(ctrl)
	result := mock_sql.NewMockResultWrapper(ctrl)

	type args struct {
		orderId     int
		orderStatus string
	}
	type want struct {
		err error
	}
	type updateOrderStatusExecCall struct {
		orderId     int
		orderStatus string
		times       int
		result      sql.ResultWrapper
		err         error
	}
	type resultCall struct {
		times        int
		rowsAffected int64
		err          error
	}
	tests := []struct {
		name string
		args
		want
		updateOrderStatusExecCall
		resultCall
	}{
		{
			name: "should fail to update order status when client fails to update",
			args: args{
				orderId:     123,
				orderStatus: "PAID",
			},
			want: want{
				err: errors.New("failed to update order status, error internal server error"),
			},
			updateOrderStatusExecCall: updateOrderStatusExecCall{
				orderId:     123,
				orderStatus: "PAID",
				times:       1,
				result:      nil,
				err:         errors.New("internal server error"),
			},
		},
		{
			name: "should fail to update order status when result returns error to check rows affected",
			args: args{
				orderId:     123,
				orderStatus: "PAID",
			},
			want: want{
				err: errors.New("failed to check order status update operation, error internal server error"),
			},
			updateOrderStatusExecCall: updateOrderStatusExecCall{
				orderId:     123,
				orderStatus: "PAID",
				times:       1,
				result:      result,
				err:         nil,
			},
			resultCall: resultCall{
				times:        1,
				rowsAffected: 0,
				err:          errors.New("internal server error"),
			},
		},
		{
			name: "should fail to update order status when order is not found",
			args: args{
				orderId:     123,
				orderStatus: "PAID",
			},
			want: want{
				err: sql.ErrNotFound,
			},
			updateOrderStatusExecCall: updateOrderStatusExecCall{
				orderId:     123,
				orderStatus: "PAID",
				times:       1,
				result:      result,
				err:         nil,
			},
			resultCall: resultCall{
				times:        1,
				rowsAffected: 0,
				err:          nil,
			},
		},
		{
			name: "should update order status successfully",
			args: args{
				orderId:     123,
				orderStatus: "PAID",
			},
			want: want{
				err: nil,
			},
			updateOrderStatusExecCall: updateOrderStatusExecCall{
				orderId:     123,
				orderStatus: "PAID",
				times:       1,
				result:      result,
				err:         nil,
			},
			resultCall: resultCall{
				times:        1,
				rowsAffected: 1,
				err:          nil,
			},
		},
	}

	for _, tt := range tests {
		sqlClient.EXPECT().
			Exec(gomock.Any(), gomock.Eq(tt.updateOrderStatusExecCall.orderId), gomock.Eq(tt.updateOrderStatusExecCall.orderStatus)).
			Times(tt.updateOrderStatusExecCall.times).
			Return(tt.updateOrderStatusExecCall.result, tt.updateOrderStatusExecCall.err)

		result.EXPECT().
			RowsAffected().
			Times(tt.resultCall.times).
			Return(tt.resultCall.rowsAffected, tt.resultCall.err)

		orderRepository := NewOrderRepositoryGateway(sqlClient)
		err := orderRepository.UpdateOrderStatus(tt.args.orderId, tt.args.orderStatus)

		if tt.want.err != nil {
			assert.EqualError(t, err, tt.want.err.Error())
		} else {
			assert.NoError(t, err)
		}
	}
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
