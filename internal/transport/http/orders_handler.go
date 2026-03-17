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
		c.JSON(400, gin.H{"error": invalidRequestMsg})
		return
	}

	o.log.Debug("CreateOrderRequest", zap.Any("body", number))

	err := o.service.CreateOrder(c.Request.Context(), number)
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

	c.Status(http.StatusAccepted)
}
