package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"monorepo-ecommerce/micro-services/warehouse/models"
)

type StockRepository interface {
	AddStockToWarehouse(productId, warehouseId int64, quantity int) error
	RemoveStockFromWarehouse(productId, warehouseId int64, quantity int) error
	GetStockByProductAndWarehouse(productId, warehouseId int64) (*models.Stock, error)
	UpdateStock(productID, warehouseID int64, newQuantity int) error
}

type stockRepository struct {
	db *sql.DB
}

func NewStockRepository(db *sql.DB) StockRepository {
	return &stockRepository{db: db}
}

func (r *stockRepository) AddStockToWarehouse(productId, warehouseId int64, quantity int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stock, err := r.GetStockByProductAndWarehouse(productId, warehouseId)
	if err != nil {
		tx.Rollback()
		return err
	}

	newStock := stock.Quantity + quantity
	_, err = tx.Exec("UPDATE stocks SET quantity = ? WHERE warehouse_id = ? AND product_id = ?", newStock, warehouseId, productId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *stockRepository) RemoveStockFromWarehouse(productId, warehouseId int64, quantity int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stock, err := r.GetStockByProductAndWarehouse(productId, warehouseId)
	if err != nil {
		tx.Rollback()
		return err
	}

	newStock := stock.Quantity - quantity
	_, err = tx.Exec("UPDATE stocks SET quantity = ? WHERE warehouse_id = ? AND product_id = ?", newStock, warehouseId, productId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *stockRepository) GetStockByProductAndWarehouse(productId, warehouseId int64) (*models.Stock, error) {
	var stock models.Stock
	row := r.db.QueryRow("SELECT id, product_id, warehouse_id, quantity FROM stocks WHERE product_id = ? AND warehouse_id = ?", productId, warehouseId)
	err := row.Scan(&stock.Id, &stock.ProductId, &stock.WarehouseId, &stock.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("stock with product Id %d and warehouse Id %d not found", productId, warehouseId)
		}

		return nil, err
	}

	return &stock, nil
}

func (r *stockRepository) UpdateStock(productID, warehouseID int64, newQuantity int) error {
	query := `UPDATE stocks
              SET quantity = ? 
              WHERE product_id = ? AND warehouse_id = ?`

	result, err := r.db.Exec(query, newQuantity, productID, warehouseID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no stock record found to update")
	}

	return nil
}
