package repository

import (
	"fmt"
	"monorepo-ecommerce/micro-services/order/models"

	"github.com/parnurzeal/gorequest"
)

type ShopRepository interface {
	ForwardOrderToShop(order models.Order) error
}

type shopRepository struct {
	baseURL string
}

func NewShopRepository(baseURL string) ShopRepository {
	return &shopRepository{baseURL: baseURL}
}

type ProceedOrderRequest struct {
	OrderID int64                 `json:"order_id"`
	Items   []ProductOrderDetails `json:"items"`
}

type ProductOrderDetails struct {
	ProductId int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

func (r *shopRepository) ForwardOrderToShop(order models.Order) error {
	url := fmt.Sprintf("%s/shop/proceed-order", r.baseURL)

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
