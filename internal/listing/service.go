package listing

import "context"

type Storage interface {
	GetUserByID(ctx context.Context, id int) (*User, error)
}
