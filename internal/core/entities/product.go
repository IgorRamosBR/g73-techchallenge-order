package entities

import "time"

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	SkuId       string    `json:"skuId" db:"sku_id"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}
