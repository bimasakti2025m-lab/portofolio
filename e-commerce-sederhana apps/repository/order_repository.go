package repository

import (
	"E-commerce-Sederhana/model"
	"database/sql"
)

type OrderRepository interface {
	GetAllOrders() ([]model.Order, error)
	GetOrderById(id int) (model.Order, error)
	CreateOrder(order model.Order) (model.Order, error)
	UpdateOrder(order model.Order) (model.Order, error)
	DeleteOrder(id int) error
}

type orderRepository struct {
	db *sql.DB
}

func (or *orderRepository) GetAllOrders() ([]model.Order, error) {
	var orders []model.Order
	rows, err := or.db.Query("SELECT id, user_id, total, status_pesanan, transaction_id_midtrans FROM orders")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.Total, &order.StatusPesanan, &order.TransactionIDMidtrans); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (or *orderRepository) GetOrderById(id int) (model.Order, error) {
	var order model.Order
	row := or.db.QueryRow("SELECT id, user_id, total, status_pesanan, transaction_id_midtrans FROM orders WHERE id = $1", id)
	if err := row.Scan(&order.ID, &order.UserID, &order.Total, &order.StatusPesanan, &order.TransactionIDMidtrans); err != nil {
		if err == sql.ErrNoRows {
			return model.Order{}, nil
		}
		return model.Order{}, err
	}
	return order, nil
}

func (or *orderRepository) CreateOrder(order model.Order) (model.Order, error) {
	err := or.db.QueryRow("INSERT INTO orders (user_id, total, status_pesanan, transaction_id_midtrans) VALUES ($1, $2, $3, $4) RETURNING id", order.UserID, order.Total, order.StatusPesanan, order.TransactionIDMidtrans).Scan(&order.ID)
	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}

func (or *orderRepository) UpdateOrder(order model.Order) (model.Order, error) {
	// Query ini hanya akan mengupdate status pesanan, yang lebih aman untuk alur notifikasi.
	// Jika Anda perlu mengupdate seluruh order, buatlah fungsi terpisah.
	_, err := or.db.Exec("UPDATE orders SET status_pesanan = $2 WHERE id = $1", order.ID, order.StatusPesanan)

	// Query lama yang berpotensi menyebabkan masalah:
	// _, err := or.db.Exec("UPDATE orders SET user_id = $2, total = $3, status_pesanan = $4, transaction_id_midtrans = $5 WHERE id = $1", order.ID, order.UserID, order.Total, order.StatusPesanan, order.TransactionIDMidtrans)
	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}

func (or *orderRepository) DeleteOrder(id int) error {
	_, err := or.db.Exec("DELETE FROM orders WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}
