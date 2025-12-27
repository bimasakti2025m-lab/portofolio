package usecase

import (
	"fmt"
	"enigmacamp.com/toko-enigma/model"
	"enigmacamp.com/toko-enigma/repository"
)

type produkUsecase struct {
	repo repository.ProductRepository
}

type ProductUseCase interface {
	CreateNewProduct(book model.Product) (model.Product, error)
	GetAllProduct() ([]model.Product, error)
	GetProductById(id int) (model.Product, error)
	UpdateProductById(book model.Product) (model.Product, error)
	DeleteProductById(id int) error
}

func (p *produkUsecase) CreateNewProduct(product model.Product) (model.Product, error) {
	return p.repo.Create(product)
}

func (p *produkUsecase) GetAllProduct() ([]model.Product, error) {
	return p.repo.FindAll()
}

func (p *produkUsecase) GetProductById(id int) (model.Product, error) {
	return p.repo.FindById(id)
}

func (p *produkUsecase) UpdateProductById(product model.Product) (model.Product, error) {
	_, err := p.repo.FindById(product.Id)

	if err != nil {
		return model.Product{}, fmt.Errorf("Product with ID : %d not found.", product.Id)
	}

	return p.repo.Update(product)
}

func (p *produkUsecase) DeleteProductById(id int) error {
	_, err := p.repo.FindById(id)

	if err != nil {
		return fmt.Errorf("Product with ID : %d not found.", id)
	}
	return p.repo.Delete(id)
}

func NewProductUseCase(repo repository.ProductRepository) ProductUseCase {
	return &produkUsecase{repo: repo}
}
