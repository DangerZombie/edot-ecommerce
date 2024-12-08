package service

import (
	"fmt"
	"monorepo-ecommerce/micro-services/order/models"
	"monorepo-ecommerce/micro-services/order/repository"

	"github.com/labstack/echo/v4"
)

type OrderService interface {
	CreateOrder(c echo.Context, orderRequest *models.OrderRequest) (*models.Order, error)
	ProcessPayment(orderId int64, paid bool) (*models.Order, error)
	CancelOrder(orderId int64) error
	ForwardOrderToShop(order models.Order) error
}

type orderService struct {
	OrderRepo   repository.OrderRepository
	ProductRepo repository.ProductRepository
	ShopRepo    repository.ShopRepository
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository, shopRepo repository.ShopRepository) OrderService {
	return &orderService{
		OrderRepo:   orderRepo,
		ProductRepo: productRepo,
		ShopRepo:    shopRepo,
	}
}

func (s *orderService) CreateOrder(c echo.Context, orderRequest *models.OrderRequest) (*models.Order, error) {
	userId := c.Get("user_id").(int64)
	var totalPrice float64
	var items []models.OrderItem

	for _, itemRequest := range orderRequest.Items {
		// Fetch detail product based on ProductId
		product, err := s.ProductRepo.GetProductStock(itemRequest.ProductId)
		if err != nil {
			return nil, fmt.Errorf("failed fetch product data: %v", err)
		}

		// Check available quantity
		if product.Stock < itemRequest.Quantity {
			return nil, fmt.Errorf("product stock %d not enough", itemRequest.ProductId)
		}

		// Reserve lock and deduction stock
		err = s.ProductRepo.DeductStock(itemRequest.ProductId, itemRequest.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to deduct stock: %v", err)
		}

		totalPrice += float64(itemRequest.Quantity) * product.Price

		items = append(items, models.OrderItem{
			ProductId: itemRequest.ProductId,
			Quantity:  itemRequest.Quantity,
			Price:     product.Price,
		})
	}

	order := &models.Order{
		UserId:     userId,
		Items:      items,
		TotalPrice: totalPrice,
		Status:     "pending",
	}

	createdOrder, err := s.OrderRepo.CreateOrder(order)
	if err != nil {
		return nil, fmt.Errorf("failed create order: %v", err)
	}

	return createdOrder, nil
}

func (s *orderService) ProcessPayment(orderId int64, paid bool) (*models.Order, error) {
	order, err := s.OrderRepo.GetOrderById(orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order: %v", err)
	}

	if order.Status != "pending" {
		return nil, fmt.Errorf("cannot process payment for order with status: %s", order.Status)
	}

	if paid {
		// send data to invoke shop service
		err = s.ForwardOrderToShop(*order)
		if err != nil {
			return nil, err
		}

		err = s.OrderRepo.UpdateOrderStatus(orderId, "success")
		if err != nil {
			return nil, fmt.Errorf("failed to update order status: %v", err)
		}

		order.Status = "success"
	} else {
		err = s.CancelOrder(orderId)
		if err != nil {
			return nil, fmt.Errorf("failed to cancel order: %v", err)
		}

		order.Status = "cancelled"
	}

	return order, nil
}

func (s *orderService) CancelOrder(orderId int64) error {
	order, err := s.OrderRepo.GetOrderById(orderId)
	if err != nil {
		return fmt.Errorf("failed to fetch order: %v", err)
	}

	err = s.OrderRepo.UpdateOrderStatus(orderId, "cancelled")
	if err != nil {
		return fmt.Errorf("failed to update order status: %v", err)
	}

	for _, item := range order.Items {
		err = s.ProductRepo.RestoreStock(item.ProductId, item.Quantity)
		if err != nil {
			return fmt.Errorf("failed to restore stock for product %d: %v", item.ProductId, err)
		}
	}

	return nil
}

func (s *orderService) ForwardOrderToShop(order models.Order) error {
	err := s.ShopRepo.ForwardOrderToShop(order)
	if err != nil {
		return fmt.Errorf("failed to forward order to shop")
	}

	return nil
}
