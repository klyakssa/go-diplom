package orders

import "context"

type Service interface {
	CreateOrder(ctx context.Context, number, userID string) error
}
