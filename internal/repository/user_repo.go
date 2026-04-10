package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gopher-95/go-merch-shop/internal/models"
)

// Репозиторий содержит БД
type Repository struct {
	db *sql.DB
}

// Конструктор репозитория
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Функция находит в бд пользователя по имени
func (r *Repository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	query := "SELECT id, username, password_hash, coins, created_at FROM users WHERE username = $1 "

	var user models.User

	err := r.db.QueryRowContext(ctx, query, username).Scan(
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

// Функция создает пользователя
func (r *Repository) CreateUser(ctx context.Context, username string, passwordHash string) (int, error) {
	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id"

	var id int

	err := r.db.QueryRowContext(ctx, query, username, passwordHash).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения id добавленного пользователя: %w", err)
	}
	return id, nil
}
