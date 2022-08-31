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
