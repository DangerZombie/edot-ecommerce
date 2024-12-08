package cron

import (
	"log"
	"monorepo-ecommerce/micro-services/order/repository"
	"time"
)

type AutoCancelJob struct {
	OrderRepo   repository.OrderRepository
	ProductRepo repository.ProductRepository
}

func NewAutoCancelJob(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) *AutoCancelJob {
	return &AutoCancelJob{OrderRepo: orderRepo, ProductRepo: productRepo}
}

func (job *AutoCancelJob) Run() {
	// check within 2 minutes orders
	cutoffTime := time.Now().Add(-2 * time.Minute)
	orders, err := job.OrderRepo.GetExpiredOrders("pending", cutoffTime)
	if err != nil {
		log.Printf("Error fetching expired orders: %v", err)
		return
	}

	for _, order := range orders {
		err := job.OrderRepo.UpdateOrderStatus(order.Id, "cancelled")
		if err != nil {
			log.Printf("Failed to cancel order ID %d: %v", order.Id, err)
			continue
		}
		log.Printf("Order Id %d successfully cancelled", order.Id)

		for _, item := range order.Items {
			err = job.ProductRepo.RestoreStock(item.ProductId, item.Quantity)
			if err != nil {
				log.Printf("failed to restore stock for product %d: %v", item.ProductId, err)
				continue
			}
		}
	}
}
