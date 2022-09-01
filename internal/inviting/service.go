package inviting

import (
	"context"
	"github.com/turbak/joom-calendar/internal/pkg/logger"
)

type Storage interface {
	// UpdateEventInviteStatus updates invite status to accepted or declined.
	// Then updates event attendees list correspondingly.
	UpdateEventInviteStatus(ctx context.Context, inviteID int, status string) error
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) AcceptInvite(ctx context.Context, inviteID int) error {
	err := s.storage.UpdateEventInviteStatus(ctx, inviteID, "accepted")
	if err != nil {
		return err
	}

	logger.Debugf("invite %d accepted", inviteID)

	return nil
}

func (s *Service) DeclineInvite(ctx context.Context, inviteID int) error {
	err := s.storage.UpdateEventInviteStatus(ctx, inviteID, "accepted")
	if err != nil {
		return err
	}

	logger.Debugf("invite %d declined", inviteID)

	return nil
}
