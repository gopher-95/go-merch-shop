package service

import (
	"context"
	"testing"

	"github.com/gopher-95/go-merch-shop/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockInfoStorage struct {
	mock.Mock
}

func (m *mockInfoStorage) GetUserBalance(ctx context.Context, userID int) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *mockInfoStorage) GetUserInventory(ctx context.Context, userID int) ([]models.InventoryItem, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.InventoryItem), args.Error(1)
}

func (m *mockInfoStorage) GetReceivedTransactions(ctx context.Context, userID int) ([]models.ReceivedTransaction, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.ReceivedTransaction), args.Error(1)
}

func (m *mockInfoStorage) GetSentTransactions(ctx context.Context, userID int) ([]models.SentTransaction, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.SentTransaction), args.Error(1)
}

func TestInfoService_GetUserInfo_Success(t *testing.T) {
	storage := new(mockInfoStorage)
	service := NewInfoService(storage)

	storage.On("GetUserBalance", mock.Anything, 1).Return(1000, nil)
	storage.On("GetUserInventory", mock.Anything, 1).Return([]models.InventoryItem{}, nil)
	storage.On("GetReceivedTransactions", mock.Anything, 1).Return([]models.ReceivedTransaction{}, nil)
	storage.On("GetSentTransactions", mock.Anything, 1).Return([]models.SentTransaction{}, nil)

	info, err := service.GetUserInfo(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, 1000, info.Coins)
	assert.NotNil(t, info.Inventory)
	assert.NotNil(t, info.CoinHistory.Received)
	assert.NotNil(t, info.CoinHistory.Sent)
	storage.AssertExpectations(t)
}
