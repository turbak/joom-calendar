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
	DaysOfWeek  []int
	DayOfMonth  int
	MonthOfYear int
	WeekOfMonth int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
