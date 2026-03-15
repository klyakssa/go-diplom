package auth

import "context"

type Repository interface {
	CreateUser(ctx context.Context, login, password string) error
	GetUserByLogin(ctx context.Context, login string) (*User, error)
}
