package postgres

import "time"

type Event struct {
	ID          int
	Title       string
	Description string
	Duration    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
