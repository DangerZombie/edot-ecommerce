package handler

import (
	"monorepo-ecommerce/micro-services/order/middleware"
	"monorepo-ecommerce/micro-services/order/models"
	"monorepo-ecommerce/micro-services/order/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	OrderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{OrderService: orderService}
}

func (h *OrderHandler) Checkout(c echo.Context) error {
	var orderRequest models.OrderRequest
	if err := c.Bind(&orderRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Checkout process
	order, err := h.OrderService.CreateOrder(c, &orderRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) Payment(c echo.Context) error {
	orderIdParam := c.Param("orderId")
	orderId, err := strconv.ParseInt(orderIdParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid order Id")
	}

	var paymentRequest struct {
		Paid bool `json:"paid"`
	}
	if err := c.Bind(&paymentRequest); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request body")
	}

	order, err := h.OrderService.ProcessPayment(orderId, paymentRequest.Paid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, order)
}

func RegisterOrderRoutes(e *echo.Echo, orderService service.OrderService) {
	handler := NewOrderHandler(orderService)
	e.POST("/order/checkout", handler.Checkout, middleware.IsAuthenticated)
	e.POST("/order/payment/:orderId", handler.Payment, middleware.IsAuthenticated)
}
