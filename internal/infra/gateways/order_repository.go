package gateways

import (
	"fmt"

	"github.com/g73-techchallenge-order/internal/core/entities"
	"github.com/g73-techchallenge-order/internal/core/usecases/dto"
	"github.com/g73-techchallenge-order/internal/infra/drivers/sql"
	"github.com/g73-techchallenge-order/internal/infra/gateways/sqlscripts"
)

type OrderRepositoryGateway interface {
	FindAllOrders(pageParams dto.PageParams) ([]entities.Order, error)
	GetOrderStatus(orderId int) (string, error)
	SaveOrder(order entities.Order) (int, error)
	UpdateOrderStatus(orderId int, orderStatus string) error
}

type orderRepositoryGateway struct {
	sqlClient sql.SQLClient
}

func NewOrderRepositoryGateway(sqlClient sql.SQLClient) OrderRepositoryGateway {
	return orderRepositoryGateway{
		sqlClient: sqlClient,
	}
}

func (r orderRepositoryGateway) FindAllOrders(pageParams dto.PageParams) ([]entities.Order, error) {
	orders := []entities.Order{}
	err := r.sqlClient.Find(&orders, sqlscripts.FindAllOrdersQuery, pageParams.GetLimit(), pageParams.GetOffset())
	if err != nil {
		return nil, fmt.Errorf("failed to find all orders, error %w", err)
	}

	for i, order := range orders {
		orderItems, err := r.getOrderItems(order.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order items, error %w", err)
		}

		orders[i].Items = orderItems
	}

	return orders, nil
}

func (r orderRepositoryGateway) GetOrderStatus(orderId int) (string, error) {
	var orderStatus string
	err := r.sqlClient.FindOne(&orderStatus, sqlscripts.FindOrderStatusByIdQuery, orderId)
	if err != nil {
		return "", fmt.Errorf("failed to find order status, error %w", err)
	}

	return orderStatus, nil
}

func (r orderRepositoryGateway) SaveOrder(order entities.Order) (int, error) {
	tx, err := r.sqlClient.Begin()
	if err != nil {
		return -1, fmt.Errorf("failed to create a transaction, error %w", err)
	}

	row := tx.ExecWithReturn(sqlscripts.InsertOrderCmd, order.Coupon, order.TotalAmount, order.CustomerCPF, order.Status, order.CreatedAt)

	var orderId int
	err = row.Scan(&orderId)
	if err != nil {
		return -1, fmt.Errorf("failed to save order, error %w", err)
	}

	for _, item := range order.Items {
		_, err := tx.Exec(sqlscripts.InsertOrderItemCmd, orderId, item.Product.ID, item.Quantity, item.Type)
		if err != nil {
			return -1, fmt.Errorf("failed to save order items associations, error %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return -1, fmt.Errorf("failed to commit the transaction, error %w", err)
	}

	return orderId, nil
}

func (r orderRepositoryGateway) UpdateOrderStatus(orderId int, orderStatus string) error {
	result, err := r.sqlClient.Exec(sqlscripts.UpdateOrderStatusCmd, orderId, orderStatus)
	if err != nil {
		return fmt.Errorf("failed to update order status, error %w", err)
	}

	rowsAffect, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check order status update operation, error %w", err)
	}

	if rowsAffect < 1 {
		return sql.ErrNotFound
	}

	return nil
}

func (r orderRepositoryGateway) getOrderItems(orderId int) ([]entities.OrderItem, error) {
	orderItems := []entities.OrderItem{}
	err := r.sqlClient.Find(&orderItems, sqlscripts.FindOrderItems, orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to find order items, error %w", err)
	}

	return orderItems, nil
}
