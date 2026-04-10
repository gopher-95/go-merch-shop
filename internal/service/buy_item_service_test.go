package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockBuyStorage struct {
	mock.Mock
}

func (m *mockBuyStorage) GetUserBalance(ctx context.Context, userID int) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *mockBuyStorage) WithdrawCoins(ctx context.Context, userID int, amount int) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

func (m *mockBuyStorage) AddToInventory(ctx context.Context, userID int, itemName string) error {
	args := m.Called(ctx, userID, itemName)
	return args.Error(0)
}

func TestBuyService_BuyMerch_Success(t *testing.T) {
	storage := new(mockBuyStorage)
	service := NewBuyService(storage)

	storage.On("GetUserBalance", mock.Anything, 1).Return(1000, nil)
	storage.On("WithdrawCoins", mock.Anything, 1, 80).Return(nil)
	storage.On("AddToInventory", mock.Anything, 1, "t-shirt").Return(nil)

	err := service.BuyMerch(context.Background(), 1, "t-shirt")

	assert.NoError(t, err)
	storage.AssertExpectations(t)
}

func TestBuyService_BuyMerch_InsufficientFunds(t *testing.T) {
	storage := new(mockBuyStorage)
	service := NewBuyService(storage)

	storage.On("GetUserBalance", mock.Anything, 1).Return(50, nil)

	err := service.BuyMerch(context.Background(), 1, "t-shirt")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "недостаточно монет")
	storage.AssertExpectations(t)
}

func TestBuyService_BuyMerch_ItemNotFound(t *testing.T) {
	storage := new(mockBuyStorage)
	service := NewBuyService(storage)

	err := service.BuyMerch(context.Background(), 1, "invalid_item")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "не удалось найти мерч в каталоге")
	storage.AssertExpectations(t)
}
