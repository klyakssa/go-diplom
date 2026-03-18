package service

import (
	"context"

	"github.com/klyakssa/go-diplom.git/internal/domain/balance"
	"github.com/klyakssa/go-diplom.git/pkg/luhn"
)

type BalanceService struct {
	repo balance.Repository
}

func NewBalanceService(repo balance.Repository) *BalanceService {
	return &BalanceService{repo: repo}
}

func (b *BalanceService) WithdrawBalance(ctx context.Context, userID string, orderNumber string, amount int) error {
	if !luhn.Valid(orderNumber) {
		return balance.ErrIncorrectOrderNumberFormat
	}

	return b.repo.WithdrawBalance(ctx, userID, orderNumber, amount)
}
