package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"monorepo-ecommerce/micro-services/shop/handler"
	mocks "monorepo-ecommerce/micro-services/shop/mocks/mock_micro-services/shop/service"
	"monorepo-ecommerce/micro-services/shop/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestProcessOrder(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockShopService := mocks.NewMockShopService(ctrl)
	h := handler.NewShopHandler(mockShopService)
	e := echo.New()

	t.Run("should success", func(t *testing.T) {
		reqBody := models.Order{
			Id:     1,
			UserId: 1,
			Items: []models.OrderItem{
				{
					ProductId: 1,
					Quantity:  10,
					Price:     100,
				},
			},
		}
		reqJSON, _ := json.Marshal(reqBody)

		mockShopService.EXPECT().
			ProcessOrder(reqBody).
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/shop/proceed-order", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.ProcessOrder(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should bad request when request invalid", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"id": "1",
		}
		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/shop/proceed-order", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.ProcessOrder(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should internal server error when error occured", func(t *testing.T) {
		reqBody := models.Order{
			Id:     1,
			UserId: 1,
			Items: []models.OrderItem{
				{
					ProductId: 1,
					Quantity:  10,
					Price:     100,
				},
			},
		}
		reqJSON, _ := json.Marshal(reqBody)

		mockShopService.EXPECT().
			ProcessOrder(reqBody).
			Return(errors.New("failed"))

		req := httptest.NewRequest(http.MethodPost, "/shop/proceed-order", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.ProcessOrder(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
