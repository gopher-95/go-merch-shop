package main

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/gopher-95/go-merch-shop/internal/config"
	"github.com/gopher-95/go-merch-shop/internal/handlers"
	"github.com/gopher-95/go-merch-shop/internal/middleware"
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

	//  Соединение с бд
	db, err := repository.NewDB(cfg.DatabaseURLString())
	if err != nil {
		log.Fatalf("ошибка соединения с бд: %v", err)
	}

	// Репозиторный слой
	repo := repository.NewRepository(db)

	// Сервисный слой
	jwt := service.NewJWT(cfg.JWTSecret)
	authService := service.NewAuthService(repo, jwt)
	buyService := service.NewBuyService(repo)

	// HTTP слой
	authHandler := handlers.NewAuthHanlder(authService)
	buyHandler := handlers.NewBuyHandler(buyService)

	// MiddleWare
	authMiddleware := middleware.NewAuthMiddleware(jwt)

	// Роутер
	router := chi.NewRouter()

	// Маршрут без авторизации
	router.Post("/api/auth", authHandler.Login)

	// Маршрут с авторизацией
	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Handler)
		r.Get("/api/buy/{item}", buyHandler.BuyMerch)
	})

	// Запускаем сервер
	server := server.NewServer(cfg.ServerPort, router)
	if err := server.Run(); err != nil {
		log.Fatal("не удалось запустить сервер")
	}
}
