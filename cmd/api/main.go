package main

import (
	"log"

	"github.com/gopher-95/go-merch-shop/internal/config"
	"github.com/gopher-95/go-merch-shop/internal/repository"
)

func main() {
	cfg := config.LoadConf()

	err := repository.RunMigrations(cfg.DatabaseURLString())
	if err != nil {
		log.Fatalf("ошибка запуска миграций: %v", err)
	}
}
