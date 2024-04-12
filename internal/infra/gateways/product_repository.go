package gateways

import (
	"fmt"

	"github.com/g73-techchallenge-order/internal/core/entities"
	"github.com/g73-techchallenge-order/internal/core/usecases/dto"
	"github.com/g73-techchallenge-order/internal/infra/drivers/sql"
	"github.com/g73-techchallenge-order/internal/infra/gateways/sqlscripts"
)

type ProductRepositoryGateway interface {
	FindAllProducts(pageParams dto.PageParams) ([]entities.Product, error)
	FindProductsByCategory(pageParams dto.PageParams, category string) ([]entities.Product, error)
	FindProductById(id int) (entities.Product, error)
	SaveProduct(product entities.Product) error
	UpdateProduct(id int, product entities.Product) error
	DeleteProduct(id int) error
}

type productRepositoryGateway struct {
	sqlClient sql.SQLClient
}

func NewProductRepositoryGateway(sqlClient sql.SQLClient) ProductRepositoryGateway {
	return productRepositoryGateway{
		sqlClient: sqlClient,
	}
}

func (r productRepositoryGateway) FindAllProducts(pageParams dto.PageParams) ([]entities.Product, error) {
	getAllProductsQuery := fmt.Sprintf(sqlscripts.GetAllProductsQuery, pageParams.GetLimit(), pageParams.GetOffset())

	products := []entities.Product{}
	err := r.sqlClient.Find(&products, getAllProductsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to find all products, error %w", err)
	}

	return products, nil
}

func (r productRepositoryGateway) FindProductsByCategory(pageParams dto.PageParams, category string) ([]entities.Product, error) {
	getProductsByCategoryQuery := fmt.Sprintf(sqlscripts.GetProductsByCategoryQuery, pageParams.GetLimit(), pageParams.GetOffset())

	products := []entities.Product{}
	err := r.sqlClient.Find(&products, getProductsByCategoryQuery, category)
	if err != nil {
		return nil, fmt.Errorf("failed to find products by category, error %w", err)
	}

	return products, nil
}

func (r productRepositoryGateway) FindProductById(id int) (entities.Product, error) {
	var product entities.Product
	err := r.sqlClient.FindOne(&product, sqlscripts.GetProductByIdQuery, id)
	if err != nil {
		return entities.Product{}, fmt.Errorf("failed to find product by id, error %w", err)
	}

	return product, nil
}

func (r productRepositoryGateway) SaveProduct(product entities.Product) error {
	inserProductCmd := fmt.Sprintf(sqlscripts.InsertProductCmd)

	_, err := r.sqlClient.Exec(inserProductCmd, product.Name, product.SkuId, product.Description, product.Category,
		product.Price, product.CreatedAt, product.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to save product, error %w", err)
	}

	return nil
}

func (r productRepositoryGateway) UpdateProduct(id int, product entities.Product) error {
	updateProductCmd := fmt.Sprintf(sqlscripts.UpdateProductCmd)

	result, err := r.sqlClient.Exec(updateProductCmd, id, product.Name, product.SkuId, product.Description, product.Category,
		product.Price, product.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update the product [%d], error %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected on updating product [%d], error %w", id, err)
	}

	if rowsAffected < 1 {
		return sql.ErrNotFound
	}

	return nil
}

func (r productRepositoryGateway) DeleteProduct(id int) error {
	deleteProductCmd := fmt.Sprintf(sqlscripts.DeleteProductCmd)

	result, err := r.sqlClient.Exec(deleteProductCmd, id)
	if err != nil {
		return fmt.Errorf("failed to delete the product [%d], error %v", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected on deleting product [%d], error %w", id, err)
	}

	if rowsAffected < 1 {
		return sql.ErrNotFound
	}
	return nil
}
