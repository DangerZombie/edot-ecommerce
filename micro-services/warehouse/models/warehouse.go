package models

type Warehouse struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"` // Active or Inactive
}

type Stock struct {
	Id          int64 `json:"id"`
	WarehouseId int64 `json:"warehouse_id"`
	ProductId   int64 `json:"product_id"`
	Quantity    int   `json:"quantity"`
}