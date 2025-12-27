package usecase

import (
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/repository"
)

type CartUseCase interface {
	CreateCart(cart model.Cart) (*model.Cart, error)
	GetAllCarts() ([]model.Cart, error)
	GetCartByID(id int) (model.Cart, error)
	UpdateCart(cart model.Cart) (*model.Cart, error)
	DeleteCart(id int) error
}

type cartUseCase struct {
	repo repository.CartRepository
}

func (cu *cartUseCase) CreateCart(cart model.Cart) (*model.Cart, error) {
	return cu.repo.CreateCart(&cart)
}

func (cu *cartUseCase) GetAllCarts() ([]model.Cart, error) {
	return cu.repo.GetAllCarts()
}

func (cu *cartUseCase) GetCartByID(id int) (model.Cart, error) {
	return cu.repo.GetCartByID(id)
}

func (cu *cartUseCase) UpdateCart(cart model.Cart) (*model.Cart, error) {
	return cu.repo.UpdateCart(&cart)
}

func (cu *cartUseCase) DeleteCart(id int) error {
	return cu.repo.DeleteCart(id)
}

func NewCartUseCase(repo repository.CartRepository) CartUseCase {
	return &cartUseCase{repo: repo}
}