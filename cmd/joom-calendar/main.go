package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/turbak/joom-calendar/internal/app"
	"github.com/turbak/joom-calendar/internal/auth"
	"github.com/turbak/joom-calendar/internal/creating"
	"github.com/turbak/joom-calendar/internal/inviting"
	"github.com/turbak/joom-calendar/internal/listing"
	"github.com/turbak/joom-calendar/internal/pkg/client/github"
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

	addingSvc := creating.NewService(storage)
	listingSvc := listing.NewService(storage)
	invitingSvc := inviting.NewService(storage)
	authSvc := auth.NewService(
		storage,
		github.NewClient(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET")),
		[]byte(os.Getenv("JWT_KEY")),
	)

	err = app.
		New(addingSvc, listingSvc, invitingSvc, authSvc).
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
