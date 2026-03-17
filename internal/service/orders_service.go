package service

import (
	"context"

	"github.com/klyakssa/go-diplom.git/internal/domain/orders"
	"github.com/klyakssa/go-diplom.git/pkg/luhn"
)

type OrdersService struct {
}

func NewOrdersService() *OrdersService {
	return &OrdersService{}
}

func (o *OrdersService) CreateOrder(ctx context.Context, number string) error {
	if !luhn.Valid(number) {
		return orders.ErrIncorrectOrderNumberFormat
	}

	return nil
}
