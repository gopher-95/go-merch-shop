package models

import "time"

type User struct {
	ID           int       `json:"id"`
	UserName     string    `json:"user_name"`
	PasswordHash string    `json:"-"`
	Coins        int       `json:"coins"`
	CreeatedAt   time.Time `json:"created_at"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
