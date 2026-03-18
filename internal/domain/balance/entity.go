package balance

import (
	"time"

	"github.com/shopspring/decimal"
)

type Balance struct {
	ID      int             `json:"id" db:"id"`
	Current decimal.Decimal `json:"current" db:"current"`
	UserID  int             `json:"user_id" db:"user_id"`
}

type WithdrawHistory struct {
	OrderNumber string          `json:"order" db:"number"`
	Sum         decimal.Decimal `json:"sum" db:"sum"`
	UserID      int             `json:"-" db:"user_id"`
	ProcessedAt time.Time       `json:"processed_at" db:"processed_at"`
}
