package service

import (
	"context"
	"fmt"

	"github.com/gopher-95/go-merch-shop/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, userName string, passwordHash string) (int, error)
	FindByUserName(ctx context.Context, username string) (*models.User, error)
}

type AuthService struct {
	userRepo UserRepository
	jwt      *JWTService
}

func NewAuthService(userRepo UserRepository, jwt *JWTService) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwt:      jwt,
	}
}

// Фнукция возвращает токен
func (authServ *AuthService) Authenticate(ctx context.Context, userName, password string) (string, error) {

	user, err := authServ.userRepo.FindByUserName(ctx, userName)
	if err != nil {
		return "", fmt.Errorf("ошибка поиска пользователя: %w", err)
	}

	if user == nil {
		hashedPassword, err := getHashedPassword(password)
		if err != nil {
			return "", fmt.Errorf("ошибка хеширования пароля: %w", err)
		}

		userID, err := authServ.userRepo.CreateUser(ctx, userName, hashedPassword)
		if err != nil {
			return "", fmt.Errorf("ошибка создания пользователя: %w", err)
		}

		token, err := authServ.jwt.GenerateToken(userID, userName)
		if err != nil {
			return "", fmt.Errorf("не удалось создать токен для пользователя: %w", err)
		}

		return token, nil
	}

	if !checkPassword(password, user.PasswordHash) {
		return "", fmt.Errorf("неверный пароль")
	}

	token, err := authServ.jwt.GenerateToken(user.ID, user.UserName)
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
