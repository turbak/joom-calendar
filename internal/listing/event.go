package listing

import "time"

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
	Attendees   []EventAttendee
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type EventAttendee struct {
	EventID   int
	Status    EventAttendeeStatus
	User      User
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EventAttendeeStatus string

const (
	EventAttendeeStatusUnconfirmed EventAttendeeStatus = "unconfirmed"
	EventAttendeeStatusConfirmed   EventAttendeeStatus = "confirmed"
	EventAttendeeStatusOrganizer   EventAttendeeStatus = "organizer"
)
