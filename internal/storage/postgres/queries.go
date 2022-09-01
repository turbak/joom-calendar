package postgres

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"time"
)

var pgQb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type Queries struct {
	querier Querier
}

type createUserParams struct {
	Name  string
	Email string
}

func (q Queries) CreateUser(ctx context.Context, params createUserParams) (int, error) {
	row := q.querier.QueryRow(ctx, "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", params.Name, params.Email)
	var ID int
	if err := row.Scan(&ID); err != nil {
		return 0, err
	}
	return ID, nil
}

func (q Queries) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := q.querier.QueryRow(ctx, "SELECT id, name, email, created_at, updated_at FROM users WHERE email = $1", email)

	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

type createEventParams struct {
	Title       string
	Description string
	Duration    int
}

func (q Queries) CreateEvent(ctx context.Context, params createEventParams) (int, error) {
	row := q.querier.QueryRow(ctx, `INSERT INTO events 
    			(title, description, duration) 
				VALUES ($1, $2, $3) RETURNING id`,
		params.Title, params.Description, params.Duration)

	var ID int
	if err := row.Scan(&ID); err != nil {
		return 0, err
	}
	return ID, nil
}

type createEventAttendeeParams struct {
	EventID int
	UserID  int
	Status  EventAttendeeStatus
}

func (q Queries) BatchCreateEventAttendees(ctx context.Context, params []createEventAttendeeParams) error {
	qb := pgQb.
		Insert("event_attendees").
		Columns("event_id", "user_id", "status")

	for _, attendee := range params {
		qb = qb.Values(attendee.EventID, attendee.UserID, attendee.Status)
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	_, err = q.querier.Exec(ctx, query, args...)

	return err
}

type createEventRepeatParams struct {
	EventID     int
	StartDate   time.Time
	DayOfWeek   string
	DayOfMonth  string
	MonthOfYear string
	WeekOfMonth string
}

func (q Queries) CreateEventRepeat(ctx context.Context, params createEventRepeatParams) error {
	const query = `INSERT INTO event_repeats
(event_id, repeat_start_date, day_of_week, day_of_month, month_of_year, week_of_month)
VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := q.querier.Exec(
		ctx,
		query,
		params.EventID,
		params.StartDate,
		params.DayOfWeek,
		params.DayOfMonth,
		params.MonthOfYear,
		params.WeekOfMonth,
	)
	if err != nil {
		return err
	}

	return err
}

func (q Queries) BatchGetUsersByIDs(ctx context.Context, IDs []int) ([]User, error) {
	query, args, err := pgQb.
		Select("id", "name", "email", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"id": IDs}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := q.querier.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (q Queries) GetEventByID(ctx context.Context, ID int) (*Event, error) {
	row := q.querier.QueryRow(ctx, "SELECT id, title, description, duration, created_at, updated_at FROM events WHERE id = $1", ID)

	var event Event
	if err := row.Scan(&event.ID, &event.Title, &event.Description, &event.Duration, &event.CreatedAt, &event.UpdatedAt); err != nil {
		return nil, err
	}

	return &event, nil
}

type createEventInviteParams struct {
	EventID int
	UserID  int
}

func (q Queries) BatchCreateEventInvites(ctx context.Context, params []createEventInviteParams) error {
	qb := pgQb.
		Insert("event_invites").
		Columns("event_id", "user_id")

	for _, invite := range params {
		qb = qb.Values(invite.EventID, invite.UserID)
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	_, err = q.querier.Exec(ctx, query, args...)

	return err
}

func (q Queries) UpdateEventInviteStatus(ctx context.Context, inviteID int, status EventInviteStatus) (*EventInvite, error) {
	const query = `UPDATE event_invites
SET status = $1,
    updated_at = NOW()
WHERE id = $2
RETURNING id, event_id, user_id, status, created_at, updated_at`

	row := q.querier.QueryRow(ctx, query, status, inviteID)

	var invite EventInvite
	if err := row.Scan(&invite.ID, &invite.EventID, &invite.UserID, &invite.Status, &invite.CreatedAt, &invite.UpdatedAt); err != nil {
		return nil, err
	}

	return &invite, nil
}

func (q Queries) UpdateEventAttendeeStatus(ctx context.Context, eventID int, userID int, status EventAttendeeStatus) (*EventAttendee, error) {
	const query = `UPDATE event_attendees
SET status = $1,
    updated_at = NOW()
WHERE event_id = $2
  AND user_id = $3
RETURNING event_id, user_id, status, created_at, updated_at`

	row := q.querier.QueryRow(ctx, query, status, eventID, userID)

	var attendee EventAttendee
	if err := row.Scan(&attendee.EventID, &attendee.UserID, &attendee.Status, &attendee.CreatedAt, &attendee.UpdatedAt); err != nil {
		return nil, err
	}

	return &attendee, nil
}

func (q Queries) DeleteEventAttendee(ctx context.Context, eventID int, userID int) error {
	const query = `DELETE FROM event_attendees
WHERE event_id = $1
  AND user_id = $2`

	_, err := q.querier.Exec(ctx, query, eventID, userID)
	return err
}
