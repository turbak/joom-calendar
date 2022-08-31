package postgres

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type Execer interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type Queries struct {
	execer Execer
}

type createUserArgs struct {
	Name  string
	Email string
}

func (q Queries) CreateUser(ctx context.Context, args createUserArgs) (int, error) {
	row := q.execer.QueryRow(ctx, "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", args.Name, args.Email)
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (q Queries) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := q.execer.QueryRow(ctx, "SELECT id, name, email, created_at, updated_at FROM users WHERE email = $1", email)

	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}
