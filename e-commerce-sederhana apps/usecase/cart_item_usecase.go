package usecase

import (
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/repository"
)

type CartItemUseCase interface {
	CreateCartItem(cartItem model.CartItem) (*model.CartItem, error)
	GetAllCartItems() ([]model.CartItem, error)
	GetCartItemByID(id int) (model.CartItem, error)
	UpdateCartItem(cartItem model.CartItem) (*model.CartItem, error)
	DeleteCartItem(id int) error
}

type cartItemUseCase struct {
	repo repository.CartItemRepository
}

func (cu *cartItemUseCase) CreateCartItem(cartItem model.CartItem) (*model.CartItem, error) {
	return cu.repo.CreateCartItem(&cartItem)
}

func (cu *cartItemUseCase) GetAllCartItems() ([]model.CartItem, error) {
	return cu.repo.GetAllCartItems()
}

func (cu *cartItemUseCase) GetCartItemByID(id int) (model.CartItem, error) {
	return cu.repo.GetCartItemByID(id)
}

func (cu *cartItemUseCase) UpdateCartItem(cartItem model.CartItem) (*model.CartItem, error) {
	return cu.repo.UpdateCartItem(&cartItem)
}

func (cu *cartItemUseCase) DeleteCartItem(id int) error {
	return cu.repo.DeleteCartItem(id)
}

func NewCartItemUseCase(repo repository.CartItemRepository) CartItemUseCase {
	return &cartItemUseCase{repo: repo}
}