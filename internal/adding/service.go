package adding

import (
	"context"
	"github.com/turbak/joom-calendar/internal/pkg/logger"
)

type User struct {
	Name  string
	Email string
}

type Storage interface {
	CreateUser(ctx context.Context, args User) (int, error)
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) CreateUser(ctx context.Context, user User) (int, error) {
	createdID, err := s.storage.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	logger.Debugf("user %d created successfully", createdID)

	return createdID, nil
}
