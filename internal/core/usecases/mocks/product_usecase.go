// Code generated by MockGen. DO NOT EDIT.
// Source: product_usecase.go
//
// Generated by this command:
//
//	mockgen -source=product_usecase.go -destination=mocks/product_usecase.go
//

// Package mock_usecases is a generated GoMock package.
package mock_usecases

import (
	reflect "reflect"

	entities "github.com/g73-techchallenge-order/internal/core/entities"
	dto "github.com/g73-techchallenge-order/internal/core/usecases/dto"
	gomock "go.uber.org/mock/gomock"
)

// MockProductUsecase is a mock of ProductUsecase interface.
type MockProductUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockProductUsecaseMockRecorder
}

// MockProductUsecaseMockRecorder is the mock recorder for MockProductUsecase.
type MockProductUsecaseMockRecorder struct {
	mock *MockProductUsecase
}

// NewMockProductUsecase creates a new mock instance.
func NewMockProductUsecase(ctrl *gomock.Controller) *MockProductUsecase {
	mock := &MockProductUsecase{ctrl: ctrl}
	mock.recorder = &MockProductUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductUsecase) EXPECT() *MockProductUsecaseMockRecorder {
	return m.recorder
}

// CreateProduct mocks base method.
func (m *MockProductUsecase) CreateProduct(productDTO dto.ProductDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProduct", productDTO)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateProduct indicates an expected call of CreateProduct.
func (mr *MockProductUsecaseMockRecorder) CreateProduct(productDTO any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProduct", reflect.TypeOf((*MockProductUsecase)(nil).CreateProduct), productDTO)
}

// DeleteProduct mocks base method.
func (m *MockProductUsecase) DeleteProduct(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProduct", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProduct indicates an expected call of DeleteProduct.
func (mr *MockProductUsecaseMockRecorder) DeleteProduct(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProduct", reflect.TypeOf((*MockProductUsecase)(nil).DeleteProduct), id)
}

// GetAllProducts mocks base method.
func (m *MockProductUsecase) GetAllProducts(pageParameters dto.PageParams) (dto.Page[entities.Product], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllProducts", pageParameters)
	ret0, _ := ret[0].(dto.Page[entities.Product])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllProducts indicates an expected call of GetAllProducts.
func (mr *MockProductUsecaseMockRecorder) GetAllProducts(pageParameters any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllProducts", reflect.TypeOf((*MockProductUsecase)(nil).GetAllProducts), pageParameters)
}

// GetProductById mocks base method.
func (m *MockProductUsecase) GetProductById(id int) (entities.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductById", id)
	ret0, _ := ret[0].(entities.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductById indicates an expected call of GetProductById.
func (mr *MockProductUsecaseMockRecorder) GetProductById(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductById", reflect.TypeOf((*MockProductUsecase)(nil).GetProductById), id)
}

// GetProductsByCategory mocks base method.
func (m *MockProductUsecase) GetProductsByCategory(pageParameters dto.PageParams, category string) (dto.Page[entities.Product], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductsByCategory", pageParameters, category)
	ret0, _ := ret[0].(dto.Page[entities.Product])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductsByCategory indicates an expected call of GetProductsByCategory.
func (mr *MockProductUsecaseMockRecorder) GetProductsByCategory(pageParameters, category any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductsByCategory", reflect.TypeOf((*MockProductUsecase)(nil).GetProductsByCategory), pageParameters, category)
}

// UpdateProduct mocks base method.
func (m *MockProductUsecase) UpdateProduct(id string, productDTO dto.ProductDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProduct", id, productDTO)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProduct indicates an expected call of UpdateProduct.
func (mr *MockProductUsecaseMockRecorder) UpdateProduct(id, productDTO any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProduct", reflect.TypeOf((*MockProductUsecase)(nil).UpdateProduct), id, productDTO)
}