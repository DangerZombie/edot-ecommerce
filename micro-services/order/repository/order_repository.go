package repository

import (
	"database/sql"
	"fmt"
	"monorepo-ecommerce/micro-services/order/models"
	"time"
)

type OrderRepository interface {
	CreateOrder(order *models.Order) (*models.Order, error)
	GetOrderById(orderId int64) (*models.Order, error)
	UpdateOrderStatus(orderId int64, status string) error
	GetExpiredOrders(status string, cutoffTime time.Time) ([]models.Order, error)
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(order *models.Order) (*models.Order, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed begin transaction: %v", err)
	}

	orderQuery := "INSERT INTO orders (user_id, total_price, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	result, err := tx.Exec(orderQuery, order.UserId, order.TotalPrice, order.Status, time.Now(), time.Now())
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed insert order: %v", err)
	}

	orderId, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed retreive Id order: %v", err)
	}

	for _, item := range order.Items {
		itemQuery := "INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)"
		_, err := tx.Exec(itemQuery, orderId, item.ProductId, item.Quantity, item.Price)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed insert item order: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed commit transaction: %v", err)
	}

	order.Id = orderId
	return order, nil
}

func (r *orderRepository) GetOrderById(orderId int64) (*models.Order, error) {
	var order models.Order
	row := r.db.QueryRow("SELECT id, user_id, status, total_price FROM orders WHERE id = ?", orderId)
	err := row.Scan(&order.Id, &order.UserId, &order.Status, &order.TotalPrice)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order with Id %d not found", orderId)
		}

		return nil, err
	}

	items, err := r.getOrderItems(order.Id)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return &order, nil
}

func (r *orderRepository) UpdateOrderStatus(orderId int64, status string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = r.GetOrderById(orderId)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE orders SET status = ? WHERE id = ?", status, orderId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) GetExpiredOrders(status string, cutoffTime time.Time) ([]models.Order, error) {
	rows, err := r.db.Query("SELECT id, user_id, total_price, status FROM orders WHERE status = ? AND created_at < ?", status, cutoffTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.Id, &order.UserId, &order.TotalPrice, &order.Status); err != nil {
			return nil, err
		}
		orders = append(orders, order)

		items, err := r.getOrderItems(order.Id)
		if err != nil {
			return nil, err
		}
		order.Items = items

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *orderRepository) getOrderItems(orderId int64) ([]models.OrderItem, error) {
	rows, err := r.db.Query("SELECT oi.id, oi.product_id, oi.quantity, oi.price FROM order_items oi WHERE oi.order_id = ?", orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.Id, &item.ProductId, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
