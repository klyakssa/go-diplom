package http_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/klyakssa/go-diplom.git/internal/domain/balance"
	mock_balance "github.com/klyakssa/go-diplom.git/internal/domain/balance/mocks"
	"github.com/klyakssa/go-diplom.git/internal/logger"
	httptransport "github.com/klyakssa/go-diplom.git/internal/transport/http"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestBalanceHandler_WithdrawBalance(t *testing.T) {

	logger := &logger.Logger{Logger: zap.NewNop()}

	tests := []struct {
		name       string
		body       string
		moke       func(service *mock_balance.MockService, user_id string)
		statusCode int
	}{
		{
			name: "success withdraw",
			body: `{"order":"2377225624", "sum":100.00}`,
			moke: func(service *mock_balance.MockService, user_id string) {
				service.EXPECT().WithdrawBalance(gomock.Any(), user_id, "2377225624", gomock.Any()).Return(nil)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "error withdraw",
			body: `{"order":"2377225624", "sum":100.00}`,
			moke: func(service *mock_balance.MockService, user_id string) {
				service.EXPECT().WithdrawBalance(gomock.Any(), user_id, "2377225624", gomock.Any()).Return(balance.ErrIncorrectOrderNumberFormat)
			},
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "error withdraw",
			body: `{"order":"2377225624", "sum":100.00}`,
			moke: func(service *mock_balance.MockService, user_id string) {
				service.EXPECT().WithdrawBalance(gomock.Any(), user_id, "2377225624", gomock.Any()).Return(balance.ErrInsufficientFunds)
			},
			statusCode: http.StatusPaymentRequired,
		},
		{
			name:       "internal request",
			body:       `{}`,
			moke:       func(service *mock_balance.MockService, user_id string) {},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			balance := mock_balance.NewMockService(c)
			tt.moke(balance, "1")

			handler := httptransport.NewBalanceHandler(logger, balance)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.POST("/balance/withdraw", func(c *gin.Context) {
				c.Set("user_id", "1")
				handler.WithdrawBalance(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/balance/withdraw", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestBalanceHandler_GetBalanceWithdrawn(t *testing.T) {
	logger := &logger.Logger{Logger: zap.NewNop()}
	decimal.MarshalJSONWithoutQuotes = true

	tests := []struct {
		name         string
		moke         func(service *mock_balance.MockService, user_id string)
		statusCode   int
		expectedBody string
	}{
		{
			name: "success get balance withdrawn",
			moke: func(service *mock_balance.MockService, user_id string) {
				service.EXPECT().GetBalanceWithdrawn(gomock.Any(), user_id).Return(decimal.Zero, decimal.Zero, nil)
			},
			statusCode:   http.StatusOK,
			expectedBody: `{"current":0,"withdrawn":0}`,
		},
		{
			name: "error get balance withdrawn",
			moke: func(service *mock_balance.MockService, user_id string) {
				service.EXPECT().GetBalanceWithdrawn(gomock.Any(), user_id).Return(decimal.Zero, decimal.Zero, fmt.Errorf("error"))
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			balance := mock_balance.NewMockService(c)
			tt.moke(balance, "1")

			handler := httptransport.NewBalanceHandler(logger, balance)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.GET("/balance", func(c *gin.Context) {
				c.Set("user_id", "1")
				handler.GetBalanceWithdrawn(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/balance", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)

			if tt.statusCode == http.StatusOK {
				assert.Equal(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestBalanceHandler_GetWithdrawals(t *testing.T) {
	logger := &logger.Logger{Logger: zap.NewNop()}
	decimal.MarshalJSONWithoutQuotes = true

	timeNow := time.Now()

	tests := []struct {
		name         string
		moke         func(service *mock_balance.MockService, user_id string)
		statusCode   int
		expectedBody string
	}{
		{
			name: "success get withdrawals",
			moke: func(service *mock_balance.MockService, user_id string) {
				withdrawls := []balance.WithdrawHistory{
					{
						OrderNumber: "2377225624",
						Sum:         decimal.Zero,
						UserID:      1,
						ProcessedAt: balance.TimeRFC3339(timeNow),
					},
				}
				service.EXPECT().GetWithdrawls(gomock.Any(), user_id).Return(withdrawls, nil)
			},
			statusCode:   http.StatusOK,
			expectedBody: `[{"order":"2377225624","sum":0,"processed_at":"` + timeNow.Format(time.RFC3339) + `"}]`,
		},
		{
			name: "error get withdrawals",
			moke: func(service *mock_balance.MockService, user_id string) {
				service.EXPECT().GetWithdrawls(gomock.Any(), user_id).Return(nil, fmt.Errorf("error"))
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			balance := mock_balance.NewMockService(c)
			tt.moke(balance, "1")

			handler := httptransport.NewBalanceHandler(logger, balance)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.GET("/withdrawals", func(c *gin.Context) {
				c.Set("user_id", "1")
				handler.GetWithdrawals(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/withdrawals", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)

			if tt.statusCode == http.StatusOK {
				assert.Equal(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}
