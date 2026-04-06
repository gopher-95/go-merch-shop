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

	signedToken, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", fmt.Errorf("ошибка генерации токена: %w", err)
	}

	return signedToken, nil
}

func (j *JWT) getSecret(token *jwt.Token) (interface{}, error) {
	return []byte(j.secret), nil
}

func (j *JWT) Validate(tokenString string) (int, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, j.getSecret)
	if err != nil || !token.Valid {
		return 0, fmt.Errorf("невалидный токен: %w", err)
	}

	userIDGFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id не найден в токене")
	}

	return int(userIDGFloat), nil
}
