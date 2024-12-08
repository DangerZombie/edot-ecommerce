package test

import (
	mocks "monorepo-ecommerce/micro-services/order/mocks/mock_micro-services/order/repository"
	"monorepo-ecommerce/micro-services/order/models"
	"monorepo-ecommerce/micro-services/order/repository"
	"monorepo-ecommerce/micro-services/order/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockShopRepo := mocks.NewMockShopRepository(ctrl)

	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo, mockShopRepo)

	orderRequest := &models.OrderRequest{
		Items: []models.OrderItem{
			{
				ProductId: 1,
				Quantity:  2,
			},
		},
	}

	getProductStock := &repository.Product{
		Id:    1,
		Stock: 10,
		Price: 100,
	}

	createOrder := &models.Order{
		UserId:     1,
		TotalPrice: 200,
		Status:     "pending",
		Items: []models.OrderItem{
			{
				ProductId: 1,
				Quantity:  2,
				Price:     100,
			},
		},
	}

	mockProductRepo.EXPECT().
		GetProductStock(gomock.Any()).
		Return(getProductStock, nil)
	mockProductRepo.EXPECT().
		DeductStock(gomock.Any(), gomock.Any()).
		Return(nil)

	mockOrderRepo.EXPECT().
		CreateOrder(gomock.Any()).
		Return(createOrder, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/orders", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", int64(1))

	order, err := orderService.CreateOrder(c, orderRequest)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, float64(200), order.TotalPrice)
	assert.Equal(t, "pending", order.Status)
}

func TestProcessPayment(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockShopRepo := mocks.NewMockShopRepository(ctrl)

	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo, mockShopRepo)

	order := &models.Order{
		Id:         1,
		Status:     "pending",
		UserId:     1,
		TotalPrice: 100,
	}

	mockOrderRepo.EXPECT().
		GetOrderById(int64(1)).
		Return(order, nil)

	mockShopRepo.EXPECT().
		ForwardOrderToShop(gomock.Any()).
		Return(nil)

	mockOrderRepo.EXPECT().
		UpdateOrderStatus(gomock.Any(), "success").
		Return(nil)

	order, err := orderService.ProcessPayment(int64(1), true)

	assert.NoError(t, err)
	assert.Equal(t, "success", order.Status)
}

func TestCancelOrder(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockShopRepo := mocks.NewMockShopRepository(ctrl)

	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo, mockShopRepo)

	order := &models.Order{
		Id:     1,
		Status: "pending",
		Items: []models.OrderItem{
			{ProductId: 1, Quantity: 2},
		},
	}

	mockOrderRepo.EXPECT().
		GetOrderById(gomock.Any()).
		Return(order, nil)

	mockOrderRepo.EXPECT().
		UpdateOrderStatus(gomock.Any(), "cancelled").
		Return(nil)

	mockProductRepo.EXPECT().
		RestoreStock(gomock.Any(), gomock.Any()).
		Return(nil)

	err := orderService.CancelOrder(int64(1))

	assert.NoError(t, err)
}
