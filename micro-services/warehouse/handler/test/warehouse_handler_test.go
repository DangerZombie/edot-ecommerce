package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"monorepo-ecommerce/micro-services/warehouse/handler"
	mocks "monorepo-ecommerce/micro-services/warehouse/mocks/mock_micro-services/warehouse/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAddStock(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockWarehouseService := mocks.NewMockWarehouseService(ctrl)
	h := handler.NewWarehouseHandler(mockWarehouseService)
	e := echo.New()

	t.Run("should success", func(t *testing.T) {
		reqBody := handler.WarehouseRequest{
			ProductId:   1,
			WarehouseId: 1,
			Quantity:    10,
		}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/warehouse/stock/add", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockWarehouseService.EXPECT().
			AddStock(reqBody.ProductId, reqBody.WarehouseId, reqBody.Quantity).
			Return(nil)

		err := h.AddStock(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should bad request when request invalid", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"product_id": "1",
		}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/warehouse/stock/add", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.AddStock(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should internal server error when failed add stock", func(t *testing.T) {
		reqBody := handler.WarehouseRequest{
			ProductId:   1,
			WarehouseId: 1,
			Quantity:    10,
		}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/warehouse/stock/add", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockWarehouseService.EXPECT().
			AddStock(reqBody.ProductId, reqBody.WarehouseId, reqBody.Quantity).
			Return(errors.New("failed"))

		err := h.AddStock(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestRemoveStock(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockWarehouseService := mocks.NewMockWarehouseService(ctrl)
	h := handler.NewWarehouseHandler(mockWarehouseService)
	e := echo.New()

	t.Run("should success", func(t *testing.T) {
		reqBody := handler.WarehouseRequest{
			ProductId:   1,
			WarehouseId: 1,
			Quantity:    10,
		}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/warehouse/stock/remove", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockWarehouseService.EXPECT().
			RemoveStock(reqBody.ProductId, reqBody.WarehouseId, reqBody.Quantity).
			Return(nil)

		err := h.RemoveStock(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should bad request when request invalid", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"product_id": "1",
		}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/warehouse/stock/remove", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.RemoveStock(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should internal server error when failed remove stock", func(t *testing.T) {
		reqBody := handler.WarehouseRequest{
			ProductId:   1,
			WarehouseId: 1,
			Quantity:    10,
		}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/warehouse/stock/remove", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockWarehouseService.EXPECT().
			RemoveStock(reqBody.ProductId, reqBody.WarehouseId, reqBody.Quantity).
			Return(errors.New("failed"))

		err := h.RemoveStock(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestTransferProduct(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockWarehouseService := mocks.NewMockWarehouseService(ctrl)
	h := handler.NewWarehouseHandler(mockWarehouseService)
	e := echo.New()

	t.Run("should success", func(t *testing.T) {
		reqBody := handler.TransferProductRequest{
			OriginWarehouseId:      1,
			DestinationWarehouseId: 2,
			ProductId:              1,
			Quantity:               10,
		}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/warehouse/stock/transfer-product", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockWarehouseService.EXPECT().
			TransferProduct(reqBody.ProductId, reqBody.OriginWarehouseId, reqBody.DestinationWarehouseId, reqBody.Quantity).
			Return(nil)

		err := h.TransferProduct(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should bad request when request invalid", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"product_id": "1",
		}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/warehouse/stock/transfer-product", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.TransferProduct(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should internal server error when failed transfer stock", func(t *testing.T) {
		reqBody := handler.TransferProductRequest{
			OriginWarehouseId:      1,
			DestinationWarehouseId: 2,
			ProductId:              1,
			Quantity:               10,
		}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/warehouse/stock/transfer-product", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockWarehouseService.EXPECT().
			TransferProduct(reqBody.ProductId, reqBody.OriginWarehouseId, reqBody.DestinationWarehouseId, reqBody.Quantity).
			Return(errors.New("failed"))

		err := h.TransferProduct(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestActiveDeactiveWarehouse(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockWarehouseService := mocks.NewMockWarehouseService(ctrl)
	h := handler.NewWarehouseHandler(mockWarehouseService)
	e := echo.New()

	t.Run("should success", func(t *testing.T) {
		reqBody := handler.AvailabilityWarehouseRequest{
			WarehouseId: 1,
		}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/warehouse/stock/active-deactive", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockWarehouseService.EXPECT().
			ActiveDeactiveWarehouseStatus(reqBody.WarehouseId).
			Return(nil)

		err := h.ActiveDeactiveWarehouse(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should bad request when request invalid", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"warehouse_id": "1",
		}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/warehouse/stock/active-deactive", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.ActiveDeactiveWarehouse(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should internal server error when failed transfer stock", func(t *testing.T) {
		reqBody := handler.AvailabilityWarehouseRequest{
			WarehouseId: 1,
		}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/warehouse/stock/active-deactive", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockWarehouseService.EXPECT().
			ActiveDeactiveWarehouseStatus(reqBody.WarehouseId).
			Return(errors.New("failed"))

		err := h.ActiveDeactiveWarehouse(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
