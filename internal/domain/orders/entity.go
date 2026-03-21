package orders

import (
	"time"

	"github.com/shopspring/decimal"
)

// Order struct
type Order struct {
	Number      string          `json:"number" db:"number"`             // number of order
	Status      string          `json:"status" db:"status"`             // status of order
	Accrual     decimal.Decimal `json:"accrual,omitempty" db:"accrual"` // accrued amount
	UploadAt    time.Time       `json:"uploaded_at" db:"uploaded_at"`   // upload date
	UserID      string          `json:"-" db:"user_id"`                 // user id
	IsAccrualed bool            `json:"-" db:"is_accrualed"`            // is accrualed
}
