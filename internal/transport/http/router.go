package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/klyakssa/go-diplom.git/internal/config"
	"github.com/klyakssa/go-diplom.git/internal/logger"
	"github.com/klyakssa/go-diplom.git/internal/middleware"
	"github.com/klyakssa/go-diplom.git/pkg/jwt"
)

type Router struct {
	cfg        *config.WebServerConfig
	engine     *gin.Engine
	server     *http.Server
	log        *logger.Logger
	jwtManager *jwt.JWTManager
}

func NewRouter(logger *logger.Logger, cfg *config.Config, jwtManager *jwt.JWTManager) *Router {
	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(
		middleware.Recovery(logger.Logger),
		middleware.LoggingMiddleware(logger.Logger),
	)
	return &Router{
		cfg:    cfg.Web,
		engine: engine,
		server: &http.Server{
			Addr:         cfg.Web.RunAddress,
			Handler:      engine,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		log:        logger,
		jwtManager: jwtManager,
	}
}

func (r *Router) Run(ctx context.Context) error {

	errChan := make(chan error, 1)
	go func() {
		r.log.Info("HTTP server started on address " + r.cfg.RunAddress)
		if err := r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
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
		api.POST("/login", authHandler.Login)
		auth := api.Group("/")
		auth.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			// Здесь можно добавлять защищенные маршруты, например:
		}
	}
}
