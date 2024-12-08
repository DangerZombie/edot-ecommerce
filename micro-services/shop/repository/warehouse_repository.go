package repository

import (
	"fmt"
	"monorepo-ecommerce/micro-services/shop/models"

	"github.com/parnurzeal/gorequest"
)

type WarehouseRepository interface {
	ForwardOrderToWarehouse(order models.Order) error
}

type warehouseRepository struct {
	baseURL string
}

func NewWarehouseRepository(baseURL string) WarehouseRepository {
	return &warehouseRepository{baseURL: baseURL}
}

func (r *warehouseRepository) ForwardOrderToWarehouse(order models.Order) error {
	url := fmt.Sprintf("%s/warehouse/stock/proceed-order", r.baseURL)

	requestBody := ProceedOrderRequest{
		OrderID: order.Id,
		Items:   make([]ProductOrderDetails, len(order.Items)),
	}
	for i, item := range order.Items {
		requestBody.Items[i] = ProductOrderDetails{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		}
	}

	request := gorequest.New()
	resp, body, errs := request.Post(url).
		SendStruct(requestBody).
		End()

	if len(errs) > 0 {
		return fmt.Errorf("failed to call warehouse service: %v", errs[0])
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("warehouse service returned error: %s", body)
	}

	return nil
}
