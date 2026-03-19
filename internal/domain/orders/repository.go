package orders

import (
	"context"

	"github.com/shopspring/decimal"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository.go

type Repository interface {
	GetOrderByNumber(ctx context.Context, number string) (*Order, error)
	CreateOrder(ctx context.Context, order *Order) error
	GetPendingOrders(ctx context.Context) ([]Order, error)
	UpdateOrder(ctx context.Context, number, status string, accrual decimal.Decimal) error
	ApplyAccrual(ctx context.Context, order *Order) error
	GetOrdersByUserID(ctx context.Context, userID string) ([]Order, error)
}
