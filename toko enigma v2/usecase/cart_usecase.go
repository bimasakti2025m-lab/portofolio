package usecase

import (
	"fmt"
	"enigmacamp.com/toko-enigma/model"
	"enigmacamp.com/toko-enigma/repository"
)

type cartUsecase struct {
	repo repository.CartRepository
}

type CartUseCase interface {
	CreateNewCart(cart model.Cart) (model.Cart, error)
	GetAllCart() ([]model.Cart, error)
	GetCartById(id int) (model.Cart, error)
	UpdateCartById(cart model.Cart) (model.Cart, error)
	DeleteCartById(id int) error
}

func (c *cartUsecase) CreateNewCart(cart model.Cart) (model.Cart, error) {
	return c.repo.Create(cart)
}

func (c *cartUsecase) GetAllCart() ([]model.Cart, error) {
	return c.repo.List()
}

func (c *cartUsecase) GetCartById(id int) (model.Cart, error) {
	return c.repo.Get(id)
}

func (c *cartUsecase) UpdateCartById(cart model.Cart) (model.Cart, error) {
	_, err := c.repo.Get(cart.Id)

	if err != nil {
		return model.Cart{}, fmt.Errorf("Cart with ID : %d not found.", cart.Id)
	}

	return c.repo.Update(cart)
}

func (c *cartUsecase) DeleteCartById(id int) error {
	_, err := c.repo.Get(id)

	if err != nil {
		return fmt.Errorf("Cart with ID : %d not found.", id)
	}
	return c.repo.Delete(id)
}

func NewCartUseCase(repo repository.CartRepository) CartUseCase {
	return &cartUsecase{repo: repo}
}
