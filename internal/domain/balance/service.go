package balance

import (
	"context"

	"github.com/shopspring/decimal"
)

//go:generate mockgen -source=service.go -destination=mocks/service.go

type Service interface {
	WithdrawBalance(ctx context.Context, userID string, orderNumber string, amount decimal.Decimal) error
	GetBalanceWithdrawn(ctx context.Context, userID string) (decimal.Decimal, decimal.Decimal, error)
	GetWithdrawls(ctx context.Context, userID string) ([]WithdrawHistory, error)
}
