package repository

import (
	"context"
	"fmt"
)

func (r *Repository) GetUserBalance(ctx context.Context, userID int) (int, error) {
	query := "SELECT coins FROM users WHERE id = $1"

	var coins int

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&coins)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения баланса пользователя: %w", err)
	}

	return coins, nil
}

func (r *Repository) WithdrawCoins(ctx context.Context, userID int, amount int) error {
	query := "UPDATE users SET coins = coins - $1 WHERE id = $2 AND coins >= $1"

	res, err := r.db.ExecContext(ctx, query, amount, userID)
	if err != nil {
		return fmt.Errorf("ошибка списания монет: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("недостаточно монет")
	}

	return nil

}

func (r *Repository) AddToInventory(ctx context.Context, userID int, itemName string) error {
	query := `INSERT INTO inventory (user_id, item_name, quantity)
	          VALUES ($1, $2, 1)
			  ON CONFLICT (user_id, item_name)
			  DO UPDATE SET quantity = inventory.quantity + 1`

	_, err := r.db.ExecContext(ctx, query, userID, itemName)
	if err != nil {
		return fmt.Errorf("ошибка добавления товара в инвентарь: %w", err)
	}

	return nil
}
