package main

import (
	"context"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB() *pgxpool.Pool {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://user:password@postgres:5432/splitwise_db"
	}
	pool, _ := pgxpool.New(context.Background(), connStr)
	return pool
}