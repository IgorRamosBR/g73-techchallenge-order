package dto

import "github.com/g73-techchallenge-order/internal/core/entities"

type OrderEventDTO struct {
	OrderId int         `json:"orderId"`
	Status  OrderStatus `json:"status"`
}

type ProductionOrderDTO struct {
	ID     int                      `json:"id"`
	Status string                   `json:"status"`
	Items  []ProductionOrderItemDTO `json:"items"`
}

type ProductionOrderItemDTO struct {
	Quantity int                  `json:"quantity"`
	Products ProductionProductDTO `json:"product"`
	Type     string               `json:"type"`
}

type ProductionProductDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

func ToProductionOrderDTO(order entities.Order) ProductionOrderDTO {
	productionOrder := ProductionOrderDTO{
		ID:     order.ID,
		Status: order.Status,
		Items:  toProductionOrderItemDTO(order.Items),
	}

	return productionOrder
}

func toProductionOrderItemDTO(orderItems []entities.OrderItem) []ProductionOrderItemDTO {
	productionOrderItems := []ProductionOrderItemDTO{}
	for _, item := range orderItems {
		productionOrderItem := ProductionOrderItemDTO{
			Quantity: item.ID,
			Type:     item.Type,
			Products: ProductionProductDTO{
				Name:        item.Product.Name,
				Description: item.Product.Description,
				Category:    item.Product.Category,
			},
		}
		productionOrderItems = append(productionOrderItems, productionOrderItem)
	}

	return productionOrderItems
}
