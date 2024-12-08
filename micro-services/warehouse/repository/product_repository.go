package repository

import (
	"encoding/json"
	"fmt"

	"github.com/parnurzeal/gorequest"
)

type ProductRepository interface {
	GetAllProducts() ([]Product, error)
	UpdateTotalProductStock(productId int64, quantity int) error
}

type productRepository struct {
	baseURL string
}

func NewProductRepository(baseURL string) ProductRepository {
	return &productRepository{
		baseURL: baseURL,
	}
}

type Product struct {
	Id          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func (r *productRepository) GetAllProducts() ([]Product, error) {
	url := fmt.Sprintf("%s/products", r.baseURL)

	request := gorequest.New()
	resp, body, errs := request.Get(url).
		End()

	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to do a request: %v", errs)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error occured: %v", resp.Status)
	}

	var products []Product
	err := json.Unmarshal([]byte(body), &products)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshall JSON: %v", err)
	}

	return products, nil
}

func (r *productRepository) UpdateTotalProductStock(productId int64, quantity int) error {
	url := fmt.Sprintf("%s/products/adjust-total-stock/%d", r.baseURL, productId)

	body := map[string]int{"quantity": quantity}

	request := gorequest.New()
	resp, _, errs := request.Post(url).
		Send(body).
		End()

	if len(errs) > 0 {
		return fmt.Errorf("failed to do a request: %v", errs)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed deduct stock: %v", resp.Status)
	}

	return nil
}
