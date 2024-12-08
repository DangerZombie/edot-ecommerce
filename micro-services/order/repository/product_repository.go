package repository

import (
	"encoding/json"
	"fmt"

	"github.com/parnurzeal/gorequest"
)

type ProductRepository interface {
	GetProductStock(productId int64) (*Product, error)
	DeductStock(productId int64, quantity int) error
	RestoreStock(productId int64, quantity int) error
}

type productRepository struct {
	baseURL string
}

type Product struct {
	Id    int64   `json:"id"`
	Stock int     `json:"stock"`
	Price float64 `json:"price"`
}

func NewProductRepository(baseURL string) ProductRepository {
	return &productRepository{
		baseURL: baseURL,
	}
}

func (r *productRepository) GetProductStock(productId int64) (*Product, error) {
	url := fmt.Sprintf("%s/products/%d", r.baseURL, productId)

	request := gorequest.New()
	resp, body, errs := request.Get(url).
		End()

	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to do a request: %v", errs)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error occured: %v", resp.Status)
	}

	var product Product
	err := json.Unmarshal([]byte(body), &product)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshall JSON: %v", err)
	}

	return &product, nil
}

func (r *productRepository) DeductStock(productId int64, quantity int) error {
	url := fmt.Sprintf("%s/products/deduct/%d", r.baseURL, productId)

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

func (r *productRepository) RestoreStock(productId int64, quantity int) error {
	url := fmt.Sprintf("%s/products/restore/%d", r.baseURL, productId)

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
