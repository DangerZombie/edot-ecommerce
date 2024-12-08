package service

import (
	"errors"
	"fmt"
	"monorepo-ecommerce/micro-services/warehouse/repository"
)

type WarehouseService interface {
	AddStock(productId, warehouseID int64, quantity int) error
	RemoveStock(productId, warehouseID int64, quantity int) error
	GetTotalStock(productId int64) (int, error)
	TransferProduct(productId int64, fromWarehouseId int64, toWarehouseId int64, quantity int) error
	ActiveDeactiveWarehouseStatus(warehouseId int64) error
	ProceedOrder(orderID int64, items []ProductOrderDetails) error
}

type warehouseService struct {
	warehouseRepo repository.WarehouseRepository
	stockRepo     repository.StockRepository
	productRepo   repository.ProductRepository
}

func NewWarehouseService(warehouseRepo repository.WarehouseRepository, stockRepo repository.StockRepository, productRepo repository.ProductRepository) WarehouseService {
	return &warehouseService{
		warehouseRepo: warehouseRepo,
		stockRepo:     stockRepo,
		productRepo:   productRepo,
	}
}

type ProductOrderDetails struct {
	ProductId int64
	Quantity  int
}

func (s *warehouseService) AddStock(productId, warehouseId int64, quantity int) error {
	err := s.stockRepo.AddStockToWarehouse(productId, warehouseId, quantity)
	if err != nil {
		return fmt.Errorf("failed to add stock to warehouse: %v", err)
	}

	// update product stock
	totalStock, err := s.GetTotalStock(productId)
	if err != nil {
		return fmt.Errorf("failed to fetch total stock: %v", err)
	}

	err = s.productRepo.UpdateTotalProductStock(productId, totalStock)
	if err != nil {
		return fmt.Errorf("failed forward update total product stock: %v", err)
	}

	return nil
}

func (s *warehouseService) RemoveStock(productId, warehouseId int64, quantity int) error {
	err := s.stockRepo.RemoveStockFromWarehouse(productId, warehouseId, quantity)
	if err != nil {
		return fmt.Errorf("failed to remove stock from warehouse: %v", err)
	}

	// update product stock
	totalStock, err := s.GetTotalStock(productId)
	if err != nil {
		return fmt.Errorf("failed to fetch total stock: %v", err)
	}

	err = s.productRepo.UpdateTotalProductStock(productId, totalStock)
	if err != nil {
		return fmt.Errorf("failed forward update total product stock: %v", err)
	}

	return nil
}

func (s *warehouseService) GetTotalStock(productId int64) (int, error) {
	warehouses, err := s.warehouseRepo.GetActiveWarehouses()
	if err != nil {
		return 0, fmt.Errorf("failed to get active warehouses: %v", err)
	}

	totalStock := 0
	for _, warehouse := range warehouses {
		stock, err := s.stockRepo.GetStockByProductAndWarehouse(productId, warehouse.Id)
		if err != nil {
			return 0, fmt.Errorf("failed to get stock for product %d in warehouse %d: %v", productId, warehouse.Id, err)
		}
		totalStock += stock.Quantity
	}

	return totalStock, nil
}

func (s *warehouseService) TransferProduct(productID int64, fromWarehouseID int64, toWarehouseID int64, quantity int) error {
	// Deduct stock from origin warehouse
	err := s.RemoveStock(productID, fromWarehouseID, quantity)
	if err != nil {
		return fmt.Errorf("failed to remove stock from source warehouse: %v", err)
	}

	// Add stock from destination warehouse
	err = s.AddStock(productID, toWarehouseID, quantity)
	if err != nil {
		return fmt.Errorf("failed to add stock to destination warehouse: %v", err)
	}

	return nil
}

func (s *warehouseService) ActiveDeactiveWarehouseStatus(warehouseId int64) error {
	warehouse, err := s.warehouseRepo.GetWarehouseById(warehouseId)
	if err != nil {
		return fmt.Errorf("failed fetch warehouse: %v", err)
	}

	if warehouse.Status == "active" {
		err = s.ActivateWarehouse(warehouse.Id)
		if err != nil {
			return fmt.Errorf("failed activated warehouse: %v", err)
		}
	} else {
		err = s.DeactivateWarehouse(warehouse.Id)
		if err != nil {
			return fmt.Errorf("failed deactivated warehouse: %v", err)
		}
	}

	return nil
}

func (s *warehouseService) ActivateWarehouse(warehouseId int64) error {
	err := s.warehouseRepo.UpdateWarehouseStatus(warehouseId, "active")
	if err != nil {
		return fmt.Errorf("failed to activate warehouse: %v", err)
	}

	return nil
}

func (s *warehouseService) DeactivateWarehouse(warehouseId int64) error {
	err := s.warehouseRepo.UpdateWarehouseStatus(warehouseId, "inactive")
	if err != nil {
		return fmt.Errorf("failed to deactivate warehouse: %v", err)
	}

	return nil
}

func (s *warehouseService) ProceedOrder(orderID int64, products []ProductOrderDetails) error {
	// Retrieve all active warehouses
	warehouses, err := s.warehouseRepo.GetActiveWarehouses()
	if err != nil {
		return err
	}

	for _, product := range products {
		remainingQuantity := product.Quantity

		// Iterate through active warehouses to fulfill the product's stock
		for _, warehouse := range warehouses {
			stock, err := s.stockRepo.GetStockByProductAndWarehouse(product.ProductId, warehouse.Id)
			if err != nil {
				return err
			}

			if stock.Quantity >= remainingQuantity {
				// Deduct remainingQuantity from this warehouse
				err = s.stockRepo.UpdateStock(product.ProductId, warehouse.Id, stock.Quantity-remainingQuantity)
				if err != nil {
					return err
				}
				remainingQuantity = 0
				break
			} else if stock.Quantity > 0 {
				// Deduct as much as possible and continue to the next warehouse
				remainingQuantity -= stock.Quantity
				err = s.stockRepo.UpdateStock(product.ProductId, warehouse.Id, 0)
				if err != nil {
					return err
				}
			}
		}

		// If there is still remaining quantity, return an error for this product
		if remainingQuantity > 0 {
			return errors.New("insufficient stock for product_id: " + string(product.ProductId))
		}
	}

	return nil
}
