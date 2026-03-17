package orders

type Order struct {
	Number   string `json:"number" db:"number"`
	Status   int    `json:"status" db:"status"`
	Accrual  string `json:"accrual" db:"accrual"`
	UploadAt string `json:"uploaded_at" db:"uploaded_at"`
}
