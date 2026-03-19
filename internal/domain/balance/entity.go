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
	ProcessedAt TimeRFC3339     `json:"processed_at" db:"processed_at"`
}

type TimeRFC3339 time.Time

func (t TimeRFC3339) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).Format(time.RFC3339) + `"`), nil
}
