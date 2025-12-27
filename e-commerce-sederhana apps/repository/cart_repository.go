package repository

import (
	"E-commerce-Sederhana/model"
	"database/sql"
)

type CartRepository interface {
	CreateCart(cart *model.Cart) (*model.Cart, error)
	GetAllCarts() ([]model.Cart, error)
	GetCartByID(id int) (model.Cart, error)
	UpdateCart(cart *model.Cart) (*model.Cart, error)
	DeleteCart(id int) error
}

type cartRepository struct {
	db *sql.DB
}

func (cr *cartRepository) CreateCart(cart *model.Cart) (*model.Cart, error) {
	err := cr.db.QueryRow("INSERT INTO carts (user_id, status) VALUES ($1, $2) RETURNING id", cart.UserID, cart.Status).Scan(&cart.ID)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (cr *cartRepository) GetAllCarts() ([]model.Cart, error) {
	var carts []model.Cart
	rows, err := cr.db.Query("SELECT id, user_id, status FROM carts")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var cart model.Cart
		if err := rows.Scan(&cart.ID, &cart.UserID, &cart.Status); err != nil {
			return nil, err
		}
		carts = append(carts, cart)
	}

	return carts, nil
}

func (cr *cartRepository) GetCartByID(id int) (model.Cart, error) {
	var cart model.Cart
	row := cr.db.QueryRow("SELECT id, user_id, status FROM carts WHERE id = $1", id)
	if err := row.Scan(&cart.ID, &cart.UserID, &cart.Status); err != nil {
		if err == sql.ErrNoRows {
			return model.Cart{}, nil
		}
		return model.Cart{}, err
	}

	return cart, nil
}

func (cr *cartRepository) UpdateCart(cart *model.Cart) (*model.Cart, error) {
	_, err := cr.db.Exec("UPDATE carts SET user_id = $2, status = $3 WHERE id = $1", cart.ID, cart.UserID, cart.Status)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (cr *cartRepository) DeleteCart(id int) error {
	_, err := cr.db.Exec("DELETE FROM carts WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func NewCartRepository(db *sql.DB) CartRepository {
	return &cartRepository{db: db}
}