package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/klyakssa/go-diplom.git/internal/domain/auth"
	"github.com/klyakssa/go-diplom.git/internal/logger"
	"go.uber.org/zap"
)

const invalidRequestMsg = "Invalid request"

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
		h.log.Error(invalidRequestMsg, zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidRequestMsg})
		return
	}

	h.log.Debug("RegisterRequest", zap.Any("body", req))

	token, err := h.service.Register(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			h.log.Warn("User already exists", zap.String("login", req.Login))
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		if errors.Is(err, auth.ErrPasswordTooLong) {
			h.log.Warn("Password too long", zap.String("login", req.Login))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password is too long"})
			return
		}
		h.log.Error("Failed to register user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.SetCookieData(
		&http.Cookie{
			Name:     "auth_token",
			Value:    token,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   false,
		},
	)
	c.Status(http.StatusOK)
}

type loginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error(invalidRequestMsg, zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidRequestMsg})
		return
	}

	h.log.Debug("LoginRequest", zap.Any("body", req))

	token, err := h.service.Login(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			h.log.Warn("User not found", zap.String("login", req.Login))
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if errors.Is(err, auth.ErrInvalidCredentials) {
			h.log.Warn("Invalid login or password", zap.String("login", req.Login))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})
			return
		}
		h.log.Error("Failed to login user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.SetCookieData(
		&http.Cookie{
			Name:     "auth_token",
			Value:    token,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   false,
		},
	)
	c.Status(http.StatusOK)
}
