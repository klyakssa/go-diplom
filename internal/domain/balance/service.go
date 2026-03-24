package balance

import (
	"context"

	"github.com/shopspring/decimal"
)

//go:generate mockgen -source=service.go -destination=mocks/service.go

// Service is an interface for balance service
type Service interface {
	WithdrawBalance(ctx context.Context, userID string, orderNumber string, amount decimal.Decimal) error // subtracts balance of user
	GetBalanceWithdrawn(ctx context.Context, userID string) (decimal.Decimal, decimal.Decimal, error)     // returns balance and withdrawls of user
	GetWithdrawls(ctx context.Context, userID string) ([]WithdrawHistory, error)                          // returns balance
}
