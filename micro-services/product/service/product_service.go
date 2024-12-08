package service

import (
	"fmt"
	"monorepo-ecommerce/micro-services/product/models"
	"monorepo-ecommerce/micro-services/product/repository"
)

type ProductService interface {
	GetAllProducts() ([]models.Product, error)
	GetProductById(productId int64) (*models.Product, error)
	DeductStock(productId int64, quantity int) error
	RestoreStock(productId int64, quantity int) error
	UpdateTotalStock(productId int64, quantity int) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetAllProducts() ([]models.Product, error) {
	return s.repo.GetAllProducts()
}

func (s *productService) GetProductById(productId int64) (*models.Product, error) {
	product, err := s.repo.GetProductStock(productId)
	if err != nil {
		return nil, fmt.Errorf("failed fetch product: %v", err)
	}
	return product, nil
}

func (s *productService) DeductStock(productId int64, quantity int) error {
	product, err := s.repo.GetProductStock(productId)
	if err != nil {
		return fmt.Errorf("product not found: %v", err)
	}

	// Validate stock
	if product.Stock < quantity {
		return fmt.Errorf("stock not enough")
	}

	// stock deduction
	newStock := product.Stock - quantity
	err = s.repo.UpdateStock(productId, newStock)
	if err != nil {
		return err
	}

	return nil
}

func (s *productService) RestoreStock(productId int64, quantity int) error {
	product, err := s.repo.GetProductStock(productId)
	if err != nil {
		return fmt.Errorf("product not found: %v", err)
	}

	// stock deduction
	newStock := product.Stock + quantity
	err = s.repo.UpdateStock(productId, newStock)
	if err != nil {
		return err
	}

	return nil
}

func (s *productService) UpdateTotalStock(productId int64, quantity int) error {
	err := s.repo.UpdateStock(productId, quantity)
	if err != nil {
		return err
	}

	return nil
}
