package balance

type Balance struct {
	ID      int `json:"id" db:"id"`
	Current int `json:"current" db:"current"`
	UserID  int `json:"user_id" db:"user_id"`
}
