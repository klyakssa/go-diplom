package auth

import "context"

//go:generate mockgen -source=service.go -destination=mocks/service.go

type Service interface {
	Register(ctx context.Context, username, password string) (string, error)
	Login(ctx context.Context, username, password string) (string, error)
}
