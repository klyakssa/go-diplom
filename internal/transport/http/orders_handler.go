package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/klyakssa/go-diplom.git/internal/domain/orders"
	"github.com/klyakssa/go-diplom.git/internal/logger"
	"go.uber.org/zap"
)

type OrdersHandler struct {
	log     *logger.Logger
	service orders.Service
}

func NewOrdersHandler(log *logger.Logger, service orders.Service) *OrdersHandler {
	return &OrdersHandler{
		log:     log,
		service: service,
	}
}

func (o *OrdersHandler) CreateOrders(c *gin.Context) {
	var number string
	if err := c.ShouldBindPlain(&number); err != nil {
		o.log.Error(invalidRequestMsg, zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidRequestMsg})
		return
	}

	o.log.Debug("CreateOrderRequest", zap.Any("body", number))

	if number == "" {
		o.log.Warn("Order number is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order number is empty"})
		return
	}

	if len(number) > 1024 {
		o.log.Warn("Order number is too long")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order number is too long"})
		return
	}

	err := o.service.CreateOrder(c.Request.Context(), number, c.GetString("user_id"))
	if err != nil {
		if errors.Is(err, orders.ErrIncorrectOrderNumberFormat) {
			o.log.Warn("Incorrect order number format", zap.String("number", number))
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Incorrect order number format"})
			return
		}
		if errors.Is(err, orders.ErrOrderAlreadyAddedByThisUser) {
			o.log.Warn("Order already added by this user", zap.String("number", number))
			c.JSON(http.StatusOK, gin.H{"error": "Order already added by this user"})
			return
		}
		if errors.Is(err, orders.ErrOrderAlreadyAddedByOtherUser) {
			o.log.Warn("Order already added by other user", zap.String("number", number))
			c.JSON(http.StatusConflict, gin.H{"error": "Order already added by other user"})
			return
		}
		o.log.Error("Failed to create order", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Order added successfully"})
}

func (o *OrdersHandler) GetOrders(c *gin.Context) {
	orders, err := o.service.GetOrders(c.Request.Context(), c.GetString("user_id"))
	if err != nil {
		o.log.Error("Failed to get orders", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if len(orders) == 0 {
		o.log.Warn("Orders not found")
		c.JSON(http.StatusNoContent, gin.H{"error": "Orders not found"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
