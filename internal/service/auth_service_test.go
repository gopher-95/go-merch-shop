package service

import (
	"context"
	"testing"

	"github.com/gopher-95/go-merch-shop/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserStorage struct {
	mock.Mock
}

func (m *mockUserStorage) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserStorage) CreateUser(ctx context.Context, username, passwordHash string) (int, error) {
	args := m.Called(ctx, username, passwordHash)
	return args.Int(0), args.Error(1)
}

func TestAuthService_Login_NewUser(t *testing.T) {
	storage := new(mockUserStorage)
	jwt := NewJWT("test-secret")
	auth := NewAuthService(storage, jwt)

	storage.On("FindByUsername", mock.Anything, "alex").Return(nil, nil)
	storage.On("CreateUser", mock.Anything, "alex", mock.Anything).Return(1, nil)

	token, err := auth.Login(context.Background(), "alex", "123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	storage.AssertExpectations(t)
}

func TestAuthService_Login_ExistingUser_WrongPassword(t *testing.T) {
	storage := new(mockUserStorage)
	jwt := NewJWT("test-secret")
	auth := NewAuthService(storage, jwt)

	existingUser := &models.User{
		ID:           1,
		UserName:     "alex",
		PasswordHash: "$2a$10$hashedpasswordfrombcrypt",
		Coins:        1000,
	}

	storage.On("FindByUsername", mock.Anything, "alex").Return(existingUser, nil)

	token, err := auth.Login(context.Background(), "alex", "wrongpassword")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "неверный пароль")

	storage.AssertExpectations(t)
}
