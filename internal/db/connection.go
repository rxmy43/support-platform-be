package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func Connect(dsn string) (*sqlx.DB, error) {
	var err error

	DB, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	log.Println("Database connected successfully")
	return DB, nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}

	return nil
}
