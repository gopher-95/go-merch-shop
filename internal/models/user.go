package models

import "time"

// Структура описывает пользователя
type User struct {
	ID           int       `json:"id"`
	UserName     string    `json:"user_name"`
	PasswordHash string    `json:"-"`
	Coins        int       `json:"coins"`
	CreeatedAt   time.Time `json:"created_at"`
}

// Структура описывает запрос пользователя на его создание
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Структура описывает ответ пользователю в виде токена
type AuthResponse struct {
	Token string `json:"token"`
}
