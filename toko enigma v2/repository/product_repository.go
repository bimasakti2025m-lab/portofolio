package repository

import (
	"database/sql"

	"enigmacamp.com/toko-enigma/model"
)

type ProductRepository interface {
	Create(product model.Product) (model.Product, error)
	FindAll() ([]model.Product, error)
	FindById(id int) (model.Product, error)
	Update(product model.Product) (model.Product, error)
	Delete(id int) error
}

type productRepository struct {
	db *sql.DB
}

// Create implements ProductRepository.
func (p *productRepository) Create(product model.Product) (model.Product, error) {
	var productID int
	err := p.db.QueryRow("INSERT INTO mst_product (name, unit, stock,price) VALUES ($1,$2,$3,$4) RETURNING id", product.Name, product.Unit, product.Stock, product.Price).Scan(&productID)
	if err != nil {
		return model.Product{}, err
	}

	product.Id = productID
	return product, nil
}

// Delete implements ProductRepository.
func (p *productRepository) Delete(id int) error {
	_, err := p.db.Exec("DELETE FROM mst_product WHERE id =$1", id)
	if err != nil {
		return err
	}
	return nil
}

// FindAll implements ProductRepository.
func (p *productRepository) FindAll() ([]model.Product, error) {
	var products []model.Product

	rows, err := p.db.Query("SELECT * FROM mst_product")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product model.Product

		if err := rows.Scan(&product.Id, &product.Name, &product.Unit, &product.Stock, &product.Price); err != nil {
			return nil, err
		}

		products = append(products, product)
	}
	return products, err
}

// FindById implements ProductRepository.
func (p *productRepository) FindById(id int) (model.Product, error) {
	var product model.Product
	err := p.db.QueryRow("SELECT * FROM mst_product WHERE id = $1", id).Scan(&product.Id, &product.Name, &product.Unit, &product.Stock, &product.Price)
	if err != nil {
		return model.Product{}, err
	}
	return product, nil
}

// Update implements ProductRepository.
func (p *productRepository) Update(product model.Product) (model.Product, error) {
	_, err := p.db.Exec("UPDATE mst_product SET id=$2, name=$3, unit=$4, stock =$5, price=$6 WHERE id = $1", product.Id, product.Id, product.Name, product.Unit, product.Stock, product.Price)
	if err != nil {
		return model.Product{}, err
	}
	return product, nil
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}
