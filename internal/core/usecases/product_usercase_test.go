package usecases

import (
	"errors"
	"testing"
	"time"

	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/entities"
	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/usecases/dto"
	mock_gateways "github.com/IgorRamosBR/g73-techchallenge-order/internal/infra/gateways/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestProductUsecase_GetAllProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	productRepository := mock_gateways.NewMockProductRepositoryGateway(ctrl)

	productUsecase := NewProductUsecase(productRepository)

	pageParams := dto.NewPageParams(20, 10)

	productRepository.EXPECT().
		FindAllProducts(gomock.Eq(pageParams)).
		Times(1).
		Return(nil, errors.New("internal server error"))

	products, err := productUsecase.GetAllProducts(pageParams)

	assert.Empty(t, products)
	assert.EqualError(t, err, "internal server error")

	returnedProducts := []entities.Product{
		{
			Name:        "Product 1",
			SkuId:       "33333",
			Description: "Description of product 1",
			Category:    "Acompanhamento",
			Price:       9.99,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		},
	}
	expectedProducts := dto.Page[entities.Product]{
		Result: returnedProducts,
		Next:   nil,
	}
	productRepository.EXPECT().
		FindAllProducts(gomock.Eq(pageParams)).
		Times(1).
		Return(returnedProducts, nil)

	products, err = productUsecase.GetAllProducts(pageParams)

	assert.Equal(t, expectedProducts, products)
	assert.NoError(t, err)
}

func TestProductUsecase_GetProductsByCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	productRepository := mock_gateways.NewMockProductRepositoryGateway(ctrl)

	productUsecase := NewProductUsecase(productRepository)

	pageParams := dto.NewPageParams(20, 10)
	category := "Acompanhamento"

	productRepository.EXPECT().
		FindProductsByCategory(gomock.Eq(pageParams), gomock.Eq(category)).
		Times(1).
		Return(nil, errors.New("internal server error"))

	products, err := productUsecase.GetProductsByCategory(pageParams, category)

	assert.Empty(t, products)
	assert.EqualError(t, err, "internal server error")

	returnedProducts := []entities.Product{
		{
			Name:        "Product 1",
			SkuId:       "33333",
			Description: "Description of product 1",
			Category:    "Acompanhamento",
			Price:       9.99,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		},
	}
	expectedProducts := dto.Page[entities.Product]{
		Result: returnedProducts,
		Next:   nil,
	}
	productRepository.EXPECT().
		FindProductsByCategory(gomock.Eq(pageParams), gomock.Eq(category)).
		Times(1).
		Return(returnedProducts, nil)

	products, err = productUsecase.GetProductsByCategory(pageParams, category)

	assert.Equal(t, expectedProducts, products)
	assert.NoError(t, err)
}

func TestProductUsecase_GetProductById(t *testing.T) {
	ctrl := gomock.NewController(t)
	productRepository := mock_gateways.NewMockProductRepositoryGateway(ctrl)

	productUsecase := NewProductUsecase(productRepository)

	id := 111

	productRepository.EXPECT().
		FindProductById(gomock.Eq(id)).
		Times(1).
		Return(entities.Product{}, errors.New("internal server error"))

	product, err := productUsecase.GetProductById(id)

	assert.Empty(t, product)
	assert.EqualError(t, err, "internal server error")

	expectedProduct := entities.Product{
		Name:        "Product 1",
		SkuId:       "33333",
		Description: "Description of product 1",
		Category:    "Acompanhamento",
		Price:       9.99,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}

	productRepository.EXPECT().
		FindProductById(gomock.Eq(id)).
		Times(1).
		Return(expectedProduct, nil)

	product, err = productUsecase.GetProductById(id)

	assert.Equal(t, expectedProduct, product)
	assert.NoError(t, err)
}

func TestProductUsecase_CreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	productRepository := mock_gateways.NewMockProductRepositoryGateway(ctrl)

	productUsecase := NewProductUsecase(productRepository)

	product := dto.ProductDTO{
		Name:        "Product 1",
		SkuId:       "33333",
		Description: "Description of product 1",
		Category:    "Acompanhamento",
		Price:       9.99,
	}

	productRepository.EXPECT().
		SaveProduct(gomock.Any()).
		Times(1).
		Return(errors.New("internal server error"))

	err := productUsecase.CreateProduct(product)

	assert.EqualError(t, err, "internal server error")

	productRepository.EXPECT().
		SaveProduct(gomock.Any()).
		Times(1).
		Return(nil)

	err = productUsecase.CreateProduct(product)

	assert.NoError(t, err)
}

func TestProductUsecase_UpdateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	productRepository := mock_gateways.NewMockProductRepositoryGateway(ctrl)

	productUsecase := NewProductUsecase(productRepository)

	type args struct {
		id      string
		product dto.ProductDTO
	}
	type want struct {
		err error
	}
	type repositoryCall struct {
		id    int
		times int
		err   error
	}
	tests := []struct {
		name string
		args
		want
		repositoryCall
	}{
		{
			name: "should fail to update product when id is not a number",
			args: args{
				id: "123abc",
				product: dto.ProductDTO{
					Name:        "Product 1",
					SkuId:       "33333",
					Description: "Description of product 1",
					Category:    "Acompanhamento",
					Price:       9.99,
				},
			},
			want: want{
				err: errors.New("strconv.Atoi: parsing \"123abc\": invalid syntax"),
			},
			repositoryCall: repositoryCall{
				id:    0,
				times: 0,
				err:   nil,
			},
		},
		{
			name: "should fail to update product when irepository returns error",
			args: args{
				id: "1111",
				product: dto.ProductDTO{
					Name:        "Product 1",
					SkuId:       "33333",
					Description: "Description of product 1",
					Category:    "Acompanhamento",
					Price:       9.99,
				},
			},
			want: want{
				err: errors.New("internal server error"),
			},
			repositoryCall: repositoryCall{
				id:    1111,
				times: 1,
				err:   errors.New("internal server error"),
			},
		},
		{
			name: "should update product succesfully",
			args: args{
				id: "1111",
				product: dto.ProductDTO{
					Name:        "Product 1",
					SkuId:       "33333",
					Description: "Description of product 1",
					Category:    "Acompanhamento",
					Price:       9.99,
				},
			},
			want: want{
				err: nil,
			},
			repositoryCall: repositoryCall{
				id:    1111,
				times: 1,
				err:   nil,
			},
		},
	}

	for _, tt := range tests {
		productRepository.EXPECT().
			UpdateProduct(gomock.Eq(tt.repositoryCall.id), gomock.Any()).
			Times(tt.repositoryCall.times).
			Return(tt.repositoryCall.err)

		err := productUsecase.UpdateProduct(tt.args.id, tt.args.product)

		if err != nil {
			assert.EqualError(t, tt.want.err, err.Error())
		} else {
			assert.NoError(t, tt.want.err)
		}
	}
}

func TestProductUsecase_DeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	productRepository := mock_gateways.NewMockProductRepositoryGateway(ctrl)

	productUsecase := NewProductUsecase(productRepository)

	type args struct {
		id string
	}
	type want struct {
		err error
	}
	type repositoryCall struct {
		id    int
		times int
		err   error
	}
	tests := []struct {
		name string
		args
		want
		repositoryCall
	}{
		{
			name: "should fail to delete product when id is not a number",
			args: args{
				id: "123abc",
			},
			want: want{
				err: errors.New("strconv.Atoi: parsing \"123abc\": invalid syntax"),
			},
			repositoryCall: repositoryCall{
				id:    0,
				times: 0,
				err:   nil,
			},
		},
		{
			name: "should fail to delete product when irepository returns error",
			args: args{
				id: "1111",
			},
			want: want{
				err: errors.New("internal server error"),
			},
			repositoryCall: repositoryCall{
				id:    1111,
				times: 1,
				err:   errors.New("internal server error"),
			},
		},
		{
			name: "should delete product succesfully",
			args: args{
				id: "1111",
			},
			want: want{
				err: nil,
			},
			repositoryCall: repositoryCall{
				id:    1111,
				times: 1,
				err:   nil,
			},
		},
	}

	for _, tt := range tests {
		productRepository.EXPECT().
			DeleteProduct(gomock.Eq(tt.repositoryCall.id)).
			Times(tt.repositoryCall.times).
			Return(tt.repositoryCall.err)

		err := productUsecase.DeleteProduct(tt.args.id)

		if err != nil {
			assert.EqualError(t, tt.want.err, err.Error())
		} else {
			assert.NoError(t, tt.want.err)
		}
	}
}
