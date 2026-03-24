package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/klyakssa/go-diplom.git/internal/domain/auth"
	mock_auth "github.com/klyakssa/go-diplom.git/internal/domain/auth/mocks"
	"github.com/klyakssa/go-diplom.git/internal/service"
	"github.com/klyakssa/go-diplom.git/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {
	jwtManager := jwt.NewJWTManager("secret", time.Hour)

	type mockBehaviour func(service *mock_auth.MockRepository, login string)

	tests := []struct {
		name     string
		login    string
		password string
		moke     mockBehaviour
		error    error
	}{
		{
			name:     "success register",
			login:    "user",
			password: "password123",
			moke: func(service *mock_auth.MockRepository, login string) {
				service.EXPECT().GetUserByLogin(gomock.Any(), login).Return(nil, fmt.Errorf("user not found"))
				service.EXPECT().CreateUser(gomock.Any(), login, gomock.Any()).Return("1", nil)
			},
			error: nil,
		},
		{
			name:     "user already exists",
			login:    "user",
			password: "password123",
			moke: func(service *mock_auth.MockRepository, login string) {
				service.EXPECT().GetUserByLogin(gomock.Any(), login).Return(nil, nil)
			},
			error: auth.ErrUserAlreadyExists,
		},
		{
			name:     "password is too long",
			login:    "user",
			password: "FASFAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
			moke: func(service *mock_auth.MockRepository, login string) {
				service.EXPECT().GetUserByLogin(gomock.Any(), login).Return(nil, fmt.Errorf("error"))
			},
			error: auth.ErrPasswordTooLong,
		},
		{
			name:     "login is too long",
			login:    "usersadddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			password: "FASFA",
			moke:     func(service *mock_auth.MockRepository, login string) {},
			error:    auth.ErrLoginTooLong,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			authRepository := mock_auth.NewMockRepository(c)
			tt.moke(authRepository, tt.login)

			authService := service.NewAuthService(authRepository, jwtManager)

			token, err := authService.Register(context.Background(), tt.login, tt.password)

			assert.ErrorIs(t, err, tt.error, "got error %v, want %v", err, tt.error)

			if token != "" {
				_, err = jwtManager.VerifyToken(token)
				assert.NoError(t, err, "verify token error %v", err)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	jwtManager := jwt.NewJWTManager("secret", time.Hour)

	type mockBehaviour func(service *mock_auth.MockRepository, login string, password string)

	tests := []struct {
		name     string
		login    string
		password string
		moke     mockBehaviour
		error    error
	}{
		{
			name:     "success login",
			login:    "user",
			password: "password123",
			moke: func(service *mock_auth.MockRepository, login string, password string) {
				hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				service.EXPECT().GetUserByLogin(gomock.Any(), login).Return(&auth.User{ID: "1", Login: "user", Password: string(hash)}, nil)
			},
			error: nil,
		},
		{
			name:     "user not found",
			login:    "user",
			password: "password123",
			moke: func(service *mock_auth.MockRepository, login string, password string) {
				service.EXPECT().GetUserByLogin(gomock.Any(), login).Return(nil, fmt.Errorf("user not found"))
			},
			error: auth.ErrUserNotFound,
		},
		{
			name:     "invalid credentials",
			login:    "user",
			password: "password123",
			moke: func(service *mock_auth.MockRepository, login string, password string) {
				hash, _ := bcrypt.GenerateFromPassword([]byte("asd"), bcrypt.DefaultCost)
				service.EXPECT().GetUserByLogin(gomock.Any(), login).Return(&auth.User{ID: "1", Login: "user", Password: string(hash)}, nil)
			},
			error: auth.ErrInvalidCredentials,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			authRepository := mock_auth.NewMockRepository(c)
			tt.moke(authRepository, tt.login, tt.password)

			authService := service.NewAuthService(authRepository, jwtManager)

			token, err := authService.Login(context.Background(), tt.login, tt.password)

			assert.ErrorIs(t, err, tt.error, "got error %v, want %v", err, tt.error)

			if token != "" {
				_, err = jwtManager.VerifyToken(token)
				assert.NoError(t, err, "verify token error %v", err)
			}
		})
	}
}
