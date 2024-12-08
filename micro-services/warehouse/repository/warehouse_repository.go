package repository

import (
	"database/sql"
	"fmt"
	"monorepo-ecommerce/micro-services/warehouse/models"
)

type WarehouseRepository interface {
	UpdateWarehouseStatus(warehouseId int64, status string) error
	GetActiveWarehouses() ([]models.Warehouse, error)
	GetWarehouseById(warehouseId int64) (*models.Warehouse, error)
}

type warehouseRepository struct {
	db *sql.DB
}

func NewWarehouseRepository(db *sql.DB) WarehouseRepository {
	return &warehouseRepository{db: db}
}

func (r *warehouseRepository) UpdateWarehouseStatus(warehouseId int64, status string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE warehouses SET status = ? WHERE warehouse_id = ?", status, warehouseId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *warehouseRepository) GetActiveWarehouses() ([]models.Warehouse, error) {
	rows, err := r.db.Query("SELECT id, name, status FROM warehouses WHERE status = ?", "active")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warehouses []models.Warehouse
	for rows.Next() {
		var warehouse models.Warehouse
		if err := rows.Scan(&warehouse.Id, &warehouse.Name, &warehouse.Status); err != nil {
			return nil, err
		}
		warehouses = append(warehouses, warehouse)
	}

	return warehouses, nil
}

func (r *warehouseRepository) GetWarehouseById(warehouseId int64) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	row := r.db.QueryRow("SELECT id, name, status FROM warehouses WHERE id = ?", warehouseId)
	err := row.Scan(&warehouse.Id, &warehouse.Name, &warehouse.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("warehouse Id %d not found", warehouseId)
		}

		return nil, err
	}

	return &warehouse, nil
}
