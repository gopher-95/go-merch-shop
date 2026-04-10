package service

import (
	"context"
	"fmt"

	"github.com/gopher-95/go-merch-shop/internal/models"
)

type BuyStorage interface {
	GetUserBalance(ctx context.Context, userID int) (int, error)
	WithdrawCoins(ctx context.Context, userID int, amount int) error
	AddToInventory(ctx context.Context, userID int, itemName string) error
}

// Cервис покупки мерча
type BuyService struct {
	storage BuyStorage
}

// Конструктор сервиса покупки мерча
func NewBuyService(storage BuyStorage) *BuyService {
	return &BuyService{
		storage: storage,
	}
}

// Функция покупки мерча
func (s *BuyService) BuyMerch(ctx context.Context, userID int, itemName string) error {
	price, ok := models.MerchCatalog[itemName]
	if !ok {
		return fmt.Errorf("не удалось найти мерч в каталоге")
	}

	balance, err := s.storage.GetUserBalance(ctx, userID)
	if err != nil {
		return fmt.Errorf("не удалось получить баланс пользователя: %w", err)
	}

	if balance < price {
		return fmt.Errorf("недостаточно монет")
	}

	err = s.storage.WithdrawCoins(ctx, userID, price)
	if err != nil {
		return fmt.Errorf("не удалось списать монеты: %w", err)
	}

	err = s.storage.AddToInventory(ctx, userID, itemName)
	if err != nil {
		return fmt.Errorf("не удалось добавить предмет в инвентарь: %w", err)
	}

	return nil
}
