package main

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/gopher-95/go-merch-shop/internal/config"
	"github.com/gopher-95/go-merch-shop/internal/handlers"
	"github.com/gopher-95/go-merch-shop/internal/repository"
	"github.com/gopher-95/go-merch-shop/internal/server"
	"github.com/gopher-95/go-merch-shop/internal/service"
)

func main() {
	// Получаем конфиг
	cfg := config.LoadConf()

	// Запускаем миграции
	err := repository.RunMigrations(cfg.DatabaseURLString())
	if err != nil {
		log.Fatalf("ошибка запуска миграций: %v", err)
	}

	// Создаем соединение с бд
	db, err := repository.NewDB(cfg.DatabaseURLString())
	if err != nil {
		log.Fatalf("ошибка соединения с бд: %v", err)
	}

	// Создаем репозиторный слой
	repo := repository.NewRepository(db)

	// Создаем JWT
	jwt := service.NewJWT(cfg.JWTSecret)

	// Создаем сервисный слой
	service := service.NewAuthService(repo, jwt)

	// Создаем слой HTTP
	authHandler := handlers.NewAuthHanlder(service)

	// Создаем роутер
	router := chi.NewRouter()

	// Регистрируем маршрут
	router.Post("/api/auth", authHandler.Login)

	// Запускаем сервер
	server := server.NewServer(cfg.ServerPort, router)
	if err := server.Run(); err != nil {
		log.Fatal("не удалось запустить сервер")
	}
}
