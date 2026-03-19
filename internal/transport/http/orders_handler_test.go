package http_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/klyakssa/go-diplom.git/internal/domain/orders"
	mock_orders "github.com/klyakssa/go-diplom.git/internal/domain/orders/mocks"
	"github.com/klyakssa/go-diplom.git/internal/logger"
	httptransport "github.com/klyakssa/go-diplom.git/internal/transport/http"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestOrdersHandler_CreateOrders(t *testing.T) {
	logger := &logger.Logger{Logger: zap.NewNop()}

	tests := []struct {
		name       string
		body       string
		moke       func(service *mock_orders.MockService, user_id string)
		statusCode int
	}{
		{
			name: "success create order",
			body: `2377225624`,
			moke: func(service *mock_orders.MockService, user_id string) {
				service.EXPECT().CreateOrder(gomock.Any(), "2377225624", user_id).Return(nil)
			},
			statusCode: http.StatusAccepted,
		},
		{
			name:       "invalid request",
			body:       ``,
			moke:       func(service *mock_orders.MockService, user_id string) {},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "invalid order number format",
			body: `4561261212345464`,
			moke: func(service *mock_orders.MockService, user_id string) {
				service.EXPECT().CreateOrder(gomock.Any(), "4561261212345464", user_id).Return(orders.ErrIncorrectOrderNumberFormat)
			},
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "order already added by this user",
			body: `2377225624`,
			moke: func(service *mock_orders.MockService, user_id string) {
				service.EXPECT().CreateOrder(gomock.Any(), "2377225624", user_id).Return(orders.ErrOrderAlreadyAddedByThisUser)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "order already added by other user",
			body: `2377225624`,
			moke: func(service *mock_orders.MockService, user_id string) {
				service.EXPECT().CreateOrder(gomock.Any(), "2377225624", user_id).Return(orders.ErrOrderAlreadyAddedByOtherUser)
			},
			statusCode: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			orders := mock_orders.NewMockService(c)
			tt.moke(orders, "1")

			handler := httptransport.NewOrdersHandler(logger, orders)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.POST("/orders", func(c *gin.Context) {
				c.Set("user_id", "1")
				handler.CreateOrders(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/orders", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestOrdersHandler_GetOrders(t *testing.T) {
	logger := &logger.Logger{Logger: zap.NewNop()}

	tests := []struct {
		name       string
		moke       func(service *mock_orders.MockService, user_id string)
		statusCode int
	}{
		{
			name: "success get orders",
			moke: func(service *mock_orders.MockService, user_id string) {
				orders := []orders.Order{
					{
						Number: "2377225624",
						UserID: user_id,
					},
				}
				service.EXPECT().GetOrders(gomock.Any(), user_id).Return(orders, nil)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "internal server error",
			moke: func(service *mock_orders.MockService, user_id string) {
				service.EXPECT().GetOrders(gomock.Any(), user_id).Return(nil, fmt.Errorf("error"))
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "orders not found",
			moke: func(service *mock_orders.MockService, user_id string) {
				orders := []orders.Order{}
				service.EXPECT().GetOrders(gomock.Any(), user_id).Return(orders, nil)
			},
			statusCode: http.StatusNoContent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			orders := mock_orders.NewMockService(c)
			tt.moke(orders, "1")

			handler := httptransport.NewOrdersHandler(logger, orders)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.GET("/orders", func(c *gin.Context) {
				c.Set("user_id", "1")
				handler.GetOrders(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/orders", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}
