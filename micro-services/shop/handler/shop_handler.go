package handler

import (
	"monorepo-ecommerce/micro-services/shop/models"
	"monorepo-ecommerce/micro-services/shop/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ShopHandler struct {
	ShopService service.ShopService
}

func NewShopHandler(shopService service.ShopService) *ShopHandler {
	return &ShopHandler{ShopService: shopService}
}

func (h *ShopHandler) ProcessOrder(c echo.Context) error {
	var order models.Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err := h.ShopService.ProcessOrder(order)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Order processed successfully")
}

func RegisterShopRoutes(e *echo.Echo, shopService service.ShopService) {
	handler := NewShopHandler(shopService)
	e.POST("/shop/proceed-order", handler.ProcessOrder)
}
