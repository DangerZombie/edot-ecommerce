package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"monorepo-ecommerce/micro-services/order/handler"
	mocks "monorepo-ecommerce/micro-services/order/mocks/mock_micro-services/order/service"
	"monorepo-ecommerce/micro-services/order/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCheckout(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockOrderService := mocks.NewMockOrderService(ctrl)
	h := handler.NewOrderHandler(mockOrderService)
	e := echo.New()

	t.Run("should success", func(t *testing.T) {
		reqBody := models.OrderRequest{
			Items: []models.OrderItem{
				{
					Id:        1,
					ProductId: 1,
					Quantity:  10,
					Price:     100,
				},
			},
		}

		mockOrder := models.Order{
			Id:     1,
			UserId: 1,
			Items: []models.OrderItem{
				{
					Id:        1,
					ProductId: 1,
					Quantity:  10,
					Price:     100,
				},
			},
			TotalPrice: 1000,
			Status:     "pending",
		}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/order/checkout", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockOrderService.EXPECT().
			CreateOrder(c, &reqBody).
			Return(&mockOrder, nil)

		err := h.Checkout(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should bad request when request invalid", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"id": "1",
				},
			},
		}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/order/checkout", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.Checkout(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should internal server error when failed to checkout", func(t *testing.T) {
		reqBody := models.OrderRequest{
			Items: []models.OrderItem{
				{
					Id:        1,
					ProductId: 1,
					Quantity:  10,
					Price:     100,
				},
			},
		}

		mockOrder := models.Order{}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/order/checkout", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockOrderService.EXPECT().
			CreateOrder(c, &reqBody).
			Return(&mockOrder, errors.New("failed"))

		err := h.Checkout(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func Test(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockOrderService := mocks.NewMockOrderService(ctrl)
	h := handler.NewOrderHandler(mockOrderService)
	e := echo.New()

	t.Run("should success", func(t *testing.T) {
		var mockId int64 = 1

		type paymentRequest struct {
			Paid bool `json:"paid"`
		}

		reqBody := paymentRequest{
			Paid: true,
		}

		mockOrder := models.Order{
			Id:     1,
			UserId: 1,
			Items: []models.OrderItem{
				{
					Id:        1,
					ProductId: 1,
					Quantity:  10,
					Price:     100,
				},
			},
			TotalPrice: 1000,
			Status:     "success",
		}

		reqJSON, _ := json.Marshal(reqBody)

		mockOrderService.EXPECT().
			ProcessPayment(mockId, reqBody.Paid).
			Return(&mockOrder, nil)

		req := httptest.NewRequest(http.MethodPost, "/order/payment/1", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// set params
		c.SetParamNames("orderId")
		c.SetParamValues("1")

		err := h.Payment(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should bad request when order id invalid", func(t *testing.T) {
		type paymentRequest struct {
			Paid bool `json:"paid"`
		}

		reqBody := paymentRequest{
			Paid: true,
		}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/order/payment/1", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// set params
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := h.Payment(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should bad request when request invalid", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"paid": "true",
		}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/order/payment/1", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// set params
		c.SetParamNames("orderId")
		c.SetParamValues("1")

		err := h.Payment(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should internal server error when failed payment", func(t *testing.T) {
		var mockId int64 = 1

		type paymentRequest struct {
			Paid bool `json:"paid"`
		}

		reqBody := paymentRequest{
			Paid: true,
		}

		mockOrder := models.Order{}

		reqJSON, _ := json.Marshal(reqBody)

		mockOrderService.EXPECT().
			ProcessPayment(mockId, reqBody.Paid).
			Return(&mockOrder, errors.New("failed"))

		req := httptest.NewRequest(http.MethodPost, "/order/payment/1", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// set params
		c.SetParamNames("orderId")
		c.SetParamValues("1")

		err := h.Payment(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
