package balance

import (
	"context"

	"github.com/shopspring/decimal"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository.go

// Repository is an interface for balance repository
type Repository interface {
	WithdrawBalance(ctx context.Context, userID string, orderNumber string, amount decimal.Decimal) error // subtracts balance of user
	GetBalanceWithdrawn(ctx context.Context, userID string) (decimal.Decimal, decimal.Decimal, error)     // returns balance and withdrawls of user
	GetWithdrawls(ctx context.Context, userID string) ([]WithdrawHistory, error)                          // returns withdrawls of user
}
