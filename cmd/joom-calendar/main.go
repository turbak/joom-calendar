package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/turbak/joom-calendar/internal/adding"
	"github.com/turbak/joom-calendar/internal/app"
	"github.com/turbak/joom-calendar/internal/pkg/logger"
	"github.com/turbak/joom-calendar/internal/storage/postgres"
	"github.com/xlab/closer"
	"os"
)

func main() {
	dbpool, err := createDatabasePool()
	if err != nil {
		logger.Fatalf("failed to create database pool: %v", err)
	}

	storage := postgres.NewStorage(dbpool)

	addingSvc := adding.NewService(storage)

	err = app.
		New(addingSvc).
		Run(":" + os.Getenv("PORT"))
	if err != nil {
		logger.Fatalf("failed to run app: %v", err)
	}
}

func createDatabasePool() (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	closer.Bind(dbpool.Close)

	return dbpool, nil
}
