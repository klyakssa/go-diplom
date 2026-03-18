package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/klyakssa/go-diplom.git/internal/domain/balance"
	"github.com/klyakssa/go-diplom.git/internal/logger"
	"go.uber.org/zap"
)

type BalanceHandler struct {
	service balance.Service
	log     *logger.Logger
}

func NewBalanceHandler(log *logger.Logger, service balance.Service) *BalanceHandler {
	return &BalanceHandler{
		service: service,
		log:     log,
	}
}

type WithdrawBalanceRequest struct {
	OrderNumber string `json:"order" binding:"required"` // TODO: validate order number
	Sum         int    `json:"sum" binding:"required"`
}

func (h *BalanceHandler) WithdrawBalance(c *gin.Context) {

	var req WithdrawBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error(invalidRequestMsg, zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidRequestMsg})
		return
	}

	h.log.Debug("WithdrawBalanceRequest", zap.Any("body", req))

	if err := h.service.WithdrawBalance(c.Request.Context(), c.GetString("user_id"), req.OrderNumber, req.Sum); err != nil {
		h.log.Error("Failed to withdraw balance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
}
