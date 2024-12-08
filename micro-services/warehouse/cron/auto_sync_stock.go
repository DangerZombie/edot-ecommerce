package cron

import (
	"log"
	"monorepo-ecommerce/micro-services/warehouse/repository"
)

type AutoSyncStock struct {
	productRepo   repository.ProductRepository
	stockRepo     repository.StockRepository
	warehouseRepo repository.WarehouseRepository
}

func NewAutoSyncStockJob(productRepo repository.ProductRepository, stockRepo repository.StockRepository, warehouseRepo repository.WarehouseRepository) *AutoSyncStock {
	return &AutoSyncStock{productRepo: productRepo, stockRepo: stockRepo, warehouseRepo: warehouseRepo}
}

func (job *AutoSyncStock) Run() {
	products, err := job.productRepo.GetAllProducts()
	if err != nil {
		log.Printf("Error fetching products: %v", err)
		return
	}

	for _, product := range products {
		warehouses, err := job.warehouseRepo.GetActiveWarehouses()
		if err != nil {
			log.Printf("Error fetching warehouses: %v", err)
			return
		}

		totalStock := 0
		for _, warehouse := range warehouses {
			stock, err := job.stockRepo.GetStockByProductAndWarehouse(product.Id, warehouse.Id)
			if err != nil {
				log.Printf("failed to get stock for product %d in warehouse %d: %v", product.Id, warehouse.Id, err)
				return
			}

			totalStock += stock.Quantity
		}

		// sync product stock
		err = job.productRepo.UpdateTotalProductStock(product.Id, totalStock)
		if err != nil {
			log.Printf("failed forward update total product stock: %v", err)
			return
		}

	}
}
