package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"monorepo-ecommerce/micro-services/product/handler"
	mocks "monorepo-ecommerce/micro-services/product/mocks/mock_micro-services/product/service"
	"monorepo-ecommerce/micro-services/product/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetProducts(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockProductService := mocks.NewMockProductService(ctrl)
	h := handler.NewProductHandler(mockProductService)
	e := echo.New()

	t.Run("should success", func(t *testing.T) {
		mockProducts := []models.Product{
			{
				Id:          1,
				Name:        "Product 1",
				Description: "Description Product 1",
				Price:       100,
				Stock:       10,
			},
		}

		mockProductService.EXPECT().
			GetAllProducts().
			Return(mockProducts, nil)

		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.GetProducts(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should success when product is empty", func(t *testing.T) {
		mockProducts := []models.Product{}

		mockProductService.EXPECT().
			GetAllProducts().
			Return(mockProducts, nil)

		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.GetProducts(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should internal server error when products failed to fetch", func(t *testing.T) {
		mockProducts := []models.Product{}

		mockProductService.EXPECT().
			GetAllProducts().
			Return(mockProducts, errors.New("failed"))

		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.GetProducts(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGetProduct(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockProductService := mocks.NewMockProductService(ctrl)
	h := handler.NewProductHandler(mockProductService)
	e := echo.New()

	t.Run("should success", func(t *testing.T) {
		var mockId int64 = 1

		mockProduct := models.Product{
			Id:          1,
			Name:        "Product 1",
			Description: "Description 1",
			Price:       100,
			Stock:       10,
		}

		mockProductService.EXPECT().
			GetProductById(mockId).
			Return(&mockProduct, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// set params
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := h.GetProduct(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should internal server error when failed to fetch product", func(t *testing.T) {
		var mockId int64 = 1

		mockProduct := models.Product{}

		mockProductService.EXPECT().
			GetProductById(mockId).
			Return(&mockProduct, errors.New("failed"))

		req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// set params
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := h.GetProduct(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestDeductStock(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockProductService := mocks.NewMockProductService(ctrl)
	h := handler.NewProductHandler(mockProductService)
	e := echo.New()

	t.Run("should success", func(t *testing.T) {
		var mockId int64 = 1

		type requestBody struct {
			Quantity int `json:"quantity"`
		}

		reqBody := requestBody{
			Quantity: 10,
		}

		reqJSON, _ := json.Marshal(reqBody)

		mockProductService.EXPECT().
			DeductStock(mockId, reqBody.Quantity).
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/products/deduct/1", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// set params
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := h.DeductStock(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should bad request when request invalid", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"quantity": "1",
		}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/products/deduct/1", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// set params
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := h.DeductStock(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should internal server error when deduct failed", func(t *testing.T) {
		var mockId int64 = 1

		type requestBody struct {
			Quantity int `json:"quantity"`
		}

		reqBody := requestBody{
			Quantity: 10,
		}

		reqJSON, _ := json.Marshal(reqBody)

		mockProductService.EXPECT().
			DeductStock(mockId, reqBody.Quantity).
			Return(errors.New("failed"))

		req := httptest.NewRequest(http.MethodPost, "/products/deduct/1", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// set params
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := h.DeductStock(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestRestoreStock(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockProductService := mocks.NewMockProductService(ctrl)
	h := handler.NewProductHandler(mockProductService)
	e := echo.New()

	t.Run("should success", func(t *testing.T) {
		var mockId int64 = 1

		type requestBody struct {
			Quantity int `json:"quantity"`
		}

		reqBody := requestBody{
			Quantity: 10,
		}

		reqJSON, _ := json.Marshal(reqBody)

		mockProductService.EXPECT().
			RestoreStock(mockId, reqBody.Quantity).
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/products/restore/1", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// set params
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := h.RestoreStock(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should bad request when request invalid", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"quantity": "1",
		}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/products/restore/1", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// set params
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := h.RestoreStock(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should internal server error when restore failed", func(t *testing.T) {
		var mockId int64 = 1

		type requestBody struct {
			Quantity int `json:"quantity"`
		}

		reqBody := requestBody{
			Quantity: 10,
		}

		reqJSON, _ := json.Marshal(reqBody)

		mockProductService.EXPECT().
			RestoreStock(mockId, reqBody.Quantity).
			Return(errors.New("failed"))

		req := httptest.NewRequest(http.MethodPost, "/products/restore/1", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// set params
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := h.RestoreStock(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestUpdateTotalProductStock(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockProductService := mocks.NewMockProductService(ctrl)
	h := handler.NewProductHandler(mockProductService)
	e := echo.New()

	t.Run("should success", func(t *testing.T) {
		var mockId int64 = 1

		type requestBody struct {
			Quantity int `json:"quantity"`
		}

		reqBody := requestBody{
			Quantity: 10,
		}

		reqJSON, _ := json.Marshal(reqBody)

		mockProductService.EXPECT().
			UpdateTotalStock(mockId, reqBody.Quantity).
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/products/adjust-total-stock/1", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// set params
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := h.UpdateTotalProductStock(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should bad request when request invalid", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"quantity": "1",
		}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/products/adjust-total-stock/1", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// set params
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := h.UpdateTotalProductStock(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should internal server error when failed to update total product stock", func(t *testing.T) {
		var mockId int64 = 1

		type requestBody struct {
			Quantity int `json:"quantity"`
		}

		reqBody := requestBody{
			Quantity: 10,
		}

		reqJSON, _ := json.Marshal(reqBody)

		mockProductService.EXPECT().
			UpdateTotalStock(mockId, reqBody.Quantity).
			Return(errors.New("failed"))

		req := httptest.NewRequest(http.MethodPost, "/products/adjust-total-stock/1", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// set params
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := h.UpdateTotalProductStock(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
