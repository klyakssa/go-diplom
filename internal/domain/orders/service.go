package orders

import "context"

//go:generate mockgen -source=service.go -destination=mocks/service.go

// Service is an interface for orders service
type Service interface {
	CreateOrder(ctx context.Context, number, userID string) error  // CreateOrder creates new order
	GetOrders(ctx context.Context, userID string) ([]Order, error) // GetOrders returns all orders
}
