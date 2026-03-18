package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/klyakssa/go-diplom.git/internal/domain/balance"
	"github.com/klyakssa/go-diplom.git/internal/logger"
	"github.com/shopspring/decimal"
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
	OrderNumber string          `json:"order" binding:"required"`
	Sum         decimal.Decimal `json:"sum"`
}

func (h *BalanceHandler) WithdrawBalance(c *gin.Context) {

	var req WithdrawBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error(invalidRequestMsg, zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidRequestMsg})
		return
	}

	h.log.Debug("WithdrawBalanceRequest", zap.Any("body", req))

	if req.Sum.IsNegative() {
		h.log.Warn("Sum is negative")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sum is negative"})
		return
	}

	if req.Sum.IsZero() {
		h.log.Warn("Sum is zero")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sum is zero"})
		return
	}

	if err := h.service.WithdrawBalance(c.Request.Context(), c.GetString("user_id"), req.OrderNumber, req.Sum); err != nil {
		if errors.Is(err, balance.ErrIncorrectOrderNumberFormat) {
			h.log.Warn("Order not found", zap.String("order", req.OrderNumber))
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Order not found"})
			return
		}
		if errors.Is(err, balance.ErrInsufficientFunds) {
			h.log.Warn("Not enough balance", zap.String("order", req.OrderNumber))
			c.JSON(http.StatusPaymentRequired, gin.H{"error": "Not enough balance"})
			return
		}
		h.log.Error("Failed to withdraw balance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Balance withdrawn successfully"})
}

type DepositBalanceResponse struct {
	CurrentBalance int             `json:"current"`
	WithDrawn      decimal.Decimal `json:"withdrawn"`
}

func (h *BalanceHandler) GetBalanceWithdrawn(c *gin.Context) {
	balance, withdrawn, err := h.service.GetBalanceWithdrawn(c.Request.Context(), c.GetString("user_id"))
	if err != nil {
		h.log.Error("Failed to get balance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, DepositBalanceResponse{
		CurrentBalance: balance,
		WithDrawn:      withdrawn,
	})
}

func (h *BalanceHandler) GetWithdrawals(c *gin.Context) {
	withdrawls, err := h.service.GetWithdrawls(c.Request.Context(), c.GetString("user_id"))
	if err != nil {
		h.log.Error("Failed to get balance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, withdrawls)
}
