package postgres

import (
	"time"
)

type Event struct {
	ID          int
	Title       string
	Description string
	Duration    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FullEvent struct {
	ID          int
	Title       string
	Description string
	Duration    int
	Repeat      EventRepeat
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type EventRepeat struct {
	StartDate   time.Time
	DayOfWeek   string
	DayOfMonth  string
	MonthOfYear string
	WeekOfMonth string
}
