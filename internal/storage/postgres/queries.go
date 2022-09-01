package postgres

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"time"
)

func init() {
	squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

type Execer interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type Queries struct {
	execer Execer
}

type createUserParams struct {
	Name  string
	Email string
}

func (q Queries) CreateUser(ctx context.Context, params createUserParams) (int, error) {
	row := q.execer.QueryRow(ctx, "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", params.Name, params.Email)
	var ID int
	if err := row.Scan(&ID); err != nil {
		return 0, err
	}
	return ID, nil
}

func (q Queries) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := q.execer.QueryRow(ctx, "SELECT id, name, email, created_at, updated_at FROM users WHERE email = $1", email)

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
	row := q.execer.QueryRow(ctx, `INSERT INTO events 
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
	//TODO: add special type for Status
	Status string
}

func (q Queries) BatchCreateEventAttendees(ctx context.Context, params []createEventAttendeeParams) error {
	qb := squirrel.Insert("event_attendees").
		Columns("event_id", "user_id", "status")

	for _, attendee := range params {
		qb = qb.Values(attendee.EventID, attendee.UserID, attendee.Status)
	}

	sql, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	_, err = q.execer.Exec(ctx, sql, args...)

	return err
}

type createEventRepeatParams struct {
	EventID     int
	StartDate   time.Time
	DaysOfWeek  []int
	DayOfMonth  int
	MonthOfYear int
}

func (q Queries) CreateEventRepeat(ctx context.Context, params createEventRepeatParams) error {
	const query = `INSERT INTO event_repeats
    			(event_id, start_date, days_of_week, day_of_month, month_of_year)
    				VALUES ($1, $2, $3, $4, $5)`

	_, err := q.execer.Exec(ctx, query, params.EventID, params.StartDate, params.DaysOfWeek, params.DayOfMonth, params.MonthOfYear)
	if err != nil {
		return err
	}

	return err
}
