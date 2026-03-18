package orders

import "time"

type Order struct {
	Number      string    `json:"number" db:"number"`
	Status      string    `json:"status" db:"status"`
	Accrual     int       `json:"accrual,omitempty" db:"accrual"`
	UploadAt    time.Time `json:"uploaded_at" db:"uploaded_at"`
	UserID      string    `json:"-" db:"user_id"`
	IsAccrualed bool      `json:"-" db:"is_accrualed"`
}
