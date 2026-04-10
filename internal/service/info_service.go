package service

import (
	"context"
	"fmt"

	"github.com/gopher-95/go-merch-shop/internal/models"
)

type InfoStorage interface {
	GetUserBalance(ctx context.Context, userID int) (int, error)
	GetUserInventory(ctx context.Context, userID int) ([]models.InventoryItem, error)
	GetReceivedTransactions(ctx context.Context, userID int) ([]models.ReceivedTransaction, error)
	GetSentTransactions(ctx context.Context, userID int) ([]models.SentTransaction, error)
}

type InfoService struct {
	storage InfoStorage
}

func NewInfoService(storage InfoStorage) *InfoService {
	return &InfoService{
		storage: storage,
	}
}

func (s *InfoService) GetUserInfo(ctx context.Context, userID int) (*models.InfoResponse, error) {
	var (
		info        models.InfoResponse
		coinHistory models.CoinHistory
	)

	balance, err := s.storage.GetUserBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить баланс пользователя: %w", err)
	}

	inventory, err := s.storage.GetUserInventory(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить инвентарь пользователя: %w", err)
	}

	if inventory == nil {
		inventory = []models.InventoryItem{}
	}

	receivedTran, err := s.storage.GetReceivedTransactions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("не удалось отобразить полученные переводы: %w", err)
	}

	if receivedTran == nil {
		receivedTran = []models.ReceivedTransaction{}
	}

	sentTran, err := s.storage.GetSentTransactions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить отправленные переводы: %w", err)
	}

	if sentTran == nil {
		sentTran = []models.SentTransaction{}
	}

	coinHistory = models.CoinHistory{
		Received: receivedTran,
		Sent:     sentTran,
	}

	info = models.InfoResponse{
		Coins:       balance,
		Inventory:   inventory,
		CoinHistory: coinHistory,
	}

	return &info, nil

}
