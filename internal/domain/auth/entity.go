package auth

// User struct
type User struct {
	ID       string `json:"user_id" db:"id"`        // ID is primary key
	Login    string `json:"login" db:"login"`       // Login is unique
	Password string `json:"password" db:"password"` // Password is hashed
}
