package balance

import (
	"context"

	"github.com/shopspring/decimal"
)

type Service interface {
	WithdrawBalance(ctx context.Context, userID string, orderNumber string, amount decimal.Decimal) error
	GetBalanceWithdrawn(ctx context.Context, userID string) (int, int, error)
	GetWithdrawls(ctx context.Context, userID string) ([]WithdrawHistory, error)
}
