// Code generated by MockGen. DO NOT EDIT.
// Source: order_repository.go
//
// Generated by this command:
//
//	mockgen -source=order_repository.go -destination=mocks/order_repository.go
//

// Package mock_gateways is a generated GoMock package.
package mock_gateways

import (
	reflect "reflect"

	entities "github.com/g73-techchallenge-order/internal/core/entities"
	dto "github.com/g73-techchallenge-order/internal/core/usecases/dto"
	gomock "go.uber.org/mock/gomock"
)

// MockOrderRepositoryGateway is a mock of OrderRepositoryGateway interface.
type MockOrderRepositoryGateway struct {
	ctrl     *gomock.Controller
	recorder *MockOrderRepositoryGatewayMockRecorder
}

// MockOrderRepositoryGatewayMockRecorder is the mock recorder for MockOrderRepositoryGateway.
type MockOrderRepositoryGatewayMockRecorder struct {
	mock *MockOrderRepositoryGateway
}

// NewMockOrderRepositoryGateway creates a new mock instance.
func NewMockOrderRepositoryGateway(ctrl *gomock.Controller) *MockOrderRepositoryGateway {
	mock := &MockOrderRepositoryGateway{ctrl: ctrl}
	mock.recorder = &MockOrderRepositoryGatewayMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderRepositoryGateway) EXPECT() *MockOrderRepositoryGatewayMockRecorder {
	return m.recorder
}

// FindAllOrders mocks base method.
func (m *MockOrderRepositoryGateway) FindAllOrders(pageParams dto.PageParams) ([]entities.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllOrders", pageParams)
	ret0, _ := ret[0].([]entities.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllOrders indicates an expected call of FindAllOrders.
func (mr *MockOrderRepositoryGatewayMockRecorder) FindAllOrders(pageParams any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllOrders", reflect.TypeOf((*MockOrderRepositoryGateway)(nil).FindAllOrders), pageParams)
}

// GetOrderStatus mocks base method.
func (m *MockOrderRepositoryGateway) GetOrderStatus(orderId int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderStatus", orderId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderStatus indicates an expected call of GetOrderStatus.
func (mr *MockOrderRepositoryGatewayMockRecorder) GetOrderStatus(orderId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderStatus", reflect.TypeOf((*MockOrderRepositoryGateway)(nil).GetOrderStatus), orderId)
}

// SaveOrder mocks base method.
func (m *MockOrderRepositoryGateway) SaveOrder(order entities.Order) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveOrder", order)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveOrder indicates an expected call of SaveOrder.
func (mr *MockOrderRepositoryGatewayMockRecorder) SaveOrder(order any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveOrder", reflect.TypeOf((*MockOrderRepositoryGateway)(nil).SaveOrder), order)
}

// UpdateOrderStatus mocks base method.
func (m *MockOrderRepositoryGateway) UpdateOrderStatus(orderId int, orderStatus string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrderStatus", orderId, orderStatus)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrderStatus indicates an expected call of UpdateOrderStatus.
func (mr *MockOrderRepositoryGatewayMockRecorder) UpdateOrderStatus(orderId, orderStatus any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrderStatus", reflect.TypeOf((*MockOrderRepositoryGateway)(nil).UpdateOrderStatus), orderId, orderStatus)
}
