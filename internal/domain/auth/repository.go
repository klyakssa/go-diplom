package auth

import "context"

type Repository interface {
	CreateUser(ctx context.Context, login, password string) (string, error) // return user ID
	GetUserByLogin(ctx context.Context, login string) (*User, error)
}
