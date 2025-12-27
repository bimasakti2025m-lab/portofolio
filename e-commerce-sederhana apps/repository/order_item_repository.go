package repository

import (
	"E-commerce-Sederhana/model"
	"database/sql"
)

type OrderItemRepository interface {
	CreateOrderItem(orderItem *model.OrderItem) (*model.OrderItem, error)
	GetAllOrderItems() ([]model.OrderItem, error)
	GetOrderItemByID(id int) (model.OrderItem, error)
	UpdateOrderItem(orderItem *model.OrderItem) (*model.OrderItem, error)
	DeleteOrderItem(id int) error
}

type orderItemRepository struct {
	db *sql.DB
}

func (or *orderItemRepository) CreateOrderItem(orderItem *model.OrderItem) (*model.OrderItem, error) {
	err := or.db.QueryRow("INSERT INTO order_items (order_id, product_id, quantity, price_snapshot) VALUES ($1, $2, $3, $4) RETURNING id", orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.PriceSnapshot).Scan(&orderItem.ID)
	if err != nil {
		return nil, err
	}

	return orderItem, nil
}

func (or *orderItemRepository) GetAllOrderItems() ([]model.OrderItem, error) {
	var orderItems []model.OrderItem
	rows, err := or.db.Query("SELECT id, order_id, product_id, quantity, price_snapshot FROM order_items")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var orderItem model.OrderItem
		if err := rows.Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.ProductID, &orderItem.Quantity, &orderItem.PriceSnapshot); err != nil {
			return nil, err
		}
		orderItems = append(orderItems, orderItem)
	}

	return orderItems, nil
}

func (or *orderItemRepository) GetOrderItemByID(id int) (model.OrderItem, error) {
	var orderItem model.OrderItem
	row := or.db.QueryRow("SELECT id, order_id, product_id, quantity, price_snapshot FROM order_items WHERE id = $1", id)
	if err := row.Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.ProductID, &orderItem.Quantity, &orderItem.PriceSnapshot); err != nil {
		if err == sql.ErrNoRows {
			return model.OrderItem{}, nil
		}
		return model.OrderItem{}, err
	}

	return orderItem, nil

}

func (or *orderItemRepository) UpdateOrderItem(orderItem *model.OrderItem) (*model.OrderItem, error) {
	_, err := or.db.Exec("UPDATE order_items SET order_id = $2, product_id = $3, quantity = $4, price_snapshot = $5 WHERE id = $1", orderItem.ID, orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.PriceSnapshot)
	if err != nil {
		return nil, err
	}

	return orderItem, nil
}

func (or *orderItemRepository) DeleteOrderItem(id int) error {
	_, err := or.db.Exec("DELETE FROM order_items WHERE id = $1", id)
	return err
}

func NewOrderItemRepository(db *sql.DB) OrderItemRepository {
	return &orderItemRepository{db: db}
}
