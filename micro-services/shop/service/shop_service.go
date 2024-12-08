package service

import (
	"fmt"
	"monorepo-ecommerce/micro-services/shop/models"
	"monorepo-ecommerce/micro-services/shop/repository"
)

type ShopService interface {
	ProcessOrder(order models.Order) error
}

type shopService struct {
	ShopRepo      repository.ShopRepository
	WarehouseRepo repository.WarehouseRepository
}

func NewShopService(shopRepo repository.ShopRepository, warehouseRepo repository.WarehouseRepository) ShopService {
	return &shopService{
		ShopRepo:      shopRepo,
		WarehouseRepo: warehouseRepo,
	}
}

func (s *shopService) ProcessOrder(order models.Order) error {
	// Forward request to warehouse
	err := s.WarehouseRepo.ForwardOrderToWarehouse(order)
	if err != nil {
		return fmt.Errorf("failed to forward order to warehouse: %v", err)
	}

	return nil
}
