package postgres

import (
	"context"
	"fmt"
	"github.com/turbak/joom-calendar/internal/adding"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

func (s *Storage) CreateUser(ctx context.Context, user adding.User) (int, error) {
	var createdID int

	err := s.withTx(ctx, func(q Queries) error {
		dbUser, err := q.GetUserByEmail(ctx, user.Email)
		if err != nil {
			return err
		}

		if dbUser != nil {
			return adding.ErrUserAlreadyExists
		}

		createdID, err = q.CreateUser(ctx, createUserParams{
			Name:  user.Name,
			Email: user.Email,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return createdID, err
}

func (s *Storage) withTx(ctx context.Context, f func(q Queries) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	q := Queries{execer: tx}

	if err = f(q); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("failed to exec queries: %v", fmt.Errorf("failed to rollback transaction: %v", rbErr))
		}
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
