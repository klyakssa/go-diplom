package orders

import "context"

//go:generate mockgen -source=service.go -destination=mocks/service.go

type Service interface {
	CreateOrder(ctx context.Context, number, userID string) error
	GetOrders(ctx context.Context, userID string) ([]Order, error)
}
