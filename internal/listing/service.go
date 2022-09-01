package listing

import "context"

type Storage interface {
	GetEventByID(ctx context.Context, ID int) (*Event, error)
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
