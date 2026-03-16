package service

import (
	"context"

	"github.com/klyakssa/go-diplom.git/internal/domain/auth"
	"github.com/klyakssa/go-diplom.git/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo       auth.Repository
	jwtManager *jwt.JWTManager
}

func NewAuthService(repo auth.Repository, jwtManager *jwt.JWTManager) *AuthService {
	return &AuthService{repo: repo, jwtManager: jwtManager}
}

func (s *AuthService) Register(ctx context.Context, login, password string) (string, error) {
	if _, err := s.repo.GetUserByLogin(ctx, login); err == nil {
		return "", auth.ErrUserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	userid, err := s.repo.CreateUser(ctx, login, string(hash))
	if err != nil {
		return "", err
	}

	token, err := s.jwtManager.GenerateToken(userid)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.repo.GetUserByLogin(ctx, login)
	if err != nil {
		return "", auth.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", auth.ErrInvalidCredentials
	}

	token, err := s.jwtManager.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
