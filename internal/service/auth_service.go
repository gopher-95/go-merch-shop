package service

import (
	"context"
	"fmt"

	"github.com/gopher-95/go-merch-shop/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage interface {
	CreateUser(ctx context.Context, username string, passwordHash string) (int, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
}

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

// Фнукция возвращает токен
func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {

	user, err := s.storage.FindByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("ошибка поиска пользователя: %w", err)
	}

	if user == nil {
		hashedPassword, err := getHashedPassword(password)
		if err != nil {
			return "", fmt.Errorf("ошибка хеширования пароля: %w", err)
		}

		userID, err := s.storage.CreateUser(ctx, username, hashedPassword)
		if err != nil {
			return "", fmt.Errorf("ошибка создания пользователя: %w", err)
		}

		token, err := s.jwt.GenerateToken(userID, username)
		if err != nil {
			return "", fmt.Errorf("не удалось создать токен для пользователя: %w", err)
		}

		return token, nil
	}

	if !checkPassword(password, user.PasswordHash) {
		return "", fmt.Errorf("неверный пароль")
	}

	token, err := s.jwt.GenerateToken(user.ID, user.UserName)
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
