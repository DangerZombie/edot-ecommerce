package models

type OrderRequest struct {
	Items []OrderItem `json:"items"`
}

type Order struct {
	Id         int64       `json:"id"`
	UserId     int64       `json:"user_id"`
	Items      []OrderItem `json:"items"`
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status"`
}

type OrderItem struct {
	ProductId int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
