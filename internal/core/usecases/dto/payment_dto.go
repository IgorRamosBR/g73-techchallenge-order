package dto

type PaymentRequest struct {
	OrderId     int                  `json:"orderId"`
	CustomerCpf string               `json:"customerCpf"`
	Items       []PaymentItemRequest `json:"items"`
	TotalAmount float64              `json:"totalAmount"`
}

type PaymentItemRequest struct {
	Quantity int                   `json:"quantity"`
	Product  PaymentProductRequest `json:"product"`
}

type PaymentProductRequest struct {
	Name        string  `json:"name"`
	SkuId       string  `json:"skuId"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Type        string  `json:"type"`
	Price       float64 `json:"price"`
}

type PaymentQRCodeResponse struct {
	QrCode string `json:"qrcode"`
}
