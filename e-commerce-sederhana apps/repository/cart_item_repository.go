package repository

import (
	"E-commerce-Sederhana/model"
	"database/sql"
)

type CartItemRepository interface {
	CreateCartItem(cartItem *model.CartItem) (*model.CartItem, error)
	GetAllCartItems() ([]model.CartItem, error)
	GetCartItemByID(id int) (model.CartItem, error)
	UpdateCartItem(cartItem *model.CartItem) (*model.CartItem, error)
	DeleteCartItem(id int) error
}

type cartItemRepository struct {
	db *sql.DB
}

func (cr *cartItemRepository) CreateCartItem(cartItem *model.CartItem) (*model.CartItem, error) {
	err := cr.db.QueryRow("INSERT INTO cart_items (cart_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING id", cartItem.CartID, cartItem.ProductID, cartItem.Quantity).Scan(&cartItem.ID)
	if err != nil {
		return nil, err
	}

	return cartItem, nil
}

func (cr *cartItemRepository) GetAllCartItems() ([]model.CartItem, error) {
	var cartItems []model.CartItem
	rows, err := cr.db.Query("SELECT id, cart_id, product_id, quantity FROM cart_items")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var cartItem model.CartItem
		if err := rows.Scan(&cartItem.ID, &cartItem.CartID, &cartItem.ProductID, &cartItem.Quantity); err != nil {
			return nil, err
		}
		cartItems = append(cartItems, cartItem)
	}

	return cartItems, nil
}

func (cr *cartItemRepository) GetCartItemByID(id int) (model.CartItem, error) {
	var cartItem model.CartItem
	row := cr.db.QueryRow("SELECT id, cart_id, product_id, quantity FROM cart_items WHERE id = $1", id)
	if err := row.Scan(&cartItem.ID, &cartItem.CartID, &cartItem.ProductID, &cartItem.Quantity); err != nil {
		if err == sql.ErrNoRows {
			return model.CartItem{}, nil
		}
		return model.CartItem{}, err
	}

	return cartItem, nil
}

func (cr *cartItemRepository) UpdateCartItem(cartItem *model.CartItem) (*model.CartItem, error) {
	_, err := cr.db.Exec("UPDATE cart_items SET cart_id = $2, product_id = $3, quantity = $4 WHERE id = $1", cartItem.ID, cartItem.CartID, cartItem.ProductID, cartItem.Quantity)
	if err != nil {
		return nil, err
	}

	return cartItem, nil
}

func (cr *cartItemRepository) DeleteCartItem(id int) error {
	_, err := cr.db.Exec("DELETE FROM cart_items WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func NewCartItemRepository(db *sql.DB) CartItemRepository {
	return &cartItemRepository{db: db}
}
