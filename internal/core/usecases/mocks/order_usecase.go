// Code generated by MockGen. DO NOT EDIT.
// Source: order_usecase.go
//
// Generated by this command:
//
//	mockgen -source=order_usecase.go -destination=mocks/order_usecase.go
//
// Package mock_usecases is a generated GoMock package.
package mock_usecases

import (
	reflect "reflect"

	entities "github.com/g73-techchallenge-order/internal/core/entities"
	dto "github.com/g73-techchallenge-order/internal/core/usecases/dto"
	gomock "go.uber.org/mock/gomock"
)

// MockOrderUseCase is a mock of OrderUseCase interface.
type MockOrderUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockOrderUseCaseMockRecorder
}

// MockOrderUseCaseMockRecorder is the mock recorder for MockOrderUseCase.
type MockOrderUseCaseMockRecorder struct {
	mock *MockOrderUseCase
}

// NewMockOrderUseCase creates a new mock instance.
func NewMockOrderUseCase(ctrl *gomock.Controller) *MockOrderUseCase {
	mock := &MockOrderUseCase{ctrl: ctrl}
	mock.recorder = &MockOrderUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderUseCase) EXPECT() *MockOrderUseCaseMockRecorder {
	return m.recorder
}

// CreateOrder mocks base method.
func (m *MockOrderUseCase) CreateOrder(orderDTO dto.OrderDTO) (dto.OrderCreationResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", orderDTO)
	ret0, _ := ret[0].(dto.OrderCreationResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockOrderUseCaseMockRecorder) CreateOrder(orderDTO any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockOrderUseCase)(nil).CreateOrder), orderDTO)
}

// GetAllOrders mocks base method.
func (m *MockOrderUseCase) GetAllOrders(pageParameters dto.PageParams) (dto.Page[entities.Order], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllOrders", pageParameters)
	ret0, _ := ret[0].(dto.Page[entities.Order])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllOrders indicates an expected call of GetAllOrders.
func (mr *MockOrderUseCaseMockRecorder) GetAllOrders(pageParameters any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllOrders", reflect.TypeOf((*MockOrderUseCase)(nil).GetAllOrders), pageParameters)
}

// GetOrderStatus mocks base method.
func (m *MockOrderUseCase) GetOrderStatus(orderId int) (dto.OrderStatusDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderStatus", orderId)
	ret0, _ := ret[0].(dto.OrderStatusDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderStatus indicates an expected call of GetOrderStatus.
func (mr *MockOrderUseCaseMockRecorder) GetOrderStatus(orderId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderStatus", reflect.TypeOf((*MockOrderUseCase)(nil).GetOrderStatus), orderId)
}

// UpdateOrderStatus mocks base method.
func (m *MockOrderUseCase) UpdateOrderStatus(orderId int, orderStatus dto.OrderStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrderStatus", orderId, orderStatus)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrderStatus indicates an expected call of UpdateOrderStatus.
func (mr *MockOrderUseCaseMockRecorder) UpdateOrderStatus(orderId, orderStatus any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrderStatus", reflect.TypeOf((*MockOrderUseCase)(nil).UpdateOrderStatus), orderId, orderStatus)
}
