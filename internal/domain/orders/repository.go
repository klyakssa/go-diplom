package orders

import (
	"context"

	"github.com/shopspring/decimal"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository.go

// Repository is an interface for orders repository
type Repository interface {
	GetOrderByNumber(ctx context.Context, number string) (*Order, error)                   // get order by number
	CreateOrder(ctx context.Context, order *Order) error                                   // create order
	GetPendingOrders(ctx context.Context) ([]Order, error)                                 // get pending orders
	UpdateOrder(ctx context.Context, number, status string, accrual decimal.Decimal) error // update order
	ApplyAccrual(ctx context.Context, order *Order) error                                  // apply accrual
	GetOrdersByUserID(ctx context.Context, userID string) ([]Order, error)                 // get orders by userID
}
