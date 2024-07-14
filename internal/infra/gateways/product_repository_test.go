package gateways

import (
	"errors"
	"testing"
	"time"

	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/entities"
	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/usecases/dto"
	"github.com/IgorRamosBR/g73-techchallenge-order/internal/infra/drivers/sql"
	mock_sql "github.com/IgorRamosBR/g73-techchallenge-order/internal/infra/drivers/sql/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestProductRepositoryGateway_FindAllProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	sqlClient := mock_sql.NewMockSQLClient(ctrl)

	type args struct {
		pageParams dto.PageParams
	}
	type want struct {
		products []entities.Product
		err      error
	}
	type findProductsCall struct {
		limit    int
		offset   int
		times    int
		products []entities.Product
		err      error
	}
	tests := []struct {
		name string
		args
		want
		findProductsCall
	}{
		{
			name: "should fail to find all products when client returns error",
			args: args{
				pageParams: dto.NewPageParams(20, 10),
			},
			want: want{
				products: nil,
				err:      errors.New("failed to find all products, error internal error"),
			},
			findProductsCall: findProductsCall{
				times:  1,
				limit:  10,
				offset: 20,
				err:    errors.New("internal error"),
			},
		},
		{
			name: "should find all products successfully",
			args: args{
				pageParams: dto.NewPageParams(20, 10),
			},
			want: want{
				products: []entities.Product{
					{
						ID:          123,
						Name:        "Product 1",
						SkuId:       "33333",
						Description: "Description of product 1",
						Category:    "Acompanhamento",
						Price:       9.99,
						CreatedAt:   time.Time{},
						UpdatedAt:   time.Time{},
					},
				},
				err: nil,
			},
			findProductsCall: findProductsCall{
				times:  1,
				limit:  10,
				offset: 20,
				products: []entities.Product{
					{
						ID:          123,
						Name:        "Product 1",
						SkuId:       "33333",
						Description: "Description of product 1",
						Category:    "Acompanhamento",
						Price:       9.99,
						CreatedAt:   time.Time{},
						UpdatedAt:   time.Time{},
					},
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		sqlClient.EXPECT().
			Find(gomock.Any(), gomock.Any()).
			SetArg(0, tt.findProductsCall.products).
			Times(tt.findProductsCall.times).
			Return(tt.findProductsCall.err)

		productRepository := NewProductRepositoryGateway(sqlClient)
		products, err := productRepository.FindAllProducts(tt.args.pageParams)

		assert.Equal(t, tt.want.products, products)
		if tt.want.err != nil {
			assert.EqualError(t, err, tt.want.err.Error())
		}
	}
}

func TestProductRepositoryGateway_FindProductsByCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	sqlClient := mock_sql.NewMockSQLClient(ctrl)

	type args struct {
		pageParams dto.PageParams
		category   string
	}
	type want struct {
		products []entities.Product
		err      error
	}
	type findProductsCall struct {
		category string
		times    int
		products []entities.Product
		err      error
	}
	tests := []struct {
		name string
		args
		want
		findProductsCall
	}{
		{
			name: "should fail to find products by category when client returns error",
			args: args{
				pageParams: dto.NewPageParams(20, 10),
				category:   "Acompanhamento",
			},
			want: want{
				products: nil,
				err:      errors.New("failed to find products by category, error internal error"),
			},
			findProductsCall: findProductsCall{
				category: "Acompanhamento",
				times:    1,
				err:      errors.New("internal error"),
			},
		},
		{
			name: "should find all products by category successfully",
			args: args{
				pageParams: dto.NewPageParams(20, 10),
				category:   "Acompanhamento",
			},
			want: want{
				products: []entities.Product{
					{
						ID:          123,
						Name:        "Product 1",
						SkuId:       "33333",
						Description: "Description of product 1",
						Category:    "Acompanhamento",
						Price:       9.99,
						CreatedAt:   time.Time{},
						UpdatedAt:   time.Time{},
					},
				},
				err: nil,
			},
			findProductsCall: findProductsCall{
				times:    1,
				category: "Acompanhamento",
				products: []entities.Product{
					{
						ID:          123,
						Name:        "Product 1",
						SkuId:       "33333",
						Description: "Description of product 1",
						Category:    "Acompanhamento",
						Price:       9.99,
						CreatedAt:   time.Time{},
						UpdatedAt:   time.Time{},
					},
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		sqlClient.EXPECT().
			Find(gomock.Any(), gomock.Any(), gomock.Eq(tt.findProductsCall.category)).
			SetArg(0, tt.findProductsCall.products).
			Times(tt.findProductsCall.times).
			Return(tt.findProductsCall.err)

		productRepository := NewProductRepositoryGateway(sqlClient)
		products, err := productRepository.FindProductsByCategory(tt.args.pageParams, tt.args.category)

		assert.Equal(t, tt.want.products, products)
		if tt.want.err != nil {
			assert.EqualError(t, err, tt.want.err.Error())
		}
	}
}

func TestProductRepositoryGateway_FindProductsById(t *testing.T) {
	ctrl := gomock.NewController(t)
	sqlClient := mock_sql.NewMockSQLClient(ctrl)

	type args struct {
		id int
	}
	type want struct {
		product entities.Product
		err     error
	}
	type findProductCall struct {
		product entities.Product
		id      int
		times   int
		err     error
	}
	tests := []struct {
		name string
		args
		want
		findProductCall
	}{
		{
			name: "should fail to find product by id when client returns error",
			args: args{
				id: 123,
			},
			want: want{
				product: entities.Product{},
				err:     errors.New("failed to find product by id, error internal error"),
			},
			findProductCall: findProductCall{
				id:    123,
				times: 1,
				err:   errors.New("internal error"),
			},
		},
		{
			name: "should find products by id successfully",
			args: args{
				id: 123,
			},
			want: want{
				product: entities.Product{
					ID:          123,
					Name:        "Product 1",
					SkuId:       "33333",
					Description: "Description of product 1",
					Category:    "Acompanhamento",
					Price:       9.99,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
				err: nil,
			},
			findProductCall: findProductCall{
				times: 1,
				id:    123,
				product: entities.Product{
					ID:          123,
					Name:        "Product 1",
					SkuId:       "33333",
					Description: "Description of product 1",
					Category:    "Acompanhamento",
					Price:       9.99,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		sqlClient.EXPECT().
			FindOne(gomock.Any(), gomock.Any(), gomock.Eq(tt.findProductCall.id)).
			SetArg(0, tt.findProductCall.product).
			Times(tt.findProductCall.times).
			Return(tt.findProductCall.err)

		productRepository := NewProductRepositoryGateway(sqlClient)
		product, err := productRepository.FindProductById(tt.args.id)

		assert.Equal(t, tt.want.product, product)
		if tt.want.err != nil {
			assert.EqualError(t, err, tt.want.err.Error())
		}
	}
}

func TestProductRepositoryGateway_SaveProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	sqlClient := mock_sql.NewMockSQLClient(ctrl)

	type args struct {
		product entities.Product
	}
	type want struct {
		err error
	}
	type insertProductCall struct {
		product entities.Product
		times   int
		err     error
	}
	tests := []struct {
		name string
		args
		want
		insertProductCall
	}{
		{
			name: "should fail to save product when client returns error",
			args: args{
				product: entities.Product{
					ID:          123,
					Name:        "Product 1",
					SkuId:       "33333",
					Description: "Description of product 1",
					Category:    "Acompanhamento",
					Price:       9.99,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
			},
			want: want{
				err: errors.New("failed to save product, error internal error"),
			},
			insertProductCall: insertProductCall{
				product: entities.Product{
					ID:          123,
					Name:        "Product 1",
					SkuId:       "33333",
					Description: "Description of product 1",
					Category:    "Acompanhamento",
					Price:       9.99,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
				times: 1,
				err:   errors.New("internal error"),
			},
		},
		{
			name: "should find products by id successfully",
			args: args{
				product: entities.Product{
					ID:          123,
					Name:        "Product 1",
					SkuId:       "33333",
					Description: "Description of product 1",
					Category:    "Acompanhamento",
					Price:       9.99,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
			},
			want: want{
				err: nil,
			},
			insertProductCall: insertProductCall{
				times: 1,
				product: entities.Product{
					ID:          123,
					Name:        "Product 1",
					SkuId:       "33333",
					Description: "Description of product 1",
					Category:    "Acompanhamento",
					Price:       9.99,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		sqlClient.EXPECT().
			Exec(gomock.Any(), gomock.Eq(tt.insertProductCall.product.Name), gomock.Eq(tt.insertProductCall.product.SkuId), gomock.Eq(tt.insertProductCall.product.Description), gomock.Eq(tt.insertProductCall.product.Category), gomock.Eq(tt.insertProductCall.product.Price), gomock.Eq(tt.insertProductCall.product.CreatedAt), gomock.Eq(tt.insertProductCall.product.UpdatedAt)).
			Times(tt.insertProductCall.times).
			Return(nil, tt.insertProductCall.err)

		productRepository := NewProductRepositoryGateway(sqlClient)
		err := productRepository.SaveProduct(tt.args.product)

		if tt.want.err != nil {
			assert.EqualError(t, err, tt.want.err.Error())
		}

	}
}

func TestProductRepositoryGateway_UpdateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	sqlClient := mock_sql.NewMockSQLClient(ctrl)
	result := mock_sql.NewMockResultWrapper(ctrl)

	type args struct {
		id      int
		product entities.Product
	}
	type want struct {
		err error
	}
	type updateProductCall struct {
		id      int
		product entities.Product
		times   int
		result  sql.ResultWrapper
		err     error
	}
	type resultCall struct {
		times        int
		rowsAffected int64
		err          error
	}
	tests := []struct {
		name string
		args
		want
		updateProductCall
		resultCall
	}{
		{
			name: "should fail to update product when client returns error",
			args: args{
				id: 123,
				product: entities.Product{
					Name:        "Product 1",
					SkuId:       "33333",
					Description: "Description of product 1",
					Category:    "Acompanhamento",
					Price:       9.99,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
			},
			want: want{
				err: errors.New("failed to update product [123], error internal error"),
			},
			updateProductCall: updateProductCall{
				id: 123,
				product: entities.Product{
					Name:        "Product 1",
					SkuId:       "33333",
					Description: "Description of product 1",
					Category:    "Acompanhamento",
					Price:       9.99,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
				times: 1,
				err:   errors.New("internal error"),
			},
		},
		{
			name: "should fail to update product when client returns error to get rows affected",
			args: args{
				id: 123,
				product: entities.Product{
					Name:        "Product 1",
					SkuId:       "33333",
					Description: "Description of product 1",
					Category:    "Acompanhamento",
					Price:       9.99,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
			},
			want: want{
				err: errors.New("failed to get rows affected on updating product [123], error internal error"),
			},
			updateProductCall: updateProductCall{
				id: 123,
				product: entities.Product{
					Name:        "Product 1",
					SkuId:       "33333",
					Description: "Description of product 1",
					Category:    "Acompanhamento",
					Price:       9.99,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
				times:  1,
				result: result,
				err:    nil,
			},
			resultCall: resultCall{
				times:        1,
				rowsAffected: 0,
				err:          errors.New("internal error"),
			},
		},
		{
			name: "should fail to update product when client does not foud the product",
			args: args{
				id: 123,
				product: entities.Product{
					Name:        "Product 1",
					SkuId:       "33333",
					Description: "Description of product 1",
					Category:    "Acompanhamento",
					Price:       9.99,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
			},
			want: want{
				err: sql.ErrNotFound,
			},
			updateProductCall: updateProductCall{
				times: 1,
				id:    123,
				product: entities.Product{
					Name:        "Product 1",
					SkuId:       "33333",
					Description: "Description of product 1",
					Category:    "Acompanhamento",
					Price:       9.99,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
				result: result,
				err:    nil,
			},
			resultCall: resultCall{
				times:        1,
				rowsAffected: 0,
				err:          nil,
			},
		},
		{
			name: "should update product successfully",
			args: args{
				id: 123,
				product: entities.Product{
					Name:        "Product 1",
					SkuId:       "33333",
					Description: "Description of product 1",
					Category:    "Acompanhamento",
					Price:       9.99,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
			},
			want: want{
				err: nil,
			},
			updateProductCall: updateProductCall{
				times: 1,
				id:    123,
				product: entities.Product{
					Name:        "Product 1",
					SkuId:       "33333",
					Description: "Description of product 1",
					Category:    "Acompanhamento",
					Price:       9.99,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
				result: result,
				err:    nil,
			},
			resultCall: resultCall{
				times:        1,
				rowsAffected: 1,
				err:          nil,
			},
		},
	}

	for _, tt := range tests {
		sqlClient.EXPECT().
			Exec(gomock.Any(), gomock.Eq(tt.updateProductCall.id), gomock.Eq(tt.updateProductCall.product.Name), gomock.Eq(tt.updateProductCall.product.SkuId), gomock.Eq(tt.updateProductCall.product.Description), gomock.Eq(tt.updateProductCall.product.Category), gomock.Eq(tt.updateProductCall.product.Price), gomock.Eq(tt.updateProductCall.product.UpdatedAt)).
			Times(tt.updateProductCall.times).
			Return(tt.updateProductCall.result, tt.updateProductCall.err)

		result.EXPECT().
			RowsAffected().
			Times(tt.resultCall.times).
			Return(tt.resultCall.rowsAffected, tt.resultCall.err)

		productRepository := NewProductRepositoryGateway(sqlClient)
		err := productRepository.UpdateProduct(tt.args.id, tt.args.product)

		if tt.want.err != nil {
			assert.EqualError(t, err, tt.want.err.Error())
		}

	}
}

func TestProductRepositoryGateway_DeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	sqlClient := mock_sql.NewMockSQLClient(ctrl)
	result := mock_sql.NewMockResultWrapper(ctrl)

	type args struct {
		id int
	}
	type want struct {
		err error
	}
	type deleteProductCall struct {
		id     int
		times  int
		result sql.ResultWrapper
		err    error
	}
	type resultCall struct {
		times        int
		rowsAffected int64
		err          error
	}
	tests := []struct {
		name string
		args
		want
		deleteProductCall
		resultCall
	}{
		{
			name: "should fail to delete product when client returns error",
			args: args{
				id: 123,
			},
			want: want{
				err: errors.New("failed to delete product [123], error internal error"),
			},
			deleteProductCall: deleteProductCall{
				id:    123,
				times: 1,
				err:   errors.New("internal error"),
			},
		},
		{
			name: "should fail to delete product when client returns error to get rows affected",
			args: args{
				id: 123,
			},
			want: want{
				err: errors.New("failed to get rows affected on deleting product [123], error internal error"),
			},
			deleteProductCall: deleteProductCall{
				id:     123,
				times:  1,
				result: result,
				err:    nil,
			},
			resultCall: resultCall{
				times:        1,
				rowsAffected: 0,
				err:          errors.New("internal error"),
			},
		},
		{
			name: "should fail to delete product when client does not foud the product",
			args: args{
				id: 123,
			},
			want: want{
				err: sql.ErrNotFound,
			},
			deleteProductCall: deleteProductCall{
				times:  1,
				id:     123,
				result: result,
				err:    nil,
			},
			resultCall: resultCall{
				times:        1,
				rowsAffected: 0,
				err:          nil,
			},
		},
		{
			name: "should delete product successfully",
			args: args{
				id: 123,
			},
			want: want{
				err: nil,
			},
			deleteProductCall: deleteProductCall{
				times:  1,
				id:     123,
				result: result,
				err:    nil,
			},
			resultCall: resultCall{
				times:        1,
				rowsAffected: 1,
				err:          nil,
			},
		},
	}

	for _, tt := range tests {
		sqlClient.EXPECT().
			Exec(gomock.Any(), gomock.Eq(tt.deleteProductCall.id)).
			Times(tt.deleteProductCall.times).
			Return(tt.deleteProductCall.result, tt.deleteProductCall.err)

		result.EXPECT().
			RowsAffected().
			Times(tt.resultCall.times).
			Return(tt.resultCall.rowsAffected, tt.resultCall.err)

		productRepository := NewProductRepositoryGateway(sqlClient)
		err := productRepository.DeleteProduct(tt.args.id)

		if tt.want.err != nil {
			assert.EqualError(t, err, tt.want.err.Error())
		}

	}
}
