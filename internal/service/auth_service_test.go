package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/klyakssa/go-diplom.git/internal/domain/auth"
	"github.com/klyakssa/go-diplom.git/internal/service"
	"github.com/klyakssa/go-diplom.git/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) GetUserByLogin(ctx context.Context, login string) (*auth.User, error) {
	args := m.Called(ctx, login)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.User), args.Error(1)
}

func (m *MockAuthRepository) CreateUser(ctx context.Context, login, password string) (string, error) {
	args := m.Called(ctx, login, password)
	return args.String(0), args.Error(1)
}

func TestAuthService_Register(t *testing.T) {

	repo := new(MockAuthRepository)
	jwtManager := jwt.NewJWTManager("secret", time.Hour)

	s := service.NewAuthService(repo, jwtManager)

	tests := []struct {
		name     string
		login    string
		password string
		mock     func()
		wantErr  bool
	}{
		{
			name:     "success register",
			login:    "user",
			password: "password123",
			mock: func() {
				repo.On("GetUserByLogin", mock.Anything, "user").
					Return(nil, auth.ErrUserNotFound)

				repo.On("CreateUser", mock.Anything, "user", mock.Anything).
					Return("", nil)
			},
			wantErr: false,
		},
		{
			name:     "user already exists",
			login:    "user",
			password: "password123",
			mock: func() {
				repo.On("GetUserByLogin", mock.Anything, "user").
					Return(&auth.User{Login: "user"}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			repo.ExpectedCalls = nil
			tt.mock()

			token, err := s.Register(context.Background(), tt.login, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, token)

		})
	}
}

func TestAuthService_Login(t *testing.T) {

	repo := new(MockAuthRepository)
	jwtManager := jwt.NewJWTManager("secret", time.Hour)

	s := service.NewAuthService(repo, jwtManager)

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	tests := []struct {
		name     string
		login    string
		password string
		mock     func()
		wantErr  bool
	}{
		{
			name:     "success login",
			login:    "user",
			password: "password123",
			mock: func() {
				repo.On("GetUserByLogin", mock.Anything, "user").
					Return(&auth.User{
						Login:    "user",
						Password: string(hash),
					}, nil)
			},
			wantErr: false,
		},
		{
			name:     "invalid credentials",
			login:    "user",
			password: "wrong",
			mock: func() {
				repo.On("GetUserByLogin", mock.Anything, "user").
					Return(&auth.User{
						Login:    "user",
						Password: string(hash),
					}, nil)
			},
			wantErr: true,
		},
		{
			name:     "user not found",
			login:    "user",
			password: "password",
			mock: func() {
				repo.On("GetUserByLogin", mock.Anything, "user").
					Return(nil, auth.ErrUserNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			repo.ExpectedCalls = nil
			tt.mock()

			token, err := s.Login(context.Background(), tt.login, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, token)

		})
	}
}
