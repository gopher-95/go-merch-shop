package repository

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия соединения с бд: %w", err)
	}

	db.SetMaxOpenConns(50)
	db.SetConnMaxIdleTime(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ошибка пинга: %w", err)
	}

	return db, nil

}
