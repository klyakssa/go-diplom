package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/klyakssa/go-diplom.git/internal/config"
	"github.com/klyakssa/go-diplom.git/internal/logger"
)

type Router struct {
	cfg    *config.WebServerConfig
	engine *gin.Engine
	server *http.Server
	log    *logger.Logger
}

func NewRouter(logger *logger.Logger, cfg *config.WebServerConfig) *Router {
	engine := gin.Default()
	engine.Use(gin.Recovery())
	return &Router{
		cfg:    cfg,
		engine: engine,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      engine,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		log: logger,
	}
}

func (r *Router) Run(ctx context.Context) error {

	errChan := make(chan error, 1)
	go func() {
		if err := r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
		r.log.Info("HTTP server started on port " + fmt.Sprintf("%d", r.cfg.Port))
	}()

	select {
	case <-ctx.Done():
		r.log.Info("HTTP server shutdown triggered by context")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return r.server.Shutdown(shutdownCtx)
	case err := <-errChan:
		return err
	}
}

func (r *Router) RegisterRoutes(authHandler *AuthHandler) {
	api := r.engine.Group("/api/user")
	{
		api.POST("/register", authHandler.Register)
		// api.POST("/login", authHandler.Login)
	}
}
