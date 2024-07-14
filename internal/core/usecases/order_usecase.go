package usecases

import (
	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/entities"
	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/usecases/dto"
	"github.com/IgorRamosBR/g73-techchallenge-order/internal/infra/gateways"
	"github.com/IgorRamosBR/g73-techchallenge-order/pkg/events"

	log "github.com/sirupsen/logrus"
)

type OrderUseCase interface {
	GetAllOrders(pageParameters dto.PageParams) (dto.Page[entities.Order], error)
	GetOrderStatus(orderId int) (dto.OrderStatusDTO, error)
	UpdateOrderStatus(orderId int, orderStatus dto.OrderStatus) error
	CreateOrder(orderDTO dto.OrderDTO) (dto.OrderCreationResponse, error)
}

type orderUseCase struct {
	authorizerUsecase AuthorizerUsecase
	paymentUsecase    PaymentUsecase
	productUsecase    ProductUsecase
	orderNotify       gateways.OrderNotify
	orderRepository   gateways.OrderRepositoryGateway
}

type OrderUseCaseConfig struct {
	AuthorizerUsecase      AuthorizerUsecase
	PaymentUseCase         PaymentUsecase
	ProductUseCase         ProductUsecase
	OrderNotify            gateways.OrderNotify
	OrderRepositoryGateway gateways.OrderRepositoryGateway
}

func NewOrderUsecase(authorizerUsecase AuthorizerUsecase, paymentUseCase PaymentUsecase, productUseCase ProductUsecase, orderNotify gateways.OrderNotify, orderRepositoryGateway gateways.OrderRepositoryGateway) OrderUseCase {
	return &orderUseCase{
		authorizerUsecase: authorizerUsecase,
		paymentUsecase:    paymentUseCase,
		productUsecase:    productUseCase,
		orderNotify:       orderNotify,
		orderRepository:   orderRepositoryGateway,
	}
}

func (u *orderUseCase) GetAllOrders(pageParams dto.PageParams) (dto.Page[entities.Order], error) {
	orders, err := u.orderRepository.FindAllOrders(pageParams)
	if err != nil {
		log.Errorf("failed to get all orders, error: %v", err)
		return dto.Page[entities.Order]{}, err
	}

	page := dto.BuildPage[entities.Order](orders, pageParams)
	return page, nil
}

func (u *orderUseCase) GetOrder(orderId int) (entities.Order, error) {
	order, err := u.orderRepository.FindOrderById(orderId)
	if err != nil {
		log.Errorf("failed to get order, error: %v", err)
	}

	return order, nil
}

func (u *orderUseCase) CreateOrder(orderDTO dto.OrderDTO) (dto.OrderCreationResponse, error) {
	// Authorize user
	_, err := u.authorizerUsecase.AuthorizeUser(orderDTO.CustomerCPF)
	if err != nil {
		log.Errorf("failed to authorize customer [%s], error: %v", orderDTO.CustomerCPF, err)
		return dto.OrderCreationResponse{}, err
	}

	// Criar um pedido a partir do DTO
	order := orderDTO.ToOrder()

	// Calcular o total dos produtos
	totalAmount, err := u.calculateProducts(order.Items)
	if err != nil {
		log.Errorf("failed to calculate products, error: %v", err)
		return dto.OrderCreationResponse{}, err
	}

	// Definir o total no pedido
	order.TotalAmount = totalAmount

	// Salvar o pedido no banco de dados
	order.ID, err = u.saveOrder(order)
	if err != nil {
		log.Errorf("failed to save order, error: %v", err)
		return dto.OrderCreationResponse{}, err
	}

	// Gerar o código QR para o pagamento
	paymentQRCode, err := u.paymentUsecase.GeneratePaymentQRCode(order)
	if err != nil {
		log.Errorf("failed to process payment order, error: %v", err)
		return dto.OrderCreationResponse{}, err
	}

	// Construir a resposta com o código QR e o ID do pedido
	response := dto.OrderCreationResponse{
		QRCode:  paymentQRCode,
		OrderID: order.ID,
	}

	return response, nil
}

func (u *orderUseCase) GetOrderStatus(orderId int) (dto.OrderStatusDTO, error) {
	status, err := u.orderRepository.GetOrderStatus(orderId)
	if err != nil {
		return dto.OrderStatusDTO{}, err
	}

	return dto.OrderStatusDTO{
		Status: dto.OrderStatus(status),
	}, nil
}

func (u *orderUseCase) UpdateOrderStatus(orderId int, status dto.OrderStatus) error {
	err := u.orderRepository.UpdateOrderStatus(orderId, string(status))
	if err != nil {
		return err
	}

	if status == dto.OrderStatusPaid {
		return u.NotifyOrderPaid(orderId)
	}

	return nil
}

func (u *orderUseCase) NotifyOrderPaid(orderId int) error {
	order, err := u.GetOrder(orderId)
	if err != nil {
		return err
	}

	productionOrder := ToProductionOrderDTO(order)
	err = u.orderNotify.NotifyPaymentOrder(productionOrder)
	if err != nil {
		return err
	}

	return nil
}

func (u *orderUseCase) calculateProducts(items []entities.OrderItem) (float64, error) {
	for i, item := range items {
		product, err := u.getProduct(item.Product.ID)
		if err != nil {
			log.Errorf("failed to find products to process order, error: %v", err)
			return 0.0, err
		}
		item.Product = product
		items[i] = item
	}

	totalAmount := u.calculateTotal(items)
	return totalAmount, nil
}

func (u *orderUseCase) getProduct(id int) (entities.Product, error) {
	product, err := u.productUsecase.GetProductById(id)
	if err != nil {
		log.Errorf("failed to find product [%d] to process order, error: %v", id, err)
		return entities.Product{}, err
	}

	return product, nil
}

func (u *orderUseCase) calculateTotal(items []entities.OrderItem) float64 {
	var total float64
	for _, item := range items {
		total += item.Product.Price * float64(item.Quantity)
	}
	return total
}

func (u *orderUseCase) saveOrder(order entities.Order) (int, error) {
	orderId, err := u.orderRepository.SaveOrder(order)
	if err != nil {
		return 0, err
	}

	return orderId, nil
}

func ToProductionOrderDTO(order entities.Order) events.OrderProductionDTO {
	productionOrder := events.OrderProductionDTO{
		ID:     order.ID,
		Status: order.Status,
		Items:  toProductionOrderItemDTO(order.Items),
	}

	return productionOrder
}

func toProductionOrderItemDTO(orderItems []entities.OrderItem) []events.OrderItemProductionDTO {
	productionOrderItems := []events.OrderItemProductionDTO{}
	for _, item := range orderItems {
		productionOrderItem := events.OrderItemProductionDTO{
			Quantity: item.ID,
			Type:     item.Type,
			Products: events.OrderProductionProductDTO{
				Name:        item.Product.Name,
				Description: item.Product.Description,
				Category:    item.Product.Category,
			},
		}
		productionOrderItems = append(productionOrderItems, productionOrderItem)
	}

	return productionOrderItems
}
