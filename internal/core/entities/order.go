package entities

import (
	"time"
)

type Order struct {
	ID          int         `json:"id"`
	Items       []OrderItem `json:"items"`
	Coupon      string      `json:"coupon"`
	TotalAmount float64     `json:"totalAmount" db:"total_amount"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"createdAt" db:"created_at"`
	CustomerCPF string      `json:"customerCPF" db:"customer_cpf"`
}

type OrderItem struct {
	ID       int    `json:"id"`
	Quantity int    `json:"quantity"`
	Type     string `json:"type"`
	Product  `json:"product"`
}
