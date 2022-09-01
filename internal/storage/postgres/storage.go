package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/turbak/joom-calendar/internal/inviting"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/turbak/joom-calendar/internal/creating"
	"github.com/turbak/joom-calendar/internal/listing"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

func (s *Storage) CreateUser(ctx context.Context, user creating.User) (int, error) {
	var createdID int

	err := s.withTx(ctx, func(q Queries) error {
		dbUser, err := q.GetUserByEmail(ctx, user.Email)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return err
			}
		}

		if dbUser != nil {
			return creating.ErrUserAlreadyExists
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

func (s *Storage) BatchGetUserByIDs(ctx context.Context, IDs []int) ([]listing.User, error) {
	users, err := Queries{querier: s.pool}.BatchGetUsersByIDs(ctx, IDs)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, listing.ErrUserNotFound
		}
		return nil, err
	}

	res := make([]listing.User, 0, len(users))
	for _, user := range users {
		res = append(res, listing.User{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return res, nil
}

func (s *Storage) CreateEvent(ctx context.Context, event creating.Event) (int, error) {
	var createdID int
	var err error

	err = s.withTx(ctx, func(q Queries) error {
		createdID, err = q.CreateEvent(ctx, createEventParams{
			Title:       event.Title,
			Description: event.Description,
			Duration:    event.Duration,
		})
		if err != nil {
			return err
		}

		eventAttendees := make([]createEventAttendeeParams, 0, len(event.InvitedUserIDs)+1)
		eventInvites := make([]createEventInviteParams, 0, len(event.InvitedUserIDs))
		for _, userID := range event.InvitedUserIDs {
			eventAttendees = append(eventAttendees, createEventAttendeeParams{
				EventID: createdID,
				UserID:  userID,
				Status:  EventAttendeeStatusUnconfirmed,
			})
			eventInvites = append(eventInvites, createEventInviteParams{
				EventID: createdID,
				UserID:  userID,
			})
		}

		eventAttendees = append(eventAttendees, createEventAttendeeParams{
			EventID: createdID,
			UserID:  event.OrganizerUserID,
			Status:  EventAttendeeStatusOrganizer,
		})

		err = q.BatchCreateEventInvites(ctx, eventInvites)
		if err != nil {
			return err
		}

		err = q.BatchCreateEventAttendees(ctx, eventAttendees)
		if err != nil {
			return err
		}

		repeat := createEventRepeatParams{
			StartDate: event.StartDate,
		}

		if event.Repeat != nil {
			repeat.EventID = createdID
			repeat.DayOfWeek = event.Repeat.DayOfWeek
			repeat.DayOfMonth = event.Repeat.DayOfMonth
			repeat.WeekOfMonth = event.Repeat.WeekOfMonth
			repeat.MonthOfYear = event.Repeat.MonthOfYear
		}

		return q.CreateEventRepeat(ctx, repeat)
	})

	return createdID, err
}

func (s *Storage) GetEventByID(ctx context.Context, ID int) (*listing.Event, error) {
	event, err := Queries{querier: s.pool}.GetEventByID(ctx, ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, listing.ErrEventNotFound
		}
		return nil, err
	}

	return &listing.Event{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Duration:    event.Duration,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}, nil
}

func (s *Storage) UpdateEventInviteStatus(ctx context.Context, inviteID int, status string) error {
	return s.withTx(ctx, func(q Queries) error {
		invite, err := q.UpdateEventInviteStatus(ctx, inviteID, EventInviteStatus(status))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return inviting.ErrInviteNotFound
			}
		}

		if invite.Status == EventInviteStatusAccepted {
			_, err = q.UpdateEventAttendeeStatus(ctx, invite.EventID, invite.UserID, EventAttendeeStatusConfirmed)
			if err != nil {
				return err
			}
		}

		if invite.Status == EventInviteStatusDeclined {
			err = q.DeleteEventAttendee(ctx, invite.EventID, invite.UserID)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Storage) withTx(ctx context.Context, f func(q Queries) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	q := Queries{querier: tx}

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
