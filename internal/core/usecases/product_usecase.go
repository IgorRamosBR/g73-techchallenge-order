package usecases

import (
	"strconv"
	"time"

	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/entities"
	"github.com/IgorRamosBR/g73-techchallenge-order/internal/infra/gateways"

	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/usecases/dto"
	log "github.com/sirupsen/logrus"
)

type ProductUsecase interface {
	GetAllProducts(pageParameters dto.PageParams) (dto.Page[entities.Product], error)
	GetProductsByCategory(pageParameters dto.PageParams, category string) (dto.Page[entities.Product], error)
	GetProductById(id int) (entities.Product, error)
	CreateProduct(productDTO dto.ProductDTO) error
	UpdateProduct(id string, productDTO dto.ProductDTO) error
	DeleteProduct(id string) error
}

type productUsecase struct {
	productRepositoryGateway gateways.ProductRepositoryGateway
}

func NewProductUsecase(productRepositoryGateway gateways.ProductRepositoryGateway) ProductUsecase {
	return productUsecase{
		productRepositoryGateway: productRepositoryGateway,
	}
}

func (u productUsecase) GetAllProducts(pageParameters dto.PageParams) (dto.Page[entities.Product], error) {
	products, err := u.productRepositoryGateway.FindAllProducts(pageParameters)
	if err != nil {
		log.Errorf("failed to get all products, error: %v", err)
		return dto.Page[entities.Product]{}, err
	}

	page := dto.BuildPage[entities.Product](products, pageParameters)
	return page, nil
}

func (u productUsecase) GetProductsByCategory(pageParameters dto.PageParams, category string) (dto.Page[entities.Product], error) {
	products, err := u.productRepositoryGateway.FindProductsByCategory(pageParameters, category)
	if err != nil {
		log.Errorf("failed to get products by category, error: %v", err)
		return dto.Page[entities.Product]{}, err
	}

	page := dto.BuildPage[entities.Product](products, pageParameters)
	return page, nil
}

func (u productUsecase) GetProductById(id int) (entities.Product, error) {
	product, err := u.productRepositoryGateway.FindProductById(id)
	if err != nil {
		log.Errorf("failed to get product by id, error: %v", err)
		return entities.Product{}, err
	}

	return product, nil
}

func (u productUsecase) CreateProduct(productDTO dto.ProductDTO) error {
	product := productDTO.ToProduct()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	err := u.productRepositoryGateway.SaveProduct(product)
	if err != nil {
		log.Errorf("failed to save product, error: %v", err)
		return err
	}

	return nil
}

func (u productUsecase) UpdateProduct(idStr string, productDTO dto.ProductDTO) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("failed to parse id [%s], error: %v", idStr, err)
		return err
	}

	product := productDTO.ToProduct()
	product.UpdatedAt = time.Now()
	err = u.productRepositoryGateway.UpdateProduct(id, product)
	if err != nil {
		log.Errorf("failed to update product, error: %v", err)
		return err
	}

	return nil
}

func (u productUsecase) DeleteProduct(idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("failed to parse id [%s], error: %v", idStr, err)
		return err
	}

	err = u.productRepositoryGateway.DeleteProduct(id)
	if err != nil {
		log.Errorf("failed to delete product, error: %v", err)
		return err
	}

	return nil
}
