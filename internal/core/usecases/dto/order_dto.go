package dto

import (
	"fmt"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/g73-techchallenge-order/internal/core/entities"
)

type OrderStatus string

const (
	OrderStatusCreated    OrderStatus = "CREATED"
	OrderStatusPaid       OrderStatus = "PAID"
	OrderStatusReceived   OrderStatus = "RECEIVED"
	OrderStatusInProgress OrderStatus = "IN_PROGRESS"
	OrderStatusExpired    OrderStatus = "EXPIRED"
	OrderStatusReady      OrderStatus = "READY"
	OrderStatusDone       OrderStatus = "DONE"
)

type OrderStatusDTO struct {
	Status OrderStatus `json:"status" valid:"in(CREATED|PAID|RECEIVED|IN_PROGRESS|READY|DONE),required~Status is invalid"`
}

func (o OrderStatusDTO) Validate() (bool, error) {
	if _, err := govalidator.ValidateStruct(o); err != nil {
		return false, err
	}

	return true, nil
}

type OrderItemType string

const (
	OrderItemTypeUnit        OrderItemType = "UNIT"
	OrderItemTypeCombo       OrderItemType = "COMBO"
	OrderItemTypeCustomCombo OrderItemType = "CUSTOM_COMBO"
)

type OrderItemDTO struct {
	ProductId int           `json:"productId"`
	Quantity  int           `json:"quantity" valid:"int,required~Quantity is required|range(1|)~Quantity greater than 0"`
	Type      OrderItemType `json:"type" valid:"in(UNIT|COMBO|CUSTOM_COMBO),required~Type is invalid"`
}

func (o OrderItemDTO) toOrderItem() entities.OrderItem {
	return entities.OrderItem{
		Product: entities.Product{
			ID: o.ProductId,
		},
		Quantity: o.Quantity,
		Type:     string(o.Type),
	}
}

type OrderDTO struct {
	Items       []OrderItemDTO `json:"items"`
	Coupon      string         `json:"coupon" valid:"length(0|100)~Description length should be less than 100 characters"`
	CustomerCPF string         `json:"customerCpf"`
	Status      OrderStatus    `json:"status" valid:"in(CREATED|PAID|RECEIVED|IN_PROGRESS|READY|DONE),required~Status is invalid"`
}

func (o OrderDTO) ToOrder() entities.Order {
	orderItems := make([]entities.OrderItem, len(o.Items))
	for i, item := range o.Items {
		orderItems[i] = item.toOrderItem()
	}

	return entities.Order{
		Items:       orderItems,
		Coupon:      o.Coupon,
		CustomerCPF: o.CustomerCPF,
		Status:      string(o.Status),
		CreatedAt:   time.Now(),
	}
}

func (o OrderDTO) ValidateOrder() (bool, error) {
	if _, err := govalidator.ValidateStruct(o); err != nil {
		return false, err
	}

	// Validate CPF using a custom function
	if !isValidCPF(o.CustomerCPF) {
		return false, fmt.Errorf("invalid CPF [%s]", o.CustomerCPF)
	}

	return true, nil
}

func isValidCPF(cpf string) bool {
	cpf = strings.Replace(cpf, ".", "", -1)
	cpf = strings.Replace(cpf, "-", "", -1)

	if len(cpf) != 11 {
		return false
	}

	if strings.Count(cpf, string(cpf[0])) == 11 {
		return false
	}

	// Check if all digits are the same
	if cpf == "00000000000" {
		return false
	}

	// Validate CPF using the standard algorithm
	var sum1, sum2 int
	for i := 0; i < 9; i++ {
		digit := int(cpf[i] - '0')
		sum1 += digit * (10 - i)
		sum2 += digit * (11 - i)
	}

	sum1 %= 11
	if sum1 < 2 {
		sum1 = 0
	} else {
		sum1 = 11 - sum1
	}

	sum2 += sum1 * 2
	sum2 %= 11
	if sum2 < 2 {
		sum2 = 0
	} else {
		sum2 = 11 - sum2
	}

	return cpf[9]-'0' == byte(sum1) && cpf[10]-'0' == byte(sum2)
}
