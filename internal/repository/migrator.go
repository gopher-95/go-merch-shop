package repository

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Функция запуска миграций
func RunMigrations(databaseURL string) error {
	const sourceURL = "file:///app/migrations"
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("ошибка создания мигратора: %w", err)
	}

	defer m.Close()

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("нет новых миграций")
			return nil
		}

		return fmt.Errorf("ошибка применения миграций: %w", err)
	}

	log.Println("все миграции успешно применены")

	return nil

}
