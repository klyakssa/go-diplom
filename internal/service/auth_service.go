package service

import (
	"context"

	"github.com/klyakssa/go-diplom.git/internal/domain/auth"
)

type AuthService struct {
	repo auth.Repository
}

func NewAuthService(repo auth.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(ctx context.Context, login, password string) (string, error) {
	if _, err := s.repo.GetUserByLogin(ctx, login); err == nil {
		return "", auth.ErrUserAlreadyExists
	}
	if err := s.repo.CreateUser(ctx, login, password); err != nil {
		return "", err
	}

	return "mocked_token", nil
}

func (s *AuthService) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.repo.GetUserByLogin(ctx, login)
	if err != nil {
		return "", auth.ErrInvalidCredentials
	}
	if user.Password != password {
		return "", auth.ErrInvalidCredentials
	}
	return "mocked_token", nil
}
