package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/klyakssa/go-diplom.git/internal/domain/orders"
	mock_orders "github.com/klyakssa/go-diplom.git/internal/domain/orders/mocks"
	"github.com/klyakssa/go-diplom.git/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestOrdersService_CreateOrder(t *testing.T) {
	type mockBehaviour func(service *mock_orders.MockRepository, user_id string, order string)

	tests := []struct {
		name    string
		user_id string
		order   string
		moke    mockBehaviour
		error   error
	}{
		{
			name:    "success create order",
			user_id: "1",
			order:   "2377225624",
			moke: func(service *mock_orders.MockRepository, user_id string, order string) {
				service.EXPECT().GetOrderByNumber(gomock.Any(), order).Return(nil, fmt.Errorf("mock error"))
				service.EXPECT().CreateOrder(gomock.Any(), &orders.Order{Number: order, UserID: user_id}).Return(nil)
			},
			error: nil,
		},
		{
			name:    "incorrect order number format",
			user_id: "1",
			order:   "4561261212345464",
			moke:    func(service *mock_orders.MockRepository, user_id string, order string) {},
			error:   orders.ErrIncorrectOrderNumberFormat,
		},
		{
			name:    "order already added by this user",
			user_id: "1",
			order:   "2377225624",
			moke: func(service *mock_orders.MockRepository, user_id string, order string) {
				service.EXPECT().GetOrderByNumber(gomock.Any(), order).Return(&orders.Order{Number: order, UserID: user_id}, nil)
			},
			error: orders.ErrOrderAlreadyAddedByThisUser,
		},
		{
			name:    "order already added by other user",
			user_id: "1",
			order:   "2377225624",
			moke: func(service *mock_orders.MockRepository, user_id string, order string) {
				service.EXPECT().GetOrderByNumber(gomock.Any(), order).Return(&orders.Order{Number: order, UserID: "2"}, nil)
			},
			error: orders.ErrOrderAlreadyAddedByOtherUser,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			ordersRepository := mock_orders.NewMockRepository(c)
			tt.moke(ordersRepository, tt.user_id, tt.order)

			ordersService := service.NewOrdersService(ordersRepository)

			err := ordersService.CreateOrder(context.Background(), tt.order, tt.user_id)

			assert.ErrorIs(t, err, tt.error, "got error %v, want %v", err, tt.error)
		})
	}
}
