package balance

import "context"

type Service interface {
	WithdrawBalance(ctx context.Context, userID string, orderNumber string, amount int) error
}
