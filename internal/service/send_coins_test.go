package service

import (
	"context"
	"testing"

	"github.com/gopher-95/go-merch-shop/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSendStorage struct {
	mock.Mock
}

func (m *mockSendStorage) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockSendStorage) GetUserBalance(ctx context.Context, userID int) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *mockSendStorage) TransferCoins(ctx context.Context, fromID, toID, amount int) error {
	args := m.Called(ctx, fromID, toID, amount)
	return args.Error(0)
}

func TestSendCoinsService_SendCoins_Success(t *testing.T) {
	storage := new(mockSendStorage)
	service := NewSendCoinsService(storage)

	toUser := &models.User{ID: 2, UserName: "ivan"}

	storage.On("FindByUsername", mock.Anything, "ivan").Return(toUser, nil)
	storage.On("GetUserBalance", mock.Anything, 1).Return(1000, nil)
	storage.On("TransferCoins", mock.Anything, 1, 2, 100).Return(nil)

	err := service.SendCoins(context.Background(), 1, "ivan", 100)

	assert.NoError(t, err)
	storage.AssertExpectations(t)
}

func TestSendCoinsService_SendCoins_SelfTransfer(t *testing.T) {
	storage := new(mockSendStorage)
	service := NewSendCoinsService(storage)

	toUser := &models.User{ID: 1, UserName: "alex"}

	storage.On("FindByUsername", mock.Anything, "alex").Return(toUser, nil)

	err := service.SendCoins(context.Background(), 1, "alex", 100)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "нельзя отправить монеты самому себе")
	storage.AssertExpectations(t)
}

func TestSendCoinsService_SendCoins_UserNotFound(t *testing.T) {
	storage := new(mockSendStorage)
	service := NewSendCoinsService(storage)

	storage.On("FindByUsername", mock.Anything, "unknown").Return(nil, nil)

	err := service.SendCoins(context.Background(), 1, "unknown", 100)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "пользователь не найден")
	storage.AssertExpectations(t)
}
