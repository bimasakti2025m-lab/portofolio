package repository

import (
	"database/sql"
	"fmt"

	"enigmacamp.com/toko-enigma/model"
)

type CartRepository interface {
	Create(todo model.Cart) (model.Cart, error)
	List() ([]model.Cart, error)
	Get(id int) (model.Cart, error)
	Update(todo model.Cart) (model.Cart, error)
	Delete(id int) error
}

type cartRepository struct {
	db *sql.DB
}

func (c *cartRepository) Create(cart model.Cart) (model.Cart, error) {
	// memulai transaksi
	tx, err := c.db.Begin()
	if err != nil {
		return model.Cart{}, err
	}

	// memasukkan data ke tabel mst_cart
	var cartID int
	err = tx.QueryRow("INSERT INTO mst_cart (user_id, total_price) VALUES ($1, $2) RETURNING id", cart.UserId, cart.TotalPrice).Scan(&cartID)
	if err != nil {
		// jika terjadi error saat memasukkan data ke tabel mst_cart, maka rollback transaksi
		tx.Rollback()
		return model.Cart{}, err
	}
	cart.Id = cartID

	// Untuk setiap item dalam keranjang:
	// 1. Buat produknya terlebih dahulu untuk mendapatkan product_id yang valid.
	// 2. Gunakan product_id tersebut untuk membuat cart_item.
	for _, item := range cart.Items {
		var productID int
		// Asumsi: setiap item di keranjang adalah produk baru.
		// Jika produk sudah ada dan hanya ingin mereferensikannya, logikanya perlu diubah.
		err = tx.QueryRow("INSERT INTO mst_product (name, unit, stock, price) VALUES ($1, $2, $3, $4) RETURNING id", item.Product.Name, item.Product.Unit, item.Product.Stock, item.Product.Price).Scan(&productID)
		if err != nil {
			tx.Rollback()
			return model.Cart{}, err
		}

		// Masukkan data ke mst_cart_item dengan cart_id dan product_id yang sudah valid
		_, err = tx.Exec("INSERT INTO mst_cart_item (cart_id, product_id, price, quantity) VALUES ($1, $2, $3, $4)", cartID, productID, item.Price, item.Quantity)
		if err != nil {
			tx.Rollback()
			return model.Cart{}, err
		}
	}

	// jika tidak ada error, maka commit transaksi
	err = tx.Commit()
	if err != nil {
		fmt.Println(err)

		return model.Cart{}, err

	}

	return cart, nil
}

func (c *cartRepository) Get(id int) (model.Cart, error) {
	// memulai transaksi
	tx, err := c.db.Begin()
	if err != nil {
		return model.Cart{}, err
	}

	var cart model.Cart
	err = tx.QueryRow("SELECT id, user_id, total_price FROM mst_cart WHERE id = $1", id).Scan(&cart.Id, &cart.UserId, &cart.TotalPrice)
	if err != nil {
		tx.Rollback()
		return model.Cart{}, err
	}

	rows, err := tx.Query("SELECT id, product_id,cart_id, price, quantity FROM mst_cart_item WHERE cart_id = $1", id)
	if err != nil {
		tx.Rollback()
		return model.Cart{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.CartItem
		err := rows.Scan(&item.Id, &item.ProductId, &item.CartId, &item.Price, &item.Quantity)
		if err != nil {
			tx.Rollback()
			return model.Cart{}, err
		}

		item.Product, err = NewProductRepository(c.db).FindById(item.ProductId)
		if err != nil {
			tx.Rollback()
			return model.Cart{}, err
		}
		cart.Items = append(cart.Items, item)
	}

	rows, err = tx.Query("SELECT id, name, unit, stock, price FROM mst_product")
	if err != nil {
		tx.Rollback()
		return model.Cart{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var product model.Product
		err := rows.Scan(&product.Id, &product.Name, &product.Unit, &product.Stock, &product.Price)
		if err != nil {
			tx.Rollback()
			return model.Cart{}, err
		}

	}

	err = tx.Commit()
	if err != nil {
		return model.Cart{}, err
	}

	return cart, nil
}

func (c *cartRepository) List() ([]model.Cart, error) {
	var carts []model.Cart
	rows, err := c.db.Query("SELECT id, user_id, total_price FROM mst_cart")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cart model.Cart
		if err := rows.Scan(&cart.Id, &cart.UserId, &cart.TotalPrice); err != nil {
			return nil, fmt.Errorf("failed to scan cart: %w", err)
		}
		carts = append(carts, cart)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	// For each cart, get its items
	for i := range carts {
		itemRows, err := c.db.Query("SELECT id, product_id, price, quantity FROM mst_cart_item WHERE cart_id = $1", carts[i].Id)
		if err != nil {
			return nil, fmt.Errorf("failed to query cart items for cart id %d: %w", carts[i].Id, err)
		}
		defer itemRows.Close()

		var items []model.CartItem
		for itemRows.Next() {
			var item model.CartItem
			if err := itemRows.Scan(&item.Id, &item.ProductId, &item.Price, &item.Quantity); err != nil {
				return nil, fmt.Errorf("failed to scan cart item: %w", err)
			}
			items = append(items, item)
		}
		carts[i].Items = items

		productRows, err := c.db.Query("SELECT id, name, unit, stock, price FROM mst_product")
		if err != nil {
			return nil, fmt.Errorf("failed to query products: %w", err)
		}
		defer productRows.Close()

		for productRows.Next() {
			var product model.Product
			if err := productRows.Scan(&product.Id, &product.Name, &product.Unit, &product.Stock, &product.Price); err != nil {
				return nil, fmt.Errorf("failed to scan product: %w", err)
			}
			for j := range items {
				if items[j].ProductId == product.Id {
					items[j].Product = product
				}
			}
		}
	}
	return carts, nil
}

func (c *cartRepository) Update(cart model.Cart) (model.Cart, error) {
	// memulai transaksi
	tx, err := c.db.Begin()
	if err != nil {
		return model.Cart{}, err
	}

	// memperbarui data di tabel mst_cart
	_, err = tx.Exec("UPDATE mst_cart SET user_id = $1, total_price = $2 WHERE id = $3", cart.UserId, cart.TotalPrice, cart.Id)
	if err != nil {
		tx.Rollback()
		return model.Cart{}, err
	}

	// memperbarui data di tabel mst_cart_item
	for _, item := range cart.Items {
		_, err = tx.Exec("UPDATE mst_cart_item SET product_id = $1, price = $2, quantity = $3 WHERE id = $4", item.ProductId, item.Price, item.Quantity, item.Id)
		if err != nil {
			tx.Rollback()
			return model.Cart{}, err
		}
	}

	// memperbarui data di tabel mst_product
	for _, item := range cart.Items {
		_, err = tx.Exec("UPDATE mst_product SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductId)
		if err != nil {
			tx.Rollback()
			return model.Cart{}, err
		}
	}

	// jika tidak ada error, maka commit transaksi
	err = tx.Commit()
	if err != nil {
		return model.Cart{}, err
	}

	return cart, nil
}

func (c *cartRepository) Delete(id int) error {
	// memulai transaksi
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	// menghapus data di tabel mst_cart_item
	_, err = tx.Exec("DELETE FROM mst_cart_item WHERE cart_id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// menghapus data di tabel mst_cart
	_, err = tx.Exec("DELETE FROM mst_cart WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// jika tidak ada error, maka commit transaksi
	return tx.Commit()
}

func NewCartRepository(db *sql.DB) CartRepository {
	return &cartRepository{db: db}
}
