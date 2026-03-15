package http

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/klyakssa/go-diplom.git/internal/domain/auth"
	"github.com/klyakssa/go-diplom.git/internal/logger"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service auth.Service
	log     *logger.Logger
}

func NewAuthHandler(log *logger.Logger, service auth.Service) *AuthHandler {
	return &AuthHandler{log: log, service: service}
}

type registerRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Invalid request", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	h.log.Debug("RegisterRequest", zap.Any("body", req))

	token, err := h.service.Register(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			h.log.Warn("User already exists", zap.String("login", req.Login))
			c.JSON(409, gin.H{"error": "User already exists"})
			return
		}
		h.log.Error("Failed to register user", zap.Error(err))
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(200, &gin.H{"token": token})
}
