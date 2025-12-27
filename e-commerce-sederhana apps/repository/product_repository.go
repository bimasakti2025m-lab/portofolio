package repository

import (
	"E-commerce-Sederhana/model"
	"database/sql"
)

type ProductRepository interface {
	GetAllProducts() ([]model.Product, error)
	GetProductByID(id int) (model.Product, error)
	CreateProduct(product *model.Product) (*model.Product, error)
	UpdateProduct(product *model.Product) (*model.Product, error)
	DeleteProduct(id int) error
}

type productRepository struct {
	db *sql.DB
}

func (pr *productRepository) GetAllProducts() ([]model.Product, error) {
	rows, err := pr.db.Query("SELECT id, name, description, price, stock FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (pr *productRepository) GetProductByID(id int) (model.Product, error) {
	var product model.Product
	err := pr.db.QueryRow("SELECT id, name, description, price, stock FROM products WHERE id = $1", id).
		Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Product{}, nil
		}
		return model.Product{}, err
	}
	return product, nil
}

func (pr *productRepository) CreateProduct(product *model.Product) (*model.Product, error) {
	err := pr.db.QueryRow("INSERT INTO products (name, description, price, stock) VALUES ($1, $2, $3, $4) RETURNING id",
		product.Name, product.Description, product.Price, product.Stock).Scan(&product.ID)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (pr *productRepository) UpdateProduct(product *model.Product) (*model.Product, error) {
	_, err := pr.db.Exec("UPDATE products SET name = $2, description = $3, price = $4, stock = $5 WHERE id = $1",
		product.ID, product.Name, product.Description, product.Price, product.Stock)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (pr *productRepository) DeleteProduct(id int) error {
	_, err := pr.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}
