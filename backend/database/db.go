package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@db:5432/musicplayer?sslmode=disable"
	}

	var err error
	Pool, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %v", err)
	}

	if err = Pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	fmt.Println("Successfully connected to database")
	return nil
}

func CloseDB() {
	if Pool != nil {
		Pool.Close()
	}
}
