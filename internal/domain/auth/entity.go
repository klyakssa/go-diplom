package auth

type User struct {
	ID       string `json:"user_id" db:"user_id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
}
