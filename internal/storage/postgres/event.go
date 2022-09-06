package postgres

import (
	"time"
)

type Event struct {
	ID          int
	Title       string
	Description string
	Duration    int
	StartDate   time.Time
	Rrule       string
	IsRepeated  bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
