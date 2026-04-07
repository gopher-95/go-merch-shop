package repository

import (
	"context"
	"fmt"
)

func (r *Repository) TransferCoins(ctx context.Context, fromID, toID, amount int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// Списываем монеты у отправителя
	withDrawQuery := "UPDATE users SET coins = coins - $1 WHERE id = $2 AND coins >= $1"

	res, err := tx.ExecContext(ctx, withDrawQuery, amount, fromID)
	if err != nil {
		return fmt.Errorf("ошибка списания монет: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения измененных строк: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("недостаточно монет для списания")
	}

	// Зачисляем монеты получателю
	addCoinsQuery := "UPDATE users SET coins = coins + $1 WHERE id = $2"
	_, err = tx.ExecContext(ctx, addCoinsQuery, amount, toID)
	if err != nil {
		return fmt.Errorf("не удалось начислить монеты получателю: %w", err)
	}

	// Записываем в таблицу транзакций
	transactionQuery := "INSERT INTO transactions (from_user_id, to_user_id, amount) VALUES ($1, $2, $3)"
	_, err = tx.ExecContext(ctx, transactionQuery, fromID, toID, amount)
	if err != nil {
		return fmt.Errorf("ошибка записи транзакций: %w", err)
	}

	return tx.Commit()

}
