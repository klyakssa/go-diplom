package balance

import "context"

type Repository interface {
	WithdrawBalance(ctx context.Context, userID string, orderNumber string, amount int) error
	GetBalanceWithdrawn(ctx context.Context, userID string) (int, int, error)
	GetWithdrawls(ctx context.Context, userID string) ([]WithdrawHistory, error)
}
