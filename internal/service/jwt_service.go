package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secret string
}

func NewJWT(secret string) *JWT {
	return &JWT{
		secret: secret,
	}
}

func (j *JWT) GenerateToken(userID int, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", fmt.Errorf("ошибка генерации токена: %w", err)
	}

	return tokenString, nil
}
