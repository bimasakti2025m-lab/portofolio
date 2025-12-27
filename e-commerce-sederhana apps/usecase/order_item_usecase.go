package usecase

import (
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/repository"
)

type OrderItemUsecase interface {
	CreateOrderItem(orderItem model.OrderItem) (*model.OrderItem, error)
	GetAllOrderItems() ([]model.OrderItem, error)
	GetOrderItemByID(id int) (model.OrderItem, error)
	UpdateOrderItem(orderItem model.OrderItem) (*model.OrderItem, error)
	DeleteOrderItem(id int) error
}

type orderItemUsecase struct {
	repo repository.OrderItemRepository
}

func (ou *orderItemUsecase) CreateOrderItem(orderItem model.OrderItem) (*model.OrderItem, error) {
	return ou.repo.CreateOrderItem(&orderItem)
}

func (ou *orderItemUsecase) GetAllOrderItems() ([]model.OrderItem, error) {
	return ou.repo.GetAllOrderItems()
}

func (ou *orderItemUsecase) GetOrderItemByID(id int) (model.OrderItem, error) {
	return ou.repo.GetOrderItemByID(id)
}

func (ou *orderItemUsecase) UpdateOrderItem(orderItem model.OrderItem) (*model.OrderItem, error) {
	return ou.repo.UpdateOrderItem(&orderItem)
}

func (ou *orderItemUsecase) DeleteOrderItem(id int) error {
	return ou.repo.DeleteOrderItem(id)
}

func NewOrderItemUsecase(repo repository.OrderItemRepository) OrderItemUsecase {
	return &orderItemUsecase{repo: repo}
}