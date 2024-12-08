package test

import (
	"fmt"
	mocks "monorepo-ecommerce/micro-services/shop/mocks/mock_micro-services/shop/repository"
	"monorepo-ecommerce/micro-services/shop/models"
	"monorepo-ecommerce/micro-services/shop/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestProcessOrder(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockShopRepo := mocks.NewMockShopRepository(ctrl)
	mockWarehouseRepo := mocks.NewMockWarehouseRepository(ctrl)

	shopService := service.NewShopService(mockShopRepo, mockWarehouseRepo)

	order := models.Order{
		Id:     1,
		UserId: 1,
		Items: []models.OrderItem{
			{
				ProductId: 1,
				Quantity:  10,
				Price:     100,
			},
		},
		TotalPrice: 1000,
		Status:     "pending",
	}

	t.Run("should success", func(t *testing.T) {
		mockWarehouseRepo.EXPECT().ForwardOrderToWarehouse(order).Return(nil)

		err := shopService.ProcessOrder(order)

		assert.NoError(t, err)
	})

	t.Run("should failed forwarding", func(t *testing.T) {
		mockWarehouseRepo.EXPECT().ForwardOrderToWarehouse(order).Return(fmt.Errorf("warehouse error"))

		err := shopService.ProcessOrder(order)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to forward order to warehouse: warehouse error")
	})
}
