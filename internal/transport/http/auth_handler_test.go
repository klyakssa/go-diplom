package http_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/klyakssa/go-diplom.git/internal/config"
	"github.com/klyakssa/go-diplom.git/internal/domain/auth"
	"github.com/klyakssa/go-diplom.git/internal/logger"
	httptransport "github.com/klyakssa/go-diplom.git/internal/transport/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Register(ctx context.Context, login, password string) (string, error) {
	args := m.Called(ctx, login, password)
	return args.String(0), args.Error(1)
}

func (m *mockAuthService) Login(ctx context.Context, login, password string) (string, error) {
	args := m.Called(ctx, login, password)
	return args.String(0), args.Error(1)
}

func setupServer(service auth.Service) *httptest.Server {
	gin.SetMode(gin.TestMode)

	log := logger.NewLogger(&config.LoggingConfiguration{
		Level:      "debug",
		Path:       "",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
	}, "test")

	handler := httptransport.NewAuthHandler(log, service)

	router := gin.New()
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)

	return httptest.NewServer(router)
}

func TestAuthHandler_Register(t *testing.T) {

	service := new(mockAuthService)
	server := setupServer(service)
	defer server.Close()

	client := resty.New().SetBaseURL(server.URL)

	tests := []struct {
		name       string
		body       string
		mock       func()
		statusCode int
	}{
		{
			name: "success register",
			body: `{"login":"user","password":"pass"}`,
			mock: func() {
				service.On("Register", mock.Anything, "user", "pass").
					Return("token123", nil).Once()
			},
			statusCode: 200,
		},
		{
			name: "user exists",
			body: `{"login":"user","password":"pass"}`,
			mock: func() {
				service.On("Register", mock.Anything, "user", "pass").
					Return("", auth.ErrUserAlreadyExists).Once()
			},
			statusCode: 409,
		},
		{
			name:       "invalid json",
			body:       `{}`,
			mock:       func() {},
			statusCode: 400,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			service.ExpectedCalls = nil
			tt.mock()

			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(tt.body).
				Post("/register")

			assert.NoError(t, err)
			assert.Equal(t, tt.statusCode, resp.StatusCode())
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {

	service := new(mockAuthService)
	server := setupServer(service)
	defer server.Close()

	client := resty.New().SetBaseURL(server.URL)

	tests := []struct {
		name       string
		body       string
		mock       func()
		statusCode int
	}{
		{
			name: "success login",
			body: `{"login":"user","password":"pass"}`,
			mock: func() {
				service.On("Login", mock.Anything, "user", "pass").
					Return("token123", nil).Once()
			},
			statusCode: 200,
		},
		{
			name: "invalid credentials",
			body: `{"login":"user","password":"pass"}`,
			mock: func() {
				service.On("Login", mock.Anything, "user", "pass").
					Return("", auth.ErrInvalidCredentials).Once()
			},
			statusCode: 401,
		},
		{
			name:       "invalid json",
			body:       `{}`,
			mock:       func() {},
			statusCode: 400,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			service.ExpectedCalls = nil
			tt.mock()

			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(tt.body).
				Post("/login")

			assert.NoError(t, err)
			assert.Equal(t, tt.statusCode, resp.StatusCode())
		})
	}
}
