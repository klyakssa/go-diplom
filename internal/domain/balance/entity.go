package balance

import (
	"time"

	"github.com/shopspring/decimal"
)

// Balance is a struct for balance
type Balance struct {
	ID      int             `json:"id" db:"id"`           // ID of balance
	Current decimal.Decimal `json:"current" db:"current"` // Current balance
	UserID  int             `json:"user_id" db:"user_id"` // ID of user
}

// WithdrawHistory is a struct for withdraw history
type WithdrawHistory struct {
	OrderNumber string          `json:"order" db:"number"`              // Order number
	Sum         decimal.Decimal `json:"sum" db:"sum"`                   // Sum
	UserID      int             `json:"-" db:"user_id"`                 // ID of user
	ProcessedAt TimeRFC3339     `json:"processed_at" db:"processed_at"` // Processed at
}

// TimeRFC3339 is a struct for time
type TimeRFC3339 time.Time

func (t TimeRFC3339) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).Format(time.RFC3339) + `"`), nil
}
