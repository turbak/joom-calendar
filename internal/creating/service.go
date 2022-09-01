package creating

import (
	"context"
	"github.com/turbak/joom-calendar/internal/listing"
	"github.com/turbak/joom-calendar/internal/pkg/logger"
)

type Storage interface {
	CreateUser(ctx context.Context, user User) (int, error)
	CreateEvent(ctx context.Context, event Event) (int, error)

	BatchGetUserByIDs(ctx context.Context, IDs []int) ([]listing.User, error)
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

func (s *Service) CreateEvent(ctx context.Context, event Event) (int, error) {
	users, err := s.storage.BatchGetUserByIDs(ctx, append(event.InvitedUserIDs, event.OrganizerUserID))
	if err != nil {
		return 0, err
	}

	if len(users) != len(event.InvitedUserIDs)+1 {
		return 0, ErrSomeUsersNotFound
	}

	createdID, err := s.storage.CreateEvent(ctx, event)
	if err != nil {
		return 0, err
	}

	logger.Debugf("event %d created successfully", createdID)

	return createdID, nil
}
