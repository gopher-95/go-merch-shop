package service

import (
	"context"
	"fmt"

	"github.com/gopher-95/go-merch-shop/internal/models"
)

type SendStorage interface {
	TransferCoins(ctx context.Context, fromID, toID, amount int) error
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserBalance(ctx context.Context, userID int) (int, error)
}

// Сервис обмена монетами
type SendCoinsService struct {
	storage SendStorage
}

// Конструктор сервиса обмена монетами
func NewSendCoinsService(storage SendStorage) *SendCoinsService {
	return &SendCoinsService{
		storage: storage,
	}
}

func (s *SendCoinsService) SendCoins(ctx context.Context, fromUserID int, toUserName string, amount int) error {
	toUser, err := s.storage.FindByUsername(ctx, toUserName)
	if err != nil {
		return fmt.Errorf("получатель не найден: %w", err)
	}

	if toUser == nil {
		return fmt.Errorf("пользователь не найден")
	}

	if fromUserID == toUser.ID {
		return fmt.Errorf("нельзя отправить монеты самому себе")
	}

	if amount <= 0 {
		return fmt.Errorf("сумма не может быть отрицательной либо равной нулю")
	}

	balance, err := s.storage.GetUserBalance(ctx, fromUserID)
	if err != nil {
		return fmt.Errorf("не удалось получить баланс отправителя: %w", err)
	}

	if balance < amount {
		return fmt.Errorf("у пользователя не хватает баланса на счету")
	}

	err = s.storage.TransferCoins(ctx, fromUserID, toUser.ID, amount)
	if err != nil {
		return fmt.Errorf("не удалось совершить транзакцию: %w", err)
	}

	return nil
}
