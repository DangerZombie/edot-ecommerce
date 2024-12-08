package handler

import (
	"monorepo-ecommerce/micro-services/product/models"
	"monorepo-ecommerce/micro-services/product/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) GetProducts(c echo.Context) error {
	products, err := h.service.GetAllProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products"})
	}

	if len(products) == 0 {
		return c.JSON(http.StatusOK, []models.Product{})
	}

	return c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetProduct(c echo.Context) error {
	id := c.Param("id")
	productId, _ := strconv.ParseInt(id, 10, 64)
	product, err := h.service.GetProductById(productId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch product"})
	}

	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeductStock(c echo.Context) error {
	id := c.Param("id")
	var requestBody struct {
		Quantity int `json:"quantity"`
	}

	// Bind JSON body to struct
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Request body invalid"})
	}

	productId, _ := strconv.ParseInt(id, 10, 64)
	err := h.service.DeductStock(productId, requestBody.Quantity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed deduct product stock"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Product stock success to deduct"})
}

func (h *ProductHandler) RestoreStock(c echo.Context) error {
	id := c.Param("id")
	var requestBody struct {
		Quantity int `json:"quantity"`
	}

	// Bind JSON body to struct
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Request body invalid"})
	}

	productId, _ := strconv.ParseInt(id, 10, 64)
	err := h.service.RestoreStock(productId, requestBody.Quantity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed deduct product stock"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Product stock success to deduct"})
}

func (h *ProductHandler) UpdateTotalProductStock(c echo.Context) error {
	id := c.Param("id")
	var requestBody struct {
		Quantity int `json:"quantity"`
	}

	// Bind JSON body to struct
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Request body invalid"})
	}

	productId, _ := strconv.ParseInt(id, 10, 64)
	err := h.service.UpdateTotalStock(productId, requestBody.Quantity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed deduct product stock"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Product stock success to deduct"})
}

func RegisterProductRoutes(e *echo.Echo, productService service.ProductService) {
	handler := NewProductHandler(productService)
	e.GET("/products", handler.GetProducts)
	e.GET("/products/:id", handler.GetProduct)
	e.POST("/products/deduct/:id", handler.DeductStock)
	e.POST("/products/restore/:id", handler.RestoreStock)
	e.POST("/products/adjust-total-stock/:id", handler.UpdateTotalProductStock)
}
