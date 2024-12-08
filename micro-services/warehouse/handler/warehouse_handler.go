package handler

import (
	"monorepo-ecommerce/micro-services/warehouse/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type WarehouseRequest struct {
	ProductId   int64 `json:"product_id"`
	WarehouseId int64 `json:"warehouse_id"`
	Quantity    int   `json:"quantity"`
}

type TransferProductRequest struct {
	OriginWarehouseId      int64 `json:"origin_warehouse_id"`
	DestinationWarehouseId int64 `json:"destination_warehouse_id"`
	ProductId              int64 `json:"product_id"`
	Quantity               int   `json:"quantity"`
}

type AvailabilityWarehouseRequest struct {
	WarehouseId int64 `json:"warehouse_id"`
}

type ProceedOrderRequest struct {
	OrderID int64                 `json:"order_id"`
	Items   []ProductOrderDetails `json:"items"`
}

type ProductOrderDetails struct {
	ProductId int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type WarehouseHandler struct {
	WarehouseService service.WarehouseService
}

func NewWarehouseHandler(warehouseService service.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{WarehouseService: warehouseService}
}

func (h *WarehouseHandler) AddStock(c echo.Context) error {
	var req WarehouseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	err := h.WarehouseService.AddStock(req.ProductId, req.WarehouseId, req.Quantity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Stock added successfully"})
}

func (h *WarehouseHandler) RemoveStock(c echo.Context) error {
	var req WarehouseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	err := h.WarehouseService.RemoveStock(req.ProductId, req.WarehouseId, req.Quantity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Stock removed successfully"})
}

func (h *WarehouseHandler) TransferProduct(c echo.Context) error {
	var req TransferProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	err := h.WarehouseService.TransferProduct(req.ProductId, req.OriginWarehouseId, req.DestinationWarehouseId, req.Quantity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Product transfered successfully"})
}

func (h *WarehouseHandler) ActiveDeactiveWarehouse(c echo.Context) error {
	var req AvailabilityWarehouseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	err := h.WarehouseService.ActiveDeactiveWarehouseStatus(req.WarehouseId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Product transfered successfully"})
}

func (h *WarehouseHandler) ProceedOrder(c echo.Context) error {
	var req ProceedOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	result := make([]service.ProductOrderDetails, len(req.Items))
	for i, item := range req.Items {
		result[i] = service.ProductOrderDetails{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		}
	}
	err := h.WarehouseService.ProceedOrder(req.OrderID, result)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Order processed successfully"})
}

func RegisterWarehouseRoutes(e *echo.Echo, warehouseService service.WarehouseService) {
	handler := NewWarehouseHandler(warehouseService)
	e.POST("/warehouse/stock/add", handler.AddStock)
	e.POST("/warehouse/stock/remove", handler.RemoveStock)
	e.POST("/warehouse/stock/transfer-product", handler.TransferProduct)
	e.POST("/warehouse/stock/active-deactive", handler.ActiveDeactiveWarehouse)
	e.POST("/warehouse/stock/proceed-order", handler.ProceedOrder)
}
