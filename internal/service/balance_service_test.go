package service_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/klyakssa/go-diplom.git/internal/domain/balance"
	mock_balance "github.com/klyakssa/go-diplom.git/internal/domain/balance/mocks"
	"github.com/klyakssa/go-diplom.git/internal/service"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestBalanceService_WithdrawBalance(t *testing.T) {

	type mockBehaviour func(service *mock_balance.MockRepository, user_id string, order string, sum decimal.Decimal)

	tests := []struct {
		name    string
		user_id string
		order   string
		sum     decimal.Decimal
		moke    mockBehaviour
		error   error
	}{
		{
			name:    "success withdraw balance",
			user_id: "1",
			order:   "2377225624",
			sum:     decimal.NewFromInt(100),
			moke: func(service *mock_balance.MockRepository, user_id string, order string, sum decimal.Decimal) {
				service.EXPECT().WithdrawBalance(gomock.Any(), user_id, order, sum).Return(nil)
			},
			error: nil,
		},
		{
			name:    "incorrect order number format",
			user_id: "1",
			order:   "4561261212345464",
			sum:     decimal.NewFromInt(100),
			moke:    func(service *mock_balance.MockRepository, user_id string, order string, sum decimal.Decimal) {},
			error:   balance.ErrIncorrectOrderNumberFormat,
		},
		{
			name:    "error insufficient funds",
			user_id: "1",
			order:   "2377225624",
			sum:     decimal.NewFromInt(100),
			moke: func(service *mock_balance.MockRepository, user_id string, order string, sum decimal.Decimal) {
				service.EXPECT().WithdrawBalance(gomock.Any(), user_id, order, sum).Return(balance.ErrInsufficientFunds)
			},
			error: balance.ErrInsufficientFunds,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			balanceRepository := mock_balance.NewMockRepository(c)
			tt.moke(balanceRepository, tt.user_id, tt.order, tt.sum)

			balanceService := service.NewBalanceService(balanceRepository)

			err := balanceService.WithdrawBalance(context.Background(), tt.user_id, tt.order, tt.sum)

			assert.ErrorIs(t, err, tt.error, "got error %v, want %v", err, tt.error)
		})
	}
}
