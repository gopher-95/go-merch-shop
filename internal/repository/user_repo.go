package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gopher-95/go-merch-shop/internal/models"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (storage *Storage) FindByUserName(ctx context.Context, username string) (*models.User, error) {
	query := "SELECT id, username, coins, created_at FROM users WHERE username = $1 "

	var user models.User

	err := storage.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.UserName,
		&user.PasswordHash,
		&user.Coins,
		&user.CreeatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователя по username: %w", err)
	}

	return &user, nil
}

func (storage *Storage) CreateUser(ctx context.Context, userName string, passwordHash string) (int, error) {
	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id"

	var id int

	err := storage.db.QueryRowContext(ctx, query, userName, passwordHash).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения id добавленного пользователя: %w", err)
	}
	return id, nil
}
