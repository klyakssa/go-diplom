package service

import (
	"context"

	"github.com/klyakssa/go-diplom.git/internal/domain/balance"
	"github.com/klyakssa/go-diplom.git/pkg/luhn"
	"github.com/shopspring/decimal"
)

// BalanceService
type BalanceService struct {
	repo balance.Repository
}

// NewBalanceService returns instance of new BalanceService
func NewBalanceService(repo balance.Repository) *BalanceService {
	return &BalanceService{repo: repo}
}

// WithdrawBalance withdraws balance
func (b *BalanceService) WithdrawBalance(ctx context.Context, userID string, orderNumber string, amount decimal.Decimal) error {
	if !luhn.Valid(orderNumber) {
		return balance.ErrIncorrectOrderNumberFormat
	}

	return b.repo.WithdrawBalance(ctx, userID, orderNumber, amount)
}

// GetBalanceWithdrawn returns balance
func (b *BalanceService) GetBalanceWithdrawn(ctx context.Context, userID string) (decimal.Decimal, decimal.Decimal, error) {
	return b.repo.GetBalanceWithdrawn(ctx, userID)
}

// GetWithdrawls returns balance
func (b *BalanceService) GetWithdrawls(ctx context.Context, userID string) ([]balance.WithdrawHistory, error) {
	return b.repo.GetWithdrawls(ctx, userID)
}
