package auth

import "context"

//go:generate mockgen -source=repository.go -destination=mocks/repository.go

type Repository interface {
	CreateUser(ctx context.Context, login, password string) (string, error) // return user ID
	GetUserByLogin(ctx context.Context, login string) (*User, error)
}
