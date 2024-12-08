package repository

import (
	"database/sql"
	"fmt"
	"monorepo-ecommerce/micro-services/product/models"
)

type ProductRepository interface {
	GetAllProducts() ([]models.Product, error)
	GetProductStock(productId int64) (*models.Product, error)
	UpdateStock(productId int64, quantity int) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetAllProducts() ([]models.Product, error) {
	rows, err := r.db.Query("SELECT id, name, description, price, stock FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Stock); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *productRepository) GetProductStock(productId int64) (*models.Product, error) {
	var product models.Product
	row := r.db.QueryRow("SELECT id, name, description, price, stock FROM products WHERE id = ?", productId)
	err := row.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Stock)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with Id %d not found", productId)
		}

		return nil, err
	}

	return &product, nil
}

func (r *productRepository) UpdateStock(productId int64, newStock int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE products SET stock = ? WHERE id = ?", newStock, productId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
