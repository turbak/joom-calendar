package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/turbak/joom-calendar/internal/creating"
	"github.com/turbak/joom-calendar/internal/inviting"
	"github.com/turbak/joom-calendar/internal/listing"
	"time"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

func (s *Storage) CreateUser(ctx context.Context, user creating.User) (int, error) {
	var createdID int

	err := s.withTx(ctx, pgx.TxOptions{}, func(q Queries) error {
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

	err := s.withTx(ctx, pgx.TxOptions{}, func(q Queries) error {
		insertedID, err := q.CreateEvent(ctx, createEventParams{
			Title:       event.Title,
			Description: event.Description,
			Duration:    event.Duration,
			StartDate:   event.StartDate,
			IsAllDay:    event.IsAllDay,
			IsRepeated:  event.Repeat != nil,
			Rrule:       event.Repeat.ToRrule(),
		})
		if err != nil {
			return err
		}

		createdID = insertedID

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

		if len(eventInvites) > 0 {
			err = q.BatchCreateEventInvites(ctx, eventInvites)
			if err != nil {
				return err
			}
		}

		return q.BatchCreateEventAttendees(ctx, eventAttendees)
	})

	return createdID, err
}

func (s *Storage) GetEventByID(ctx context.Context, ID int) (*listing.Event, error) {
	var event Event
	var attendees []FullEventAttendee

	err := s.withTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly}, func(q Queries) error {
		foundEvent, err := q.GetEventByID(ctx, ID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return listing.ErrEventNotFound
			}
			return err
		}

		event = *foundEvent

		attendees, err = q.BatchGetFullEventAttendees(ctx, event.ID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	res := toListingEvent(event, toListingEventAttendees(attendees))
	return &res, nil
}

func (s *Storage) UpdateEventInviteStatus(ctx context.Context, inviteID int, status string) error {
	return s.withTx(ctx, pgx.TxOptions{}, func(q Queries) error {
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

func (s *Storage) withTx(ctx context.Context, options pgx.TxOptions, f func(q Queries) error) error {
	return s.pool.BeginTxFunc(ctx, options, func(tx pgx.Tx) error {
		return f(Queries{querier: tx})
	})
}

func (s *Storage) ListUsersEvents(ctx context.Context, userID int, from, to time.Time) ([]listing.Event, error) {
	var events []Event
	var attendees []FullEventAttendee

	err := s.withTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly}, func(q Queries) error {
		var err error
		events, err = q.ListUsersEvents(ctx, from, to, userID)
		if err != nil {
			return err
		}

		attendees, err = q.BatchGetFullEventAttendees(ctx, pluckEventIDs(events)...)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return toListingEvents(events, attendees), nil
}

func (s *Storage) BatchGetEventsByUserIDs(ctx context.Context, userIDs []int) ([]listing.Event, error) {
	var events []Event
	var attendees []FullEventAttendee

	err := s.withTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly}, func(q Queries) error {
		var err error
		events, err = q.BatchGetEventsByUserIDs(ctx, userIDs)
		if err != nil {
			return err
		}

		attendees, err = q.BatchGetFullEventAttendees(ctx, pluckEventIDs(events)...)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return toListingEvents(events, attendees), nil
}

func pluckEventIDs(events []Event) []int {
	ids := make([]int, 0, len(events))
	for _, event := range events {
		ids = append(ids, event.ID)
	}
	return ids
}
