package usecase

import (
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/repository"
)

type ProductUsecase interface {
	GetAllProducts() ([]model.Product, error)
	GetProductByID(id int) (model.Product, error)
	CreateProduct(product model.Product) (*model.Product, error)
	UpdateProduct(product *model.Product) (*model.Product, error)
	DeleteProduct(id int) error
}

type productUsecase struct {
	productRepository repository.ProductRepository
}

func (pu *productUsecase) GetAllProducts() ([]model.Product, error) {
	return pu.productRepository.GetAllProducts()
}

func (pu *productUsecase) GetProductByID(id int) (model.Product, error) {
	return pu.productRepository.GetProductByID(id)
}

func (pu *productUsecase) CreateProduct(product model.Product) (*model.Product, error) {
	return pu.productRepository.CreateProduct(&product)
}

func (pu *productUsecase) UpdateProduct(product *model.Product) (*model.Product, error) {
	return pu.productRepository.UpdateProduct(product)
}

func (pu *productUsecase) DeleteProduct(id int) error {
	return pu.productRepository.DeleteProduct(id)
}

func NewProductUsecase(productRepository repository.ProductRepository) ProductUsecase {
	return &productUsecase{productRepository: productRepository}
}