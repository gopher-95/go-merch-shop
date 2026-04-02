package service

import (
	"context"
	"fmt"
	"log"

	"github.com/gopher-95/go-merch-shop/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Сервис авторизации
type AuthService struct {
	storage UserStorage
	jwt     *JWT
}

// Конструктор сервиса авторизации
func NewAuthService(storage UserStorage, jwt *JWT) *AuthService {
	return &AuthService{
		storage: storage,
		jwt:     jwt,
	}
}

type UserStorage interface {
	CreateUser(ctx context.Context, username string, passwordHash string) (int, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
}

// Фнукция возвращает токен
func (a *AuthService) Login(ctx context.Context, username, password string) (string, error) {

	log.Printf("🔵 Login attempt: username=%s", username)
	user, err := a.storage.FindByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("ошибка поиска пользователя: %w", err)
	}

	if user == nil {
		hashedPassword, err := getHashedPassword(password)
		if err != nil {
			return "", fmt.Errorf("ошибка хеширования пароля: %w", err)
		}

		userID, err := a.storage.CreateUser(ctx, username, hashedPassword)
		if err != nil {
			return "", fmt.Errorf("ошибка создания пользователя: %w", err)
		}

		token, err := a.jwt.GenerateToken(userID, username)
		if err != nil {
			return "", fmt.Errorf("не удалось создать токен для пользователя: %w", err)
		}

		return token, nil
	}

	if !checkPassword(password, user.PasswordHash) {
		return "", fmt.Errorf("неверный пароль")
	}

	token, err := a.jwt.GenerateToken(user.ID, user.UserName)
	if err != nil {
		return "", fmt.Errorf("ошибка генерации токена: %w", err)
	}

	return token, nil

}

// Генерирует хэшированный пароль
func getHashedPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", fmt.Errorf("не удалось хешироват пароль: %w", err)
	}

	return string(hashedPass), nil
}

// Проверяет пароль пользователя
func checkPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}

	return true
}
