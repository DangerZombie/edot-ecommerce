package repository

import (
	"database/sql"
	"monorepo-ecommerce/micro-services/shop/models"
)

type ShopRepository interface {
	GetAllShops() ([]models.Shop, error)
	GetShopById(id int64) (*models.Shop, error)
}

type shopRepository struct {
	db *sql.DB
}

func NewShopRepository(db *sql.DB) ShopRepository {
	return &shopRepository{db: db}
}

type ProceedOrderRequest struct {
	OrderID int64                 `json:"order_id"`
	Items   []ProductOrderDetails `json:"items"`
}

type ProductOrderDetails struct {
	ProductId int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

func (r *shopRepository) GetAllShops() ([]models.Shop, error) {
	rows, err := r.db.Query("SELECT id, name, description FROM shops")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shops []models.Shop
	for rows.Next() {
		var shop models.Shop
		if err := rows.Scan(&shop.Id, &shop.Name, &shop.Description); err != nil {
			return nil, err
		}
		shops = append(shops, shop)
	}

	return shops, nil
}

func (r *shopRepository) GetShopById(id int64) (*models.Shop, error) {
	row := r.db.QueryRow("SELECT id, name, description FROM shops WHERE id = ?", id)
	var shop models.Shop
	if err := row.Scan(&shop.Id, &shop.Name, &shop.Description); err != nil {
		return nil, err
	}

	return &shop, nil
}
