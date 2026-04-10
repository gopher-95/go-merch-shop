package repository

import (
	"context"
	"fmt"

	"github.com/gopher-95/go-merch-shop/internal/models"
)

// Функция получения инвентаря пользователя
func (r *Repository) GetUserInventory(ctx context.Context, userID int) ([]models.InventoryItem, error) {
	var (
		inventoryItem  models.InventoryItem
		inventoryItems []models.InventoryItem
	)
	query := "SELECT item_name, quantity from inventory WHERE user_id = $1"

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения инвентаря: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&inventoryItem.Type, &inventoryItem.Quantity); err != nil {
			return nil, fmt.Errorf("не удалось просканировать результат: %w", err)
		}

		inventoryItems = append(inventoryItems, inventoryItem)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %w", err)
	}

	return inventoryItems, nil
}

// Функция получения транзакций
func (r *Repository) GetReceivedTransactions(ctx context.Context, userID int) ([]models.ReceivedTransaction, error) {
	var (
		recievedTransaction  models.ReceivedTransaction
		recievedTransactions []models.ReceivedTransaction
	)
	query := `SELECT u.username, t.amount 
			  FROM transactions t
			  JOIN users u ON u.id =t.from_user_id
			  WHERE t.to_user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения транзакции: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&recievedTransaction.FromUser, &recievedTransaction.Amount); err != nil {
			return nil, fmt.Errorf("не удалось просканировать результат: %w", err)
		}

		recievedTransactions = append(recievedTransactions, recievedTransaction)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %w", err)
	}

	return recievedTransactions, nil
}

// Функция отправки транзакций
func (r *Repository) GetSentTransactions(ctx context.Context, userID int) ([]models.SentTransaction, error) {
	var (
		sentTransaction  models.SentTransaction
		sentTransactions []models.SentTransaction
	)

	query := `SELECT u.username, t.amount
			  FROM transactions t
			  JOIN users u ON t.to_user_id = u.id
			  WHERE t.from_user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения транзакции: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&sentTransaction.ToUser, &sentTransaction.Amount); err != nil {
			return nil, fmt.Errorf("не удалось просканировать результат: %w", err)
		}

		sentTransactions = append(sentTransactions, sentTransaction)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %w", err)
	}

	return sentTransactions, nil

}
