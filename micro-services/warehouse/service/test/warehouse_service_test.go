package test

import (
	"errors"
	mocks "monorepo-ecommerce/micro-services/warehouse/mocks/mock_micro-services/warehouse/repository"
	"monorepo-ecommerce/micro-services/warehouse/models"
	"monorepo-ecommerce/micro-services/warehouse/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestWarehouseService_AddStock(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockStockRepo := mocks.NewMockStockRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockWarehouseRepo := mocks.NewMockWarehouseRepository(ctrl)

	warehouseService := service.NewWarehouseService(mockWarehouseRepo, mockStockRepo, mockProductRepo)

	productID := int64(1)
	warehouseID := int64(1)
	quantity := 10

	t.Run("should success add stock", func(t *testing.T) {
		warehouses := []models.Warehouse{
			{
				Id:     warehouseID,
				Status: "active",
				Name:   "Warehouse name",
			},
		}

		stock := &models.Stock{Quantity: 20}

		mockStockRepo.EXPECT().
			AddStockToWarehouse(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)

		mockWarehouseRepo.EXPECT().
			GetActiveWarehouses().
			Return(warehouses, nil)

		mockStockRepo.EXPECT().
			GetStockByProductAndWarehouse(gomock.Any(), gomock.Any()).
			Return(stock, nil)

		mockProductRepo.EXPECT().
			UpdateTotalProductStock(gomock.Any(), gomock.Any()).
			Return(nil)

		err := warehouseService.AddStock(productID, warehouseID, quantity)

		assert.NoError(t, err)
	})

	t.Run("should failed to add stock to warehouse", func(t *testing.T) {
		mockStockRepo.EXPECT().
			AddStockToWarehouse(productID, warehouseID, quantity).
			Return(errors.New("database error"))

		err := warehouseService.AddStock(productID, warehouseID, quantity)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to add stock to warehouse")
	})
}

func TestWarehouseService_RemoveStock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStockRepo := mocks.NewMockStockRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockWarehouseRepo := mocks.NewMockWarehouseRepository(ctrl)

	warehouseService := service.NewWarehouseService(mockWarehouseRepo, mockStockRepo, mockProductRepo)

	productID := int64(1)
	warehouseID := int64(1)
	quantity := 5

	t.Run("should success remove stock", func(t *testing.T) {
		warehouses := []models.Warehouse{
			{Id: warehouseID, Status: "active"},
		}

		stock := &models.Stock{Quantity: 15}

		mockStockRepo.EXPECT().
			RemoveStockFromWarehouse(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)

		mockWarehouseRepo.EXPECT().
			GetActiveWarehouses().
			Return(warehouses, nil)

		mockStockRepo.EXPECT().
			GetStockByProductAndWarehouse(gomock.Any(), gomock.Any()).
			Return(stock, nil)

		mockProductRepo.EXPECT().
			UpdateTotalProductStock(gomock.Any(), gomock.Any()).
			Return(nil)

		err := warehouseService.RemoveStock(productID, warehouseID, quantity)

		assert.NoError(t, err)
	})

	t.Run("should failed to remove stock from warehouse", func(t *testing.T) {
		mockStockRepo.EXPECT().
			RemoveStockFromWarehouse(productID, warehouseID, quantity).
			Return(errors.New("database error"))

		err := warehouseService.RemoveStock(productID, warehouseID, quantity)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to remove stock from warehouse")
	})
}

func TestWarehouseService_TransferProduct(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockStockRepo := mocks.NewMockStockRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockWarehouseRepo := mocks.NewMockWarehouseRepository(ctrl)

	warehouseService := service.NewWarehouseService(mockWarehouseRepo, mockStockRepo, mockProductRepo)

	productID := int64(1)
	fromWarehouseID := int64(1)
	toWarehouseID := int64(2)
	quantity := 10

	t.Run("should success transfer product", func(t *testing.T) {
		mockStockRepo.EXPECT().
			RemoveStockFromWarehouse(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)

		mockStockRepo.EXPECT().
			AddStockToWarehouse(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)

		err := warehouseService.TransferProduct(productID, fromWarehouseID, toWarehouseID, quantity)

		assert.NoError(t, err)
	})

	t.Run("should failed to remove stock from source warehouse", func(t *testing.T) {
		mockStockRepo.EXPECT().
			RemoveStockFromWarehouse(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(errors.New("insufficient stock"))

		err := warehouseService.TransferProduct(productID, fromWarehouseID, toWarehouseID, quantity)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to remove stock from source warehouse")
	})
}
