package listing

import (
	"context"
	"time"
)

type Storage interface {
	GetEventByID(ctx context.Context, ID int) (*Event, error)
	ListUsersEvents(ctx context.Context, userID int, from, to time.Time) ([]Event, error)
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) GetEventByID(ctx context.Context, ID int) (*Event, error) {
	event, err := s.storage.GetEventByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *Service) ListUsersEvents(ctx context.Context, userID int, from, to time.Time) ([]Event, error) {
	events, err := s.storage.ListUsersEvents(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}

	j := 0
	for i := range events {
		if events[i].IsRepeated {
			isEventInPeriod := len(events[i].Rrule.Between(from, to, true)) > 0
			if isEventInPeriod {
				events[j] = events[i]
				j++
			}
		}
	}

	res := events[:j]
	if len(res) == 0 {
		return nil, ErrEventNotFound
	}

	return res, nil
}
