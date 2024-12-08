package test

import (
	mocks "monorepo-ecommerce/micro-services/product/mocks/mock_micro-services/product/repository"
	"monorepo-ecommerce/micro-services/product/models"
	"monorepo-ecommerce/micro-services/product/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetAllProducts(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockProductRepository(ctrl)

	productService := service.NewProductService(mockRepo)

	mockProducts := []models.Product{
		{Id: 1, Name: "Product 1", Stock: 10, Price: 100},
		{Id: 2, Name: "Product 2", Stock: 20, Price: 200},
	}

	mockRepo.EXPECT().GetAllProducts().Return(mockProducts, nil)

	products, err := productService.GetAllProducts()

	assert.NoError(t, err)
	assert.Equal(t, mockProducts, products)
}

func TestGetProductById(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockProductRepository(ctrl)

	productService := service.NewProductService(mockRepo)

	mockProduct := &models.Product{Id: 1, Name: "Product 1", Stock: 10, Price: 100}

	mockRepo.EXPECT().GetProductStock(int64(1)).Return(mockProduct, nil)

	product, err := productService.GetProductById(1)

	assert.NoError(t, err)
	assert.Equal(t, mockProduct, product)
}

func TestDeductStock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProductRepository(ctrl)

	productService := service.NewProductService(mockRepo)

	mockProduct := &models.Product{Id: 1, Name: "Product 1", Stock: 10, Price: 100}
	mockRepo.EXPECT().GetProductStock(int64(1)).Return(mockProduct, nil)
	mockRepo.EXPECT().UpdateStock(int64(1), 8).Return(nil)

	err := productService.DeductStock(1, 2)

	assert.NoError(t, err)
}

func TestRestoreStock(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockProductRepository(ctrl)

	productService := service.NewProductService(mockRepo)

	mockProduct := &models.Product{Id: 1, Name: "Product 1", Stock: 10, Price: 100}

	mockRepo.EXPECT().GetProductStock(int64(1)).Return(mockProduct, nil)
	mockRepo.EXPECT().UpdateStock(int64(1), 12).Return(nil) // Adding 2 to stock

	err := productService.RestoreStock(1, 2)

	assert.NoError(t, err)
}

func TestUpdateTotalStock(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockProductRepository(ctrl)

	productService := service.NewProductService(mockRepo)

	mockRepo.EXPECT().UpdateStock(int64(1), 15).Return(nil) // Directly setting stock to 15

	err := productService.UpdateTotalStock(1, 15)

	assert.NoError(t, err)
}
