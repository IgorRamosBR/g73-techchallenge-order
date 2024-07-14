package events

type OrderStatusEventDTO struct {
	OrderId int    `json:"orderId"`
	Status  string `json:"status"`
}

type OrderProductionDTO struct {
	ID     int                      `json:"id"`
	Status string                   `json:"status"`
	Items  []OrderItemProductionDTO `json:"items"`
}

type OrderItemProductionDTO struct {
	Quantity int                       `json:"quantity"`
	Products OrderProductionProductDTO `json:"product"`
	Type     string                    `json:"type"`
}

type OrderProductionProductDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}
