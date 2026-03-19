package http_test

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/klyakssa/go-diplom.git/internal/domain/auth"
	mock_auth "github.com/klyakssa/go-diplom.git/internal/domain/auth/mocks"
	"github.com/klyakssa/go-diplom.git/internal/logger"
	httptransport "github.com/klyakssa/go-diplom.git/internal/transport/http"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAuthHandler_Register(t *testing.T) {

	type mockBehavior func(s *mock_auth.MockService, login, password string)

	logger := &logger.Logger{Logger: zap.NewNop()}

	tests := []struct {
		name           string
		body           string
		mock           mockBehavior
		statusCode     int
		expectedCookie string
	}{
		{
			name: "success register",
			body: `{"login":"user","password":"pass"}`,
			mock: func(s *mock_auth.MockService, login, password string) {
				s.EXPECT().Register(gomock.Any(), login, password).Return("1", nil)
			},
			statusCode:     200,
			expectedCookie: "1",
		},
		{
			name: "user exists",
			body: `{"login":"user","password":"pass"}`,
			mock: func(s *mock_auth.MockService, login, password string) {
				s.EXPECT().Register(gomock.Any(), login, password).Return("", auth.ErrUserAlreadyExists)
			},
			statusCode:     409,
			expectedCookie: "",
		},
		{
			name:           "invalid json",
			body:           `{}`,
			mock:           func(s *mock_auth.MockService, login, password string) {},
			statusCode:     400,
			expectedCookie: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_auth.NewMockService(c)
			tt.mock(auth, "user", "pass")

			handler := httptransport.NewAuthHandler(logger, auth)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.POST("/register", handler.Register)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)

			if tt.expectedCookie != "" {
				cookies := w.Result().Cookies()
				var token string
				for _, c := range cookies {
					if c.Name == "auth_token" {
						token = c.Value
					}
				}

				assert.Equal(t, tt.expectedCookie, token)
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {

	type mockBehavior func(s *mock_auth.MockService, login, password string)

	logger := &logger.Logger{Logger: zap.NewNop()}

	tests := []struct {
		name           string
		body           string
		mock           mockBehavior
		statusCode     int
		expectedCookie string
	}{
		{
			name: "success login",
			body: `{"login":"user","password":"pass"}`,
			mock: func(s *mock_auth.MockService, login, password string) {
				s.EXPECT().Login(gomock.Any(), login, password).Return("1", nil)
			},
			statusCode:     200,
			expectedCookie: "1",
		},
		{
			name: "invalid credentials",
			body: `{"login":"user","password":"pass"}`,
			mock: func(s *mock_auth.MockService, login, password string) {
				s.EXPECT().Login(gomock.Any(), login, password).Return("", auth.ErrInvalidCredentials)
			},
			statusCode:     401,
			expectedCookie: "",
		},
		{
			name:           "invalid json",
			body:           `{}`,
			mock:           func(s *mock_auth.MockService, login, password string) {},
			statusCode:     400,
			expectedCookie: "",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_auth.NewMockService(c)
			tt.mock(auth, "user", "pass")

			handler := httptransport.NewAuthHandler(logger, auth)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.POST("/login", handler.Login)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)

			if tt.expectedCookie != "" {
				cookies := w.Result().Cookies()
				var token string
				for _, c := range cookies {
					if c.Name == "auth_token" {
						token = c.Value
					}
				}

				assert.Equal(t, tt.expectedCookie, token)
			}
		})
	}
}
