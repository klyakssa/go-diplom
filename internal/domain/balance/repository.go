package balance

import "context"

type Repository interface {
	AddBalance(ctx context.Context, userID string, amount int) error
}
